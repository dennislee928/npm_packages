package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/currencies"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	seed.Register(&SeederUsersWallets{})
}

// SeederUsers is a kind of ISeeder, You can set it as same as model.
type SeederUsersWallets migrations.Migration1689728805

// SeederName for ISeeder
func (*SeederUsersWallets) SeederName() string {
	return "UsersWallets"
}

// TableName for gorm
func (*SeederUsersWallets) TableName() string {
	return (*migrations.Migration1689728805)(nil).TableName()
}

// Default for ISeeder
func (*SeederUsersWallets) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (s *SeederUsersWallets) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	var currencyRecords []currencies.Model
	if err := db.Find(&currencyRecords).Error; err != nil {
		return err
	} else if len(currencyRecords) == 0 {
		return errors.New("currency records is nil")
	}

	for _, userID := range userIDs {
		records := make([]SeederUsersWallets, 0, len(currencyRecords))

		for _, currency := range currencyRecords {
			newRecord := SeederUsersWallets{
				UsersID:          userID,
				CurrenciesSymbol: currency.Symbol,
				Type:             int32(currency.Type),
			}

			switch currency.Type {
			case currencies.TypeFiat:
			default:
				newRecord.FreeAmount = decimal.New(rand.Intn[int64](200000000)+1000, -8)
			}

			records = append(records, newRecord)
		}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error; err != nil {
			return err
		}
	}

	return nil
}
