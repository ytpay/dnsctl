package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "show example config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`dnskey: /gdns
etcd:
  cert: /etc/etcd/ssl/etcd.pem
  key: /etc/etcd/ssl/etcd-key.pem
  ca: /etc/etcd/ssl/etcd-root-ca.pem
  endpoints:
    - https://172.16.10.11:2379
    - https://172.16.10.12:2379
    - https://172.16.10.13:2379
`)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
