package main

import (
	"net/http"
)

type Todo struct {
	Todo        string `json:"todo"`
	Description string `json:"description"`
}

func (app *application) ListTodos(w http.ResponseWriter, r *http.Request) {

	todos := []Todo{{"queacer", "Hacer queacer"}, {"Comprar cosas", "ir al super a comprar cosas"}}
	err := app.writeJSON(w, http.StatusOK, envelope{"todos": todos}, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}
