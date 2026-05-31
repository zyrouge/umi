package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func CreateUser(user *UmiUser, userKey []byte) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	encryptedEmail, err := EncryptUserEmail(user.Email, userKey)
	if err != nil {
		return err
	}
	_, err = connection.Exec(
		`INSERT INTO umi_user (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		user.Id, user.Username, encryptedEmail, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
