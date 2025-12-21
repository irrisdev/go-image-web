package handlers

import (
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"html/template"
	"net/http"
	"path"
	"strings"
)

type BoardHandler struct {
	srv *services.BoardService
}

func NewBoardHandler(srv *services.BoardService) *BoardHandler {
	return &BoardHandler{
		srv: srv,
	}
}

func (h *BoardHandler) Default(w http.ResponseWriter, r *http.Request) {

	tpl := template.Must(template.Must(baseLayout.Clone()).ParseFiles(path.Join(publicDir, "boards.html")))

	boards, err := h.srv.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "layout", boards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {

	// parse post form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate form inputs
	slug := r.PostFormValue("slug")
	if slug == "" {
		http.Error(w, "slug parameter required", http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	if slug == "" {
		http.Error(w, "name parameter required", http.StatusBadRequest)
		return
	}

	// create board using service method
	board, err := h.srv.Create(r.Context(), models.BoardParams{
		Slug: strings.ToLower(slug),
		Name: name,
	})

	// evaluate errors and nil results
	if err != nil || board == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Refresh", fmt.Sprintf("1, /%s", board.Slug))

	http.Redirect(w, r, "/boards", http.StatusSeeOther)

}
