package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tactical-elements/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var puzzle = models.Puzzle{
	ID:        0,
	FEN:       "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
	Solutions: []string{"a1", "a8", "e4", "e5", "h1", "h8"},
}

var puzzleHandler = PuzzleHandler{
	Collection: models.PuzzleCollection{
		Puzzles: []models.Puzzle{puzzle},
	},
}

func TestGuessCorrect(t *testing.T) {
	response := requestGuess(t, puzzle.ID, map[string]interface{}{
		"browserId": "test",
		"solution":  puzzle.Solutions[0],
	})

	assert.True(t, response.Correct)
}

func TestGuessIncorrect(t *testing.T) {
	response := requestGuess(t, puzzle.ID, map[string]interface{}{
		"browserId": "test",
		"solution":  "incorrect",
	})

	assert.False(t, response.Correct)
}

func TestAdvanceCorrect(t *testing.T) {
	response := requestAdvance(t, puzzle.ID, map[string]interface{}{
		"browserId": "test",
		"solutions": puzzle.Solutions,
	})

	assert.True(t, response.Correct)
}

func TestAdvanceIncorrect(t *testing.T) {
	response := requestAdvance(t, puzzle.ID, map[string]interface{}{
		"browserId": "test",
		"solutions": []string{"incorrect"},
	})

	assert.False(t, response.Correct)
}

func makeContext(t *testing.T, body map[string]interface{}) (echo.Context, *httptest.ResponseRecorder) {
	marshaledBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshaledBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	return ctx, rec
}

func readResponse(t *testing.T, err error, rec *httptest.ResponseRecorder, response any) {
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, json.Valid(rec.Body.Bytes()))
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
}

func requestGuess(t *testing.T, puzzleId int, body map[string]interface{}) GuessResponse {
	ctx, rec := makeContext(t, body)

	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(puzzleId))

	response := GuessResponse{}
	readResponse(t, puzzleHandler.Guess(ctx), rec, &response)

	return response
}

func requestAdvance(t *testing.T, puzzleId int, body map[string]interface{}) AdvanceResponse {
	ctx, rec := makeContext(t, body)

	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(puzzleId))

	response := AdvanceResponse{}
	readResponse(t, puzzleHandler.Advance(ctx), rec, &response)

	return response
}

func requestGiveUp(t *testing.T, puzzleId int, body map[string]interface{}) GiveUpResponse {
	ctx, rec := makeContext(t, body)

	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(puzzleId))

	response := GiveUpResponse{}
	readResponse(t, puzzleHandler.Advance(ctx), rec, &response)

	return response
}
