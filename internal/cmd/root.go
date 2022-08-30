package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
)

type commandLineOptions struct {
	logLevel string
}

var commandLineOpts = &commandLineOptions{}

func Execute() error {
	rootCmd := &cobra.Command{
		// Use:               "tejolote",
		SilenceUsage:      false,
		PersistentPreRunE: initLogging,
	}

	rootCmd.PersistentFlags().StringVar(
		&commandLineOpts.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)

	addStats(rootCmd)
	addList(rootCmd)
	return rootCmd.Execute()
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}
