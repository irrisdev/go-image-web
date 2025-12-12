package main

import (
	"context"
	"go-image-web/handlers"
	"go-image-web/store"
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
	db := openDB()
	defer db.Close()

	// intialise mux router
	router := handlers.SetupRouter()

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
	db, err := sqlx.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// check sqlite version
	var version string
	err = db.Get(&version, "SELECT SQLITE_VERSION()")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("running with sqlite version: %s", version)

	db.MustExec(`PRAGMA foreign_keys = ON;`)
	db.MustExec(`PRAGMA busy_timeout = 5000;`)

	return db
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
