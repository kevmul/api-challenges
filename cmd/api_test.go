package main

import (
	"api-challenges/internal/database"
	"api-challenges/internal/todo"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

	db, err := database.New(":memory:")
	if err != nil {
		s.FailNow("failed to open test db", err)
	}
	if err := database.Migrate(db); err != nil {
		s.FailNow("failed to migrate test db", err)
	}

	api := Application{port: ":8080", DB: db}
	s.router = api.mount()
}

// helper: POST a todo and return the created Todo.
// Fails the test immediately if creation doesn't return 201.
func (s *ApiTestSuite) createTodo(title string) todo.Todo {
	body := strings.NewReader(`{"title": "` + title + `"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todos/", body)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)
	s.Require().Equal(http.StatusCreated, w.Code, "setup: failed to create todo %q", title)

	var created todo.Todo
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &created))
	return created
}

// Test for getting all Todos
// r.GET("/todos/", handler.GetTodos)
func (s *ApiTestSuite) TestApi_Get_Todos() {
	s.createTodo("Learn net/http")
	s.createTodo("Learn JSON encoding")

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
	created := s.createTodo("Learn JSON encoding")
	id := strconv.FormatInt(created.ID, 10)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/"+id, nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var got todo.Todo
	json.Unmarshal(w.Body.Bytes(), &got)
	s.Equal("Learn JSON encoding", got.Title)
	s.False(got.Done)
}

// Test for getting a Todo that doesn't exist
// r.GET("/todos/:id", handler.GetTodo)
func (s *ApiTestSuite) TestApi_Get_Todo_NotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/999", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("todo not found", response["error"])
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

	var created todo.Todo
	json.Unmarshal(w.Body.Bytes(), &created)
	s.Equal("Buy Groceries", created.Title)
	s.False(created.Done)
	s.NotZero(created.ID)
}

// Test that POST with no title returns 400
func (s *ApiTestSuite) TestApi_Post_Todo_MissingTitle() {
	body := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todos/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("title is required", response["error"])
}

// Test for Toggling a Todo's Done status
// r.PATCH("/todos/:id", handler.ToggleTodoState)
func (s *ApiTestSuite) TestApi_Patch_Todo_Toggle() {
	created := s.createTodo("Learn net/http")
	id := strconv.FormatInt(created.ID, 10)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/todos/"+id, nil)
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
// r.PUT("/todos/:id", handler.UpdateTodo)
func (s *ApiTestSuite) TestApi_Put_Todo_Update() {
	created := s.createTodo("Original Title")
	id := strconv.FormatInt(created.ID, 10)

	body := `{"title": "Updated Title", "done": true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/todos/"+id, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	updated, ok := response["updated"].(map[string]any)
	s.True(ok, "expected 'updated' key in response")
	s.Equal("Updated Title", updated["title"])
	s.Equal(true, updated["done"])
}

// Test for Deleting a Todo
// r.DELETE("/todos/:id", handler.DestroyTodo)
func (s *ApiTestSuite) TestApi_Delete_Todo() {
	created := s.createTodo("To Be Deleted")
	id := strconv.FormatInt(created.ID, 10)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todos/"+id, nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	_, ok := response["deleted"]
	s.True(ok, "expected 'deleted' key in response")

	// Verify it's actually gone
	getW := httptest.NewRecorder()
	getReq, _ := http.NewRequest("GET", "/todos/"+id, nil)
	s.router.ServeHTTP(getW, getReq)
	s.Equal(http.StatusNotFound, getW.Code)
}

// Test deleting a Todo that doesn't exist
// r.DELETE("/todos/:id", handler.DestroyTodo)
func (s *ApiTestSuite) TestApi_Delete_Todo_NotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todos/999", nil)
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("todo not found", response["error"])
}

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
