package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	authctrl "asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type getOnboardingStateResponse struct {
	Status      string     `json:"status"`
	Step        int16      `json:"step"`
	IsOnboarded bool       `json:"is_onboarded"`
	OnboardedAt *time.Time `json:"onboarded_at,omitempty"`
}

// GetOnboardingState handles GET /api/v1/me/onboarding.
// It returns the authenticated user's onboarding status snapshot.
func (h Handler) GetOnboardingState(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[GetOnboardingState] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	userID, ok := val.(int64)
	if !ok {
		logger.ERROR.Printf("[GetOnboardingState] invalid userID type in context: %T", val)
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}

	state, err := h.ctrl.GetOnboardingState(c.Request.Context(), userID)
	if err != nil {
		logger.ERROR.Printf("[GetOnboardingState] failed to get onboarding state for user %d: %+v", userID, err)
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

	logger.INFO.Printf("[GetOnboardingState] onboarding state retrieved for user %d", userID)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		getOnboardingStateResponse{
			Status:      string(state.Status),
			Step:        state.Step,
			IsOnboarded: state.IsOnboarded,
			OnboardedAt: state.OnboardedAt,
		},
	))
}
