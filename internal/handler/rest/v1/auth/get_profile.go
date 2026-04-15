package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	authctrl "asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type GetProfileResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Image       string `json:"image,omitempty"`
	IsOnboarded bool   `json:"is_onboarded"`
}

// Profile handles GET /auth/profile
// @Summary      Get user profile
// @Description  Get the profile of the currently authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} response.Response{data=GetProfileResponse}
// @Failure      401  {object} response.Response
// @Router       /auth/profile [get]
func (h Handler) Profile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[Profile] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	userID, ok := val.(int64)
	if !ok {
		logger.ERROR.Printf("[Profile] invalid userID type in context: %T", val)
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	user, err := h.ctrl.GetProfile(c.Request.Context(), userID)
	if err != nil {
		logger.ERROR.Printf("[Profile] get profile failed for user %d: %+v", userID, err)
		if errors.Is(err, authctrl.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.NewResponse(
				constants.UserNotFound.Code,
				constants.UserNotFound.Message,
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

	res := GetProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Image:       user.Image,
		IsOnboarded: user.OnboardedAt != nil,
	}

	logger.INFO.Printf("[Profile] user profile retrieved successfully: %s", user.Email)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		res,
	))
}
