package main

import (
	"github.com/030/logging/pkg/logging"
	"github.com/030/p2iwd/internal/app/p2iwd/push"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an image",
	Long: `Push an image.

Examples:
  # Push images:
  p2iwd push
`,
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.Logging{File: "p2iwd-push.log", Level: logLevel, Syslog: syslog}
		if _, err := l.Setup(); err != nil {
			log.Fatal(err)
		}

		pdr := push.DockerRegistry{Dir: dir, Host: host, Pass: pass, User: user}
		if err := pdr.All(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
