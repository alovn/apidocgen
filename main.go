package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/alovn/apidocgen/gen"
)

func main() {
	var searchDir string
	var outputDir string
	var templateDir string
	var excludesDir string
	var isSingle bool
	var isHelp bool
	flag.StringVar(&searchDir, "dir", ".", "--dir")
	flag.StringVar(&outputDir, "output", "./docs/", "--output")
	flag.StringVar(&templateDir, "template", "", "--template")
	flag.StringVar(&excludesDir, "excludes", "", "--excludes")
	flag.BoolVar(&isSingle, "single", false, "--single")
	flag.BoolVar(&isHelp, "help", false, "--help")
	flag.Parse()
	if isHelp {
		fmt.Println(`apidocgen is a tool for Go to generate apis markdown docs.

Usage:
  apidocgen --dir= --excludes= --output= --template= --single

Flags:
	--dir:		Search apis dir, comma separated, default .
	--excludes:	Exclude directories and files when searching, comma separated
	--output: 	Generate markdown files dir, default ./docs/
	--template:	Custom template files dir.
	--single: 	If true, generate a single markdown file, default false`)
		return
	}
	g := gen.New(&gen.Config{
		SearchDir:       searchDir,
		OutputDir:       outputDir,
		TemplateDir:     templateDir,
		ExcludesDir:     excludesDir,
		IsGenSingleFile: isSingle,
	})
	if err := g.Build(); err != nil {
		log.Fatal(err)
	}
}
