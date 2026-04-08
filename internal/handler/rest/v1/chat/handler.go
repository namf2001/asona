package chat

import (
	"asona/internal/controller/chat/channel"
	"asona/internal/controller/chat/message"
)

// Handler aggregates chat-related controllers to handle HTTP requests.
type Handler struct {
	channelCtrl channel.Controller
	messageCtrl message.Controller
}

// New creates a new chat handler with necessary controllers.
func New(channelCtrl channel.Controller, messageCtrl message.Controller) Handler {
	return Handler{
		channelCtrl: channelCtrl,
		messageCtrl: messageCtrl,
	}
}
