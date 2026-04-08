package auth

import (
	"asona/internal/controller/auth"
	"asona/internal/constants"
	"asona/internal/handler/response"
	"net/http"

	"asona/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image,omitempty"`
}

// Register handles user registration
// @Summary      User Registration
// @Description  Register a new user in the system
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration details"
// @Success      201      {object}  response.Response{data=RegisterResponse}
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/register [post]
func (h Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[Register] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	user, err := h.ctrl.Register(c.Request.Context(), auth.RegisterInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.ERROR.Printf("[Register] registration failed for user %s: %+v", req.Email, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.RegisterUserFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[Register] user registered successfully: %s", user.Username)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.RegisterUserSuccess.Code,
		constants.RegisterUserSuccess.Message,
		RegisterResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
		},
	))
}
