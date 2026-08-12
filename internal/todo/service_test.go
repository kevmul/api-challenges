package todo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockStore implements TodoStorer for testing.
// Each field is a function so individual tests can control exactly what the store returns.
type mockStore struct {
	getAllFn   func() ([]Todo, error)
	getByIdFn func(id int64) (*Todo, error)
	createFn  func(title string) (*Todo, error)
	updateFn  func(id int64, updatedTodo Todo) (*Todo, error)
	patchFn   func(id int64) (*Todo, error)
	destroyFn func(id int64) error
}

func (m *mockStore) GetAll() ([]Todo, error)                        { return m.getAllFn() }
func (m *mockStore) GetById(id int64) (*Todo, error)                { return m.getByIdFn(id) }
func (m *mockStore) Create(title string) (*Todo, error)             { return m.createFn(title) }
func (m *mockStore) Update(id int64, t Todo) (*Todo, error)         { return m.updateFn(id, t) }
func (m *mockStore) Patch(id int64) (*Todo, error)                  { return m.patchFn(id) }
func (m *mockStore) Destroy(id int64) error                         { return m.destroyFn(id) }

// --- GetAll ---

func TestGetAll(t *testing.T) {
	store := &mockStore{
		getAllFn: func() ([]Todo, error) {
			return []Todo{
				{ID: 1, Title: "Learn net/http", Done: false},
				{ID: 2, Title: "Learn JSON encoding", Done: false},
			}, nil
		},
	}
	svc := NewTodoService(store)

	todos, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, "Learn net/http", todos[0].Title)
}

func TestGetAll_StoreError(t *testing.T) {
	store := &mockStore{
		getAllFn: func() ([]Todo, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewTodoService(store)

	todos, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, todos)
}

// --- GetById ---

func TestGetById(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		storeTodo  *Todo
		storeErr   error
		wantErr    bool
		wantErrVal error
	}{
		{
			name:      "existing todo",
			id:        1,
			storeTodo: &Todo{ID: 1, Title: "Learn net/http", Done: false},
			storeErr:  nil,
			wantErr:   false,
		},
		{
			name:       "non-existing todo",
			id:         999,
			storeTodo:  nil,
			storeErr:   ErrNotFound,
			wantErr:    true,
			wantErrVal: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				getByIdFn: func(id int64) (*Todo, error) {
					return tt.storeTodo, tt.storeErr
				},
			}
			svc := NewTodoService(store)

			todo, err := svc.GetById(tt.id)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErrVal, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.id, todo.ID)
			}
		})
	}
}

// --- Create ---

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		storeTodo *Todo
		storeErr  error
		wantErr   bool
	}{
		{
			name:      "successful create",
			title:     "New Todo",
			storeTodo: &Todo{ID: 1, Title: "New Todo", Done: false},
			storeErr:  nil,
			wantErr:   false,
		},
		{
			name:      "store returns error",
			title:     "New Todo",
			storeTodo: nil,
			storeErr:  errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				createFn: func(title string) (*Todo, error) {
					return tt.storeTodo, tt.storeErr
				},
			}
			svc := NewTodoService(store)

			todo, err := svc.Create(tt.title)

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.title, todo.Title)
				assert.False(t, todo.Done)
			}
		})
	}
}

// --- Update ---

func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		updatedTodo Todo
		storeTodo   *Todo
		storeErr    error
		wantErr     bool
	}{
		{
			name: "update existing todo",
			id:   1,
			updatedTodo: Todo{Title: "Updated Title", Done: true},
			storeTodo:   &Todo{ID: 1, Title: "Updated Title", Done: true},
			storeErr:    nil,
			wantErr:     false,
		},
		{
			name:        "update non-existing todo",
			id:          999,
			updatedTodo: Todo{Title: "Updated Title"},
			storeTodo:   nil,
			storeErr:    ErrNotFound,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				updateFn: func(id int64, t Todo) (*Todo, error) {
					return tt.storeTodo, tt.storeErr
				},
			}
			svc := NewTodoService(store)

			got, err := svc.Update(tt.id, tt.updatedTodo)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, ErrNotFound, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.updatedTodo.Title, got.Title)
				assert.Equal(t, tt.updatedTodo.Done, got.Done)
			}
		})
	}
}

// --- Patch (toggle done) ---

func TestPatch(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		storeTodo  *Todo
		storeErr   error
		wantDone   bool
		wantErr    bool
	}{
		{
			name:      "toggle false to true",
			id:        1,
			storeTodo: &Todo{ID: 1, Title: "Test Todo", Done: true},
			storeErr:  nil,
			wantDone:  true,
			wantErr:   false,
		},
		{
			name:      "toggle true to false",
			id:        2,
			storeTodo: &Todo{ID: 2, Title: "Done Todo", Done: false},
			storeErr:  nil,
			wantDone:  false,
			wantErr:   false,
		},
		{
			name:      "non-existing todo",
			id:        99,
			storeTodo: nil,
			storeErr:  ErrNotFound,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				patchFn: func(id int64) (*Todo, error) {
					return tt.storeTodo, tt.storeErr
				},
			}
			svc := NewTodoService(store)

			got, err := svc.Patch(tt.id)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, ErrNotFound, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.wantDone, got.Done)
			}
		})
	}
}

// --- Destroy ---

func TestDestroy(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		storeErr error
		wantErr  bool
	}{
		{
			name:     "delete existing todo",
			id:       1,
			storeErr: nil,
			wantErr:  false,
		},
		{
			name:     "delete non-existing todo",
			id:       999,
			storeErr: ErrNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{
				destroyFn: func(id int64) error {
					return tt.storeErr
				},
			}
			svc := NewTodoService(store)

			err := svc.Destroy(tt.id)

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
