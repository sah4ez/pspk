package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/sah4ez/pspk/pkg/config"
	environment "github.com/sah4ez/pspk/pkg/evnironment"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

const (
	baseURL = "https://pspk.now.sh"
)

var (
	//Version current tools
	Version string
	// Hash revision number from git
	Hash string
	// BuildDate when building this utilitites
	BuildDate string
)

func main() {
	var (
		err error
	)

	var api pspk.PSPK
	{
		api = pspk.New(baseURL)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config", err.Error())
		return
	}
	path := environment.LoadDataPath()

	app := cli.NewApp()
	app.Name = "pspk"
	app.Version = Version + "." + Hash
	app.Metadata = map[string]interface{}{"builded": BuildDate}
	app.Description = "Console tool for encyption/decription data through pspk.now.sh"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "key name",
		}}

	app.Commands = []cli.Command{
		{
			Name:    "publish",
			Usage:   `Generate x25519 pair to pspk`,
			Aliases: []string{"p"},
			Action: func(c *cli.Context) error {
				name := c.GlobalString("name")
				if name == "" {
					if cfg.CurrentName == "" {
						return fmt.Errorf("empty current name, set to config or use --name")
					}
					name = cfg.CurrentName
				}
				path = path + "/" + name

				pub, priv, err := keys.GenereateDH()
				if err != nil {
					return err
				}
				err = utils.Write(path, "pub.bin", pub[:])
				if err != nil {
					return err
				}
				err = api.Publish(name, pub[:])
				if err != nil {
					return err
				}

				err = utils.Write(path, "key.bin", priv[:])
				if err != nil {
					return err
				}

				fmt.Println("Generate key pair on x25519")
				return nil
			},
		},
		{
			Name:    "secret",
			Aliases: []string{"s"},
			Usage:   `Generate shared secret key by private and public keys from pspk by name`,
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
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   `Encrypt input message with shared key`,
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
		},
		{
			Name:    "ephemeral-encrypt",
			Aliases: []string{"ee"},
			Usage:   `Encrypt input message with ephemeral key`,
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
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   `Decrypt input message with shared key`,
			Action: func(c *cli.Context) error {
				pubName := c.Args().Get(0)
				message := c.Args().Get(1)
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
		},
		{
			Name:    "ephemeral-decrypt",
			Aliases: []string{"ed"},
			Usage:   `Decrypt input message with ephemral shared key`,
			Action: func(c *cli.Context) error {
				message := c.Args().Get(0)
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
		},
		{
			Name:  "use-current",
			Usage: `Set currnet name by default`,
			Action: func(c *cli.Context) error {
				name := c.GlobalString("name")
				if name == "" {
					return fmt.Errorf("empty name use  --name")
				}
				cfg.CurrentName = name
				return cfg.Save()
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run has error:", err.Error())
	}

}
