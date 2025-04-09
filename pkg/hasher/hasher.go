package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func CompareStringHash(str, hash string) bool {
	return GenerateHashString(str) == hash
}
