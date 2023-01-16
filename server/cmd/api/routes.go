package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Todos
	router.HandlerFunc(http.MethodGet, "/v1/todos", app.listTodos)
	router.HandlerFunc(http.MethodPost, "/v1/todos", app.createTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todos/:id", app.showTodoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/todos/:id", app.updateTodoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todos/:id", app.deleteTodoHandler)

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.enableCORS(router)
}
