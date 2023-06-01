package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(w, "get movie %d handler", id)
}
