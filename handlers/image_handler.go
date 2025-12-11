package handlers

import (
	"fmt"
	"go-image-web/services"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("imageFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	id, err := services.SaveImage(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = id

	http.Redirect(w, r, "/", http.StatusCreated)

}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if v, ok := vars["id"]; !ok || v == "" {
		http.Error(w, "invalid image id", http.StatusBadRequest)
		return
	}

	varient, err := services.GetImage(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error loading image: %v", err), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, varient.Path)

}
