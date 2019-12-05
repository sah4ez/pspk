package main

import (
	"fmt"
	"strings"

	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/urfave/cli"
)

func Group() cli.Command {
	return cli.Command{
		Name:        "group",
		Aliases:     []string{"g"},
		Description: "create prime base point and publish to pspk.now.sh",
		Usage:       "--name base_name group",
		Action: func(c *cli.Context) error {
			name := c.GlobalString("name")
			if name == "" {
				return fmt.Errorf("empty name use  --name")
			}
			pub, priv, err := keys.GenerateDH()
			if err != nil {
				return err
			}
			base := keys.Secret(priv[:], pub[:])
			err = api.Publish(name, base[:])
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func StartGroup() cli.Command {
	return cli.Command{
		Name:        "start-group",
		Aliases:     []string{"sg"},
		Usage:       `start-group groupName [pubName1 pubName2 ...]`,
		Description: "calculate intermediate keys",
		Action: func(c *cli.Context) error {
			groupName := c.Args().Get(0)
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
			base, err := api.Load(groupName)
			if err != nil {
				return err
			}
			publicGroup := keys.Secret(priv, base)
			err = api.Publish(name+groupName, publicGroup[:])
			if err != nil {
				return err
			}

			names := make([]string, len(c.Args()[1:]))
			copy(names, c.Args()[1:])

			for i, _ := range names {
				n := []string{}
				n = append(n, names[:i]...)
				n = append(n, names[i+1:]...)
				n = append(n, groupName)
				if len(n) > 0 {
					intermediate := strings.Join(n, "")
					pub, err := api.Load(intermediate)
					if err != nil {
						fmt.Println("start-join-group load error: ", err.Error())
						return err
					}
					dh := keys.Secret(priv, pub)
					err = api.Publish(name+intermediate, dh[:])
					if err != nil {
						fmt.Println("start-join-group publish error: ", err.Error())
						return err
					}
				}
			}
			if len(names) > 0 {
				intermediate := strings.Join(names, "") + groupName
				pub, err := api.Load(intermediate)
				if err != nil {
					fmt.Println("start-join-group load error: ", err.Error())
					return err
				}
				dh := keys.Secret(priv, pub)
				err = api.Publish(name+intermediate, dh[:])
				if err != nil {
					fmt.Println("start-join-group publish error: ", err.Error())
					return err
				}
			}

			return nil
		},
	}
}

func FinishGroup() cli.Command {
	return cli.Command{
		Name:        "finish-group",
		Aliases:     []string{"fg"},
		Usage:       `finish-group groupName pubName1 [pubName2 ...]`,
		Description: "calculate shared group keys",
		Action: func(c *cli.Context) error {
			groupName := c.Args().Get(0)
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
			base, err := api.Load(groupName)
			if err != nil {
				return err
			}
			publicGroup := keys.Secret(priv, base)
			err = api.Publish(name+groupName, publicGroup[:])
			if err != nil {
				return err
			}

			names := make([]string, len(c.Args()[1:]))
			copy(names, c.Args()[1:])

			for i, _ := range names {
				n := []string{}
				n = append(n, names[:i]...)
				n = append(n, names[i+1:]...)
				n = append(n, groupName)
				if len(n) > 0 {
					intermediate := strings.Join(n, "")
					pub, err := api.Load(intermediate)
					if err != nil {
						fmt.Println("start-join-group load error: ", err.Error())
						return err
					}
					dh := keys.Secret(priv, pub)
					err = api.Publish(name+intermediate, dh[:])
					if err != nil {
						fmt.Println("start-join-group publish error: ", err.Error())
						return err
					}
				}
			}
			return nil
		},
	}
}

func SecretGroup() cli.Command {
	return cli.Command{
		Name:        "secret-group",
		Aliases:     []string{"seg"},
		Usage:       `secret-group groupName pubName1 [pubName2 ...]`,
		Description: "calculate shared group keys",
		Action: func(c *cli.Context) error {
			groupName := c.Args().Get(0)
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
			intermediate := strings.Join(c.Args()[1:], "") + groupName
			pub, err := api.Load(intermediate)
			if err != nil {
				return err
			}
			publicGroup := keys.Secret(priv, pub)
			err = utils.Write(path, groupName+".secret", publicGroup[:])
			if err != nil {
				return err
			}
			return nil
		},
	}
}
