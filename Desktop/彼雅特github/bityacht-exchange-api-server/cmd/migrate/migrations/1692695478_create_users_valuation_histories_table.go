package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1692695478{})
}

// Migration1692695478 is a kind of IMigration, You can define schema here.
type Migration1692695478 struct {
	ID        int64
	UsersID   int64           `gorm:"not null;uniqueIndex:idx_user_date"`
	Date      time.Time       `gorm:"type:date;not null;uniqueIndex:idx_user_date"`
	Valuation decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	CreatedAt time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1692695478) TableName() string {
	return "users_valuation_histories"
}

// Version for IMigration
func (*Migration1692695478) Version() int64 {
	return 1692695478
}

// Up for IMigration
func (m *Migration1692695478) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1692695478) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
