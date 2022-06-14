/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version version
	Version = "0.1.0"
	// BuildDate build date
	BuildDate string
	// GitCommit
	GitCommit string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Long:  `print version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", getVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion() string {
	ret := Version
	if GitCommit != "" && len(GitCommit) >= 7 {
		ret += fmt.Sprintf("\ngit-commit: %s", GitCommit[0:7])
	}
	if BuildDate != "" {
		ret += fmt.Sprintf("\nbuild-at: %s", BuildDate)
	}
	return ret
}
