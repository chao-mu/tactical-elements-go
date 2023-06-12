package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/puzzles/next", app.nextPuzzleHandler)
	router.HandlerFunc(http.MethodPost, "/v1/puzzles/guess", app.guessHandler)
	router.HandlerFunc(http.MethodPost, "/v1/puzzles/advance", app.advanceHandler)
	router.HandlerFunc(http.MethodPost, "/v1/puzzles/give-up", app.giveUpHandler)

	return router
}
