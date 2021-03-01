package main

import (
	"github.com/aagu/go-i18n/cmd/extract"
	"github.com/spf13/cobra"
	"os"
)

var command = cobra.Command{
	Use:   "go-i18n",
	Short: "Command line tool for go-i18n",
}

func init() {
	command.AddCommand(extract.NewExtractCommand())
}

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
