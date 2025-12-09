package handlers

import (
	"github.com/gorilla/mux"
)

var publicDir string = "./public"
var staticDir string = "./public/assets"

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/", IndexPage).Methods("GET")
	r.HandleFunc("/upload", UploadAsset).Methods("POST")

	return r

}
