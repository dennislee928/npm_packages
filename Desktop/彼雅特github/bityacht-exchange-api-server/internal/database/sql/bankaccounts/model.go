package bankaccounts

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	dbsql "database/sql"
	"net/http"
	"time"
)

// TableName of bank_accounts table
const TableName = "bank_accounts"

// Model of bank_accounts table
type Model struct {
	ID          int64
	UsersID     int64     `gorm:"not null"`
	BanksCode   string    `gorm:"not null"`
	BranchsCode string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Account     string    `gorm:"not null"`
	CoverImage  []byte    `gorm:"type:mediumblob"`
	Status      Status    `gorm:"not null;default:1"`
	CreatedAt   time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime   dbsql.NullTime
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetBankAccountByID(id int64) (BankAccount, *errpkg.Error) {
	var record BankAccount

	if err := sql.DB().Table(TableName).Where("`id` = ?", id).Take(&record).Error; err != nil {
		return BankAccount{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}
