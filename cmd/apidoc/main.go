package main

import (
	"log"

	"github.com/alovn/apidoc/gen"
)

func main() {
	g := gen.New(&gen.Config{
		SearchDir:      "../../examples/svc-user/",
		MainFile:       "main.go",
		OutputDir:      "../../examples/docs",
		IsGenGroupFile: true,
	})
	if err := g.Build(); err != nil {
		log.Fatal(err)
	}
}
