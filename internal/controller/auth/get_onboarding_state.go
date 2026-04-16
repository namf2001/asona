package auth

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
	"asona/internal/repository/users"
)

// OnboardingState represents the current onboarding progress for a user.
type OnboardingState struct {
	Status      model.OnboardingStatus
	Step        int16
	IsOnboarded bool
	OnboardedAt *time.Time
}

// GetOnboardingState returns onboarding status fields for the authenticated user.
func (i impl) GetOnboardingState(ctx context.Context, userID int64) (OnboardingState, error) {
	user, err := i.repo.User().GetByID(ctx, userID)
	if err != nil {
		if pkgerrors.Is(err, users.ErrUserNotFound) {
			return OnboardingState{}, pkgerrors.WithStack(ErrUserNotFound)
		}
		return OnboardingState{}, err
	}

	status := user.OnboardingStatus
	if status == "" {
		if user.OnboardedAt != nil {
			status = model.OnboardingStatusCompleted
		} else {
			status = model.OnboardingStatusPending
		}
	}

	step := user.OnboardingStep
	if step == 0 && user.OnboardedAt != nil {
		step = 3
	}

	return OnboardingState{
		Status:      status,
		Step:        step,
		IsOnboarded: user.OnboardedAt != nil,
		OnboardedAt: user.OnboardedAt,
	}, nil
}
