package main

import (
	"flag"
	"log"
	"os"

	"github.com/tenkoh/tdd-moshi-share/ranking"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("invalid command.")
	}
	filepath := args[0]
	w := os.Stdout
	err := ranking.GetRank(filepath, w)
	if err != nil {
		log.Fatal("operation failed")
	}
}
