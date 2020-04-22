package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli"
	"os/user"
)

// Publish a public key of pair x25519
func Publish() cli.Command {
	return cli.Command{
		Name:        "publish",
		Description: "Generate x25519 pair to pspk",
		Usage:       "--name <NAME> publish",
		Aliases:     []string{"p"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "qr",
				Usage: "Generate QR",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Specify path",
			},
		},
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

			if c.Bool("qr") {
				path := c.String("path")
				if path != "" {
					path = fmt.Sprintf("%s/%s.png", path, name)
				} else {
					usr, err := user.Current()
					if err != nil {
						return errors.Wrap(err, "can not get information about use")
					}
					path = fmt.Sprintf("%s/.local/share/pspk/%[2]s/%[2]s.png", usr.HomeDir, name)
				}
				err = qrcode.WriteFile(string(pub[:]), qrcode.Medium, 256, path)
				if err != nil {
					return errors.Wrap(err, "can not create file")
				}
			}

			fmt.Println("Generate key pair on x25519")
			return nil
		},
	}
}
