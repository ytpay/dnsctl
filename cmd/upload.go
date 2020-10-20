package cmd

import (
	"github.com/ytpay/dnsctl/etcdhosts"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:    "upload FILE",
	Short:  "upload hosts from file",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
		} else {
			etcdhosts.Upload(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
