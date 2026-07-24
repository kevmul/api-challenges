package main

import (
	"api-challenges/internal/todo"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	suite.Suite
	router http.Handler
}

func (s *ApiTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	api := application{port: ":8080"}
	s.router = api.mount()
}

// Test for getting all Todos
// r.GET("/todos/", handler.GetTodos)
func (s *ApiTestSuite) TestApi_Get_Todos() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var todos []todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todos)
	s.Len(todos, 2)
	s.Equal("Learn net/http", todos[0].Title)
	s.False(todos[0].Done)
}

// Test for getting a Single Todo
// r.GET("/todos/:id", handler.GetTodo)
func (s *ApiTestSuite) TestApi_Get_Todo() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/2", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var todo todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todo)
	s.Equal("Learn JSON encoding", todo.Title)
	s.False(todo.Done)
}

// Test for Creating a Todo
// r.POST("/todos/", handler.CreateTodo)
func (s *ApiTestSuite) TestApi_Post_Todo() {
	body := `{"title": "Buy Groceries"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todos/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var todo todo.Todo
	json.Unmarshal(w.Body.Bytes(), &todo)
	s.Equal("Buy Groceries", todo.Title)
	s.False(todo.Done)
}

// Test for Toggling a Todo's Done status
// r.PUT("/todos/:id", handler.ToggleTodoState) // Toggle done state
func (s *ApiTestSuite) TestApi_Put_Todo() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	updated, ok := response["updated"].(map[string]any)
	s.True(ok, "expected 'updated' key in response")
	s.Equal("Learn net/http", updated["title"])
	s.Equal(true, updated["done"])
}

// Test for Updating a whole Todo
// r.PATCH("/todos/:id", handler.UpdateTodo)    // Update all status
func (s *ApiTestSuite) TestApi_Patch_Todo() {
	body := `{"title": "Updated", "done": true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/todos/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	_, ok := response["updated"].(map[string]any)
	s.True(ok, "expected 'updated' key in response")
}

// Test for Deleting a Todo
// r.DELETE("/todos/:id", handler.DestroyTodo)
func (s *ApiTestSuite) TestApi_Delete_Todo() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	_, ok := response["deleted"]
	s.True(ok, "expected 'deleted' key in response")
}

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
