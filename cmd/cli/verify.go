package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

func Verify() cli.Command {
	return cli.Command{
		Name:        "verify",
		Aliases:     []string{"v"},
		Description: "Verify the message through ed25519",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file",
				Usage: "path to file which will signed",
			},
		},
		Usage: "verify <KEY_NAME> <SIGNATURE_IN_BASE64> <MESSAGE>\nor\n verify <KEY_NAME> <SIGNATURE_IN_BASE64> --file <PATH_TO_FILE>",
		Action: func(c *cli.Context) error {
			var (
				data []byte
				err  error
			)

			keyName := c.Args()[0]
			signature := c.Args()[1]

			if path := c.String("file"); path != "" {
				data, err = utils.ReadPath(path)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("could read file from path %s", path))
				}
			} else {
				message := c.Args()[2:]
				data = []byte(strings.Join(message, " "))
			}

			return pcli.Verify(keyName, signature, data)
		},
	}
}
