package pspk

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/sah4ez/pspk/pkg/config"
	environment "github.com/sah4ez/pspk/pkg/evnironment"
	"github.com/sah4ez/pspk/pkg/keys"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/skip2/go-qrcode"
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

	Secret(name, pubName string) (err error)
	Publish(name string) (err error)
	PublishAndGenerateQR(name string, path string) (err error)

	Group(name string) (err error)
	StartGroup(name, groupName string, names ...string) (err error)
	FinishGroup(name, groupName string, names ...string) (err error)
	SecretGroup(name, groupName string, names ...string) (err error)

	Verify(name string, signature string, data []byte) (err error)
}

type PSPKcli struct {
	cfg     *config.Config
	api     PSPK
	path    string
	baseURL string
	out     io.Writer
	fs      utils.FS
}

// Decrypt decryp message by name key, public name of recipeint and message
func (p *PSPKcli) Decrypt(name, message, pubName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}
	pub := p.api.Load(pubName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load the public name")
	}
	chain := keys.Secret(priv, pub.Key)
	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load message key")
	}
	bytesMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "bytesMessage has error")
	}

	b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage)
	if err != nil {
		return errors.Wrap(err, "decrypt has error")
	}
	fmt.Println(string(b))
	return nil
}

// EphemeralDecrypt decryp message by name key, ephemeral key and message
func (p *PSPKcli) EphemeralDecrypt(name, message string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}
	bytesMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "error when decode in base64 the message")
	}
	chain := keys.Secret(priv, bytesMessage[:32])
	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load message key")
	}

	b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage[32:])
	if err != nil {
		return errors.Wrap(err, "can not decryp message key")
	}
	fmt.Println(string(b))
	return nil

}

func (p *PSPKcli) DecryptGroup(name, message, groupName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, groupName+".secret")
	if err != nil {
		return errors.Wrap(err, "can not read group secret")
	}
	pub := p.api.Load(groupName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load group name")
	}
	chain := keys.Secret(priv, pub.Key)
	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load message key")
	}
	bytesMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "can not decode in base64 the message")
	}

	b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage)
	if err != nil {
		return errors.Wrap(err, "can not decrypt message")
	}
	fmt.Println(string(b))
	return nil
}

func (p *PSPKcli) EphemeralDecryptGroup(name, message, groupName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, groupName+".secret")
	if err != nil {
		return errors.Wrap(err, "can not read group secret")
	}
	bytesMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "can not decode in base64 the message")
	}
	chain := keys.Secret(priv, bytesMessage[:32])
	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load keys")
	}

	b, err := utils.Decrypt(messageKey[64:], messageKey[:32], bytesMessage[32:])
	if err != nil {
		return errors.Wrap(err, "can not decrypt message")
	}
	fmt.Println(string(b))
	return nil
}

func (p *PSPKcli) Encrypt(name, message, pubName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}
	pub := p.api.Load(pubName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load public name")
	}
	chain := keys.Secret(priv, pub.Key)

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load message key")
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return errors.Wrap(err, "can not encrypt message")
	}
	data := base64.StdEncoding.EncodeToString(b)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)
}

func (p *PSPKcli) EphemeralEncrypt(message, pubName string, link bool) (err error) {
	pubEphemeral, privEphemeral, err := keys.GenerateDH()
	if err != nil {
		return errors.Wrap(err, "can not generate keys")
	}
	pub := p.api.Load(pubName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load public name")
	}
	chain := keys.Secret(privEphemeral[:], pub.Key)

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load keys")
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return errors.Wrap(err, "can not encrypt message")
	}
	m := append(pubEphemeral[:], b...)
	data := base64.StdEncoding.EncodeToString(m)
	fmt.Fprintln(p.out, data)

	return p.generateLink(link, data)
}

func (p *PSPKcli) EncryptGroup(name, message, groupName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, groupName+".secret")
	if err != nil {
		return errors.Wrap(err, "can not read group name")
	}
	pub := p.api.Load(groupName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "api load can not group name")
	}
	chain := keys.Secret(priv, pub.Key)

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load keys")
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return errors.Wrap(err, "can not Encrypt message")
	}
	data := base64.StdEncoding.EncodeToString(b)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)
}

func (p *PSPKcli) EphemeralEncrypGroup(name, message, groupName string, link bool) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, groupName+".secret")
	if err != nil {
		return errors.Wrap(err, "can not read group name")
	}

	pubEphemeral, _, err := keys.GenerateDH()
	if err != nil {
		return errors.Wrap(err, "can not generate keys")
	}
	chain := keys.Secret(priv[:], pubEphemeral[:])

	messageKey, err := keys.LoadMaterialKey(chain)
	if err != nil {
		return errors.Wrap(err, "can not load keys")
	}

	b, err := utils.Encrypt(messageKey[64:], messageKey[:32], []byte(message))
	if err != nil {
		return errors.Wrap(err, "can not encrypt message")
	}
	m := append(pubEphemeral[:], b...)
	data := base64.StdEncoding.EncodeToString(m)
	fmt.Fprintln(p.out, data)
	return p.generateLink(link, data)
}

func (p *PSPKcli) Secret(name, pubName string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not find the path")
	}
	pub := p.api.Load(pubName)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load the public name")
	}
	dh := keys.Secret(priv, pub.Key)
	fmt.Fprintln(p.out, base64.StdEncoding.EncodeToString(dh))

	err = p.fs.Write(path, pubName+".secret.bin", dh[:])
	if err != nil {
		return errors.Wrap(err, "can not write")
	}
	return nil
}

func (p *PSPKcli) Publish(name string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	pub, priv, err := keys.GenerateDH()
	if err != nil {
		return errors.Wrap(err, "can not generate keys")
	}
	err = p.fs.Write(path, "pub.bin", pub[:])
	if err != nil {
		return errors.Wrap(err, "can not write in pub.bin")
	}
	err = p.api.Publish(name, pub[:])
	if err != nil {
		return errors.Wrap(err, "can not publish")
	}

	err = p.fs.Write(path, "key.bin", priv[:])
	if err != nil {
		return errors.Wrap(err, "can not find the path")
	}

	fmt.Fprintln(p.out, "Generate key pair on x25519")
	return nil
}

func (p *PSPKcli) PublishAndGenerateQR(name string, qrPath string) (err error) {
	if err = p.Publish(name); err != nil {
		return
	}

	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}

	pub, err := p.fs.Read(path, "pub.bin")
	if err != nil {
		return errors.Wrap(err, "can't read pub.bin")
	}

	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}

	if qrPath != "" {
		qrPath = fmt.Sprintf("%s/%s.png", qrPath, name)
	} else {
		qrPath = fmt.Sprintf("%s/%[2]s/", environment.LoadDataPath(), name)
	}
	err = qrcode.WriteFile(string(pub[:]), qrcode.Medium, 256, qrPath+"pub.png")
	if err != nil {
		return errors.Wrap(err, "can not create qrcode pub file")
	}
	err = qrcode.WriteFile(string(priv[:]), qrcode.Medium, 256, qrPath+"key.png")
	if err != nil {
		return errors.Wrap(err, "can not create qrcode key file")
	}

	return
}

func (p *PSPKcli) Group(name string) (err error) {
	if name == "" {
		return fmt.Errorf("empty name use  --name")
	}
	pub, priv, err := keys.GenerateDH()
	if err != nil {
		return errors.Wrap(err, "can not generate keys")
	}
	base := keys.Secret(priv[:], pub[:])
	if err = p.api.Publish(name, base[:]); err != nil {
		return errors.Wrap(err, "can not publish name")
	}

	return nil
}

func (p *PSPKcli) StartGroup(name, groupName string, names ...string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return fmt.Errorf("load path to keys: %w", err)
	}
	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}
	base := p.api.Load(groupName)
	if base.Error != nil {
		return errors.Wrap(base.Error, "can not load group name")
	}
	publicGroup := keys.Secret(priv, base.Key)
	err = p.api.Publish(name+groupName, publicGroup[:])
	if err != nil {
		return errors.Wrap(err, "can not publish group name")
	}

	local_names, err := p.processNames(name, groupName, priv, names...)
	if err != nil {
		return errors.Wrap(err, "can not process names of group")
	}
	// TODO add print the remaining users
	if len(local_names) > 0 {
		intermediate := strings.Join(local_names, "") + groupName

		pub := p.api.Load(intermediate)
		if pub.Error != nil {
			return errors.Wrap(err, "start-join-group load error")
		}
		dh := keys.Secret(priv, pub.Key)
		err = p.api.Publish(name+intermediate, dh[:])
		if err != nil {
			return errors.Wrap(err, "start-join-group publish error")
		}
	}

	return nil
}

func (p *PSPKcli) FinishGroup(name, groupName string, names ...string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}
	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}

	base := p.api.Load(groupName)
	if base.Error != nil {
		return errors.Wrap(base.Error, "can not load group name")
	}
	publicGroup := keys.Secret(priv, base.Key)
	err = p.api.Publish(name+groupName, publicGroup[:])
	if err != nil {
		return errors.Wrap(err, "can not publish")
	}
	if _, err = p.processNames(name, groupName, priv, names...); err != nil {
		return errors.Wrap(err, "can not process name users of group")
	}
	return nil
}

func (p *PSPKcli) SecretGroup(name, groupName string, names ...string) (err error) {
	path, err := p.loadPath(name)
	if err != nil {
		return errors.Wrap(err, "load path to keys")
	}
	priv, err := p.fs.Read(path, "key.bin")
	if err != nil {
		return errors.Wrap(err, "can not read key.bin")
	}
	intermediate := strings.Join(names, "") + groupName

	pub := p.api.Load(intermediate)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load group name")
	}
	publicGroup := keys.Secret(priv, pub.Key)
	err = p.fs.Write(path, groupName+".secret", publicGroup[:])
	if err != nil {
		return errors.Wrap(err, "can not write in group name")
	}
	return nil

}

func (p *PSPKcli) processNames(name, groupName string, priv []byte, names ...string) (local_names []string, err error) {
	local_names = make([]string, len(names))
	copy(local_names, names)

	for i, _ := range local_names {
		n := []string{}
		n = append(n, local_names[:i]...)
		n = append(n, local_names[i+1:]...)
		n = append(n, groupName)
		if len(n) > 0 {
			intermediate := strings.Join(n, "")

			pub := p.api.Load(intermediate)
			if pub.Error != nil {
				return nil, fmt.Errorf("failed load intermediate key: %w", pub.Error)
			}
			dh := keys.Secret(priv, pub.Key)
			err = p.api.Publish(name+intermediate, dh[:])
			if err != nil {
				return nil, fmt.Errorf("failed publish intermediate key: %w", err)
			}
		}
	}

	return
}

func (p *PSPKcli) loadPath(name string) (path string, err error) {

	if p.cfg == nil {
		return "/" + name, nil
	}
	if name, err = p.cfg.LoadCurrentName(name); err != nil {
		return
	}

	return p.path + "/" + name, nil
}

func (p *PSPKcli) generateLink(isLink bool, data string) error {
	if isLink {
		link := p.api.GenerateLink(data)
		if link.Error != nil {
			return link.Error
		}
		fmt.Fprintln(p.out, p.baseURL+"/?link="+link.Link)
	}
	return nil
}

func (p *PSPKcli) Verify(name string, signature string, data []byte) (err error) {

	pub := p.api.Load(name)
	if pub.Error != nil {
		return errors.Wrap(pub.Error, "can not load public name")
	}

	pubArray := utils.Slice2Array32(pub.Key)

	signatureBinary, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.Wrap(err, "can decode signature from base64")
	}
	signatureArrya := utils.Slice2Array64(signatureBinary)

	verify := keys.Verify(pubArray, data, &signatureArrya)
	fmt.Fprintln(p.out, verifyMessage(signature, verify))
	return
}

func verifyMessage(signature string, verify bool) string {
	if verify {
		return fmt.Sprintf("Signature %s is valid.", signature)
	}
	return fmt.Sprintf("Signature %s is NOT valid.", signature)
}

// NewPSPKcli return new API client for CLI interface
func NewPSPKcli(api PSPK, cfg *config.Config, basePath string, baseURL string, out io.Writer, fs utils.FS) *PSPKcli {
	if cfg != nil {
		cfg.Init()
	}

	return &PSPKcli{
		cfg:     cfg,
		api:     api,
		path:    basePath,
		baseURL: baseURL,
		out:     out,
		fs:      fs,
	}
}
