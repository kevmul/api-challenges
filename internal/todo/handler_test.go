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

// mockService implements TodoService for handler tests.
// Each field is a function so individual tests control exactly what the service returns.
type mockService struct {
	getAllFn   func() ([]Todo, error)
	getByIdFn func(id int64) (Todo, error)
	createFn  func(title string) (Todo, error)
	updateFn  func(id int64, t Todo) (Todo, error)
	patchFn   func(id int64) (Todo, error)
	destroyFn func(id int64) error
}

func (m *mockService) GetAll() ([]Todo, error)                   { return m.getAllFn() }
func (m *mockService) GetById(id int64) (Todo, error)            { return m.getByIdFn(id) }
func (m *mockService) Create(title string) (Todo, error)         { return m.createFn(title) }
func (m *mockService) Update(id int64, t Todo) (Todo, error)     { return m.updateFn(id, t) }
func (m *mockService) Patch(id int64) (Todo, error)              { return m.patchFn(id) }
func (m *mockService) Destroy(id int64) error                    { return m.destroyFn(id) }

func TestGetTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockService{
		getAllFn: func() ([]Todo, error) {
			return []Todo{
				{ID: 1, Title: "Learn net/http", Done: false},
				{ID: 2, Title: "Learn JSON encoding", Done: false},
			}, nil
		},
	}
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

	svc := &mockService{
		getAllFn: func() ([]Todo, error) {
			return []Todo{}, nil
		},
	}
	h := NewTodoHandler(svc)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/todos", nil)

	h.GetTodos(c)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got []Todo
	json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Len(t, got, 0)
}

func TestGetTodoById(t *testing.T) {}

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockService{
		createFn: func(title string) (Todo, error) {
			return Todo{ID: 1, Title: title, Done: false}, nil
		},
	}
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

	deleted := Todo{ID: 1, Title: "Learn net/http", Done: false}
	svc := &mockService{
		getByIdFn: func(id int64) (Todo, error) {
			return deleted, nil
		},
		destroyFn: func(id int64) error {
			return nil
		},
	}
	h := NewTodoHandler(svc)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/todos/1", nil)

	h.DestroyTodo(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"deleted": {"id": 1, "title": "Learn net/http", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}}`, rec.Body.String())
}

func TestTodoHandler_GetTodo(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		serviceTodo    Todo
		serviceErr     error
		expectedStatus int
		expectedBody   string
		wantErr        bool
	}{
		{
			name:           "existing todo",
			id:             1,
			serviceTodo:    Todo{ID: 1, Title: "Learn net/http", Done: false},
			serviceErr:     nil,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "non-existing todo",
			id:             999,
			serviceTodo:    Todo{},
			serviceErr:     ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"todo not found"}`,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := &mockService{
				getByIdFn: func(id int64) (Todo, error) {
					return tt.serviceTodo, tt.serviceErr
				},
			}
			h := NewTodoHandler(svc)

			id := strconv.Itoa(tt.id)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodGet, "/todos/"+id, nil)

			h.GetTodo(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.wantErr {
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
			} else {
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
				var got Todo
				json.Unmarshal(rec.Body.Bytes(), &got)
				assert.Equal(t, int64(tt.id), got.ID)
				assert.Equal(t, tt.serviceTodo.Title, got.Title)
			}
		})
	}
}

func TestTodoHandler_UpdateTodo(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		request     string
		serviceTodo Todo
		serviceErr  error
		errCode     int
		err         error
		wantErr     bool
	}{
		{
			name:        "update existing todo",
			id:          1,
			request:     `{"title": "Updated Todo"}`,
			serviceTodo: Todo{ID: 1, Title: "Updated Todo", Done: false},
			serviceErr:  nil,
			errCode:     http.StatusOK,
			wantErr:     false,
		},
		{
			name:       "update non-existing todo",
			id:         999,
			request:    `{"title": "Should Not Exist"}`,
			serviceErr: ErrNotFound,
			errCode:    http.StatusNotFound,
			err:        errors.New(`{"error":"todo not found"}`),
			wantErr:    true,
		},
		{
			name:       "update with empty title",
			id:         1,
			request:    `{"title": ""}`,
			serviceErr: nil,
			errCode:    http.StatusBadRequest,
			err:        errors.New(`{"error":"title is required"}`),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := &mockService{
				getByIdFn: func(id int64) (Todo, error) {
					return tt.serviceTodo, tt.serviceErr
				},
				updateFn: func(id int64, t Todo) (Todo, error) {
					return tt.serviceTodo, tt.serviceErr
				},
			}
			h := NewTodoHandler(svc)

			id := strconv.Itoa(tt.id)
			body := strings.NewReader(tt.request)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodPatch, "/todos/"+id, body)

			h.UpdateTodo(c)

			if tt.wantErr {
				assert.Equal(t, tt.errCode, rec.Code)
				assert.Equal(t, tt.err.Error(), strings.TrimSpace(rec.Body.String()))
			} else {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func TestTodoHandler_PatchTodo(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		serviceTodo    Todo
		serviceErr     error
		expected       string
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "toggle false to true",
			id:             1,
			serviceTodo:    Todo{ID: 1, Title: "Undone Todo", Done: true},
			serviceErr:     nil,
			expected:       `{"updated": {"id": 1, "title": "Undone Todo", "done": true, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}}`,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "toggle true to false",
			id:             2,
			serviceTodo:    Todo{ID: 2, Title: "Done Todo", Done: false},
			serviceErr:     nil,
			expected:       `{"updated": {"id": 2, "title": "Done Todo", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}}`,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "non-existing todo",
			id:             99,
			serviceTodo:    Todo{},
			serviceErr:     ErrNotFound,
			expectedStatus: http.StatusNotFound,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := &mockService{
				getByIdFn: func(id int64) (Todo, error) {
					return tt.serviceTodo, tt.serviceErr
				},
				patchFn: func(id int64) (Todo, error) {
					return tt.serviceTodo, tt.serviceErr
				},
			}
			h := NewTodoHandler(svc)

			id := strconv.Itoa(tt.id)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
			c.Request = httptest.NewRequest(http.MethodPut, "/todos/"+id, nil)

			h.ToggleTodoState(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if !tt.wantErr {
				assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
				assert.JSONEq(t, tt.expected, rec.Body.String())
			}
		})
	}
}
