package workplaces

import (
	"context"
	"time"

	"asona/internal/model"
	"asona/internal/repository"
)

// CreateWorkplaceInput represents the input data for creating a workplace.
type CreateWorkplaceInput struct {
	Name    string             `json:"name"`
	IconURL string             `json:"icon_url"`
	Size    model.WorkplaceSize `json:"size"`
}

// WorkplaceResponse represents the public workplace data.
type WorkplaceResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	IconURL   string             `json:"icon_url"`
	Size      model.WorkplaceSize `json:"size"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// Create generates a new workplace in the system and adds the creator as the initial owner.
func (i impl) Create(ctx context.Context, userID int64, input CreateWorkplaceInput) (WorkplaceResponse, error) {
	var res WorkplaceResponse

	err := i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
		// 1. Create workplace
		wp, err := txRepo.Workplace().Create(ctx, model.Workplace{
			Name:      input.Name,
			IconURL:   input.IconURL,
			Size:      input.Size,
			CreatedBy: userID,
		})
		if err != nil {
			return err
		}

		// 2. Add creator as admin
		err = txRepo.Workplace().AddMember(ctx, model.WorkplaceMember{
			WorkplaceID: wp.ID,
			UserID:      userID,
			Role:        model.WorkplaceRoleAdmin,
		})
		if err != nil {
			return err
		}

		res = WorkplaceResponse{
			ID:        wp.ID,
			Name:      wp.Name,
			IconURL:   wp.IconURL,
			Size:      wp.Size,
			CreatedAt: wp.CreatedAt,
			UpdatedAt: wp.UpdatedAt,
		}
		return nil
	}, nil)

	return res, err
}
