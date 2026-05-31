package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetUserById(id string, userKey []byte) (*UmiUser, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, username, email, password_hash, created_at, updated_at FROM umi_user WHERE id = ?`, id)
	user, err := SqlScanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Email, err = DecryptUserEmail(user.Email, userKey)
	return user, err
}

func GetUserByUsername(username string, userKey []byte) (*UmiUser, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, username, email, password_hash, created_at, updated_at FROM umi_user WHERE username = ?`, username)
	user, err := SqlScanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Email, err = DecryptUserEmail(user.Email, userKey)
	return user, err
}

func CountUsers() (int, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}
	var count int
	err = connection.QueryRow(`SELECT COUNT(*) FROM umi_user`).Scan(&count)
	return count, err
}
