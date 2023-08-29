package hash256

import (
	"crypto/sha256"
)

func New(b []byte) []byte {
	h1 := sha256.New()
	h1.Write(b)
	intermediateHash := h1.Sum(nil)
	h2 := sha256.New()
	h2.Write(intermediateHash)
	return h2.Sum(nil)
}
