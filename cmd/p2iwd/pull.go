package main

import (
	"github.com/030/logging/pkg/logging"
	"github.com/030/p2iwd/internal/app/p2iwd/pull"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull an image",
	Long:  `Pull an image.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.Logging{File: "p2iwd-pull.log", Level: logLevel, Syslog: syslog}
		if _, err := l.Setup(); err != nil {
			log.Fatal(err)
		}

		pdr := pull.DockerRegistry{Dir: dir, Host: host, Pass: pass, User: user}
		if repo != "" {
			if err := pdr.Image(repo, tag); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := pdr.All(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
