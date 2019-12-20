package pspk

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/sah4ez/pspk/pkg/config"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
)

// TODO split to Decrypter, Encryper and etc.
type CLI interface {
	// Decrypter
	Decrypt(name, message, pubName string) (err error)
	EphemeralDecrypt(name, message string) (err error)
	DecryptGroup(name, message, groupName string) (err error)
	EphemeralDecryptGroup(name, message, groupName string) (err error)
	// Encrypter
	Encrypt(name, message, pubName string, link bool) (err error)
	EphemeralEncrypt(message, pubName string, link bool) (err error)
	EncryptGroup(name, message, groupName string, link bool) (err error)
	EphemeralEncrypGroup(name, message, groupName string, link bool) (err error)
}

type PSPKcli struct {
	cfg     *config.Config
	api     PSPK
	path    string
	baseURL string
	out     io.Writer
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

func (p *PSPKcli) Encrypt(name, message, pubName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, "key.bin")
	if err != nil {
		return fmt.Errorf("read key.bin %w", err)
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

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return err
	}
	data := base64.StdEncoding.EncodeToString(b)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)
}

func (p *PSPKcli) EphemeralEncrypt(message, pubName string, link bool) (err error) {
	pubEphemeral, privEphemeral, err := keys.GenerateDH()
	if err != nil {
		return err
	}
	pub, err := p.api.Load(pubName)
	if err != nil {
		return err
	}
	chain := keys.Secret(privEphemeral[:], pub)

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return err
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return err
	}
	m := append(pubEphemeral[:], b...)
	data := base64.StdEncoding.EncodeToString(m)
	fmt.Fprintln(p.out, data)

	return p.generateLink(link, data)
}

func (p *PSPKcli) EncryptGroup(name, message, groupName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, groupName+".secret")
	if err != nil {
		return err
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

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return err
	}
	data := base64.StdEncoding.EncodeToString(b)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)
}

func (p *PSPKcli) EphemeralEncrypGroup(name, message, groupName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}

	priv, err := utils.Read(path, groupName+".secret")
	if err != nil {
		return err
	}

	pubEphemeral, _, err := keys.GenerateDH()
	if err != nil {
		return err
	}
	chain := keys.Secret(priv[:], pubEphemeral[:])

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return err
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return err
	}
	m := append(pubEphemeral[:], b...)
	data := base64.StdEncoding.EncodeToString(m)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)

}

func (p *PSPKcli) loadPath(name string) (path string, err error) {

	if name, err = p.cfg.LoadCurrentName(name); err != nil {
		return
	}

	return p.path + "/" + name, nil
}

func (p *PSPKcli) generateLink(isLink bool, data string) error {
	if isLink {
		id, err := p.api.GenerateLink(data)
		if err != nil {
			return err
		}
		fmt.Fprintln(p.out, p.baseURL+"/?link="+id)
	}
	return nil
}

// NewPSPKcli return new API client for CLI interface
func NewPSPKcli(api PSPK, cfg *config.Config, basePath string, baseURL string, out io.Writer) *PSPKcli {
	return &PSPKcli{
		cfg:     cfg,
		api:     api,
		path:    basePath,
		baseURL: baseURL,
		out:     out,
	}
}
