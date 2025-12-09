package handlers

import (
	"html/template"
	"net/http"
	"path"

	"github.com/irrisdev/go-image-web/services"
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
