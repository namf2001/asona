package projects

import (
	"context"
	"time"

	"asona/internal/model"
	"asona/internal/repository"
)

// CreateProjectInput represents the data required to create a new project.
type CreateProjectInput struct {
	WorkplaceID int64  `json:"workplace_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectResponse represents the public project data.
type ProjectResponse struct {
	ID          int64     `json:"id"`
	WorkplaceID int64     `json:"workplace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create creates a new project and adds the creator as a member.
func (i impl) Create(ctx context.Context, userID int64, input CreateProjectInput) (ProjectResponse, error) {
	var res ProjectResponse

	err := i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
		// 1. Create project
		p, err := txRepo.Project().Create(ctx, model.Project{
			WorkplaceID: input.WorkplaceID,
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   userID,
			Access:      model.ProjectAccessPublic,
		})
		if err != nil {
			return err
		}

		// 2. Add creator as owner/member (if repo supports it, using AddMember for now)
		err = txRepo.Project().AddMember(ctx, model.ProjectMember{
			ProjectID: p.ID,
			UserID:    userID,
			Role:      model.ProjectRoleOwner,
		})
		if err != nil {
			return err
		}

		res = ProjectResponse{
			ID:          p.ID,
			WorkplaceID: p.WorkplaceID,
			Name:        p.Name,
			Description: p.Description,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
		return nil
	}, nil)

	return res, err
}
