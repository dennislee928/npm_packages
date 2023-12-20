package banks

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
)

// TableName of banks table
const TableName = "banks"

// Model of banks table
type Model struct {
	Code    string `gorm:"primaryKey"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null;default:''"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetAllBanks() ([]Bank, *errpkg.Error) {
	var records []Bank

	if err := sql.DB().Table(TableName).Order("`code` ASC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}
