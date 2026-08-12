package main

import (
	"api-challenges/internal/database"
	"log"
	"os"
)

func main() {
	// Setup
	db, err := database.New("todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	api := application{
		port: ":8080",
		db:   db,
	}

	// Mount the API
	if err := api.serve(api.mount()); err != nil {
		log.Printf("server failed to start: %s", err)
		os.Exit(1)
	}
}
