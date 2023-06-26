package data

import (
	"database/sql"
	"github.com/pkg/errors"
)

type PuzzleCollection struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PuzzleCollectionModel struct {
	db *sql.DB
}

func (m PuzzleCollectionModel) Upsert(col *PuzzleCollection) error {
	row := m.db.QueryRow("select id from puzzle_collections where slug = ?", col.Slug)

	err := row.Scan(&col.ID)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = m.db.Exec("insert into puzzle_collections (slug) values (?)", col.Slug)

		return err
	}

	return err
}

type Puzzle struct {
	ID           int    `json:"id"`
	FEN          string `json:"fen"`
	CollectionID int    `json:"collectionId"`
}

type PuzzleModel struct {
	db *sql.DB
}

func (m PuzzleModel) Upsert(p *Puzzle) error {
	row := m.db.QueryRow("select id from puzzles where fen = ? and collection_id = ?", p.FEN, p.CollectionID)

	values := []any{
		p.FEN,
		p.CollectionID,
	}

	err := row.Scan(&p.ID)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = m.db.Exec("insert into puzzles (fen, collection_id) values (?, ?)", values...)

		return err
	}

	values = append(values, p.ID)
	_, err = m.db.Exec("update puzzles set fen=?, collection_id=? where id = ?", values...)

	return err
}

type PuzzleSolution struct {
	ID             int    `json:"id"`
	PuzzleID       int    `json:"puzzleId"`
	Solution       string `json:"solution"`
	SolutionPretty string `json:"solutionPretty"`
}

type PuzzleSolutionModel struct {
	db *sql.DB
}

func (m PuzzleSolutionModel) DeleteAll(pid int) error {
	_, err := m.db.Exec("delete from puzzle_solutions where puzzle_id = ?", pid)

	return err
}

func (m PuzzleSolutionModel) Insert(p *PuzzleSolution) error {
	values := []any{
		p.PuzzleID,
		p.Solution,
		p.SolutionPretty,
	}

	_, err := m.db.Exec("insert into puzzle_solutions (puzzle_id, solution, solution_pretty) values (?, ?, ?)", values...)

	return err
}

type PuzzleGuesses struct {
	ID        int    `json:"id"`
	Guess     string `json:"guess"`
	Correct   bool   `json:"correct"`
	SessionID string `json:"sessionId"`
}

type PuzzleSolutionSet []PuzzleSolution
