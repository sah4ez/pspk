package main

import (
	"github.com/urfave/cli"
)

// Secret generate shared key by
// your private key and their pulic key from pspk
func Secret() cli.Command {
	return cli.Command{
		Name:        "secret",
		Aliases:     []string{"s"},
		Description: "Generate shared secret key by private and public keys from pspk by name",
		Usage:       "secret public_name",
		Action: func(c *cli.Context) error {
			pubName := c.Args().Get(0)
			name := c.GlobalString("name")
			return pcli.Secret(name, pubName)
		},
	}
}
