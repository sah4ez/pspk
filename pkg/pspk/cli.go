package pspk

import (
	"encoding/base64"
	"fmt"

	"github.com/sah4ez/pspk/pkg/config"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
)

type CLI interface {
	Decrypt(name, message, pubName string) (err error)
	EphemeralDecrypt(name, message string) (err error)
	DecryptGroup(name, message, groupName string) (err error)
	EphemeralDecryptGroup(name, message, groupName string) (err error)
}

type PSPKcli struct {
	cfg  *config.Config
	api  PSPK
	path string
}

func (p *PSPKcli) loadPath(name string) (path string, err error) {

	if name, err = p.cfg.LoadCurrentName(name); err != nil {
		return
	}

	return p.path + "/" + name, nil
}

// Decrypt decryp message by name key, public name of recipeint and message
func (p *PSPKcli) Decrypt(name, message, pubName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, "key.bin")
	if err != nil {
		return fmt.Errorf("read key.bin: %w", err)
	}
	pub, err := p.api.Load(pubName)
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
}

// EphemeralDecrypt decryp message by name key, ephemeral key and message
func (p *PSPKcli) EphemeralDecrypt(name, message string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, "key.bin")
	if err != nil {
		return fmt.Errorf("read key.bin: %w", err)
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

}

func (p *PSPKcli) DecryptGroup(name, message, groupName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, groupName+".secret")
	if err != nil {
		return fmt.Errorf("read key.bin: %w", err)
	}
	pub, err := p.api.Load(groupName)
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
}

func (p *PSPKcli) EphemeralDecryptGroup(name, message, groupName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, groupName+".secret")
	if err != nil {
		return fmt.Errorf("read group secret: %w", err)
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

}

// NewPSPKcli return new API client for CLI interface
func NewPSPKcli(api PSPK, cfg *config.Config, basePath string) *PSPKcli {
	return &PSPKcli{
		cfg:  cfg,
		api:  api,
		path: basePath,
	}
}
