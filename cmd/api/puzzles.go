package main

import (
	"encoding/json"
	"net/http"

	"github.com/chao-mu/tactical-elements-go/internal/data"
)

func (app *application) guessHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PuzzleID     int    `json:"puzzleId"`
		CollectionID string `json:"collectionId"`
		Solution     string `json:"solution"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, err)
		return
	}

	puzzle, err := data.GetPuzzleWithSolutions(input.CollectionID, input.PuzzleID)
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	correct := puzzle.GuessOne(input.Solution)

	app.writeSuccess(w, r, responseData{
		"correct": correct,
	})
}

func (app *application) advanceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PuzzleID     int      `json:"puzzleId"`
		CollectionID string   `json:"collectionID"`
		Solutions    []string `json:"solutions"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, err)
		return
	}

	puzzle, err := data.GetPuzzleWithSolutions(input.CollectionID, input.PuzzleID)
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	correct := puzzle.GuessAll(input.Solutions)

	app.writeSuccess(w, r, responseData{
		"correct": correct,
	})
}

func (app *application) nextPuzzleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CollectionID string `json:"collectionID"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, err)
		return
	}

	col, err := data.GetPuzzleCollection(input.CollectionID)
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	p := col.RandomPuzzle().Puzzle

	app.writeSuccess(w, r, responseData{
		"puzzle": p,
	})
}

func (app *application) giveUpHandler(w http.ResponseWriter, r *http.Request) {
	app.writeSuccess(w, r, responseData{})
}
