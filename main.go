package main

import (
	"context"
	"go-image-web/internal/db"
	"go-image-web/internal/handlers"
	"go-image-web/internal/repo"
	"go-image-web/internal/services"
	"go-image-web/internal/store"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	AssetsFolder string = "public/assets"

	DbDir  string = "data/db"
	DbPath string = "data/db/storage.db"
)

func main() {

	// initialise database
	xdb := openDB()
	defer xdb.Close()

	// intialise mux router
	router := handlers.SetupRouter()

	// create post repo
	postRepo := repo.NewRepo(xdb)

	// create post service
	postService := services.NewPostService(postRepo)

	// create index handler
	indexHandler := handlers.NewIndexHandler(postService)

	// register routes with handler functions
	router.HandleFunc("/", indexHandler.Home).Methods("GET")
	router.HandleFunc("/upload", indexHandler.Upload).Methods("POST")

	// serve static server
	fs := http.FileServer(http.Dir(AssetsFolder))
	router.PathPrefix("/public/assets/").Handler(http.StripPrefix("/public/assets/", fs))

	server := &http.Server{
		Addr:         ":9991",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go handleGracefulShutdown(server)

	// listen and serve on port
	log.Printf("started on port :9991")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func openDB() *sqlx.DB {
	// create/open database
	store.CheckCreateDir(DbDir)
	xdb, err := sqlx.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// check sqlite version
	var version string
	err = xdb.Get(&version, "SELECT SQLITE_VERSION()")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("running with sqlite version: %s", version)

	xdb.MustExec(`PRAGMA foreign_keys = ON;`)
	xdb.MustExec(`PRAGMA busy_timeout = 5000;`)

	if err := db.EnsureSchema(xdb); err != nil {
		log.Fatal(err)
	}

	return xdb
}

func handleGracefulShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
