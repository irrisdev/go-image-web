package main

import (
	"log"
	"net/http"

	"github.com/irrisdev/go-image-web/handlers"
)

func main() {

	router := handlers.SetupRouter()

	log.Printf("web server started on port :9991")
	err := http.ListenAndServe(":9991", router)
	if err != nil {
		log.Fatal(err)
	}

}
