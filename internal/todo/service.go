package todo

import "errors"

//	type TodoService interface {
//		GetAll() ([]Todo, error)
//		GetByID(id int) (Todo, error)
//		Create(todoTitle string) (Todo, error)
//	}
type TodoService struct {
	Todos []Todo
}

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var ErrNotFound = errors.New("todo not found")

func NewTodoService() *TodoService {
	return &TodoService{
		Todos: []Todo{
			{ID: 1, Title: "Learn net/http", Done: false},
			{ID: 2, Title: "Learn JSON encoding", Done: false},
		},
	}
}

func (s *TodoService) getAll() ([]Todo, error) {
	if len(s.Todos) == 0 {
		return nil, ErrNotFound
	}
	return s.Todos, nil
}

func (s *TodoService) getById(id int) (Todo, error) {

	for _, todo := range s.Todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return Todo{}, ErrNotFound
}

func (s *TodoService) create(todoTitle string) (Todo, error) {

	todo := Todo{
		ID:    len(s.Todos) + 1,
		Title: todoTitle,
		Done:  false,
	}

	s.Todos = append(s.Todos, todo)

	return todo, nil
}
