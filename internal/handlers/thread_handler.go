package handlers

import (
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"go-image-web/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

type ThreadHandler struct {
	srv          *services.ThreadService
	boardService *services.BoardService
}

func NewThreadHandler(srv *services.ThreadService, boardService *services.BoardService) *ThreadHandler {
	return &ThreadHandler{
		srv:          srv,
		boardService: boardService,
	}
}

const (
	MaxThreadBytes  int64 = 15 << 20
	MaxSubjectChars int   = 70
	MaxMessageChars int   = 1500
)

func (h *ThreadHandler) NewThread(w http.ResponseWriter, r *http.Request) {
	// hard limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, MaxThreadBytes)

	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// check if board exists
	board, err := h.boardService.GetBySlug(r.Context(), slug)
	if err != nil || board == nil {
		http.Error(w, "404 board not found", http.StatusNotFound)
		return
	}

	// read multipart file and header
	file, header, fileErr := r.FormFile("imageFile")
	if fileErr == http.ErrMissingFile {
		redirectErr(w, r, slug, services.ErrMissingImage)
		return
	}

	// if file detected but has error and isn't missing file
	if fileErr != nil && fileErr != http.ErrMissingFile {
		redirectErr(w, r, slug, fileErr)
		return
	}
	defer file.Close()

	// validate idempotency key
	idempotencyKey := r.FormValue("idempotency_key")
	entry, ok := h.srv.GetUploadEntry(idempotencyKey)
	if idempotencyKey == "" || !ok {
		redirectErr(w, r, slug, services.ErrBadIdempotencyKey)
	}

	if entry.State != store.Created {
		redirectThreadByUUID(w, r, slug, entry.ThreadUUID)
	}

	subject, message := r.FormValue("subject"), r.FormValue("message")
	inputs := &models.NewThreadInputs{
		File:           file,
		Header:         header,
		Subject:        subject,
		Message:        message,
		BoardID:        board.ID,
		IdempotencyKey: idempotencyKey,
	}

	id, err := h.srv.Create(r.Context(), inputs)
	// retireve the board information

	if err != nil {
		redirectErr(w, r, slug, err)
	}

	redirectThreadByID(w, r, slug, id)

}

func (h *ThreadHandler) Default(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}
