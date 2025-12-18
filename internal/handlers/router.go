package handlers

import (
	"go-image-web/internal/store"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	publicDirPrefix string = "/public/assets/"
	publicDir       string = "public"
	assetsDir       string = "public/assets"
)

type RouterHandlers struct {
	Post   *IndexHandler
	Board  *BoardHandler
	Thread *ThreadHandler
}

func SetupRouter(handlers *RouterHandlers) *mux.Router {
	// create mux router
	r := mux.NewRouter()

	// serve static server
	fs := http.FileServer(http.Dir(assetsDir))
	r.PathPrefix(publicDirPrefix).Handler(http.StripPrefix(publicDirPrefix, fs))

	// register routes with handler functions
	if handlers.Post != nil {
		r.HandleFunc("/", handlers.Post.Home).Methods("GET")
		r.HandleFunc("/upload", handlers.Post.Upload).Methods("POST")
	}
	log.Printf("post handlers registered")

	return r
}

// setup router with some static routes
func SetupRouterz() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))
	r.HandleFunc("/img/{id}", GetImageHandler).Methods("GET")

	return r
}

func init() {
	store.CheckCreateDir(assetsDir)
}
