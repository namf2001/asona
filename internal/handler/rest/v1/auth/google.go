package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"asona/config"
	"asona/internal/constants"
	authctrl "asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

const oauthStateSessionKey = "google_oauth_state"

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

// GoogleCallback completes the Google OAuth flow and returns the local session token.
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

	isOnboarded := "false"
	if user.OnboardedAt != nil {
		isOnboarded = "true"
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

	callbackURL := frontendURL.ResolveReference(&url.URL{Path: "/api/auth/callback"})
	query := callbackURL.Query()
	query.Set("token", token)
	query.Set("is_onboarded", isOnboarded)
	callbackURL.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, callbackURL.String())
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
