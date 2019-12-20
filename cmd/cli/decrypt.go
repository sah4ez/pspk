package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func Decrypt() cli.Command {
	return cli.Command{
		Name:        "decrypt",
		Aliases:     []string{"d"},
		Description: "Decrypt input message with shared key",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "<URL> - for loading decode data by link instead of text input data",
			},
		},
		Usage: "decrypt pub_name base64==",
		Action: func(c *cli.Context) error {
			pubName := c.Args().Get(0)
			var message string
			if link := c.String("link"); len(link) > 0 {
				m, err := api.DownloadByLink(link)
				if err != nil {
					return errors.Wrap(err, "download by link failed")
				}
				message = m
			} else {
				message = c.Args().Get(1)
			}
			name := c.GlobalString("name")

			return pcli.Decrypt(name, message, pubName)
		},
	}
}

func EphemeralDecrypt() cli.Command {
	return cli.Command{
		Name:        "ephemeral-decrypt",
		Aliases:     []string{"ed"},
		Description: `Decrypt input message with ephemral shared key`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "<URL> - for loading decode data by link instead of text input data",
			},
		},
		Usage: "ephemeral-decryp pub_name base64==",
		Action: func(c *cli.Context) error {
			var message string
			if link := c.String("link"); len(link) > 0 {
				m, err := api.DownloadByLink(link)
				if err != nil {
					return errors.Wrap(err, "download by link failed")
				}
				message = m
			} else {
				message = c.Args().Get(0)
			}
			name := c.GlobalString("name")
			return pcli.EphemeralDecrypt(name, message)
		},
	}
}

func DecryptGroup() cli.Command {
	return cli.Command{
		Name:    "decrypt-group",
		Aliases: []string{"dg"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "<URL> - for loading decode data by link instead of text input data",
			},
		},
		Usage: "dg <GROUP_NAME> base64",
		Action: func(c *cli.Context) error {
			groupName := c.Args().Get(0)
			var message string
			if link := c.String("link"); len(link) > 0 {
				m, err := api.DownloadByLink(link)
				if err != nil {
					return errors.Wrap(err, "download by link failed")
				}
				message = m
			} else {
				message = c.Args().Get(1)
			}
			name := c.GlobalString("name")

			return pcli.DecryptGroup(name, message, groupName)
		},
	}
}

func EphemeralDecryptGroup() cli.Command {
	return cli.Command{
		Name:    "ephemeral-decrypt-group",
		Aliases: []string{"edg"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "<URL> - for loading decode data by link instead of text input data",
			},
		},
		Usage: `Decrypt input message with ephemral shared key`,
		Action: func(c *cli.Context) error {
			groupName := c.Args().Get(0)
			var message string
			if link := c.String("link"); len(link) > 0 {
				m, err := api.DownloadByLink(link)
				if err != nil {
					return errors.Wrap(err, "download by link failed")
				}
				message = m
			} else {
				message = c.Args().Get(1)
			}
			name := c.GlobalString("name")
			return pcli.EphemeralDecryptGroup(name, message, groupName)
		},
	}
}
