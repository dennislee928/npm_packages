package modelpkg

import (
	"database/sql"
)

func GetSqlNullString(column sql.NullString, defaultVal string) string {
	if column.Valid {
		return column.String
	}

	return defaultVal
}

func GetSqlNullInt64(column sql.NullInt64, defaultVal int64) int64 {
	if column.Valid {
		return column.Int64
	}

	return defaultVal
}
