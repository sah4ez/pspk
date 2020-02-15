package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

// Publish a public key of pair x25519
func Publish() cli.Command {
	return cli.Command{
		Name:        "publish",
		Description: "Generate x25519 pair to pspk",
		Usage:       "--name <NAME> publish",
		Aliases:     []string{"p"},
		Action: func(c *cli.Context) error {
			name := c.GlobalString("name")
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			pub, priv, err := keys.GenerateDH()
			if err != nil {
				return errors.Wrap(err, "can not generate keys")
			}
			err = utils.Write(path, "pub.bin", pub[:])
			if err != nil {
				return errors.Wrap(err, "can not write in pub.bin")
			}
			err = api.Publish(name, pub[:])
			if err != nil {
				return errors.Wrap(err, "can not publish")
			}

			err = utils.Write(path, "key.bin", priv[:])
			if err != nil {
				return errors.Wrap(err, "can not find the path")
			}

			fmt.Println("Generate key pair on x25519")
			return nil
		},
	}
}
