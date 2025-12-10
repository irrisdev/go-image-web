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

func ServeStaticAsset(w http.ResponseWriter, r *http.Request) {

}

func UploadAsset(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("imageFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	_ = header
	// sanitise filename

	// save file

	// redirect to index pager / new upload
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
