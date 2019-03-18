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
