package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flag struct {
	name         string
	shortName    string
	defaultValue interface{}
	usage        string
}

func getFlags(cmd *cobra.Command, persistent bool) *pflag.FlagSet {
	flags := cmd.Flags()
	if persistent {
		flags = cmd.PersistentFlags()
	}
	return flags
}
