package main

import (
	"strings"

	"github.com/urfave/cli"
)

func Encrypt() cli.Command {
	return cli.Command{
		Name:    "encrypt",
		Aliases: []string{"e"},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "link",
				Usage: "for generation 24hr link for loading data",
			},
		},
		Usage:       "ecnrypt pub_name some message will encrypt",
		Description: `Encrypt input message with shared key`,
		Action: func(c *cli.Context) error {
			pubName := c.Args()[0]
			message := c.Args()[1:]
			name := c.GlobalString("name")
			link := c.Bool("link")

			return pcli.Encrypt(name, strings.Join(message, " "), pubName, link)
		},
	}
}

func EphemeralEncrypt() cli.Command {
	return cli.Command{
		Name:    "ephemeral-encrypt",
		Aliases: []string{"ee"},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "link",
				Usage: "for generation 24hr link for loading data",
			},
		},
		Usage:       "ee pub_name some message will encrypt",
		Description: `Encrypt input message with ephemeral key`,
		Action: func(c *cli.Context) error {
			pubName := c.Args()[0]
			message := c.Args()[1:]
			link := c.Bool("link")
<<<<<<< HEAD
			return pcli.EphemeralEncrypt(strings.Join(message, " "), pubName, link)
		},
	}
}

func EncryptGroup() cli.Command {
	return cli.Command{
		Name:    "encrypt-group",
		Aliases: []string{"eg"},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "link",
				Usage: "for generation 24hr link for loading data",
			},
		},
		Usage:       "eg <GROUP_NAME> message",
		Description: "Encrypt message for group",
		Action: func(c *cli.Context) error {
			groupName := c.Args()[0]
			message := c.Args()[1:]
			name := c.GlobalString("name")
			link := c.Bool("link")

			return pcli.EncryptGroup(name, strings.Join(message, " "), groupName, link)
		},
	}
}

func EphemeralEncrypGroup() cli.Command {
	return cli.Command{
		Name:    "ephemeral-encrypt-group",
		Aliases: []string{"eeg"},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "link",
				Usage: "for generation 24hr link for loading data",
			},
		},
		Usage: `Encrypt input message with ephemeral key`,
		Action: func(c *cli.Context) error {
			groupName := c.Args()[0]
			message := c.Args()[1:]
			name := c.GlobalString("name")
			link := c.Bool("link")

			return pcli.EphemeralEncrypGroup(name, strings.Join(message, " "), groupName, link)
		},
	}
}
