package cmd

import (
	"github.com/gozap/dnsctl/etcdhosts"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:    "edit",
	Short:  "edit hosts",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		etcdhosts.Edit()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
