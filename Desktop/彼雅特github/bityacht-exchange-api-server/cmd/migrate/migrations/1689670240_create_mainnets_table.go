package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689670240{})
}

// Ref: https://binance-docs.github.io/apidocs/spot/en/#all-coins-39-information-user_data

// Migration1689670240 is a kind of IMigration, You can define schema here.
type Migration1689670240 struct {
	CurrenciesSymbol         string          `gorm:"primaryKey"` // binance: coin
	Mainnet                  string          `gorm:"primaryKey"` // binance: network
	Name                     string          `gorm:"not null;default:''"`
	AddressRegex             string          `gorm:"not null;default:''"`
	WithdrawDecimalPrecision int32           `gorm:"not null;default:0"` // parse from binance: withdrawIntegerMultiple
	WithdrawFee              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	WithdrawMin              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	WithdrawMax              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Status                   int32           `gorm:"not null;default:0"`

	CurrenciesSymbolFK Migration1689663148 `gorm:"foreignKey:CurrenciesSymbol;references:Symbol;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689670240) TableName() string {
	return "mainnets"
}

// Version for IMigration
func (*Migration1689670240) Version() int64 {
	return 1689670240
}

// Up for IMigration
func (m *Migration1689670240) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689670240) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
