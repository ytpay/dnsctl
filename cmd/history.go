package cmd

import (
	"github.com/gozap/dnsctl/etcdhosts"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:    "history",
	Short:  "show hosts history",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		etcdhosts.History()
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
