package cmd

import (
	"errors"

	"chainguard.dev/fagin/pkg/unpack"
	"github.com/spf13/cobra"
)

type listOptions struct {
	set string
}

func addList(parentCmd *cobra.Command) {
	opts := listOptions{}

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
			if err := u.List(args[0], opts.set); err != nil {
				return err
			}
			return nil
		},
	}
	listCmd.PersistentFlags().StringVar(
		&opts.set,
		"set",
		"untracked",
		"set of files to output (all | tracked | untracked)",
	)
	parentCmd.AddCommand(listCmd)
}
