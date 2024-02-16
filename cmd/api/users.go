package main

import (
	"net/http"

	"github.com/athifirshad/eucalyptus/internal/data"
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
		err = app.models.Users.CreateUser(user,"patient")
		// } else {
		//     // Handle other user types if needed
		//     // For example, if there are other user types like "doctor" or "administrator",
		//     // you would call a different method here, e.g., app.models.Users.CreateDoctor(user)
		//     // or app.models.Users.CreateAdministrator(user)
		// }

		// if err != nil {
		// 	app.badRequestResponse(w, r, err)
		// 	return
		// }
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}

	}

