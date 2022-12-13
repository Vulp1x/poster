package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Переводит string в sql.NullString.
func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: len(s) > 0}
}

// SQLInterval переводит time.Duration в SQL интервал.
func SQLInterval(duration time.Duration) string {
	return fmt.Sprintf("%v milliseconds", duration.Milliseconds())
}
