package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	syslog                                                       bool
	cfgFile, dir, logLevel, host, pass, repo, tag, user, Version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p2iwd",
	Short: "Pull and Push Images Without Docker (P2IWD).",
	Long: `Pull and Push Images Without Docker (P2IWD) allows a user to pull an
image without docker and push it to a registry. By default it will
download all images that reside in a Docker registry.

Examples:
  # Pull images:
  p2iwd pull

  # Pull an individual image
  p2iwd pull -r a -t b

  # Push images:
  p2iwd push
`,

	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			dir = viper.GetString("dir")
		}

		if host == "" {
			host = viper.GetString("host")
		}

		if logLevel == "info" {
			value := viper.GetString("logLevel")
			if value != "" {
				logLevel = value
			}
		}

		if pass == "" {
			pass = viper.GetString("pass")
		}

		if repo == "" {
			repo = viper.GetString("repo")
		}

		if !syslog {
			syslog = viper.GetBool("syslog")
		}

		if tag == "" {
			tag = viper.GetString("tag")
		}

		if user == "" {
			user = viper.GetString("user")
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file (default: "${HOME}/.p2iwd/config.yml")`)
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "describe dir")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "describe host")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "change the log level (default: info, options: trace, debug, info, warn, error or none)")
	rootCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "describe pass")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "describe repo")
	rootCmd.PersistentFlags().BoolVar(&syslog, "syslog", false, "write the log to syslog")
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "describe tag")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "describe user")
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

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
