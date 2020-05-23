package main

import (
	"github.com/urfave/cli"
)

// Publish a public key of pair x25519
func Publish() cli.Command {
	return cli.Command{
		Name:        "publish",
		Description: "Generate x25519 pair to pspk",
		Usage:       "--name <NAME> publish",
		Aliases:     []string{"p"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "qr",
				Usage: "Generate QR",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Specify path",
			},
		},
		Action: func(c *cli.Context) error {
			name := c.GlobalString("name")
			if c.GlobalBool("qr") {
				qrPath := c.GlobalString("path")
				return pcli.PublishAndGenerateQR(name, qrPath)
			}
			return pcli.Publish(name)
		},
	}
}
