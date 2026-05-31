package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func HexToBytes(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

func BytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64ToBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
