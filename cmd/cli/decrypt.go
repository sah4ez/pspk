package main

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
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
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			priv, err := utils.Read(path, "key.bin")
			if err != nil {
				return errors.Wrap(err, "read key.bin")
			}
			pub, err := api.Load(pubName)
			if err != nil {
				return err
			}
			chain := keys.Secret(priv, pub)
			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}
			bytesMessage, err := base64.StdEncoding.DecodeString(message)
			if err != nil {
				return fmt.Errorf("bytesMessage has error: %s", err.Error())
			}

			b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage)
			if err != nil {
				return fmt.Errorf("decrypt has error: %s", err.Error())
			}
			fmt.Println(string(b))
			return nil
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
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			priv, err := utils.Read(path, "key.bin")
			if err != nil {
				return errors.Wrap(err, "read key.bin")
			}
			bytesMessage, err := base64.StdEncoding.DecodeString(message)
			if err != nil {
				return err
			}
			chain := keys.Secret(priv, bytesMessage[:32])
			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}

			b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage[32:])
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
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
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			priv, err := utils.Read(path, groupName+".secret")
			if err != nil {
				return errors.Wrap(err, "read group secret")
			}
			pub, err := api.Load(groupName)
			if err != nil {
				return err
			}
			chain := keys.Secret(priv, pub)
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
			fmt.Println(string(b))
			return nil
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
			if name == "" {
				if cfg.CurrentName == "" {
					return fmt.Errorf("empty current name, set to config or use --name")
				}
				name = cfg.CurrentName
			}
			path = path + "/" + name

			priv, err := utils.Read(path, groupName+".secret")
			if err != nil {
				return errors.Wrap(err, "read group secret")
			}
			bytesMessage, err := base64.StdEncoding.DecodeString(message)
			if err != nil {
				return err
			}
			chain := keys.Secret(priv, bytesMessage[:32])
			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}

			b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage[32:])
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}
}
