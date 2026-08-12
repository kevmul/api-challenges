package todo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service TodoService
}

// functions declared in an interface MUST have the
// uppercased first letter to be Public. This allows
// the service to have access to the functions.
type TodoService interface {
	GetAll() ([]Todo, error)
	GetById(id int64) (Todo, error)
	Create(todoTitle string) (Todo, error)
	Update(id int64, updatedTodo Todo) (Todo, error)
	Patch(id int64) (Todo, error)
	Destroy(id int64) error
}

func NewTodoHandler(service TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// GetTodos handles the GET /todos/ endpoint and
// returns a list of todos in JSON format.
func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.service.GetAll() // This is the error
	if err != nil {
		log.Printf("error getting todos: %s", err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

// Get a single todo item by ID and return it in JSON format.
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id. ID Must be an integer"})
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		log.Printf("error getting todo by id: %s", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Create a new todo item and return it in JSON format.
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req CreateTodoRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Printf("error decoding json: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		log.Printf("error: title is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	todo, err := h.service.Create(req.Title)
	if err != nil {
		log.Printf("error creating todo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) ToggleTodoState(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id. ID Must be an integer"})
		return
	}

	_, err = h.service.GetById(id)
	if err != nil {
		log.Printf("error getting todo by id: %s", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	response, err := h.service.Patch(id)
	if err != nil {
		log.Printf("error updating todo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": response})

}

type PatchTodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// UpdateTodo handles the PATCH /todos/{id} endpoint and updates a todo item.
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id. ID Must be an integer"})
		return
	}

	var req PatchTodoRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Printf("error decoding json: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		log.Print("error: title is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	log.Printf("updating todo with id: %d and title: %s", id, req.Title)

	_, err = h.service.GetById(id)
	if err != nil {
		log.Printf("error getting todo by id: %s", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	response, err := h.service.Update(id, Todo{ID: id, Title: req.Title, Done: req.Done})
	if err != nil {
		log.Printf("error updating todo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": response})
}

// DestroyTodo handles the DELETE /todos/{id} endpoint and
func (h *TodoHandler) DestroyTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("id invalid. Must be an integeger")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id. ID Must be an integer"})
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

	c.JSON(http.StatusOK, response)
}
