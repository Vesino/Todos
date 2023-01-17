package main

import (
	"errors"
	"net/http"

	"github.com/Vesino/todos/internal/data"
	"github.com/Vesino/todos/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// create anonymous struct to hold the expected data from the request body
	var input struct {
		Name string `json:"name"`
		Email string `jason:"email"`
		Password string `json:"password"`
	}

	// parse the request body into the anonymouse input
	err := app.readJSON(w,r,&input)
	if err != nil {
		app.badRequestResponse(w,r,err)
		return
	}

	user := &data.User{
		Name: input.Name,
		Email: input.Email,
		Activated: false,
	}

	// Password.Set() method to generate and store the hashed and plaintext passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w,r,err)	
		return
	}
	v := validator.New()

	// Validate the user struct and return the error mesaage to the client if any of the checks fails
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w,r,v.Errors)
		return
	}

	// Insert the user into the DB
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicatedEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w,r,v.Errors)
		default:
			app.serverErrorResponse(w,r,err)
		}
		return
	}

	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(w,r,err)
		return
	}
	// write a JSON response containing the user data along with the 201 Created Status code
	err = app.writeJSON(w, 201, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w,r,err)
	}

}
