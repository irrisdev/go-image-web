package main

import (
	"go-image-web/handlers"
	"go-image-web/store"

	"log"
	"net/http"
)

func main() {

	// initialise database
	storage := store.SetupDatabase()
	defer storage.DB.Close()

	// intialise mux router
	router := handlers.SetupRouter()

	// serve static server
	fs := http.FileServer(http.Dir(store.AssetsFolder))
	router.PathPrefix("/public/assets/").Handler(http.StripPrefix("/public/assets/", fs))

	// listen and serve on port
	log.Printf("web server started on port :9991")
	err := http.ListenAndServe(":9991", router)
	if err != nil {
		log.Fatal(err)
	}

}
