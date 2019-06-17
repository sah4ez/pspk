package main

import (
	"fmt"
	"os"

	"github.com/sah4ez/pspk/pkg/config"
	environment "github.com/sah4ez/pspk/pkg/evnironment"
	"github.com/sah4ez/pspk/pkg/pspk"
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
	// BuildDate when building this utilities
	BuildDate string
)

var (
	app  *cli.App
	api  pspk.PSPK
	cfg  *config.Config
	path string
	err  error
)

func init() {
	api = pspk.New(baseURL)
	app = cli.NewApp()
	app.Name = "pspk"
	app.Usage = "encrypt you message and send through open communication channel"
	app.Metadata = map[string]interface{}{"builded": BuildDate}
	app.Version = Version + "." + Hash
	app.Description = "Console tool for encyption/decription data through pspk.now.sh"
	cfg, err = config.Load()
	if err != nil {
		fmt.Println("load config has error", err.Error())
		os.Exit(2)
	}
	path = environment.LoadDataPath()
}

func main() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "key name",
		}}

	app.Commands = []cli.Command{
		Publish(),
		Secret(),
		Encrypt(),
		EphemeralEncrypt(),
		Decrypt(),
		EphemeralDecrypt(),
		UseCurrent(),
		Group(),
		StartGroup(),
		FinishGroup(),
		SecretGroup(),
		EncryptGroup(),
		EphemeralEncrypGroup(),
		DecryptGroup(),
		EphemeralDecryptGroup(),
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run has error:", err.Error())
	}
}
