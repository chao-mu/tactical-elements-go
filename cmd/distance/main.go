package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/notnil/chess"
)

var PieceDiffPoints float64 = 5
var Ranks = map[chess.Rank]float64{
	chess.Rank1: 0,
	chess.Rank2: 1,
	chess.Rank3: 2,
	chess.Rank4: 3,
	chess.Rank5: 4,
	chess.Rank6: 5,
	chess.Rank7: 6,
	chess.Rank8: 7,
}

var Files = map[chess.File]float64{
	chess.FileA: 0,
	chess.FileB: 1,
	chess.FileC: 2,
	chess.FileD: 3,
	chess.FileE: 4,
	chess.FileF: 5,
	chess.FileG: 6,
	chess.FileH: 7,
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("Usage: <fen> <fen>")
	}

	fenA := flag.Arg(0)
	fenB := flag.Arg(1)

	bA, err := toBoard(fenA)
	if err != nil {
		log.Fatal(err)
	}

	bB, err := toBoard(fenB)
	if err != nil {
		log.Fatal(err)
	}

	score := 0.0
	for p, squaresA := range bA {
		squaresB := bB[p]

		captures := len(squaresA) - len(squaresB)
		if captures > 0 {
			score += float64(captures) * PieceDiffPoints
		}

		if len(squaresB) == 0 {
			continue
		}

		for _, sA := range squaresA {
			minDist := -1.0
			for _, sB := range squaresB {
				dist := distance(sA, sB)
				if minDist < 0 || dist < minDist {
					minDist = dist
				}
			}

			score += minDist
		}
	}

	for p, squaresB := range bB {
		squaresA := bA[p]

		additions := len(squaresB) - len(squaresA)
		if additions > 0 {
			score += float64(additions) * PieceDiffPoints
		}
	}

	fmt.Println(score)

	/*
	  https://pkg.go.dev/github.com/notnil/chess
	  f, err := os.Open("lichess_db_standard_rated_2013-01.pgn")
	  if err != nil {
	    panic(err)
	  }
	  defer f.Close()

	  scanner := chess.NewScanner(f)
	  for scanner.Scan() {
	    game := scanner.Next()
	    fmt.Println(game.GetTagPair("Site"))
	    // Output &{Site https://lichess.org/8jb5kiqw}
	  }
	*/
}

func distance(a chess.Square, b chess.Square) float64 {
	// Chebyshev Distance
	r1 := Ranks[a.Rank()]
	r2 := Ranks[b.Rank()]

	f1 := Files[a.File()]
	f2 := Files[b.File()]

	return math.Max(math.Abs(r2-r1), math.Abs(f2-f1))
}

func toBoard(fen string) (map[chess.Piece][]chess.Square, error) {
	b := make(map[chess.Piece][]chess.Square)

	f, err := chess.FEN(fen)
	if err != nil {
		return b, err
	}

	game := chess.NewGame(f)

	m := game.Position().Board().SquareMap()

	for s, p := range m {
		b[p] = append(b[p], s)
	}

	return b, nil
}
