package handlers

import (
	"go-image-web/internal/services"
	"net/http"
)

type ThreadHandler struct {
	srv *services.ThreadService
}

func NewThreadService(srv *services.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		srv: srv,
	}
}
func (h *ThreadHandler) Default(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}
