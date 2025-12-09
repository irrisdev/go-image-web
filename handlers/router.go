package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var publicDir string = "./public"
var staticDir string = "./public/assets"

func SetupRouter() *mux.Router {

	CheckAndCreateDir(publicDir)
	CheckAndCreateDir(staticDir)

	fs := http.FileServer(http.Dir(staticDir))

	r := mux.NewRouter()

	r.PathPrefix("/public/assets/").Handler(http.StripPrefix("/public/assets/", fs))

	r.HandleFunc("/", IndexPage).Methods("GET")
	r.HandleFunc("/upload", UploadAsset).Methods("POST")

	return r

}

func CheckAndCreateDir(path string) {
	if _, err := os.Stat(staticDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(staticDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
