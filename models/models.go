package models

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
)

type Player struct {
	ID        int    `json:"id"`
	Health    int    `json:"health"`
	Exp       int    `json:"exp"`
	Gold      int    `json:"gold"`
	BrowserID string `json:"browserId"`
}

type PuzzleCollection struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Puzzles []Puzzle `json:"puzzles"`
}

func (c PuzzleCollection) RandomPuzzle() Puzzle {
	id := rand.Intn(len(c.Puzzles))
	return c.Puzzles[id]
}

type PuzzlePublic struct {
	ID  int    `json:"id"`
	FEN string `json:"fen"`
}

type Puzzle struct {
	ID              int               `json:"id"`
	FEN             string            `json:"fen"`
	Solutions       []string          `json:"solution"`
	SolutionAliases map[string]string `json:"solutionAliases"`
}

func (p Puzzle) Public() PuzzlePublic {
	return PuzzlePublic{
		ID:  p.ID,
		FEN: p.FEN,
	}
}

func (p Puzzle) GuessOne(candidate string) bool {
	for _, solution := range p.Solutions {
		if solution == candidate {
			return true
		}
	}

	for solution, alias := range p.SolutionAliases {
		if solution == candidate || alias == candidate {
			return true
		}
	}

	return false
}

func (p Puzzle) GuessAll(solutions []string) bool {
	if len(solutions) != len(p.Solutions) {
		return false
	}

	for _, solution := range solutions {
		if !p.GuessOne(solution) {
			return false
		}
	}

	return true
}

func ReadPuzzles(jsonPath string) ([]Puzzle, error) {
	puzzles := make([]Puzzle, 0)
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return puzzles, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &puzzles)
	if err != nil {
		return puzzles, err
	}

	for idx, _ := range puzzles {
		puzzles[idx].ID = idx
	}

	return puzzles, nil
}
