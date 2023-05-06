package main

import (
	"github.com/030/logging/pkg/logging"
	"github.com/030/p2iwd/internal/app/p2iwd/pull"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
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
		if errs := validator.Validate(pdr); errs != nil {
			log.Fatal(errs)
		}
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
}
