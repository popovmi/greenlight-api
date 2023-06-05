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

	router.HandlerFunc(http.MethodPost, "/v1/users/", self.signupUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", self.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", self.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", self.requirePermission("movies:read", self.getMoviesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", self.requirePermission("movies:read", self.getMovieHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", self.requirePermission("movies:write", self.createMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", self.requirePermission("movies:write", self.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", self.requirePermission("movies:write", self.deleteMovieHandler))

	return self.recoverPanic(self.rateLimit(self.authenticate(router)))
}
