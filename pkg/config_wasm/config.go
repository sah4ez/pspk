// +build js,wasm

package config

type Config struct {
	CurrentName string `json:"current_name,omitempty"`
}

var (
	configName = "config.json"
	path       = ""
)

func Load() (c *Config, err error) {
	return &Config{
		CurrentName: "",
	}, nil
}

func (c *Config) Save() (err error) {
	return nil
}

func (c *Config) LoadCurrentName(name string) (string, error) {
	return name, nil
}
