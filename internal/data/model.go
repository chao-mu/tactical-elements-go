package data

import "database/sql"

type Models struct {
	PuzzleCollections interface {
		Upsert(col *PuzzleCollection) error
	}
	PuzzleSolutions interface {
		Insert(p *PuzzleSolution) error
		DeleteAll(pid int) error
	}
	Puzzles interface {
		Upsert(p *Puzzle) error
	}
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		PuzzleCollections: PuzzleCollectionModel{db: db},
		PuzzleSolutions:   PuzzleSolutionModel{db: db},
		Puzzles:           PuzzleModel{db: db},
	}
}
