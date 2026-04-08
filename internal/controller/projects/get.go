package projects

import (
	"context"
)

// GetByID retrieves a single project's details by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (ProjectResponse, error) {
	p, err := i.repo.Project().GetByID(ctx, id)
	if err != nil {
		return ProjectResponse{}, err
	}

	return ProjectResponse{
		ID:          p.ID,
		WorkplaceID: p.WorkplaceID,
		Name:        p.Name,
		Description: p.Description,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}
