package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (self *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(self.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(self.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", self.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", self.getMoviesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", self.getMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", self.createMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", self.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", self.deleteMovieHandler)

	return self.recoverPanic(router)
}
