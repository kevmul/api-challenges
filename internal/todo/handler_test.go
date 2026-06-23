package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestTodoHandler_GetTodo(t *testing.T) {
	tests := []struct {
		name           string // description of this test case
		id             int
		expectedStatus int
		expectedBody   string
		wantErr        bool
	}{
		{
			name:           "existing todo",
			id:             1,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id" : 1, "title": "Learn net/http", "done": false}`,
			wantErr:        false,
		},
		{
			name:           "non-existing todo",
			id:             999,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `todo not found`,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService()
			h := NewTodoHandler(svc)

			// Convert the integer ID to a string for the request URL
			id := strconv.Itoa(tt.id)

			req := httptest.NewRequest(http.MethodDelete, "/todos/"+id, nil)
			req.SetPathValue("id", id) // manually inject what the mux would have set
			rec := httptest.NewRecorder()

			h.GetTodo(rec, req)

			var got Todo

			err := json.Unmarshal(rec.Body.Bytes(), &got)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
				return
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				assert.Nil(t, err)
				assert.JSONEq(t, rec.Body.String(), `{"id" : 1, "title": "Learn net/http", "done": false}`)
			}
		})
	}
}
