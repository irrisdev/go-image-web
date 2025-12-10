package handlers

import (
	"fmt"
	"go-image-web/services"
	"go-image-web/store"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var baseLayout = template.Must(template.ParseFiles(path.Join(publicDir, "layout.html")))

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.Must(baseLayout.Clone()).
		ParseFiles(path.Join(publicDir, "index.html")))

	data := services.IndexService()

	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UploadAsset(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("imageFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	_ = header

	dst, err := os.Create(path.Join(store.ImageAssetsFolder, fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sanitise filename

	// redirect to index pager / new upload
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
