package levellimits

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/shopspring/decimal"
)

// TableName of level_limits table
const TableName = "level_limits"

// Model of level_limits table
type Model struct {
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

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetAllLevelLimit() ([]LevelLimit, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Table(TableName).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{}
	}

	output := make([]LevelLimit, len(records))
	for i, v := range records {
		output[i] = LevelLimit{
			Type:  v.Type,
			Level: v.Level,
			TWD: LimitByCurrency{
				Deposit: LimitByAction{
					PerTransaction: Limit{Max: v.MaxDepositTwdPerTransaction},
					PerDay:         Limit{Max: v.MaxDepositTwdPerDay},
					PerMonth:       Limit{Max: v.MaxDepositTwdPerMonth},
				},
				Withdraw: LimitByAction{
					PerTransaction: Limit{Min: v.MinWithdrawTwdPerTransaction, Max: v.MaxWithdrawTwdPerTransaction},
					PerDay:         Limit{Max: v.MaxWithdrawTwdPerDay},
					PerMonth:       Limit{Max: v.MaxWithdrawTwdPerMonth},
				},
			},
			USDT: LimitByCurrency{
				Deposit: LimitByAction{
					PerTransaction: Limit{Max: v.MaxDepositUsdtPerTransaction},
					PerDay:         Limit{Max: v.MaxDepositUsdtPerDay},
					PerMonth:       Limit{Max: v.MaxDepositUsdtPerMonth},
				},
				Withdraw: LimitByAction{
					PerTransaction: Limit{Min: v.MinWithdrawUsdtPerTransaction, Max: v.MaxWithdrawUsdtPerTransaction},
					PerDay:         Limit{Max: v.MaxWithdrawUsdtPerDay},
					PerMonth:       Limit{Max: v.MaxWithdrawUsdtPerMonth},
				},
			},
		}
	}

	return output, nil
}
