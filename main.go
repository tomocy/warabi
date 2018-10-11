package main

import (
	"os"

	"github.com/tomocy/warabi/repl"
)

func main() {
	repler := repl.NewStandard(os.Stdin, os.Stdout)
	// repler := repl.NewWarabi(os.Stdin, os.Stdout)
	repler.REPL()
}
