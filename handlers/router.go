package handlers

import (
	"go-image-web/store"
	"net/http"

	"github.com/gorilla/mux"
)

var publicDir string = "public"
var assetsDir string = "public/assets"

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/img/{id}", GetImageHandler).Methods("GET")
	r.HandleFunc("/upload", UploadImageHandler).Methods("POST")

	return r
}

func init() {
	store.CheckCreateDir(assetsDir)
}
