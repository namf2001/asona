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

type completeOnboardResponse struct {
	IsOnboarded bool `json:"is_onboarded"`
}

// CompleteOnboard handles PATCH /api/v1/me/onboard.
// It marks the authenticated user's onboarding as complete
// and returns the updated is_onboarded flag.
func (h Handler) CompleteOnboard(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CompleteOnboard] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	userID, ok := val.(int64)
	if !ok {
		logger.ERROR.Printf("[CompleteOnboard] invalid userID type in context: %T", val)
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	if err := h.ctrl.CompleteOnboard(c.Request.Context(), userID); err != nil {
		logger.ERROR.Printf("[CompleteOnboard] failed to complete onboard for user %d: %+v", userID, err)
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

	logger.INFO.Printf("[CompleteOnboard] user %d completed onboarding", userID)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.OnboardSuccess.Code,
		constants.OnboardSuccess.Message,
		completeOnboardResponse{IsOnboarded: true},
	))
}
