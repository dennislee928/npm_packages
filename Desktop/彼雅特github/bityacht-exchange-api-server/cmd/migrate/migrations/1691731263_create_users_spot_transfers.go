package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1691731263{})
}

// Migration1691731263 is a kind of IMigration, You can define schema here.
type Migration1691731263 struct {
	TransfersID      string          `gorm:"primaryKey"`
	Type             int32           `gorm:"not null"`
	UsersID          int64           `gorm:"not null;index:idx_user_time"`
	CurrenciesSymbol string          `gorm:"not null"`
	Mainnet          string          `gorm:"not null"`
	FromAddress      string          `gorm:"not null;default:''"`
	ToAddress        string          `gorm:"not null;default:''"`
	Status           int32           `gorm:"not null;default:0"`
	Action           int32           `gorm:"not null;default:0"`
	TxID             string          `gorm:"not null;default:''"`
	Serial           sql.NullInt64   `gorm:"uniqueIndex"`
	Amount           decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Valuation        decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"` // in USDT
	HandlingCharge   decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Extra            string          `gorm:"type:json;not null;default:'{}'"`
	FinishedAt       *time.Time
	CreatedAt        time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Users     Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	MainnetFK Migration1689670240 `gorm:"foreignKey:CurrenciesSymbol,Mainnet;references:CurrenciesSymbol,Mainnet;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1691731263) TableName() string {
	return "users_spot_transfers"
}

// Version for IMigration
func (*Migration1691731263) Version() int64 {
	return 1691731263
}

// Up for IMigration
func (m *Migration1691731263) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1691731263) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
