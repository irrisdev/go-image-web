package handlers

import (
	"go-image-web/services"
	"html/template"
	"net/http"
	"path"
)

var indexTpl = template.Must(template.ParseFiles(path.Join(publicDir, "index.html")))

func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := services.IndexService()
	if err := indexTpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeStaticAsset(w http.ResponseWriter, r *http.Request) {

}

func UploadAsset(w http.ResponseWriter, r *http.Request) {

}
