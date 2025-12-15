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

	// do nothing with error at the moment, however in future display error message
	postData, _ := h.PostService.GetPosts()
	for _, post := range postData {

		if post.ImageUUID == "" {
			viewModel = append(viewModel, &models.PostViewModel{
				Post: post,
			})
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

	if err := tpl.ExecuteTemplate(w, "layout", viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *IndexHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// hard limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, store.MaxUploadBytes)

	var uuid string

	// read multipart file and header
	file, header, fileErr := r.FormFile("imageFile")

	if fileErr == http.ErrMissingFile {
		http.Error(w, fmt.Errorf("must provide an image").Error(), http.StatusBadRequest)
		return
	}

	// if file detected but has error and isn't missing file
	if fileErr != nil && fileErr != http.ErrMissingFile {
		http.Error(w, fileErr.Error(), http.StatusBadRequest)
		return
	}

	if fileErr == nil {

		defer file.Close()

		// check if size in header is too big
		if header.Size > store.MaxUploadBytes {
			http.Error(w, "file too big", http.StatusRequestEntityTooLarge)
			return
		}

		// save image to system and return uuid
		var saveErr error
		uuid, saveErr = services.SaveImage(file, header.Filename)
		if saveErr != nil {
			http.Error(w, saveErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	subject, message := r.FormValue("subject"), r.FormValue("message")

	// validate required fields
	// if message == "" && fileErr != nil {
	// 	http.Error(w, fmt.Errorf("must provide a message or an image").Error(), http.StatusBadRequest)
	// 	return
	// }

	// Create post model
	postModel := &models.PostModel{
		Name:      services.DefaultPostName,
		Subject:   subject,
		Message:   message,
		ImageUUID: uuid,
	}

	_, saveErr := h.PostService.SavePost(postModel)
	if saveErr != nil {
		log.Println(saveErr)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
