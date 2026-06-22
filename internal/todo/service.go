package todo

import "errors"

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var ErrNotFound = errors.New("todo not found")

func getAll() ([]Todo, error) {
	todos := []Todo{
		{ID: 1, Title: "Learn net/http", Done: false},
		{ID: 2, Title: "Learn JSON encoding", Done: false},
	}

	return todos, nil
}

func getById(id int) (Todo, error) {
	return Todo{}, ErrNotFound
}
