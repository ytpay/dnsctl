package cmd

import (
	"github.com/gozap/dnsctl/etcdhosts"

	"github.com/spf13/cobra"
)

var dumpFile string

var dumpCmd = &cobra.Command{
	Use:    "dump",
	Short:  "dump hosts",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		etcdhosts.Dump(dumpFile)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().StringVarP(&dumpFile, "output", "o", "", "output file")
}
