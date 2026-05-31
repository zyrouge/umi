package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func UpdateUserEmail(id string, email string, userKey []byte) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	encryptedEmail, err := EncryptUserEmail(email, userKey)
	if err != nil {
		return err
	}
	_, err = connection.Exec(`UPDATE umi_user SET email = ?, updated_at = ? WHERE id = ?`, encryptedEmail, time.Now().Unix(), id)
	return err
}

func UpdateUserUsername(id, username string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`UPDATE umi_user SET username = ?, updated_at = ? WHERE id = ?`, username, time.Now().Unix(), id)
	return err
}

func UpdateUserPasswordHash(id, hash string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`UPDATE umi_user SET password_hash = ?, updated_at = ? WHERE id = ?`, hash, time.Now().Unix(), id)
	return err
}
