package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/jwt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles manual login
// @Summary      User login
// @Description  Authenticate user and return token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body auth.LoginRequest true "Login credentials"
// @Success      200  {object} response.Response
// @Failure      400  {object} response.Response
// @Failure      500  {object} response.Response
// @Router       /login [post]
func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	// In a complete app, DB check/OAuth logic happens here.
	// For this stateless demo, we trust the input and generate the token directly.
	// Hardcoding UserID 1.
	token, err := jwt.GenerateToken(1, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.GenerateTokenFailed.Code,
			constants.GenerateTokenFailed.Message,
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		LoginResponse{
			Token: token,
		},
	))
}


