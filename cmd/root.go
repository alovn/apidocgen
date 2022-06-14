/*
Copyright © 2022 alovn

*/
package cmd

import (
	"log"
	"os"

	"github.com/alovn/apidocgen/gen"
	"github.com/spf13/cobra"
)

var (
	searchDir        string
	outputDir        string
	outputIndexName  string
	templateDir      string
	excludesDir      string
	isSingle         bool
	isGenMocks       bool
	isMockServer     bool
	mockServerListen string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "apidocgen",
	Long: `apidocgen is a tool for Go to generate apis markdown docs and mocks. For example:

apidocgen --dir= --excludes= --output= --output-index= --template= --single --gen-mocks
apidocgen --mock --mock-listen=:8001
apidocgen mock --data= --listen=`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		g := gen.New(&gen.Config{
			SearchDir:        searchDir,
			OutputDir:        outputDir,
			OutputIndexName:  outputIndexName,
			TemplateDir:      templateDir,
			ExcludesDir:      excludesDir,
			IsGenSingleFile:  isSingle,
			IsGenMocks:       isGenMocks,
			IsMockServer:     isMockServer,
			MockServerListen: mockServerListen,
		})
		if err := g.Build(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&searchDir, "dir", ".", "Search apis dir, comma separated")
	rootCmd.Flags().StringVar(&outputDir, "output", "./docs/", "Generate markdown files dir")
	rootCmd.Flags().StringVar(&outputIndexName, "output-index", "", "Generate index file name")
	rootCmd.Flags().StringVar(&templateDir, "template", "", "Template name or custom template directory, built-in includes markdown and apidocs")
	rootCmd.Flags().StringVar(&excludesDir, "excludes", "", "Exclude directories and files when searching, comma separated")
	rootCmd.Flags().BoolVar(&isSingle, "single", false, "Is generate a single doc")
	rootCmd.Flags().BoolVar(&isGenMocks, "gen-mocks", false, "Is generate the mock files")
	rootCmd.Flags().BoolVar(&isMockServer, "mock", false, "Serve a mock server")
	rootCmd.Flags().StringVar(&mockServerListen, "mock-listen", "localhost:8001", "Mock Server listen address")
}
