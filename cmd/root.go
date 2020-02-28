package cmd

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var bannerBase64 = "ICAgICAgICAgICAgICAgICAgICAgICAgICAgXyAgCiAgIHwgICAgICAgICAgICAgICAgICAgICAgfCB8IAogX198ICAgXyAgXyAgICAsICAgX18gX3xfIHwgfCAKLyAgfCAgLyB8LyB8ICAvIFxfLyAgICB8ICB8LyAgClxfL3xfLyAgfCAgfF8vIFwvIFxfX18vfF8vfF9fLwo="

var versionTpl = `%s
Name: dnsctl
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
`

var (
	Version   string
	BuildDate string
	CommitID  string
)

var rootCmd = &cobra.Command{
	Use:     "dnsctl",
	Version: Version,
	Short:   "dnsctl for etcdhosts plugin",
	Run:     func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dnsctl.yaml)")
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	rootCmd.SetVersionTemplate(fmt.Sprintf(versionTpl, banner, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID))
}

func initConfig(_ *cobra.Command, _ []string) {
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
