package handlers

import (
	"github.com/gorilla/mux"
)

var publicDir string = "./public"
var staticDir string = "./public/assets"

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/img/{id}", GetImageHandler).Methods("GET")
	r.HandleFunc("/upload", UploadImageHandler).Methods("POST")

	return r

}
