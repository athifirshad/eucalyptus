package main

import (
	"context"
	"net/http"

	"github.com/athifirshad/eucalyptus/internal/data"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userModel := &data.UserModel{
		User: &data.User{
			Name:      input.Name,
			Email:     input.Email,
			UserType:  input.UserType,
			Activated: false,
		},
	}

	err = userModel.SetPassword(input.Password)
	app.Users = &data.UserModel{} // Initialize the Users field with the appropriate value

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.Queries.CreatePatientUser(context.Background(), data.CreatePatientUserParams{
		Name:         userModel.Name,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
	})
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": userModel}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
