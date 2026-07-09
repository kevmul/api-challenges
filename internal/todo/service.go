package todo

import (
	"errors"
	"slices"
)

type todoService struct {
	Todos []Todo
}

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var ErrNotFound = errors.New("todo not found")

func NewTodoService() *todoService {
	return &todoService{
		Todos: []Todo{
			{ID: 1, Title: "Learn net/http", Done: false},
			{ID: 2, Title: "Learn JSON encoding", Done: false},
		},
	}
}

func (s *todoService) GetAll() ([]Todo, error) {
	if len(s.Todos) == 0 {
		return nil, ErrNotFound
	}
	return s.Todos, nil
}

func (s *todoService) GetById(id int) (Todo, error) {

	for _, todo := range s.Todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return Todo{}, ErrNotFound
}

func (s *todoService) Create(todoTitle string) (Todo, error) {

	todo := Todo{
		ID:    len(s.Todos) + 1,
		Title: todoTitle,
		Done:  false,
	}

	s.Todos = append(s.Todos, todo)

	return todo, nil
}

func (s *todoService) Update(id int, updatedTodo Todo) (Todo, error) {
	for i, todo := range s.Todos {
		if todo.ID == id {
			s.Todos[i].Title = updatedTodo.Title
			s.Todos[i].Done = updatedTodo.Done
			return s.Todos[i], nil
		}
	}

	return Todo{}, ErrNotFound
}

func (s *todoService) Destroy(id int) error {
	s.Todos = slices.DeleteFunc(s.Todos, func(t Todo) bool {
		return t.ID == id
	})

	return nil
}
