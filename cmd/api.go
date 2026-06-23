package main

import (
	"log"
	"net/http"
	"time"

	"api-challenges/internal/todo"
)

type application struct {
	port string
}

func (app *application) mount() http.Handler {
	mux := http.NewServeMux()

	service := todo.NewTodoService()
	handler := todo.NewTodoHandler(service)
	mux.HandleFunc("GET /todos/", handler.GetTodos)
	mux.HandleFunc("GET /todos/{id}", handler.GetTodo)
	mux.HandleFunc("POST /todos/", handler.CreateTodo)
	mux.HandleFunc("PATCH /todos/{id}", handler.UpdateTodo) // might move to PUT and leaeve PATCH for updating Done status only
	mux.HandleFunc("DELETE /todos/{id}", handler.DestroyTodo)

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
