/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
)

type commandLineOptions struct {
	logLevel string
	distro   string
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

	rootCmd.PersistentFlags().StringVar(
		&commandLineOpts.distro,
		"distro",
		"",
		fmt.Sprintf("distro format of the scanned image (alpine | debian)"),
	)

	addStats(rootCmd)
	addList(rootCmd)
	return rootCmd.Execute()
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}
