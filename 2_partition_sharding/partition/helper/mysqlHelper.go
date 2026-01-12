package helper

import "fmt"

// Helper to generate table names like "users_2025"
func GetTableName(year string) string {
	return fmt.Sprintf("users_%s", year)
}
