package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689670823{})
}

// Ref: https://binance-docs.github.io/apidocs/spot/cn/#e7746f7d60

// Migration1689670823 is a kind of IMigration, You can define schema here.
type Migration1689670823 struct {
	BaseCurrenciesSymbol            string          `gorm:"primaryKey"`
	QuoteCurrenciesSymbol           string          `gorm:"primaryKey"`
	SpreadsOfBuy                    decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	SpreadsOfSell                   decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	HandlingChargeRate              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BaseCurrenciesDecimalPrecision  int32           `gorm:"not null;default:0"`
	QuoteCurrenciesDecimalPrecision int32           `gorm:"not null;default:0"`
	Status                          int32           `gorm:"not null;default:0"`

	BaseCurrenciesSymbolFK  Migration1689663148 `gorm:"foreignKey:BaseCurrenciesSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	QuoteCurrenciesSymbolFK Migration1689663148 `gorm:"foreignKey:QuoteCurrenciesSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689670823) TableName() string {
	return "transaction_pairs"
}

// Version for IMigration
func (*Migration1689670823) Version() int64 {
	return 1689670823
}

// Up for IMigration
func (m *Migration1689670823) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689670823) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
