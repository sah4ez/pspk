package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func UseCurrent() cli.Command {
	return cli.Command{
		Name:        "use-current",
		Aliases:     []string{"uc"},
		Description: `Set currnet name by default`,
		Usage:       "--name name_pub_key use-current",
		Action: func(c *cli.Context) error {
			name := c.GlobalString("name")
			if name == "" {
				return fmt.Errorf("empty name use  --name")
			}
			cfg.CurrentName = name
			return cfg.Save()
		},
	}
}
