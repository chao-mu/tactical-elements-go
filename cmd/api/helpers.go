package main

import (
	"encoding/json"
	"net/http"
)

type responseData map[string]any

func (app *application) writeSuccess(w http.ResponseWriter, r *http.Request, data responseData) {
	resp := struct {
		Status string       `json:"status"`
		Data   responseData `json:"data"`
	}{
		Status: "success",
		Data:   data,
	}

	err := app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// For friendlier terminal output
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
