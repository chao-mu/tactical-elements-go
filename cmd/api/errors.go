package main

import (
	"errors"
	"net/http"
)

func (app *application) writeServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := errors.New("the server encountered a problem and could not process your request")
	app.writeError(w, r, http.StatusInternalServerError, message)
}

func (app *application) writeError(w http.ResponseWriter, r *http.Request, status int, err error) {
	resp := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "error",
		Message: err.Error(),
	}

	err = app.writeJSON(w, status, resp)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}
