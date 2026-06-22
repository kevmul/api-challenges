package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	port string
}

func (app *application) mount() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /todos/", func() {
		return
	})
	return mux
}

func (app *application) serve(h http.Handler) error {
	srv := &http.Server{
		Addr:              app.port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Minute,
	}

	log.Printf("server started on port: %s", app.port)

	return srv.ListenAndServe()
}
