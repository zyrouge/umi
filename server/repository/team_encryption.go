package repository

import (
	"encoding/base64"
	"fmt"

	"zyrouge.me/umi/utils"
)

func EncryptTeamKey(teamKey []byte, masterKey []byte) (string, error) {
	return utils.EncryptAESGCM(masterKey, base64.StdEncoding.EncodeToString(teamKey))
}

func DecryptTeamKey(encryptedKey string, masterKey []byte) ([]byte, error) {
	if encryptedKey == "" {
		return nil, fmt.Errorf("team has no encryption key")
	}
	keyB64, err := utils.DecryptAESGCM(masterKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt team key: %w", err)
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode team key: %w", err)
	}
	return key, nil
}

func DecryptTeamEncryptionKey(team *UmiTeam, masterKey []byte) ([]byte, error) {
	return DecryptTeamKey(team.EncryptionKey, masterKey)
}
