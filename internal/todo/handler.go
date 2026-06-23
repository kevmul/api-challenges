package todo

import (
	"api-challenges/internal/helpers"
	"encoding/json"
	"log"
	"net/http"
)

type TodoHandler struct {
	service TodoService
}

// functions declared in an interface MUST have the
// uppercased first letter to be Public. This allows
// the service to have access to the functions.
type TodoService interface {
	GetAll() ([]Todo, error)
	GetById(id int) (Todo, error)
	Create(todoTitle string) (Todo, error)
}

func NewTodoHandler(service TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// GetTodos handles the GET /todos/ endpoint and
// returns a list of todos in JSON format.
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAll() // This is the error
	if err != nil {
		log.Printf("error getting todos: %s", err)
		return
	}

	helpers.WriteJson(w, http.StatusOK, todos)
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

// Create a new todo item and return it in JSON format.
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding json: %s", err)
		return
	}

	if req.Title == "" {
		log.Printf("error: title is required")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	todo, err := h.service.Create(req.Title)
	if err != nil {
		log.Printf("error creating todo: %s", err)
		return
	}

	helpers.WriteJson(w, http.StatusCreated, todo)
}
