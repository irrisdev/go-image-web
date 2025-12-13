package handlers

import (
	"fmt"
	"go-image-web/internal/models"
	"go-image-web/internal/services"
	"go-image-web/internal/store"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {

	// hard limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, store.MaxUploadBytes)

	post := models.PostUploadModel{
		Name:    r.FormValue("name"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	fmt.Println(post)

	// read multipart file and header
	file, header, err := r.FormFile("imageFile")
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

	id, err := services.SaveImage(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = id

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if v, ok := vars["id"]; !ok || v == "" {
		http.Error(w, "invalid image id", http.StatusBadRequest)
		return
	}

	varient, err := services.GetImage(vars["id"])
	if err != nil {

		http.Error(w, fmt.Sprintf("image not found: %v", err), http.StatusNotFound)
		return
	}

	if cType := imageContentType(varient.Ext); cType != "" {
		w.Header().Set("Content-Type", cType)
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	http.ServeFile(w, r, varient.Path)

}

func imageContentType(ext string) string {
	switch strings.ToLower(ext) {
	case "gif":
		return "image/gif"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return ""
	}
}
