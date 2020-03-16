package cmd

import (
	"github.com/gozap/dnsctl/etcdhosts"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "show hosts version",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		etcdhosts.Version()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
