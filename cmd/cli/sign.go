package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
	"strings"
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

			if err != nil {
				return errors.Wrap(err, "can not read key.bin")
			}

			sign := keys.Sign(pub, []byte(strings.Join(message, " ")), utils.Random())
			data := base64.StdEncoding.EncodeToString(sign[:])

			fmt.Fprintln(out, data)
			return nil
		},
	}
}
