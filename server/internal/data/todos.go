package data

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Todo        string    `json:"todo"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Define a TodoModel struct type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (todo, description)
		VALUES  ($1, $2)
		RETURNING id, created_at
	`

	args := []interface{}{todo.Todo, todo.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt)
}
