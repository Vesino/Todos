package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Todo        string    `json:"todo"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsDone      bool      `json:"is_done"`
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

func (m TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT * FROM todos
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println("Error when quering", err)
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}
	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.CreatedAt,
			&todo.Todo,
			&todo.Description,
			&todo.IsDone,
		)
		if err != nil {
			log.Println("Error when scanning rows", err)
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
