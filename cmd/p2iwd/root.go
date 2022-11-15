package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/030/logging/pkg/logging"
	"github.com/030/p2iwd/internal/app/p2iwd/pull"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logFile                                                                 bool
	cfgFile, logLevel, host, pass, protocol, repository, tag, user, Version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p2iwd",
	Short: "Pull and Push Images Without Docker (P2IWD).",
	Long: `Pull and Push Images Without Docker (P2IWD) allows a user to pull an
image without docker and push it to a registry. By default it will
download all images that reside in a Docker registry.

Examples:
  # Pull an individual image:
  p2iwd -r utrecht/dip -t 2.0.0
`,

	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.Logging{File: "p2iwd.log", Level: logLevel, Syslog: logFile}
		if _, err := l.Setup(); err != nil {
			log.Fatal(err)
		}

		dr := pull.DockerRegistry{Dir: "", Host: host, Protocol: protocol, Pass: pass, Repo: repository, Tag: tag, User: user}

		if tag == "" {
			log.Info("hello")
			if err := dr.AllTags(); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := dr.Image(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "describe repository")
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "describe tag")
	rootCmd.PersistentFlags().StringVar(&host, "host", "docker.io", "describe host")
	rootCmd.PersistentFlags().StringVarP(&protocol, "protocol", "p", "https", "describe protocol")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file (default: "${HOME}/.p2iwd/config.yml")`)
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "change the log level (default: info, options: trace, debug, info, warn, error or none)")
	rootCmd.PersistentFlags().BoolVar(&logFile, "logFile", false, "write the log to syslog")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		home = filepath.Join(home, ".p2iwd")

		// Search config in home directory with name ".p2iwd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
