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

### Group encryption exchange
For encryption/decryption need generate shared secret in group. 
Use this algorithm (CLIQUES) [IV.A](https://pdfs.semanticscholar.org/dc45/970a9c43aaff17295c3769fdd0af9bded855.pdf) 

1. *Creat group*.
Create prime base point for group `base` and publish to pspk.now.sh
```
pspk --name base group
```
2. Decide number of members in group and select order for generation secret.
As example Alice, Bob, Carol and Daron want creage shared secret in group `base`.
```bash
pspk --name alice start-group base 
pspk --name bob start-group base alice
pspk --name carol start-group base bob alice
```
The last members finis generate intermediate secrets.
```bash
pspk --name daron finish-group base carol bob alice
```
Members can start generate shared secret keys via intermediate keys.
```bash
pspk --name daron secret-group base carol bob alice
pspk --name carol secret-group base daron bob alice
pspk --name bob secret-group base daron carol alice
pspk --name alice secret-group base daron carol bob
```
3. *Encryption* Encrypt some messages for `base` group members
```bash
pspk --name alice ephemeral-encrypt-group base Super secret message
```
4. *Decription* Decrypt the message from member's `base` group
```bash
pspk --name bob ephemeral-decryp-group base base64
```

**NOTE** All intermediate secrets would saved in pspk storage!

## pspk config

pspk use `$XDG_CONFIG_HOME` for saving configuration or default value `$HOME/.config/pspk`
Use `config.json` file for saving configuration:
```
{"current_name":"name"}
```

Also pspk use `$XDG_DATA_HOME` for saving appication data ro default value `$HOME/.local/share/pspk`
