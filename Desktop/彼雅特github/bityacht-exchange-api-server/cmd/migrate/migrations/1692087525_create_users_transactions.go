package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1692087525{})
}

// Migration1692087525 is a kind of IMigration, You can define schema here.
type Migration1692087525 struct {
	TransactionsID        string          `gorm:"primaryKey"`
	UsersID               int64           `gorm:"not null;index:idx_user_time"`
	BaseSymbol            string          `gorm:"not null;default:''"`
	QuoteSymbol           string          `gorm:"not null;default:''"`
	Status                int32           `gorm:"not null;default:0"`
	Side                  int32           `gorm:"not null;default:0"`
	Quantity              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Price                 decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Amount                decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	TwdExchangeRate       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	TwdTotalValue         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	HandlingCharge        decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceID             int64           `gorm:"not null;default:0"`
	BinanceQuantity       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinancePrice          decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceAmount         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceHandlingCharge decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Extra                 string          `gorm:"type:json;not null;default:'{}'"`
	CreatedAt             time.Time       `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Users         Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	BaseSymbolFK  Migration1689663148 `gorm:"foreignKey:BaseSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	QuoteSymbolFK Migration1689663148 `gorm:"foreignKey:QuoteSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1692087525) TableName() string {
	return "users_transactions"
}

// Version for IMigration
func (*Migration1692087525) Version() int64 {
	return 1692087525
}

// Up for IMigration
func (m *Migration1692087525) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1692087525) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
