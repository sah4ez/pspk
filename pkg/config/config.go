package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sah4ez/pspk/pkg/evnironment"
	"github.com/sah4ez/pspk/pkg/utils"
)

type Config struct {
	CurrentName string `json:"current_name,omitempty"`
}

var (
	configName = "config.json"
	path       = ""
)

func init() {
	path = environment.LoadConfigPath()

	os.OpenFile(path+"/"+configName, os.O_RDONLY|os.O_CREATE, 0666)
}

func Load() (c *Config, err error) {
	b, err := utils.Read(path, configName)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return &Config{}, nil
	}

	c = &Config{}

	err = json.Unmarshal(b, c)
	return
}

func (c *Config) Save() (err error) {
	b, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = utils.Write(path, configName, b)
	if err != nil {
		return
	}
	return
}

func (c *Config) LoadCurrentName(name string) (string, error) {
	if name == "" {
		if c.CurrentName == "" {
			return "", fmt.Errorf("empty current name, set to config or use --name")
		}
	}
	return c.CurrentName, nil
}
