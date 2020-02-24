package utils

import (
	"crypto/rand"
	"io"
)

func Random() [64]byte {
	var random [64]byte
	r := rand.Reader
	io.ReadFull(r, random[:])
	return random
}

func Slice2Array32(src []byte) [32]byte {
	var result [32]byte
	copy(result[:], src[:32])
	return result
}

func Slice2Array64(src []byte) [64]byte {
	var result [64]byte
	copy(result[:], src[:64])
	return result
}
