package main

import (
	"net/http"
)

func (self *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"environment": self.config.env,
			"version":     version},
	}

	err := self.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		self.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
