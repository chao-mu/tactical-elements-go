package main

import (
	"flag"
)

type config struct {
	sourcesPath string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.sourcesPath, "sources", "./assets/sources.pgn", "Path to PGN of games")

	flag.Parse()

	/*
	  https://pkg.go.dev/github.com/notnil/chess
	  f, err := os.Open("lichess_db_standard_rated_2013-01.pgn")
	  if err != nil {
	  }
	    panic(err)
	  defer f.Close()

	  scanner := chess.NewScanner(f)
	  for scanner.Scan() {
	    game := scanner.Next()
	    fmt.Println(game.GetTagPair("Site"))
	    // Output &{Site https://lichess.org/8jb5kiqw}
	  }
	*/
}
