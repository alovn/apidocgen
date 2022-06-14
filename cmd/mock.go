/*
Copyright © 2022 alovn

*/
package cmd

import (
	"log"

	"github.com/alovn/apidocgen/mock"
	"github.com/spf13/cobra"
)

var (
	mockDataDir string
	mockListen  string
)

// mockCmd represents the mock command
var mockCmd = &cobra.Command{
	Use: "mock",
	Long: `Serve a mock serve with generated mock files. For example:

apidocgen mock --data=./docs/mocks --listen=localhost:8001`,
	Run: func(cmd *cobra.Command, args []string) {
		mockServer := mock.NewMockServer(mockListen)
		if err := mockServer.InitFiles(mockDataDir); err != nil {
			log.Fatal(err)
		}
		if err := mockServer.Serve(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mockCmd)
	mockCmd.Flags().StringVar(&mockDataDir, "data", "./docs/mocks/", "mocks data directory.")
	mockCmd.Flags().StringVar(&mockListen, "listen", "localhost:8001", "mock server listen address.")
}
