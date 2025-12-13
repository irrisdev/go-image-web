package handlers

import (
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"go-image-web/internal/store"
	"html/template"
	"log"
	"net/http"
	"path"
	"sort"
)

type IndexHandler struct {
	PostService *services.PostService
}

func NewIndexHandler(postService *services.PostService) *IndexHandler {
	return &IndexHandler{
		PostService: postService,
	}
}

var baseLayout = template.Must(template.ParseFiles(path.Join(publicDir, "layout.html")))

func (h *IndexHandler) Home(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.Must(baseLayout.Clone()).ParseFiles(path.Join(publicDir, "index.html")))

	var viewModel []*models.PostViewModel

	// imageData := services.IndexService()

	// do nothing with error at the moment, however in future display error message
	postData, _ := h.PostService.GetPosts()
	for _, post := range postData {
		if post.ImageUUID == "" {
			continue
		}
		meta := store.GetGuidImageMetadata(post.ImageUUID)
		if meta != nil {
			viewModel = append(viewModel, &models.PostViewModel{
				Image: &models.ImageModel{
					ID:        meta.UUID,
					Path:      meta.OriginalPath,
					Extension: meta.OriginalExt,
					Width:     meta.OriginalWidth,
					Height:    meta.OriginalHeight,
					Timestamp: meta.ModifiedTime,
					Size:      meta.OriginalSize,
				},
				Post: post,
			})
		}

	}

	sort.Slice(viewModel, func(i, j int) bool {
		return viewModel[i].Image.Timestamp.After(viewModel[j].Image.Timestamp)
	})

	if err := tpl.ExecuteTemplate(w, "layout", viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *IndexHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// hard limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, store.MaxUploadBytes)

	var uuid string

	// read multipart file and header
	file, header, err := r.FormFile("imageFile")
	if err != http.ErrMissingFile {

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// check if size in header is too big
		if header.Size > store.MaxUploadBytes {
			http.Error(w, "file too big", http.StatusRequestEntityTooLarge)
			return
		}

		// save image to system and return uuid
		uuid, err = services.SaveImage(file, header.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// extract name value and set to default if not exist
	name := r.FormValue("name")
	if name == "" {
		name = services.DefaultPostName
	}

	postModel := &models.PostModel{
		Name:      name,
		Subject:   r.FormValue("subject"),
		Message:   r.FormValue("message"),
		ImageUUID: uuid,
	}

	_, err = h.PostService.SavePost(postModel)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
