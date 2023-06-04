package main

import (
	"fmt"
	"net/http"
)

func (self *application) logError(r *http.Request, err error) {
	self.logger.Println(err)
}

func (self *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := self.writeJSON(w, status, env, nil)
	if err != nil {
		self.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (self *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	self.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	self.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (self *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	self.errorResponse(w, r, http.StatusNotFound, message)
}

func (self *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	self.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (self *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	self.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (self *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	self.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (self *application) editConfilctResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	self.errorResponse(w, r, http.StatusConflict, message)
}