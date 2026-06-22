package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	s := NewTodoService()
	todos, err := s.getAll()

	assert.Nil(t, err, "getAll should not return an error")
	assert.Len(t, todos, 2, "getAll should return 2 todos")
	assert.Equal(t, "Learn net/http", todos[0].Title, "First todo title should be 'Learn net/http'")
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		id      int
		wantErr bool
	}{
		{
			name:    "existing todo",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non-existing todo",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTodoService()
			todo, err := s.getById(tt.id)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, ErrNotFound, err)
			} else {
				assert.Nil(t, err, "getByID should not return an error")
				assert.Equal(t, tt.id, todo.ID)
			}
		})
	}
}

func Test_create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		todoTitle string
		wantErr   bool
	}{
		{
			todoTitle: "New Todo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTodoService()
			assert.Len(t, s.Todos, 2, "Initial number of todos should be 2")
			todo, err := s.create(tt.todoTitle)
			assert.Nil(t, err, "create should not return an error")
			assert.Equal(t, tt.todoTitle, todo.Title)
			assert.False(t, todo.Done, "Newly created todo should not be done")
		})
	}
}

func TestCreateTodoIncrementsID(t *testing.T) {
	// Create a new todo and check if the ID is incremented correctly.
	s := NewTodoService()
	newTodo, err := s.create("Another Todo")
	assert.Nil(t, err, "create should not return an error")
	assert.Equal(t, 3, newTodo.ID, "Newly created todo should have ID 3")
	todos, _ := s.getAll()
	assert.Len(t, todos, 3, "There should now be 3 todos in total")
}
