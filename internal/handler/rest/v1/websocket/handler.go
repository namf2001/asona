package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"

	wsservice "asona/internal/service/websocket"
)

// Handler handles WebSocket requests.
type Handler struct {
	wsService wsservice.Service
}

// New creates a new WebSocket handler.
func New(wsService wsservice.Service) *Handler {
	return &Handler{
		wsService: wsService,
	}
}

// HandleWebSocket upgrades the HTTP connection to a WebSocket connection.
// @Summary      WebSocket Connection
// @Description  Establishes a WebSocket connection for real-time communication
// @Tags         WebSocket
// @Accept       json
// @Produce      json
// @Param        userId query string true "User ID"
// @Success      101 {string} string "Switching Protocols"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /ws [get]
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// Get user ID from query or context (from auth middleware)
	userID := c.Query("userId")
	if userID == "" {
		// Try to get from auth context
		if id, exists := c.Get("userID"); exists {
			userID = id.(string)
		}
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.wsService.Upgrader().Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	// Create new client and register to hub
	client := wsservice.NewClient(userID, conn, h.wsService.Hub())
	h.wsService.Hub().Register(client)
}
