package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	AssetsFolder string = "public/assets"
	DbDir        string = "data/db"
	DbPath       string = "data/db/storage.db"
)

type Store struct {
	DB *sql.DB
}

func SetupDatabase() *Store {
	// create database if not exists
	CheckCreateDir(DbDir)
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// check sqlite version
	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running with sqlite version: %s", version)

	return &Store{
		DB: db,
	}
}
