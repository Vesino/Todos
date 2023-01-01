package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Vesino/todos/internal/data"
)

type Todo struct {
	Todo        string `json:"todo"`
	Description string `json:"description"`
}

func (app *application) listTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := app.models.Todos.GetAll()

	err = app.writeJSON(w, http.StatusOK, envelope{"todos": todos}, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Todo        string `json:"todo"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	todo := data.Todo{
		Todo:        input.Todo,
		Description: input.Description,
	}

	err = app.models.Todos.Insert(&todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", todo.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	todo, err := app.models.Todos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	todo, err := app.models.Todos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// declaration of an input to hold the data from the user
	var input struct {
		ID          int64     `json:"id"`
		Todo        string    `json:"todo"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		IsDone      bool      `json:"is_done"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	todo.Todo = input.Todo
	todo.CreatedAt = input.CreatedAt
	todo.Description = input.Description
	todo.IsDone = input.IsDone

	err = app.models.Todos.Update(todo)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.notFoundResponse(w, r)
	}

}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Todos.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, envelope{}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
