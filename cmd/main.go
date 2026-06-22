package main

import (
	"log"
	"os"
)

func main() {
	api := application{
		port: ":8080",
	}

	if err := api.serve(api.mount()); err != nil {
		log.Printf("server failed to start: %s", err)
		os.Exit(1)
	}
}
