package main

import (
	"os"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/beldpro-ci/subscriber/server"
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
		commonutils.AssertIntFlagsSet("_app", c, mainLog,
			"port")
		commonutils.AssertStringFlagsSet("_app", c, mainLog,
			"api-key",
			"url",
			"list-id")
		return nil
	}
	app.Action = func(c *cli.Context) error {
		mc, err := mailchimp.New(mailchimp.Config{
			APIKey: c.String("api-key"),
			ListId: c.String("list-id"),
			URL:    c.String("url"),
		})
		if err != nil {
			return err
		}

		httpserver, err := server.New(server.Config{
			Port:      c.Int("port"),
			MailChimp: &mc,
		})
		if err != nil {
			return err
		}

		return httpserver.Run()
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port",
			EnvVar: "SUBSCRIBER_PORT",
			Usage:  "Port to listen to HTTP requests",
			Value:  8080,
		},

		cli.StringFlag{
			Name:   "api-key",
			EnvVar: "SUBSCRIBER_MAILCHIMP_API_KEY",
			Usage:  "MailChimp API Key",
		},

		cli.StringFlag{
			Name:   "url",
			EnvVar: "SUBSCRIBER_MAILCHIMP_URL",
			Usage:  "MailChimp's datacenter URL",
		},

		cli.StringFlag{
			Name:   "list-id",
			EnvVar: "SUBSCRIBER_MAILCHIMP_LIST_ID",
			Usage:  "Id of the list which users should subscribe to",
		},
	}

	app.Run(os.Args)
}
