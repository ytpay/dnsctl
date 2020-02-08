package cmd

import (
	"github.com/gozap/gdnsctl/gdns"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload hosts from file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
		} else {
			gdns.Upload(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
