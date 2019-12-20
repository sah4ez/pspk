package main

import (
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
			return pcli.Group(name)
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

			names := make([]string, len(c.Args()[1:]))
			copy(names, c.Args()[1:])

			return pcli.StartGroup(name, groupName, names...)
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

			names := make([]string, len(c.Args()[1:]))
			copy(names, c.Args()[1:])

			return pcli.FinishGroup(name, groupName, names...)
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

			names := make([]string, len(c.Args()[1:]))
			copy(names, c.Args()[1:])

			return pcli.SecretGroup(name, groupName, names...)
		},
	}
}
