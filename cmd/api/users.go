package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/athifirshad/eucalyptus/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"usertype"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		UserType:  data.UserType(input.UserType),
		Activated: false,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//	userType, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// if userType == 1 {	//1 - patient, 2 - doctor, 3 - admin
	err = app.models.Users.CreateUser(user, "patient")
	// } else {
	//
	//
	//
	//
	// }

	// if err != nil {
	// 	app.badRequestResponse(w, r, err)
	// 	return
	// }
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"activationToken": token.Plaintext,
		"userID":          user.ID,
	}

	err = app.mailer.Send(user.Email, "user_welcome.htm", data)
	if err != nil {
		app.logger.Sugar().Error(err, nil)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getLoggedInUserHandler(w http.ResponseWriter, r *http.Request) {
    user := app.contextGetUser(r)

    response := struct {
        ID        int64     `json:"id"`
        Name      string    `json:"name"`
        Email     string    `json:"email"`
        Activated bool      `json:"activated"`
        UserType data.UserType `json:"userType"`
    }{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Activated: user.Activated,
        UserType: user.UserType,
    }

    // Use the writeJSON helper function to write the response.
    err := app.writeJSON(w, http.StatusOK, envelope{"user": response}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
