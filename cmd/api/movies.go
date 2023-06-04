package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.aenkas.org/internal/data"
	"greenlight.aenkas.org/internal/validator"
)

func (self *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = self.models.Movies.Insert(movie)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	self.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}

func (self *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := self.readIDParam(r)
	if err != nil || id < 1 {
		self.notFoundResponse(w, r)
		return
	}

	movie, err := self.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			self.notFoundResponse(w, r)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}

func (self *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := self.readIDParam(r)
	if err != nil {
		self.notFoundResponse(w, r)
		return
	}

	movie, err := self.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			self.notFoundResponse(w, r)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err = self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = self.models.Movies.Update(movie)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}

func (self *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := self.readIDParam(r)
	if err != nil {
		self.notFoundResponse(w, r)
		return
	}

	err = self.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			self.notFoundResponse(w, r)
		default:
			self.serverErrorResponse(w, r, err)
		}
		return
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"message": "movie succesfully deleted"}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}
