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
	router.HandlerFunc(http.MethodGet, "/v1/todos", app.requireActivatedUser(app.listTodos))
	router.HandlerFunc(http.MethodPost, "/v1/todos", app.requireActivatedUser(app.createTodoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/todos/:id", app.requireActivatedUser(app.showTodoHandler))
	router.HandlerFunc(http.MethodPut, "/v1/todos/:id", app.requireActivatedUser(app.updateTodoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/todos/:id", app.requireActivatedUser(app.deleteTodoHandler))

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Use the authenticate() middleware on all requests.
	return app.enableCORS(app.recoverPanic(app.rateLimit(app.authenticate(router))))
}
