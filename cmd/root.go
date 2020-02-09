package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dnsctl",
	Short: "dnsctl for etcdhosts plugin",
	Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dnsctl.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".dnsctl")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
