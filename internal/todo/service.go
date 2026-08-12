package todo

import (
	"errors"
	"time"
)

type TodoStorer interface {
	GetAll() ([]Todo, error)
	GetById(id int64) (*Todo, error)
	Create(todoTitle string) (*Todo, error)
	Update(id int64, updatedTodo Todo) (*Todo, error)
	Patch(id int64) (*Todo, error)
	Destroy(id int64) error
}

type todoService struct {
	store TodoStorer // interface, not *TodoStore
}

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNotFound = errors.New("todo not found")

func NewTodoService(store TodoStorer) *todoService {
	return &todoService{store: store}
}

func (s *todoService) GetAll() ([]Todo, error) {
	return s.store.GetAll()
}

func (s *todoService) GetById(id int64) (Todo, error) {
	todo, err := s.store.GetById(id)
	if err != nil {
		return Todo{}, err
	}
	return *todo, nil
}

func (s *todoService) Create(title string) (Todo, error) {
	todo, err := s.store.Create(title)
	if err != nil {
		return Todo{}, err
	}
	return *todo, nil
}

func (s *todoService) Update(id int64, updatedTodo Todo) (Todo, error) {
	todo, err := s.store.Update(id, updatedTodo)
	if err != nil {
		return Todo{}, err
	}
	return *todo, nil
}

func (s *todoService) Patch(id int64) (Todo, error) {
	todo, err := s.store.Patch(id)
	if err != nil {
		return Todo{}, err
	}
	return *todo, nil
}

func (s *todoService) Destroy(id int64) error {
	return s.store.Destroy(id)
}
