package organizations

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/model"
	"asona/internal/pkg/logger"
)

type OrganizationWithRoleResponse struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	LogoURL  string        `json:"logo_url"`
	Role     model.OrgRole `json:"role"`
	JoinedAt time.Time     `json:"joined_at"`
}

// ListByUser handles GET /api/v1/organizations
// @Summary      List Organizations by User
// @Description  Retrieve all organizations the current user is a member of
// @Tags         organizations
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.Response{data=[]OrganizationWithRoleResponse}
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /organizations [get]
func (h Handler) ListByUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[ListByUser] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	orgs, err := h.ctrl.ListByUser(c.Request.Context(), uid)
	if err != nil {
		logger.ERROR.Printf("[ListByUser] list failed for user %d: %+v", uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	result := make([]OrganizationWithRoleResponse, 0, len(orgs))
	for _, o := range orgs {
		result = append(result, OrganizationWithRoleResponse{
			ID:       o.ID,
			Name:     o.Name,
			LogoURL:  o.LogoURL,
			Role:     o.Role,
			JoinedAt: o.JoinedAt,
		})
	}

	logger.INFO.Printf("[ListByUser] organizations listed successfully for user %d: count %d", uid, len(result))
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		result,
	))
}
