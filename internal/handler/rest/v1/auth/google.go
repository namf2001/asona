package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"asona/config"
	"asona/internal/constants"
	authctrl "asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

const oauthStateSessionKey = "google_oauth_state"

// oauthCodePayload holds the token and onboarding state stored under a one-time code.
type oauthCodePayload struct {
	Token       string `json:"token"`
	IsOnboarded bool   `json:"is_onboarded"`
}

// exchangeOAuthCodeResponse is the JSON response body for the ExchangeOAuthCode endpoint.
type exchangeOAuthCodeResponse struct {
	SessionToken string `json:"session_token"`
	IsOnboarded  bool   `json:"is_onboarded"`
}

// GoogleLoginResponse represents the payload containing the OAuth URL
type GoogleLoginResponse struct {
	URL string `json:"url"`
}

// GoogleLogin starts the Google OAuth flow and returns the Google authorization URL.
// @Summary      Google OAuth Login
// @Description  Redirect the user to Google OAuth consent screen
// @Tags         auth
// @Produce      json
// @Success      200  {object} response.Response{data=GoogleLoginResponse}
// @Failure      500  {object} response.Response
// @Router       /auth/google [get]
func (h Handler) GoogleLogin(c *gin.Context) {
	state, err := randomState()
	if err != nil {
		logger.ERROR.Printf("[GoogleLogin] failed to generate oauth state: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	sess := sessions.Default(c)
	sess.Set(oauthStateSessionKey, state)
	if err := sess.Save(); err != nil {
		logger.ERROR.Printf("[GoogleLogin] failed to save oauth state: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	authURL, err := h.ctrl.GoogleAuthURL(c.Request.Context(), state)
	if err != nil {
		logger.ERROR.Printf("[GoogleLogin] failed to build auth url: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	logger.INFO.Printf("[GoogleLogin] returning google oauth url")
	c.JSON(http.StatusOK, response.NewResponse(
		constants.LoginSuccess.Code,
		constants.LoginSuccess.Message,
		GoogleLoginResponse{URL: authURL},
	))
}

// GoogleCallback completes the Google OAuth flow.
// Instead of embedding the JWT token in the redirect URL (which leaks it via
// browser history and server logs), it stores the token in Redis under a
// short-lived one-time authorization code, and redirects the frontend to
// exchange that code for the real token via a secure POST endpoint.
// @Summary      Google OAuth Callback
// @Description  Handle the Google OAuth callback and create a local session
// @Tags         auth
// @Produce      json
// @Success      302
// @Failure      400  {object} response.Response
// @Failure      401  {object} response.Response
// @Failure      500  {object} response.Response
// @Router       /auth/google/callback [get]
func (h Handler) GoogleCallback(c *gin.Context) {
	if oauthErr := c.Query("error"); oauthErr != "" {
		logger.ERROR.Printf("[GoogleCallback] oauth error: %s", oauthErr)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		logger.ERROR.Printf("[GoogleCallback] missing state or code")
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	sess := sessions.Default(c)
	savedState, _ := sess.Get(oauthStateSessionKey).(string)
	if savedState == "" || savedState != state {
		logger.ERROR.Printf("[GoogleCallback] invalid oauth state")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	sess.Delete(oauthStateSessionKey)
	_ = sess.Save()

	user, token, err := h.ctrl.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		logger.ERROR.Printf("[GoogleCallback] oauth callback failed: %+v", err)
		if errors.Is(err, authctrl.ErrOAuthEmailNotVerified) {
			c.JSON(http.StatusUnauthorized, response.NewResponse(
				constants.PermissionDenied.Code,
				constants.PermissionDenied.Message,
				nil,
			))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	logger.INFO.Printf("[GoogleCallback] user authenticated successfully: %s", user.Email)

	// Generate a one-time authorization code (not the JWT token itself).
	// The JWT token is stored in Redis behind this code for 60 seconds.
	// The frontend must exchange this code for the real token via POST /api/v1/auth/exchange.
	// This prevents the JWT from appearing in browser history, server logs, or Referer headers.
	authCode, err := randomState()
	if err != nil {
		logger.ERROR.Printf("[GoogleCallback] failed to generate auth code: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	payload := oauthCodePayload{
		Token:       token,
		IsOnboarded: user.OnboardedAt != nil,
	}
	payloadBytes, _ := json.Marshal(payload)
	redisKey := constants.OAuthCodePrefix + authCode

	// Store payload with 60-second TTL — short enough to prevent replay attacks.
	if err := h.rdb.Set(c.Request.Context(), redisKey, string(payloadBytes), 60*time.Second); err != nil {
		logger.ERROR.Printf("[GoogleCallback] failed to store oauth code in redis: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	frontendURL, err := url.Parse(config.GetConfig().FrontendURL)
	if err != nil {
		logger.ERROR.Printf("[GoogleCallback] invalid frontend url config: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	// Redirect with the one-time code (not the JWT token).
	callbackURL := frontendURL.ResolveReference(&url.URL{Path: "/api/auth/callback"})
	query := callbackURL.Query()
	query.Set("code", authCode)
	callbackURL.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, callbackURL.String())
}

// ExchangeOAuthCode exchanges a one-time authorization code for a JWT session token.
// The code is stored in Redis for 60 seconds and deleted immediately after use.
// @Summary      Exchange OAuth Code
// @Description  Exchange a one-time OAuth authorization code for a session token
// @Tags         auth
// @Produce      json
// @Param        code  query  string  true  "One-time authorization code"
// @Success      200  {object} response.Response{data=LoginResponse}
// @Failure      400  {object} response.Response
// @Failure      401  {object} response.Response
// @Router       /auth/exchange [post]
func (h Handler) ExchangeOAuthCode(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[ExchangeOAuthCode] missing code: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	redisKey := constants.OAuthCodePrefix + req.Code
	payloadStr, err := h.rdb.Get(c.Request.Context(), redisKey)
	if err != nil {
		logger.ERROR.Printf("[ExchangeOAuthCode] code not found or expired: %+v", err)
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	// Delete immediately — one-time use only.
	_ = h.rdb.Del(c.Request.Context(), redisKey)

	var payload oauthCodePayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		logger.ERROR.Printf("[ExchangeOAuthCode] failed to parse payload: %+v", err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	logger.INFO.Printf("[ExchangeOAuthCode] oauth code exchanged successfully")
	c.JSON(http.StatusOK, response.NewResponse(
		constants.LoginSuccess.Code,
		constants.LoginSuccess.Message,
		exchangeOAuthCodeResponse{
			SessionToken: payload.Token,
			IsOnboarded:  payload.IsOnboarded,
		},
	))
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
