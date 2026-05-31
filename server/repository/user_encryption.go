package repository

import (
	"fmt"

	"zyrouge.me/umi/utils"
)

func EncryptUserEmail(email string, key []byte) (string, error) {
	encrypted, err := utils.EncryptAESGCM(key, email)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt email: %w", err)
	}
	return encrypted, nil
}

func DecryptUserEmail(email string, key []byte) (string, error) {
	decrypted, err := utils.DecryptAESGCM(key, email)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt email: %w", err)
	}
	return decrypted, nil
}
