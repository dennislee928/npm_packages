package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689651014{})
}

// Migration1689651014 is a kind of IMigration, You can define schema here.
type Migration1689651014 struct {
	Type  int32 `gorm:"primaryKey;autoIncrement:false"`
	Level int32 `gorm:"primaryKey;autoIncrement:false"`

	MaxDepositTwdPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"` // if value < 0 mean unlimited.
	MaxDepositTwdPerDay         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxDepositTwdPerMonth       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`

	MinWithdrawTwdPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawTwdPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawTwdPerDay         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawTwdPerMonth       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`

	MaxDepositUsdtPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxDepositUsdtPerDay         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxDepositUsdtPerMonth       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`

	MinWithdrawUsdtPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawUsdtPerTransaction decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawUsdtPerDay         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	MaxWithdrawUsdtPerMonth       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
}

// TableName for gorm
func (*Migration1689651014) TableName() string {
	return "level_limits"
}

// Version for IMigration
func (*Migration1689651014) Version() int64 {
	return 1689651014
}

// Up for IMigration
func (m *Migration1689651014) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689651014) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
