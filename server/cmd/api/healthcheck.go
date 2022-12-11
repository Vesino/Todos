package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status":     "availble",
		"enviroment": app.config.env,
		"version":    version,
	}
	fmt.Println(data)
	if err := app.writeJson(w, http.StatusOK, data, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
