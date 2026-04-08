package projects

import (
	"context"
	"asona/internal/repository/projects"
)

// ListByWorkplace retrieves all projects belonging to a specific workplace.
func (i impl) ListByWorkplace(ctx context.Context, workplaceID int64) ([]ProjectResponse, error) {
	// Assuming repository has a List method that takes a workplace filter
	// If not, we'd need to add it. For now, matching the interface.
	list, err := i.repo.Project().List(ctx, projects.ListFilter{
		WorkplaceID: workplaceID,
	})
	if err != nil {
		return nil, err
	}

	res := make([]ProjectResponse, 0, len(list))
	for _, p := range list {
		res = append(res, ProjectResponse{
			ID:          p.ID,
			WorkplaceID: p.WorkplaceID,
			Name:        p.Name,
			Description: p.Description,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return res, nil
}
