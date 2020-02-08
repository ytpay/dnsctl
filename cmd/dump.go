package cmd

import (
	"github.com/gozap/gdnsctl/gdns"

	"github.com/spf13/cobra"
)

var dumpFile string

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump hosts",
	Run: func(cmd *cobra.Command, args []string) {
		gdns.Dump(dumpFile)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().StringVarP(&dumpFile, "output", "o", "", "output file")
}
