package handlers

import (
	"go-image-web/internal/services"
	"net/http"
)

type BoardHandler struct {
	srv *services.BoardService
}

func NewBoardService(srv *services.BoardService) *BoardHandler {
	return &BoardHandler{
		srv: srv,
	}
}

func (h *BoardHandler) Default(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}
