package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/alovn/apidoc/gen"
)

func main() {
	var searchDir string
	var outputDir string
	var mainFile string
	var isSingle bool
	var isHelp bool
	flag.StringVar(&searchDir, "dir", ".", "--dir")
	flag.StringVar(&outputDir, "output", "./docs/", "--output")
	flag.StringVar(&mainFile, "main", "main.go", "--main")
	flag.BoolVar(&isSingle, "single", false, "--single")
	flag.BoolVar(&isHelp, "help", false, "--help")
	flag.Parse()
	if isHelp {
		fmt.Println(`apidoc is a tool for Go to generate apis markdown docs.

Usage:
  apidoc --dir= --output= --main= --single

Flags:
	--dir:		search apis dir, default .
	--output: 	generate markdown files dir, default ./docs/
	--main: 	the path of main go file, default main.go
	--single: 	generate single markdown file, default multi group files`)
		return
	}
	g := gen.New(&gen.Config{
		SearchDir:      searchDir,
		MainFile:       mainFile,
		OutputDir:      outputDir,
		IsGenGroupFile: !isSingle,
	})
	if err := g.Build(); err != nil {
		log.Fatal(err)
	}
}
