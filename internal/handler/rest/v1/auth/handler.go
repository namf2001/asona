package auth

import (
	authctrl "asona/internal/controller/auth"
)

type Handler struct {
	authCtrl authctrl.Controller
}

func New(ctrl authctrl.Controller) *Handler {
	return &Handler{
		authCtrl: ctrl,
	}
}
