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

func Verify() cli.Command {
	return cli.Command{
		Name:        "verify",
		Aliases:     []string{"v"},
		Description: "Verify the message through ed25519",
		Usage:       "verify <KEY_NAME> <SIGNATURE_IN_BASE64> <MESSAGE>",
		Action: func(c *cli.Context) error {
			keyName := c.Args()[0]
			signature := c.Args()[1]
			message := c.Args()[2:]

			pub, err := api.Load(keyName)
			if err != nil {
				return errors.Wrap(err, "can not load public name")
			}

			pubArray := utils.Slice2Array32(pub)


			data, err := base64.StdEncoding.DecodeString(signature)
			sig := utils.Slice2Array64(data)

			verify := keys.Verify(pubArray, []byte(strings.Join(message, " ")), &sig)
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