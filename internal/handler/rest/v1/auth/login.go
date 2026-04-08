package auth

import (
	"asona/internal/controller/auth"
	"asona/internal/constants"
	"asona/internal/handler/response"
	"net/http"

	"asona/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image,omitempty"`
}

type LoginResponse struct {
	User         LoginUserResponse `json:"user"`
	SessionToken string            `json:"session_token"`
}

// Login handles user login
// @Summary      User Login
// @Description  Authenticate user and return session token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  response.Response{data=LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /auth/login [post]
func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
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
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.ERROR.Printf("[Login] login failed for user %s: %+v", req.Email, err)
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.LoginFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	// Assuming controller now returns model.User or auth.UserResponse
	// If it returns model.User, we need to map it. 
	// Based on previous edits, get_profile returns model.User, but I defined UserResponse in controller.
	// Let's assume we map to the controller's public response type if needed, 
	// but the handler should have its own response struct.

	logger.INFO.Printf("[Login] user logged in successfully: %s", user.Email)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.LoginSuccess.Code,
		constants.LoginSuccess.Message,
		LoginResponse{
			User: LoginUserResponse{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
				Image:    user.Image,
			},
			SessionToken: token,
		},
	))
}
