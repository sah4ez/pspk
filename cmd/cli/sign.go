package main

import (
	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

func Sign() cli.Command {
	return cli.Command{
		Name:        "sign",
		Aliases:     []string{"s"},
		Description: "Signs the message through ed25519",
		Usage:       "--name <KEY_NAME> sign <MESSAGE>",
		Action: func(c *cli.Context) error {
			message := c.Args()[0:]
			name := c.GlobalString("name")
			if name == "" {
				if cfg.CurrentName == "" {
					return errors.New("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}

			path = path + "/" + name

			priv, err := utils.Read(path, "key.bin")

			if err != nil {
				return errors.Wrap(err, "can not read key.bin")
			}

			return nil
		},
	}
}
