package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1693795930{})
}

// Migration1693795930 is a kind of IMigration, You can define schema here.
type Migration1693795930 struct {
	ID             int64
	UsersID        int64 `gorm:"not null"`
	FromUsersID    sql.NullInt64
	TransactionsID sql.NullString
	Action         int32           `gorm:"not null;default:0"`
	Amount         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	CreatedAt      time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`

	Users        Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Transactions Migration1692087525 `gorm:"references:TransactionsID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	FromUsers    Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1693795930) TableName() string {
	return "users_commissions"
}

// Version for IMigration
func (*Migration1693795930) Version() int64 {
	return 1693795930
}

// Up for IMigration
func (m *Migration1693795930) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1693795930) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
