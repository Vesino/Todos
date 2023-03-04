package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Todo        string    `json:"todo"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsDone      bool      `json:"is_done"`
	UserId      int64     `json:"user_id"`
}

// Define a TodoModel struct type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (todo, description, user_id)
		VALUES  ($1, $2, $3)
		RETURNING id, created_at
	`

	args := []interface{}{todo.Todo, todo.Description, todo.UserId}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt)
}

func (m TodoModel) Get(id int64) (*Todo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, todo, description, is_done
	FROM todos
	WHERE id = $1
	`

	var todo Todo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Todo,
		&todo.Description,
		&todo.IsDone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &todo, nil
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

func (t TodoModel) Update(todo *Todo) error {

	query := `
	UPDATE todos
	SET todo = $1, description = $2, created_at = $3, is_done = $4
	WHERE id = $5
	RETURNING todo
	`
	args := []interface{}{
		todo.Todo,
		todo.Description,
		todo.CreatedAt,
		todo.IsDone,
		todo.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Todo)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (t TodoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// SQL query to delete the record
	query := `
	DELETE FROM todos
	WHERE id = $1
	`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := t.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
