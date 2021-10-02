package main

import (
	"database/sql"
	"log"

	"github.com/pscompsci/cbot/internal/cbot"
)

func main() {
	db, err := openDB("host=localhost port=5432 user=postgres password=password dbname=cbot sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := cbot.New(db)
	app.Run()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
