package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

func Encrypt() cli.Command {
	return cli.Command{
		Name:        "encrypt",
		Aliases:     []string{"e"},
		Usage:       "ecnrypt pub_name some message will encrypt",
		Description: `Encrypt input message with shared key`,
		Action: func(c *cli.Context) error {
			pubName := c.Args()[0]
			message := c.Args()[1:]
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
			chain := keys.Secret(priv, pub)

			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}

			b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(strings.Join(message, " ")))
			if err != nil {
				return err
			}
			fmt.Println(base64.StdEncoding.EncodeToString(b))
			return nil
		},
	}
}

func EphemeralEncrypt() cli.Command {
	return cli.Command{
		Name:        "ephemeral-encrypt",
		Aliases:     []string{"ee"},
		Usage:       "ee pub_name some message will encrypt",
		Description: `Encrypt input message with ephemeral key`,
		Action: func(c *cli.Context) error {
			pubName := c.Args()[0]
			message := c.Args()[1:]

			pubEphemeral, privEphemeral, err := keys.GenereateDH()
			if err != nil {
				return err
			}
			pub, err := api.Load(pubName)
			if err != nil {
				return err
			}
			chain := keys.Secret(privEphemeral[:], pub)

			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}

			b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(strings.Join(message, " ")))
			if err != nil {
				return err
			}
			m := append(pubEphemeral[:], b...)
			fmt.Println(base64.StdEncoding.EncodeToString(m))
			return nil
		},
	}
}

func EncryptGroup() cli.Command {
	return cli.Command{
		Name:        "encrypt-group",
		Aliases:     []string{"eg"},
		Usage:       "eg <GROUP_NAME> message",
		Description: "Encrypt message for group",
		Action: func(c *cli.Context) error {
			groupName := c.Args()[0]
			message := c.Args()[1:]
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
				return err
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

			b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(strings.Join(message, " ")))
			if err != nil {
				return err
			}
			fmt.Println(base64.StdEncoding.EncodeToString(b))
			return nil
		},
	}
}

func EphemeralEncrypGroup() cli.Command {
	return cli.Command{
		Name:    "ephemeral-encrypt-group",
		Aliases: []string{"eeg"},
		Usage:   `Encrypt input message with ephemeral key`,
		Action: func(c *cli.Context) error {
			groupName := c.Args()[0]
			message := c.Args()[1:]

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
				return err
			}

			pubEphemeral, _, err := keys.GenereateDH()
			if err != nil {
				return err
			}
			chain := keys.Secret(priv[:], pubEphemeral[:])

			messageKey, err := keys.LoadMaterialKey(chain)
			if err != nil {
				return err
			}

			b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(strings.Join(message, " ")))
			if err != nil {
				return err
			}
			m := append(pubEphemeral[:], b...)
			fmt.Println(base64.StdEncoding.EncodeToString(m))
			return nil
		},
	}
}
