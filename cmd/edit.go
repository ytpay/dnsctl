package cmd

import (
	"github.com/gozap/gdnsctl/gdns"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit hosts",
	Run: func(cmd *cobra.Command, args []string) {
		gdns.Edit()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
