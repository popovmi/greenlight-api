package main

import (
	"errors"
	"net/http"

	"greenlight.aenkas.org/internal/data"
	"greenlight.aenkas.org/internal/validator"
)

func (self *application) signupUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := self.readJSON(w, r, &input)

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = self.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			self.failedValidationResponse(w, r, v.Errors)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	self.background(func() {
		err = self.mailer.Send(user.Email, "user_welcome.html", user)
		if err != nil {
			self.logger.PrintError(err, nil)
		}
	})

	err = self.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}
