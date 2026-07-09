package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := NewTodoService()
	h := NewTodoHandler(svc)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/todos", nil)

	h.GetTodos(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var got []Todo
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "Learn net/http", got[0].Title)
}

func TestGetTodosHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &todoService{Todos: []Todo{}} // explicitly empty
	h := NewTodoHandler(svc)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/todos", nil)

	h.GetTodos(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", strings.TrimSpace(rec.Body.String()))
}

func TestGetTodoById(t *testing.T) {}

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := NewTodoService()
	h := NewTodoHandler(svc)

	body := strings.NewReader(`{"title": "New Todo"}`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/todos", body)

	h.CreateTodo(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var got Todo
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, "New Todo", got.Title)
}

func TestDestroyTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := NewTodoService()
	h := NewTodoHandler(svc)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/todos/1", nil)

	h.DestroyTodo(c)

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
			expectedBody:   `{"error":"todo not found"}`,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := NewTodoService()
			h := NewTodoHandler(svc)

			// Convert the integer ID to a string for the request URL
			id := strconv.Itoa(tt.id)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodDelete, "/todos/"+id, nil)

			h.GetTodo(c)

			var got Todo

			json.Unmarshal(rec.Body.Bytes(), &got)
			if tt.wantErr {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
				return
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
				assert.JSONEq(t, rec.Body.String(), `{"id" : 1, "title": "Learn net/http", "done": false}`)
			}
		})
	}
}

func TestTodoHandler_UpdateTodo(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		id      int
		request string
		errCode int
		err     error
		wantErr bool
	}{
		// {
		// 	name:    "update existing todo",
		// 	id:      1,
		// 	request: `{"title": "Updated Todo"}`,
		// 	errCode: http.StatusOK,
		// 	err:     nil,
		// 	wantErr: false,
		// },
		{
			name:    "update non-existing todo",
			id:      999,
			request: `{"title": "Should Not Exist"}`,
			errCode: http.StatusNotFound,
			err:     errors.New(`{"error":"todo not found"}`),
			wantErr: true,
		},
		// {
		// 	name:    "update with empty title",
		// 	id:      1,
		// 	request: `{"title": ""}`,
		// 	errCode: http.StatusBadRequest,
		// 	err:     errors.New(`{"error":"title is required"}`),
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := NewTodoService()
			h := NewTodoHandler(svc)

			// Convert the integer ID to a string for the request URL
			id := strconv.Itoa(tt.id)

			body := strings.NewReader(tt.request)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodPatch, "/todos/"+id, body)

			h.UpdateTodo(c)

			var got Todo

			json.Unmarshal(rec.Body.Bytes(), &got)

			if tt.wantErr {
				assert.Equal(t, tt.errCode, rec.Code)
				assert.Equal(t, tt.err.Error(), strings.TrimSpace(rec.Body.String()))
			} else {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
				assert.JSONEq(t, `{"id" : 1, "title": "Updated Todo", "done": false}`, rec.Body.String())
			}
		})
	}
}

func TestTodoHandler_PatchTodo(t *testing.T) {
	tests := []struct {
		name           string // description of this test case
		id             int
		expected       string
		expectedStatus int
		err            error
		wantErr        bool
	}{
		{
			name:           "Patch an existing todo that is not done",
			id:             1,
			expected:       `{"updated": {"id": 1, "title": "Undone Todo", "done": true}}`,
			expectedStatus: http.StatusOK,
			err:            nil,
			wantErr:        false,
		},
		{
			name:           "Patch an existing todo that is done",
			id:             2,
			expected:       `{"updated": {"id": 2, "title": "Done Todo", "done": false}}`,
			expectedStatus: http.StatusOK,
			err:            nil,
			wantErr:        false,
		},
		{
			name:           "Patch an invalid todo",
			id:             99,
			expected:       ``,
			expectedStatus: http.StatusNotFound,
			err:            ErrNotFound,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := &todoService{
				Todos: []Todo{
					{ID: 1, Title: "Undone Todo", Done: false},
					{ID: 2, Title: "Done Todo", Done: true},
				},
			}
			h := NewTodoHandler(svc)

			// Convert the integer ID to a string for the request URL
			id := strconv.Itoa(tt.id)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodPatch, "/todos/"+id, nil)

			h.ToggleTodoState(c)

			var got Todo

			json.Unmarshal(rec.Body.Bytes(), &got)

			if tt.wantErr {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, `{"error":"`+tt.err.Error()+`"}`, strings.TrimSpace(rec.Body.String()))
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
				assert.JSONEq(t, tt.expected, rec.Body.String())
			}
		})
	}
}
