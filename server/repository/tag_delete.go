package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func DeleteTagById(id string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`DELETE FROM umi_tag WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
