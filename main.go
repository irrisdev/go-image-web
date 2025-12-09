package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/irrisdev/go-image-web/handlers"
)

var assetsFolder string = "public/assets"

func main() {

	CheckAndCreateDir(assetsFolder)

	router := handlers.SetupRouter()

	// Serve static server
	fs := http.FileServer(http.Dir(assetsFolder))

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
