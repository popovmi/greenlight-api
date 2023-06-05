package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.aenkas.org/internal/data"
	"greenlight.aenkas.org/internal/validator"
)

func (self *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := self.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			self.invalidCredentialsResponse(w, r)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		self.invalidCredentialsResponse(w, r)
		return
	}

	token, err := self.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	err = self.writeJSON(w, http.StatusCreated, envelope{"token": token}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}
