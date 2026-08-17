package todo

import (
	"database/sql"
	"errors"
	"time"
)

type TodoStore struct {
	DB *sql.DB
}

func NewTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{DB: db}
}

func (s *TodoStore) Create(title string) (*Todo, error) {
	query := `INSERT INTO todos (title, done, created_at, updated_at) VALUES (?, ?, ?, ?)`
	now := time.Now()

	result, err := s.DB.Exec(query, title, false, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Todo{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *TodoStore) GetAll() ([]Todo, error) {
	rows, err := s.DB.Query(`SELECT id, title, done, created_at, updated_at FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *TodoStore) GetById(id int64) (*Todo, error) {
	var t Todo
	err := s.DB.QueryRow(`SELECT id, title, done, created_at, updated_at FROM todos WHERE id = ?`, id).
		Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TodoStore) Update(id int64, updated Todo) (*Todo, error) {
	now := time.Now().UTC()
	result, err := s.DB.Exec(
		`UPDATE todos SET title = ?, done = ?, updated_at = ? WHERE id = ?`,
		updated.Title, updated.Done, now, id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("no rows affected")
	}

	return s.GetById(id)
}

func (s *TodoStore) Patch(id int64) (*Todo, error) {
	now := time.Now().UTC()
	_, err := s.DB.Exec(
		`UPDATE todos SET done = NOT done, updated_at = ? WHERE id = ?`,
		now, id,
	)
	if err != nil {
		return nil, err
	}
	return s.GetById(id)
}

func (s *TodoStore) Destroy(id int64) error {
	result, err := s.DB.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
