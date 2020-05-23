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

			pub, err := api.Load(keyName)
			if err != nil {
				return errors.Wrap(err, "can not load public name")
			}

			pubArray := utils.Slice2Array32(pub)

			signatureBinary, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				return errors.Wrap(err, "can decode signature from base64")
			}
			signatureArrya := utils.Slice2Array64(signatureBinary)

			verify := keys.Verify(pubArray, data, &signatureArrya)
			fmt.Fprintln(out, verifyMessage(signature, verify))
			return nil
		},
	}
}

func verifyMessage(signature string, verify bool) string {
	if verify {
		return fmt.Sprintf("Signature %s is valid.", signature)
	}
	return fmt.Sprintf("Signature %s is NOT valid.", signature)
}
