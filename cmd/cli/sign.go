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

			priv, err := utils.Read(path, "key.bin")
			if err != nil {
				return errors.Wrap(err, "can not read key.bin")
			}

			privArray := utils.Slice2Array32(priv)

			sign := keys.Sign(&privArray, []byte(strings.Join(message, " ")), utils.Random())
			data := base64.StdEncoding.EncodeToString(sign[:])

			fmt.Fprintln(out, data)
			return nil
		},
	}
}
