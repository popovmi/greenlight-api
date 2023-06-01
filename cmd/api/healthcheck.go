package main

import (
	"net/http"
)

func (self *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available", "environment": self.config.env, "version": version,
	}

	err := self.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		if err != nil {
			self.logger.Println(err)
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
	}
}
