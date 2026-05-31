package utils

import (
	"strings"
)

type SqlScannable interface {
	Scan(dest ...any) error
}

func GenerateSqlPlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	placeholders := strings.Repeat("?,", n)
	placeholders = placeholders[:len(placeholders)-1]
	return placeholders
}
