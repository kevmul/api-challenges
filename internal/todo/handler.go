package todo

import (
	"api-challenges/internal/helpers"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	Update(id int, updatedTodo Todo) (Todo, error)
	Destroy(id int) error
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

// Get a single todo item by ID and return it in JSON format.
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		http.Error(w, "invalid id. ID Must be an integer", http.StatusBadRequest)
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		log.Printf("error getting todo by id: %s", err)
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	helpers.WriteJson(w, http.StatusOK, todo)
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

type PatchTodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// UpdateTodo handles the PATCH /todos/{id} endpoint and updates a todo item.
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		http.Error(w, "invalid id. ID Must be an integer", http.StatusBadRequest)
		return
	}

	var req PatchTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding json: %s", err)
		return
	}

	if req.Title == "" {
		log.Printf("error: title is required")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	log.Printf("updating todo with id: %d and title: %s", id, req.Title)

	_, err = h.service.GetById(id)
	if err != nil {
		log.Printf("error getting todo by id: %s", err)
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	response, err := h.service.Update(id, Todo{ID: id, Title: req.Title, Done: req.Done})
	if err != nil {
		log.Printf("error updating todo: %s", err)
		http.Error(w, "error updating todo", http.StatusInternalServerError)
		return
	}

	helpers.WriteJson(w, http.StatusOK, response)
}

// DestroyTodo handles the DELETE /todos/{id} endpoint and
func (h *TodoHandler) DestroyTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		http.Error(w, "invalid id. ID Must be an integer", http.StatusBadRequest)
		return
	}

	todo, _ := h.service.GetById(id)

	err = h.service.Destroy(id)
	if err != nil {
		log.Printf("could not delete todo: %s", todo.Title)
	}

	type DeletedResponse struct {
		Deleted Todo `json:"deleted"`
	}

	response := DeletedResponse{Deleted: todo}

	helpers.WriteJson(w, http.StatusOK, response)
}
