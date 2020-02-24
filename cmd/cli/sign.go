package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

func Sign() cli.Command {
	return cli.Command{
		Name:        "sign",
		Aliases:     []string{"s"},
		Description: "Signs the message through ed25519",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file",
				Usage: "path to file which will signed",
			},
		},
		Usage: "--name <KEY_NAME> sign <MESSAGE>\nor\n --name <KEY_NAME> sign --file <PATH_TO_FILE>",
		Action: func(c *cli.Context) error {
			var (
				data []byte
				err  error
			)

			if path := c.String("file"); path != "" {
				data, err = utils.ReadPath(path)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("could read file from path %s", path))
				}
			} else {
				message := c.Args()[0:]
				data = []byte(strings.Join(message, " "))
			}

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

			sign := keys.Sign(&privArray, data, utils.Random())
			sginature := base64.StdEncoding.EncodeToString(sign[:])

			fmt.Fprintln(out, sginature)
			return nil
		},
	}
}
