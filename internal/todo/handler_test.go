package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	svc := NewTodoService()
	h := NewTodoHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()

	h.GetTodos(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got []Todo
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "Learn net/http", got[0].Title)
}

func TestGetTodosHandler_Empty(t *testing.T) {
	svc := &todoService{Todos: []Todo{}} // explicitly empty
	h := NewTodoHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()

	h.GetTodos(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", strings.TrimSpace(rec.Body.String()))
}

func TestGetTodoById(t *testing.T) {}

func TestCreateTodo(t *testing.T) {
	svc := NewTodoService()
	h := NewTodoHandler(svc)

	body := strings.NewReader(`{"title": "New Todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	rec := httptest.NewRecorder()

	h.CreateTodo(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got Todo
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, "New Todo", got.Title)
}

func TestDestroyTodo(t *testing.T) {
	svc := NewTodoService()
	h := NewTodoHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req.SetPathValue("id", "1") // manually inject what the mux would have set
	rec := httptest.NewRecorder()

	h.DestroyTodo(rec, req)

	var got Todo
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.JSONEq(t, rec.Body.String(), `{"deleted": {"id" : 1, "title": "Learn net/http", "done": false}}`)

	assert.Len(t, svc.Todos, 1)
}
