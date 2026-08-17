package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"api-challenges/internal/todo"
)

type Application struct {
	port string
	DB   *sql.DB
}

func (app *Application) mount() http.Handler {
	r := gin.Default()

	todo.RegisterRoutes(r, app.DB)

	return r
}

func (app *Application) serve(h http.Handler) error {
	srv := &http.Server{
		Addr:              app.port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Minute,
	}

	log.Printf("server started on port: %s", app.port)

	return srv.ListenAndServe()
}
