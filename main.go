package main

import (
	"errors"
	"go-image-web/handlers"
	"go-image-web/store"
	"log"
	"net/http"
	"os"
)

func main() {

	CheckAndCreateDir(store.AssetsFolder)

	router := handlers.SetupRouter()

	// Serve static server
	fs := http.FileServer(http.Dir(store.AssetsFolder))

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
