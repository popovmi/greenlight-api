package main

import (
	"errors"
	"net/http"
	"time"

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
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
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

	token, err := self.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}

	self.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = self.mailer.Send(user.Email, "user_welcome.html", data)
		if err != nil {
			self.logger.PrintError(err, nil)
		}
	})

	err = self.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}

func (self *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlainText); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := self.models.Users.GetByToken(input.TokenPlainText, data.ScopeActivation)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			self.failedValidationResponse(w, r, v.Errors)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	err = self.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			self.editConflictResponse(w, r)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	err = self.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}
