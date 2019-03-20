package keys

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"io"

	"github.com/agl/ed25519"
	"github.com/agl/ed25519/edwards25519"
	"golang.org/x/crypto/curve25519"
	ecdh "golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/hkdf"
)

var (
	info = []byte("pspk_info")
)

func GenereateDH() (*[32]byte, *[32]byte, error) {
	random := rand.Reader

	// Create a byte array for our public and private keys.
	var private, public [32]byte

	// Generate some random data
	_, err := io.ReadFull(random, private[:])
	if err != nil {
		return nil, nil, err
	}

	// Documented at: http://cr.yp.to/ecdh.html
	private[0] &= 248
	private[31] &= 127
	private[31] |= 64

	curve25519.ScalarBaseMult(&public, &private)

	return &public, &private, nil
}

func Secret(priv ecdh.PrivateKey, pub ecdh.PublicKey) []byte {
	secret := new([32]byte)
	private := new([32]byte)
	public := new([32]byte)
	copy(private[:], priv[:32])
	copy(public[:], pub[:32])
	curve25519.ScalarMult(secret, private, public)
	return secret[:]
}

func HKDF(secret, info []byte, outputLength int) ([]byte, error) {
	// Underlying hash function for HMAC.
	hash := sha256.New

	// Non-secret salt, optional (can be nil).
	// Recommended: hash-length random value.
	salt := make([]byte, hash().Size())

	kdf := hkdf.New(hash, secret, salt, info)

	secrets := make([]byte, outputLength)
	length, err := io.ReadFull(kdf, secrets)
	if err != nil {
		return nil, err
	}
	if length != outputLength {
		return nil, err
	}

	return secrets, nil
}

func Sign(privateKey *[32]byte, message []byte, random [64]byte) *[64]byte {

	// Calculate Ed25519 public key from Curve25519 private key
	var A edwards25519.ExtendedGroupElement
	var publicKey [32]byte
	edwards25519.GeScalarMultBase(&A, privateKey)
	A.ToBytes(&publicKey)

	// Calculate r
	diversifier := [32]byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	var r [64]byte
	hash := sha512.New()
	hash.Write(diversifier[:])
	hash.Write(privateKey[:])
	hash.Write(message)
	hash.Write(random[:])
	hash.Sum(r[:0])

	// Calculate R
	var rReduced [32]byte
	edwards25519.ScReduce(&rReduced, &r)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &rReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	// Calculate S = r + SHA2-512(R || A_ed || msg) * a  (mod L)
	var hramDigest [64]byte
	hash.Reset()
	hash.Write(encodedR[:])
	hash.Write(publicKey[:])
	hash.Write(message)
	hash.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, privateKey, &rReduced)

	signature := new([64]byte)
	copy(signature[:], encodedR[:])
	copy(signature[32:], s[:])
	signature[63] |= publicKey[31] & 0x80

	return signature
}

func Verify(publicKey [32]byte, message []byte, signature *[64]byte) bool {

	publicKey[31] &= 0x7F

	/* Convert the Curve25519 public key into an Ed25519 public key.  In
	particular, convert Curve25519's "montgomery" x-coordinate into an
	Ed25519 "edwards" y-coordinate:

	ed_y = (mont_x - 1) / (mont_x + 1)

	NOTE: mont_x=-1 is converted to ed_y=0 since fe_invert is mod-exp

	Then move the sign bit into the pubkey from the signature.
	*/

	var edY, one, montX, montXMinusOne, montXPlusOne edwards25519.FieldElement
	edwards25519.FeFromBytes(&montX, &publicKey)
	edwards25519.FeOne(&one)
	edwards25519.FeSub(&montXMinusOne, &montX, &one)
	edwards25519.FeAdd(&montXPlusOne, &montX, &one)
	edwards25519.FeInvert(&montXPlusOne, &montXPlusOne)
	edwards25519.FeMul(&edY, &montXMinusOne, &montXPlusOne)

	var A_ed [32]byte
	edwards25519.FeToBytes(&A_ed, &edY)

	A_ed[31] |= signature[63] & 0x80
	signature[63] &= 0x7F

	return ed25519.Verify(&A_ed, message, signature)
}

func LoadMaterialKey(chain []byte) ([]byte, error) {

	mac := hmac.New(sha256.New, chain[:])
	mac.Write([]byte{0x01})
	mk := mac.Sum(nil)

	messageKey, err := HKDF(mk, info, 80)
	if err != nil {
		return nil, err
	}
	return messageKey, nil
}
