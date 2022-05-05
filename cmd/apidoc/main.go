package main

import (
	"log"

	"github.com/alovn/apidoc/gen"
)

func main() {
	g := gen.New()
	if err := g.Build(&gen.Config{
		SearchDir: "./examples",
		MainFile:  "main.go",
		OutputDir: "./examples/docs",
	}); err != nil {
		log.Fatal(err)
	}
}
