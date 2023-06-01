package main

import (
	"net/http"
	"time"

	"greenlight.aenkas.org/internal/data"
)

func (self *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}
	self.writeJSON(w, http.StatusOK, envelope{"movie": input}, nil)
}

func (self *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := self.readIDParam(r)
	if err != nil || id < 1 {
		self.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}
