package main

import (
	"errors"
	"go-image-web/handlers"
	"log"
	"net/http"
	"os"
)

var AssetsFolder string = "public/assets"

func main() {

	CheckAndCreateDir(AssetsFolder)

	router := handlers.SetupRouter()

	// Serve static server
	fs := http.FileServer(http.Dir(AssetsFolder))

	router.PathPrefix("/public/assets/").Handler(http.StripPrefix("/public/assets/", fs))

	log.Printf("web server started on port :9991")
	err := http.ListenAndServe(":9991", router)
	if err != nil {
		log.Fatal(err)
	}

}

func CheckAndCreateDir(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
