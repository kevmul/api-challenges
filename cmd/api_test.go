package main

import (
	"api-challenges/internal/todo"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// Test for getting all Todos
// r.GET("/todos/", handler.GetTodos)
func TestApi_Get_Todos(t *testing.T) {
	api := application{port: ":8080"}

	router := api.mount()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todos []todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todos)
	assert.Len(t, todos, 2)
	assert.Equal(t, "Learn net/http", todos[0].Title)
	assert.False(t, todos[0].Done)
}

// Test for getting a Single Todo
// r.GET("/todos/:id", handler.GetTodo)
func TestApi_Get_Todo(t *testing.T) {
	api := application{port: ":8080"}

	router := api.mount()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/2", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todo todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todo)
	assert.Equal(t, "Learn JSON encoding", todo.Title)
	assert.False(t, todo.Done)
}

// Test for Creating a Todo
// r.POST("/todos/", handler.CreateTodo)
func TestApi_Post_Todo(t *testing.T) {
	api := application{port: ":8080"}

	router := api.mount()

	body := `{"title": "Buy Groceries"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todos/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var todo todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todo)
	assert.Equal(t, "Buy Groceries", todo.Title)
	assert.False(t, todo.Done)
}

// Test for Toggling a Todo's Done status
// r.PUT("/todos/:id", handler.ToggleTodoState) // Toggle done state
func TestApi_Put_Todo(t *testing.T) {
	api := application{port: ":8080"}

	router := api.mount()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	updated, ok := response["updated"].(map[string]any)
	assert.True(t, ok, "expected 'updated' key in response")
	assert.Equal(t, "Learn net/http", updated["title"])
	assert.Equal(t, true, updated["done"])
}

// Test for Updating a whole Todo
// r.PATCH("/todos/:id", handler.UpdateTodo)    // Update all status
func TestApi_Patch_Todo(t *testing.T) {
	api := application{port: ":8080"}
	router := api.mount()

	body := `{"title": "Updated", "done": true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/todos/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	_, ok := response["updated"].(map[string]any)
	assert.True(t, ok, "expected 'updated' key in response")
}

// Test for Deleting a Todo
// r.DELETE("/todos/:id", handler.DestroyTodo)
func TestApi_Delete_Todo(t *testing.T) {
	api := application{port: ":8080"}

	router := api.mount()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	_, ok := response["deleted"]
	assert.True(t, ok, "expected 'deleted' key in response")
}
