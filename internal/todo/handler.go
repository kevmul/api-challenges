package todo

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	todos, err := getAll()
	if err != nil {
		log.Printf("error getting todos: %s", err)
		return
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("error encoding json: %s", err)
	}
}
