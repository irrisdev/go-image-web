package handlers

import "go-image-web/internal/services"

type BoardHandler struct {
	srv *services.BoardService
}

func NewBoardService(srv *services.BoardService) *BoardHandler {
	return &BoardHandler{
		srv: srv,
	}
}
