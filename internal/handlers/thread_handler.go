package handlers

import (
	"fmt"
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"go-image-web/internal/store"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

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

	// if idempotency has already been used redirect to the thread page
	if entry.State != store.Created {
		redirectThreadByID(w, r, slug, entry.ThreadID)
	}

	// craft thread inputs for service
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

	if err != nil {
		redirectErr(w, r, slug, err)
	}
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "state_key",
	// 	Value:    idempotencyKey,
	// 	Path:     fmt.Sprintf("/%s/%d", slug, id),
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	SameSite: http.SameSiteLaxMode,
	// 	Expires:  time.Now().Add(1 * time.Hour),
	// })
	redirectThreadByID(w, r, slug, id)
}

var threadTmpl = template.Must(template.Must(baseLayout.Clone()).ParseFiles(path.Join(publicDir, "thread.html")))

func (h *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// retrieve board slug part of url path
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// retrieve thread id part of url path
	sid, ok := vars["id"]
	id, err := strconv.Atoi(sid)
	if !ok || err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// retireve boards by slug if exists
	board, err := h.boardService.GetBySlug(r.Context(), slug)
	if err != nil || board == nil {
		http.Error(w, "404 board not found", http.StatusNotFound)
		return
	}

	// craft view model
	var view = models.ThreadView{
		Slug:  board.Slug,
		Error: r.URL.Query().Get("error"),
	}

	// find thread by ID
	thread, err := h.srv.GetByID(r.Context(), id)
	if err != nil {
		view.Error += fmt.Sprintf("; %s", err.Error())
		if err := threadTmpl.ExecuteTemplate(w, "layout", view); err != nil {
			log.Println(err)
			http.Error(w, "internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	// set view thread
	view.Thread = thread

	// check if image state store via image uuid
	// if recently created and still processing add refresh header to response
	imageState, _ := h.srv.GetUploadEntry(thread.UUID)
	if ok && imageState.State == store.Processing {
		w.Header().Set("Refresh", "1")
	}

	// retireve image metadata regardless upload state as upload state only persists for 1 hour
	view.Image = h.srv.GetImageByUUID(thread.UUID)

	if err := threadTmpl.ExecuteTemplate(w, "layout", view); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
