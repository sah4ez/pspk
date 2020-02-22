package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
	"strings"
)

func Verify() cli.Command {
	return cli.Command{
		Name:        "verify",
		Aliases:     []string{"v"},
		Description: "Verify the message through ed25519",
		Usage:       "verify <KEY_NAME> <MESSAGE>",
		Action: func(c *cli.Context) error {
			keyName := c.Args()[0]
			message := c.Args()[1:]


			path = path + "/" + keyName

			priv, err := utils.Read(path, "key.bin")

			pub, err := api.Load(keyName)
			if err != nil {
				return errors.Wrap(err, "can not load public name")
			}

			pubArray := utils.Slice2Array32(pub)
			privArray := utils.Slice2Array64(priv)

			verify := keys.Verify(pubArray, []byte(strings.Join(message, " ")), &privArray)
			fmt.Fprintln(out, verify)
			return nil
		},
	}
}
