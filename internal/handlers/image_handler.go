package handlers

import (
	"fmt"
	"go-image-web/internal/services"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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
