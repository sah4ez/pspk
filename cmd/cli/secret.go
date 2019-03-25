package main

import (
	"encoding/base64"
	"fmt"

	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

// Secret generate shared key by
// your private key and their pulic key from pspk
func Secret() cli.Command {
	return cli.Command{
		Name:        "secret",
		Aliases:     []string{"s"},
		Description: "Generate shared secret key by private and public keys from pspk by name",
		Usage:       "secret public_name",
		Action: func(c *cli.Context) error {
			pubName := c.Args().Get(0)
			name := c.GlobalString("name")
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			priv, err := utils.Read(path, "key.bin")
			if err != nil {
				return err
			}
			pub, err := api.Load(pubName)
			if err != nil {
				return err
			}
			dh := keys.Secret(priv, pub)
			fmt.Println(base64.StdEncoding.EncodeToString(dh))

			err = utils.Write(path, pubName+".secret.bin", dh[:])
			if err != nil {
				return err
			}
			return nil
		},
	}
}
