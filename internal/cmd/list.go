package cmd

import (
	"errors"

	"chainguard.dev/fagin/pkg/unpack"
	"github.com/spf13/cobra"
)

type listOptions struct {
	set string
}

var listOpts = listOptions{}

func addList(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Short: "Lists sets of files in a container image",
		Long: `list [options] imageref
	
The list subcommand outputs to stdout lists of files in an image.
The list of files can be one of three:

		all       → All files in the container image
		tracked   → Files tracked by the OS package manager
		untracked → Files found in the image not tracked by the package manager

By default list will output all non-tracked files.
	
	`,
		Use:               "list",
		SilenceUsage:      false,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				return errors.New("no image reference specified")
			}
			u := unpack.New()
			u.Options.Distro = commandLineOpts.distro
			if err := u.List(args[0], listOpts.set); err != nil {
				return err
			}
			return nil
		},
	}
	listCmd.PersistentFlags().StringVar(
		&listOpts.set,
		"set",
		"untracked",
		"set of files to output (all | tracked | untracked)",
	)
	parentCmd.AddCommand(listCmd)
}
