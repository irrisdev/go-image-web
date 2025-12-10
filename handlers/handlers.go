package handlers

import (
	"go-image-web/services"
	"html/template"
	"net/http"
	"path"
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

	result, err := services.AddImage(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(result.ErrorMessages) > 0 {
		// http.Error(w, result.ErrorMessages, http.StatusBadRequest)
		return
	}

	// sanitise filename

	// redirect to index pager / new upload
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
