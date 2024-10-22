package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func TokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
