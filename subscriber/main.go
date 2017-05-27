package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"

	log "github.com/Sirupsen/logrus"
	commonutils "github.com/beldpro-ci/common/utils"
)

var mainLog = log.WithField("from", "fserve/main")

func main() {
	app := cli.NewApp()
	app.Name = "subscriber"
	app.Usage = "Subscribes user to mailchimp"
	app.Before = func(c *cli.Context) error {
		commonutils.AssertStringFlagsSet("_app", c, mainLog,
			"api-key",
			"list-id")
		return nil
	}

	app.Action = func (c *cli.Context) error {
i		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			EnvVar: "SUBSCRIBER_MAILCHIMP_API_KEY",
			Usage:  "MailChimp API Key",
		},

		cli.StringFlag{
			Name:   "list-id",
			EnvVar: "SUBSCRIBER_MAILCHIMP_LIST_ID",
			Usage:  "Id of the list which users should subscribe to",
		},
	}

	app.Run(os.Args)
}
