package cmd

import (
	"errors"

	"chainguard.dev/fagin/pkg/unpack"
	"github.com/spf13/cobra"
)

type statsOptions struct {
	format string
}

func addStats(parentCmd *cobra.Command) {
	opts := statsOptions{}
	statsCmd := &cobra.Command{
		Short: "Output some stats about files in an image",
		Long: `stats imageref
	
The stats subcommand checks files in an image and returns some stats
about files tracked (and not tracked) in a container image.
	
	`,
		Use:               "stats",
		SilenceUsage:      false,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				return errors.New("no image reference specified")
			}
			u := unpack.New()
			u.Options.Distro = commandLineOpts.distro
			if err := u.PrintStats(args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	statsCmd.PersistentFlags().StringVar(
		&opts.format,
		"format",
		"text",
		"format for the stats (text or json)",
	)
	parentCmd.AddCommand(statsCmd)
}
