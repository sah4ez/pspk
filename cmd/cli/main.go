package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

const (
	baseURL = "https://pspk.now.sh"
)

func main() {
	var (
		err error
	)
	app := cli.NewApp()
	app.Name = "pspk"
	app.Version = "0.0.1"
	app.Description = "Console tool for encyption/decription data through pspk.now.sh"
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Usage:   `Generate x25519 pair`,
			Aliases: []string{"g"},
			Action: func(c *cli.Context) error {
				pub, priv, err := keys.GenereateDH()
				if err != nil {
					return err
				}
				err = utils.Write("./", "pub.bin", pub[:])
				if err != nil {
					return err
				}
				err = utils.Write(".", "key.bin", priv[:])
				if err != nil {
					return err
				}

				fmt.Println("Generate key pair on x25519")
				return nil
			},
		},
		{
			Name:  "secret",
			Usage: `Generate shared secret key by private and public keys`,
			Action: func(c *cli.Context) error {
				privPath := c.Args().Get(0)
				pubPath := c.Args().Get(1)
				priv, err := utils.ReadPath(privPath)
				if err != nil {
					return err
				}
				pub, err := utils.ReadPath(pubPath)
				if err != nil {
					return err
				}
				dh := keys.Secret(priv, pub)
				fmt.Println("secret:", base64.StdEncoding.EncodeToString(dh))
				if len(c.Args()) == 3 {
					err = utils.Write("./", c.Args().Get(2), dh[:])
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name:  "encrypt",
			Usage: `Encrypt input message with shared key`,
			Action: func(c *cli.Context) error {
				key := c.Args()[0]
				message := c.Args()[1:]
				chain, err := utils.Read("./", key)
				if err != nil {
					return err
				}
				messageKey, err := keys.LoadMaterialKey(chain)
				if err != nil {
					return err
				}

				b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(strings.Join(message, " ")))
				if err != nil {
					return err
				}
				fmt.Println("encrypted:", base64.StdEncoding.EncodeToString(b))
				return nil
			},
		},
		{
			Name:  "decrypt",
			Usage: `Decrypt input message with shared key`,
			Action: func(c *cli.Context) error {
				key := c.Args()[0]
				message := c.Args()[1]
				chain, err := utils.Read("./", key)
				if err != nil {
					return err
				}
				messageKey, err := keys.LoadMaterialKey(chain)
				if err != nil {
					return err
				}
				bytesMessage, err := base64.StdEncoding.DecodeString(message)
				if err != nil {
					return err
				}

				b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage)
				if err != nil {
					return err
				}
				fmt.Println("decoded:", string(b))
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run has error:", err.Error())
	}

}
