package manager

import (
	"crypto/rand"
	"encoding/hex"
)

const ServiceName = "wallet-manager"

func GenerateId(size int) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
