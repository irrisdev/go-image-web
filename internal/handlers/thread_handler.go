package handlers

import "go-image-web/internal/services"

type ThreadHandler struct {
	srv *services.ThreadService
}

func NewThreadService(srv *services.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		srv: srv,
	}
}
