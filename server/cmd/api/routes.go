package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/todos/", app.ListTodos)
	router.HandlerFunc(http.MethodPost, "/v1/todos/", app.createTodoHandler)

	return app.enableCORS(router)
}
