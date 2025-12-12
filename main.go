package main

import (
	"go-image-web/handlers"
	"go-image-web/store"
	"log"
	"net/http"
)

func main() {

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
