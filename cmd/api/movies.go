package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"greenlight.aenkas.org/internal/data"
	"greenlight.aenkas.org/internal/validator"
)

func (self *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.ListParams
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Title = self.readString(qs, "title", "")
	input.Genres = self.readCSV(qs, "genres", []string{})

	input.ListParams.Page = self.readInt(qs, "page", 1, v)
	input.ListParams.PageSize = self.readInt(qs, "pageSize", 20, v)
	input.ListParams.Sort = self.readString(qs, "sort", "id")
	input.ListParams.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateListParams(v, input.ListParams); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := self.models.Movies.GetMany(input.Title, input.Genres, input.ListParams)
	if err != nil {
		self.serverErrorResponse(w, r, err)
		return
	}

	err = self.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		self.serverErrorResponse(w, r, err)
	}
}

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

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
			self.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = self.readJSON(w, r, &input)
	if err != nil {
		self.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		self.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = self.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			self.editConflictResponse(w, r)
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
