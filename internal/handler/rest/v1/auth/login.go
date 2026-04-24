package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

// loginRequest holds the credentials submitted by the user.
type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// loginUserResponse is the user sub-object returned inside the login payload.
type loginUserResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Image       string `json:"image,omitempty"`
	IsOnboarded bool   `json:"is_onboarded"`
}

// loginResponse is the full JSON body returned on successful login or registration.
type loginResponse struct {
	User         loginUserResponse `json:"user"`
	SessionToken string            `json:"session_token"`
}

// Login handles user login
// @Summary      User Login
// @Description  Authenticate user and return session token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "Login credentials"
// @Success      200      {object}  response.Response{data=loginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /auth/login [post]
func (h Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[Login] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	user, token, err := h.ctrl.Login(c.Request.Context(), auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: c.Request.UserAgent(),
		IPAddress: clientIP(c),
	})
	if err != nil {
		logger.ERROR.Printf("[Login] login failed for user %s: %+v", req.Email, err)
		if errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, response.NewResponse(
				constants.LoginFail.Code,
				constants.LoginFail.Message,
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

	logger.INFO.Printf("[Login] user logged in successfully: %s", user.Email)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.LoginSuccess.Code,
		constants.LoginSuccess.Message,
		loginResponse{
			User: loginUserResponse{
				ID:          user.ID,
				Name:        user.Name,
				Username:    user.Username,
				Email:       user.Email,
				Image:       user.Image,
				IsOnboarded: user.OnboardedAt != nil,
			},
			SessionToken: token,
		},
	))
}

func clientIP(c *gin.Context) string {
	ip := strings.TrimSpace(c.ClientIP())
	if ip == "" {
		return ""
	}
	return ip
}
