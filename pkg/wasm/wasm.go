// +build js,wasm

package wasm

import (
	"os"
	"syscall/js"

	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/utils"
)

var (
	cbs = map[string]js.Func{
		"NewPspkAPI": js.FuncOf(JsNewPspkAPI),
		"NewPspkCLI": js.FuncOf(JsNewPspkCLI),
	}
)

func Load() (release func()) {
	for key, cb := range cbs {
		js.Global().Set(key, cb)
	}
	return func() {
		for _, cb := range cbs {
			cb.Release()
		}
	}
}

func JsDecrypt(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		//Decrypt(name, message, pubName string) (err error)
		return nil
	}
}

func JsEphemeralDecrypt(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// EphemeralDecrypt(name, message string) (err error)
		return nil
	}
}

func JsDecryptGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// DecryptGroup(name, message, groupName string) (err error)
		return nil
	}
}

func JsEphemeralDecryptGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// EphemeralDecryptGroup(name, message, groupName string) (err error)
		return nil
	}
}

func JsEncrypt(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Encrypt(name, message, pubName string, link bool) (err error)
		return nil
	}
}

func JsEphemeralEncrypt(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// EphemeralEncrypt(message, pubName string, link bool) (err error)
		return nil
	}
}

func JsEncryptGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// EncryptGroup(name, message, groupName string, link bool) (err error)
		return nil
	}
}

func JsEphemeralEncrypGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// EphemeralEncrypGroup(name, message, groupName string, link bool) (err error)
		return nil
	}
}

func JsSecret(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Secret(name, pubName string) (err error)
		return nil
	}
}

func JsPublish(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Publish(name string) (err error)
		jsName := args[0]

		name := jsName.String()

		var result error
		result = cli.Publish(name)

		if result != nil {
			return map[string]interface{}{
				"error": result.Error(),
			}
		}
		obj := map[string]interface{}{}
		return obj
	}
}

func JsGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Group(name string) (err error)
		return nil
	}
}

func JsStartGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// StartGroup(name, groupName string, names ...string) (err error)
		return nil
	}
}

func JsFinishGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// FinishGroup(name, groupName string, names ...string) (err error)
		return nil
	}
}

func JsSecretGroup(cli *pspk.PSPKcli) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// SecretGroup(name, groupName string, names ...string) (err error)
		return nil
	}
}

func JsApiPublish(api pspk.PSPK) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Publish(name string, key []byte) (err error)
		return nil
	}
}

func JsApiLoad(api pspk.PSPK) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// Load(name string) (key *PublicKey)
		return nil
	}
}

func JsApiGenerateLink(api pspk.PSPK) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// GenerateLink(data string) (link *Link)
		return nil
	}
}

func JsApiDownloadByLink(api pspk.PSPK) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// DownloadByLink(link string) (data *DownloadData)
		return nil
	}
}

func JsApiGetAll(api pspk.PSPK) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		// GetAll(opts GetAllOptions) (keys *Keys)
		return nil
	}
}

func JsNewPspkCLI(this js.Value, args []js.Value) interface{} {
	api := pspk.New("https://pspk.now.sh/")
	fs := utils.NewWasmStorage()
	cli := pspk.NewPSPKcli(api, nil, "/", "https://pspk.now.sh", os.Stdout, fs)
	wrapper := map[string]interface{}{
		"Decrypt":               js.FuncOf(JsDecrypt(cli)),
		"EphemeralDecrypt":      js.FuncOf(JsEphemeralDecrypt(cli)),
		"DecryptGroup":          js.FuncOf(JsDecryptGroup(cli)),
		"EphemeralDecryptGroup": js.FuncOf(JsEphemeralDecryptGroup(cli)),
		"Encrypt":               js.FuncOf(JsEncrypt(cli)),
		"EphemeralEncrypt":      js.FuncOf(JsEphemeralEncrypt(cli)),
		"EncryptGroup":          js.FuncOf(JsEncryptGroup(cli)),
		"EphemeralEncrypGroup":  js.FuncOf(JsEphemeralEncrypGroup(cli)),
		"Secret":                js.FuncOf(JsSecret(cli)),
		"Publish":               js.FuncOf(JsPublish(cli)),
		"Group":                 js.FuncOf(JsGroup(cli)),
		"StartGroup":            js.FuncOf(JsStartGroup(cli)),
		"FinishGroup":           js.FuncOf(JsFinishGroup(cli)),
		"SecretGroup":           js.FuncOf(JsSecretGroup(cli)),
	}
	return wrapper
}

func JsNewPspkAPI(this js.Value, args []js.Value) interface{} {
	api := pspk.New("https://pspk.now.sh/")
	wrapper := map[string]interface{}{
		"Publish":        js.FuncOf(JsApiPublish(api)),
		"Load":           js.FuncOf(JsApiLoad(api)),
		"GenerateLink":   js.FuncOf(JsApiGenerateLink(api)),
		"DownloadByLink": js.FuncOf(JsApiDownloadByLink(api)),
		"GetAll":         js.FuncOf(JsApiGenerateLink(api)),
	}
	return wrapper
}
func BytesToJS(b []byte) js.Value {
	jsB := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(jsB, b)
	return jsB
}

func BytesToGo(jsB js.Value) []byte {
	b := make([]byte, jsB.Length())
	js.CopyBytesToGo(b, jsB)
	return b
}
