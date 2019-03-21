package environment

import "os"

var (
	defaultDataPath   = os.Getenv("HOME") + "/.local/share"
	defaultConfigPath = os.Getenv("HOME") + "/.config"
	name              = "/pspk"
)

func LoadDataPath() string {
	env, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		return defaultDataPath + name
	}
	return env + name
}

func LoadConfigPath() string {
	env, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		return defaultConfigPath + name
	}
	return env + name
}
