package main

import (
	"tactical-elements/handlers"
)

func main() {
	e := handlers.GetEcho()
	e.Logger.Fatal(e.Start(":1323"))
}
