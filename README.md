## Public storage public key

Simple in-memory public storage public key.

HTTP-base key-value storage.

## Model

You should use this model:
```json
{
	"name":"Some Name",
	"key":"base64=="
}
```

## Request

Save public key:
```bash
curl -X POST "https://pspk.now.sh" -d '{"name":"Some.Name","key":"E7+TL112lj1GmJRHf9jT5MZJDgYIhUbtBLc4/ZFMZ5c="}'
```

Read public key:
```bash
curl -X POST "https://pspk.now.sh" -d '{"name":"Some.Name"}'
```

## pspk cli usage

### Generation key pair.
Will generation private and public keys and publish public pice to pspk.now.sh.
```bash
pspk --name <NAME_YOUR_KEY> publish
```

### Encryption some text message. 
Will encryption message through your private key and public key name from pspk.now.sh.
```bash
pspk --name <NAME_YOUR_KEY> encrypt <PUBLIC_PART> <SOME_MESSAGE_WITH_SPACES>
```

### Decription some text message. 
Will decription message through your private key and public key name from pspk.now.sh.
```bash
pspk --name <NAME_YOUR_KEY> decrypt <PUBLIC_PART> <SOME_BASE64_WITH_SPACES>
```
