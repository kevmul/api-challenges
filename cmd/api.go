package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"api-challenges/internal/todo"
)

type application struct {
	port string
}

func (app *application) mount() http.Handler {
	r := gin.Default()

	todo.RegisterRoutes(r)

	return r
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
