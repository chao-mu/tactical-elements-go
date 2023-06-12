package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

type PuzzleCollection struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	PuzzlesWithSolutions []PuzzleWithSolutions `json:"puzzles"`
}

func GetPuzzleCollection(unsafeID string) (*PuzzleCollection, error) {
	allowedIDs := []string{
		"checkables",
		"checks-captures",
		"knight-forkable",
		"undefended",
	}
	for _, allowedID := range allowedIDs {
		if allowedID == unsafeID {
			path := fmt.Sprintf("assets/puzzles/%s.json", allowedID)
			puzzles, err := ReadPuzzlesWithSolutions(path)
			if err != nil {
				return nil, err
			}

			return &PuzzleCollection{
				ID:      allowedID,
				PuzzlesWithSolutions: puzzles,
			}, nil
		}
	}

  return nil, fmt.Errorf("Unable to find puzzle collection of id %s", unsafeID)
}

func (c PuzzleCollection) RandomPuzzle() PuzzleWithSolutions {
	id := rand.Intn(len(c.PuzzlesWithSolutions))
	return c.PuzzlesWithSolutions[id]
}

type Puzzle struct {
	ID  int    `json:"id"`
	FEN string `json:"fen"`
}

type PuzzleWithSolutions struct {
  Puzzle
	Solutions       []string          `json:"solution"`
	SolutionAliases map[string]string `json:"solutionAliases"`
}

func GetPuzzleWithSolutions(colID string, puzzleID int) (PuzzleWithSolutions, error) {
  col, err := GetPuzzleCollection(colID)
  if err != nil {
    return PuzzleWithSolutions{}, err
  }

  if puzzleID >= len(col.PuzzlesWithSolutions) || puzzleID < 0 {
    return PuzzleWithSolutions{}, fmt.Errorf("Invalid puzzle id %d", puzzleID)
  }

  return col.PuzzlesWithSolutions[puzzleID], nil
}

func (p *PuzzleWithSolutions) GuessOne(candidate string) bool {
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

func (p *PuzzleWithSolutions) GuessAll(solutions []string) bool {
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

func ReadPuzzlesWithSolutions(jsonPath string) ([]PuzzleWithSolutions, error) {
	puzzles := make([]PuzzleWithSolutions, 0)
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

	for idx := range puzzles {
		puzzles[idx].ID = idx
	}

	return puzzles, nil
}
