package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func HashText(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Base64(v []byte) string {
	return base64.URLEncoding.EncodeToString(v)
}

func DecodeBase64(str string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(str)
}
