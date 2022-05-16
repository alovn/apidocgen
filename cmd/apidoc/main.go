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
	var isSingle bool
	var isHelp bool
	flag.StringVar(&searchDir, "dir", ".", "--dir")
	flag.StringVar(&outputDir, "output", "./docs/", "--output")
	flag.BoolVar(&isSingle, "single", false, "--single")
	flag.BoolVar(&isHelp, "help", false, "--help")
	flag.Parse()
	if isHelp {
		fmt.Println(`apidoc is a tool for Go to generate apis markdown docs.

Usage:
  apidoc --dir= --output= --single

Flags:
	--dir:		search apis dir, default .
	--output: 	generate markdown files dir, default ./docs/
	--single: 	generate single markdown file, default multi group files`)
		return
	}
	g := gen.New(&gen.Config{
		SearchDir:       searchDir,
		OutputDir:       outputDir,
		IsGenSingleFile: isSingle,
	})
	if err := g.Build(); err != nil {
		log.Fatal(err)
	}
}
