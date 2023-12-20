package bankbranchs

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
)

// TableName of bank_branchs table
const TableName = "bank_branchs"

// Model of bank_branchs table
type Model struct {
	BanksCode string `gorm:"primaryKey"`
	Code      string `gorm:"primaryKey"`
	Chinese   string `gorm:"not null"`
	English   string `gorm:"not null;default:''"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetBranchListMapByBank() (map[string][]Branch, *errpkg.Error) {
	var records []Branch

	if err := sql.DB().Table(TableName).Order("`banks_code` ASC, `code` ASC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make(map[string][]Branch)
	for _, record := range records {
		output[record.BanksCode] = append(output[record.BanksCode], record)
	}

	return output, nil
}
