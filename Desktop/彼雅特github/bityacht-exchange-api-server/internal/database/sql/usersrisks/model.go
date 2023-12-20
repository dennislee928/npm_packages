package usersrisks

import (
	"net/http"

	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

// TableName of managers table
const TableName = "users_risks"

// Model of managers table
type Model struct {
	UsersID int64 `gorm:"primaryKey"`
	RisksID int64 `gorm:"primaryKey"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetRisksIDsByUserID(userID int64) ([]int64, *errpkg.Error) {
	var records []int64
	if err := sql.DB().Table(TableName).Where("users_id = ?", userID).Pluck("risks_id", &records).Error; err != nil {
		return records, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}
