package main

import (
	"fmt"
	"net/http"

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
