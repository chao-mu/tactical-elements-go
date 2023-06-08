package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"tactical-elements/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GuessResponse struct {
	Correct bool `json:"correct"`
}

type AdvanceResponse struct {
	Correct bool                `json:"correct"`
	Puzzle  models.PuzzlePublic `json:"puzzle"`
}

type GiveUpResponse struct {
	Puzzle models.PuzzlePublic `json:"puzzle"`
}

type PuzzleHandler struct {
	Collection models.PuzzleCollection
}

type PuzzleAction struct {
	BrowserID string `json:"browserId"`
	PuzzleID  int    `json:"puzzleId"`
}

type GuessRequest struct {
	PuzzleAction
	Solution string `json:"solution"`
}

type GiveUpRequest struct {
	PuzzleAction
}

type AndvanceRequest struct {
	PuzzleAction
	Solutions []string `json:"solutions"`
}

func (h PuzzleHandler) Guess(c echo.Context) error {
	id := c.Param("id")
	puzzleID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	puzzle := h.Collection.Puzzles[puzzleID]
	action := new(GuessRequest)
	if err := c.Bind(action); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	correct := puzzle.GuessOne(action.Solution)

	if correct {
		// Add exp
		// Add gold
	} else {
		// Remove health
	}

	return c.JSON(http.StatusOK, GuessResponse{
		Correct: correct,
	})
}

func (h PuzzleHandler) Advance(c echo.Context) error {
	id := c.Param("id")
	puzzleID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	puzzle := h.Collection.Puzzles[puzzleID]
	action := new(AndvanceRequest)
	if err := c.Bind(action); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if puzzle.GuessAll(action.Solutions) {
		return c.JSON(http.StatusOK, AdvanceResponse{
			Correct: true,
			Puzzle:  h.Collection.RandomPuzzle().Public(),
		})
	}

	return c.JSON(http.StatusOK, AdvanceResponse{
		Correct: false,
		Puzzle:  puzzle.Public(),
	})
}

func (h PuzzleHandler) NextPuzzle(c echo.Context) error {
	puzzle := h.Collection.RandomPuzzle()
	return c.JSON(http.StatusOK, puzzle.Public())
}

func (h PuzzleHandler) GiveUp(c echo.Context) error {
	nextPuzzle := h.Collection.RandomPuzzle()
	return c.JSON(http.StatusOK, GiveUpResponse{Puzzle: nextPuzzle.Public()})
}

func GetEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	puzzlesGroup := e.Group("/puzzles")

	collectionIDs := []string{
		"checkables",
		"checks-captures",
		"knight-forkable",
		"undefended",
	}

	for _, collectionID := range collectionIDs {
		path := fmt.Sprintf("assets/puzzles/%s.json", collectionID)

		puzzles, err := models.ReadPuzzles(path)
		if err != nil {
			e.Logger.Fatalf("Failed to read puzzles from %s: %s", path, err)
		}

		collection := models.PuzzleCollection{
			ID:      collectionID,
			Puzzles: puzzles,
		}

		handler := PuzzleHandler{
			Collection: collection,
		}

		puzzlesGroup.GET("/"+collection.ID, handler.NextPuzzle)
		puzzlesGroup.POST("/"+collection.ID+"/:id/guess", handler.Guess)
		puzzlesGroup.POST("/"+collection.ID+"/:id/advance", handler.Advance)
		puzzlesGroup.POST("/"+collection.ID+"/:id/give-up", handler.GiveUp)
	}

	return e
}

/* Vercel integration */
func Handler(w http.ResponseWriter, r *http.Request) {
	GetEcho().ServeHTTP(w, r)
}
