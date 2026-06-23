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
