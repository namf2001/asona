package projects

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlProj "asona/internal/controller/projects"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type CreateProjectRequest struct {
	WorkplaceID int64  `json:"workplace_id" binding:"required"`
	Name        string `json:"name"         binding:"required"`
	Description string `json:"description"`
}

type CreateProjectResponse struct {
	ID          int64     `json:"id"`
	WorkplaceID int64     `json:"workplace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProject handles POST /api/v1/projects
// @Summary      Create Project
// @Description  Create a new project in a specific workplace
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateProjectRequest  true  "Project details"
// @Success      201      {object}  response.Response{data=CreateProjectResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /projects [post]
func (h Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[CreateProject] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CreateProject] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	proj, err := h.ctrl.Create(c.Request.Context(), uid, ctrlProj.CreateProjectInput{
		WorkplaceID: req.WorkplaceID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		logger.ERROR.Printf("[CreateProject] create failed for user %d in workplace %d: %+v", uid, req.WorkplaceID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[CreateProject] project created successfully: %s (ID: %d)", proj.Name, proj.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Project created",
		CreateProjectResponse{
			ID:          proj.ID,
			WorkplaceID: proj.WorkplaceID,
			Name:        proj.Name,
			Description: proj.Description,
			CreatedBy:   proj.CreatedBy,
			CreatedAt:   proj.CreatedAt,
			UpdatedAt:   proj.UpdatedAt,
		},
	))
}
