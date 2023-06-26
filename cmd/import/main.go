package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/chao-mu/tactical-elements-go/internal/data"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type config struct {
	collectionSlug string
	path           string
	db             data.DBConfig
}

type puzzleWithSolutions struct {
	FEN              string            `json:"fen"`
	MoveNumber       int               `json:"moveNumber"`
	ChessmenCount    int               `json:"chessmenCount"`
	SolutionsAliases map[string]string `json:"solutionAliases"`
	Solutions        []string          `json:"solutions"`
}

func main() {
	var cfg config

	flag.StringVar(&cfg.collectionSlug, "slug", "", "Puzzle Collection reference slug")
	flag.StringVar(&cfg.path, "path", "", "Path to puzzle collection")
	flag.IntVar(&cfg.db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.db.DSN, "db-dsn", os.Getenv("TE_DB_DSN"), "Database DSN")

	flag.Parse()

	db, err := data.OpenDB(cfg.db)
	if err != nil {
		log.Fatal(err)
	}

	puzzles, err := readPuzzlesWithSolutions(cfg.path)
	if err != nil {
		log.Fatal(err)
	}

	models := data.NewModels(db)

	puzCol := &data.PuzzleCollection{
		Slug: cfg.collectionSlug,
		Name: cfg.collectionSlug,
	}

	err = models.PuzzleCollections.Upsert(puzCol)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range puzzles {
		dbPuz := &data.Puzzle{
			FEN:          p.FEN,
			CollectionID: puzCol.ID,
		}

		err = models.Puzzles.Upsert(dbPuz)
		if err != nil {
			log.Fatal(err)
		}

		models.PuzzleSolutions.DeleteAll(dbPuz.ID)

		for _, s := range p.Solutions {
			pretty, ok := p.SolutionsAliases[s]
			if !ok {
				pretty = s
			}

			models.PuzzleSolutions.Insert(&data.PuzzleSolution{
				PuzzleID:       dbPuz.ID,
				Solution:       s,
				SolutionPretty: pretty,
			})
		}
	}
}

func readPuzzlesWithSolutions(path string) ([]puzzleWithSolutions, error) {
	puzzles := make([]puzzleWithSolutions, 0)
	jsonFile, err := os.Open(path)
	if err != nil {
		return puzzles, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &puzzles)
	if err != nil {
		return puzzles, err
	}

	return puzzles, nil
}
