package cmd

import (
	"github.com/spf13/cobra"
)

type rootCmd struct {
	*baseCmd
}

func newRootCmd() *rootCmd {
	rc := &rootCmd{newBaseCmd()}
	rc.cmd = &cobra.Command{
		Use:   "",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				panic(err)
			}
		},
	}
	rc.cmd.SetHelpCommand(&cobra.Command{})
	rc.addCommands(
		newRunCmd(),
		newJobCmd(),
		newCronCmd(),
		newGenModelCmd(),
	)

	return rc
}
