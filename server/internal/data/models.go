package data

import "database/sql"

type Models struct {
	Todos TodoModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Todos: TodoModel{DB: db},
	}
}
