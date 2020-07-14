package cmd

import (
	"github.com/gozap/dnsctl/etcdhosts"
	"time"

	"github.com/spf13/cobra"
)

var dumpFile string
var revision int64

var dumpCmd = &cobra.Command{
	Use:    "dump",
	Short:  "dump hosts",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		etcdhosts.Dump(dumpFile, revision)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().StringVarP(&dumpFile, "output", "o", "etcdhosts-"+time.Now().Format("20060102150405"), "output file")
	dumpCmd.PersistentFlags().Int64VarP(&revision, "revision", "v", -1, "hosts etcd revision")
}
