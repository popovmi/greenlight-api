package main

import (
	"fmt"
	"net/http"
)

func (self *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status: available\n")
	fmt.Fprintf(w, "environment: %s\n", self.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
