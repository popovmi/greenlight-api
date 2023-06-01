package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.aenkas.org/internal/data"
)

func (self *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create movie handler")
}

func (self *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := self.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
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
		self.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
