package handlers

import (
	"fmt"
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type BoardHandler struct {
	srv       *services.BoardService
	threadSrv *services.ThreadService
}

func NewBoardHandler(srv *services.BoardService, threadSrv *services.ThreadService) *BoardHandler {
	return &BoardHandler{
		srv:       srv,
		threadSrv: threadSrv,
	}
}

var boardTmpl = template.Must(template.Must(baseLayout.Clone()).ParseFiles(path.Join(publicDir, "board.html")))

func (h *BoardHandler) Default(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	board, err := h.srv.GetBySlug(r.Context(), slug)
	if err != nil || board == nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// threads, err :=

	view := models.BoardView{
		Meta:           board,
		IdempotencyKey: h.threadSrv.NewUploadToken(),
	}

	if err := boardTmpl.ExecuteTemplate(w, "layout", view); err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

}

var boardsTpl = template.Must(template.Must(baseLayout.Clone()).ParseFiles(path.Join(publicDir, "boards.html")))

func (h *BoardHandler) Boards(w http.ResponseWriter, r *http.Request) {

	boards, err := h.srv.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view := models.BoardCreateView{
		Boards: boards,
		Error:  r.URL.Query().Get("error"),
	}

	if err := boardsTpl.ExecuteTemplate(w, "layout", view); err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok || sid == "" {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape("invalid board identifier")), http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape("invalid board identifier")), http.StatusSeeOther)
		return
	}

	if err := h.srv.DeleteByID(r.Context(), id); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape("failed to delete board")), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/boards?success=%s", url.QueryEscape("deleted board")), http.StatusSeeOther)

}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {

	// parse post form
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape(err.Error())), http.StatusBadRequest)
		return
	}

	// validate form inputs
	slug := r.PostFormValue("slug")
	if slug == "" {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape("slug parameter required")), http.StatusSeeOther)
		return
	}

	name := r.PostFormValue("name")
	if slug == "" {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape("name parameter required")), http.StatusSeeOther)
		return
	}

	// create board using service method
	board, err := h.srv.Create(r.Context(), models.BoardParams{
		Slug: strings.ToLower(slug),
		Name: name,
	})

	// evaluate errors and nil results
	if err != nil || board == nil {
		http.Redirect(w, r, fmt.Sprintf("/boards?error=%s", url.QueryEscape(err.Error())), http.StatusSeeOther)
		return
	}

	// w.Header().Set("Refresh", fmt.Sprintf("1, /%s", board.Slug))

	http.Redirect(w, r, "/boards", http.StatusSeeOther)

}
