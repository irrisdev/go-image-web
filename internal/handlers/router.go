package handlers

import (
	"go-image-web/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

var publicDir string = "public"
var assetsDir string = "public/assets"

// setup router with some static routes
func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))
	r.HandleFunc("/img/{id}", GetImageHandler).Methods("GET")

	return r
}

func init() {
	store.CheckCreateDir(assetsDir)
}
