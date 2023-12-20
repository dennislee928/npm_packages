package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689728805{})
}

// Migration1689728805 is a kind of IMigration, You can define schema here.
type Migration1689728805 struct {
	ID               int64
	UsersID          int64           `gorm:"not null;uniqueIndex:idx_users_currency,priority:1"`
	CurrenciesSymbol string          `gorm:"not null;uniqueIndex:idx_users_currency,priority:2"`
	Type             int32           `gorm:"not null;default:0"`
	FreeAmount       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	LockedAmount     decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`

	DeletedAt *time.Time

	Users              Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CurrenciesSymbolFK Migration1689663148 `gorm:"foreignKey:CurrenciesSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689728805) TableName() string {
	return "users_wallets"
}

// Version for IMigration
func (*Migration1689728805) Version() int64 {
	return 1689728805
}

// Up for IMigration
func (m *Migration1689728805) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689728805) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
