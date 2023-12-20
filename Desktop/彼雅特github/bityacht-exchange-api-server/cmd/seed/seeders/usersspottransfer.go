package seeders

import (
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederUsersSpotTransfer{})
}

// SeederUsersSpotTransfer is a kind of ISeeder, You can set it as same as model.
type SeederUsersSpotTransfer usersspottransfers.Model

// SeederName for ISeeder
func (*SeederUsersSpotTransfer) SeederName() string {
	return "UsersSpotTransfer"
}

// TableName for gorm
func (*SeederUsersSpotTransfer) TableName() string {
	return usersspottransfers.TableName
}

func (s *SeederUsersSpotTransfer) BeforeCreate(tx *gorm.DB) error {
	s.TransfersID = modelpkg.GetOrderID(modelpkg.TypeSpot, modelpkg.Action(s.Action))
	s.TransfersID = s.TransfersID[:len(s.TransfersID)-1] + "F"
	return nil
}

// Default for ISeeder
func (*SeederUsersSpotTransfer) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (*SeederUsersSpotTransfer) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	options := []struct {
		Currency wallet.CurrencyType
		Mainnet  wallet.Mainnet
	}{
		{wallet.CurrencyTypeBTC, wallet.MainnetBTC},
		{wallet.CurrencyTypeETH, wallet.MainnetETH},
		{wallet.CurrencyTypeUSDC, wallet.MainnetERC20},
		{wallet.CurrencyTypeUSDT, wallet.MainnetERC20},
		{wallet.CurrencyTypeUSDT, wallet.MainnetTRC20},
	}

	startAt := time.Now().AddDate(0, 0, -7)
	for _, userID := range userIDs {
		recordCount := 25 + rand.Intn(75)
		createdAt := startAt
		halfInterval := 84 * time.Hour / time.Duration(recordCount) // 7*24/2

		for i := 0; i < recordCount; i++ {
			option := options[rand.Intn(len(options))]
			newRecord := SeederUsersSpotTransfer{
				Type:             usersspottransfers.TypeCybavoAPI,
				UsersID:          userID,
				CurrenciesSymbol: option.Currency.String(),
				Mainnet:          option.Mainnet.String(),
				FromAddress:      "fake_from_address_" + rand.LetterAndNumberString(20),
				ToAddress:        "fake_to_address_" + rand.LetterAndNumberString(20),
				Status:           rand.Intn[usersspottransfers.Status](4) + usersspottransfers.StatusProcessing,
				Action:           usersspottransfers.ActionDeposit + rand.Intn[usersspottransfers.Action](2),
				Amount:           decimal.New(rand.Intn[int64](10000)+1, -3),
				CreatedAt:        createdAt,
			}

			if newRecord.Action == usersspottransfers.ActionDeposit {
				newRecord.Type += rand.Intn[usersspottransfers.Type](2)
			}
			if newRecord.Status == usersspottransfers.StatusFinished {
				newRecord.TxID = "fake_tx_id_" + rand.LetterAndNumberString(20)
			}
			if newRecord.Status != usersspottransfers.StatusProcessing {
				finishedAt := createdAt.Add((rand.Intn[time.Duration](24) + 1) * time.Hour)
				newRecord.FinishedAt = &finishedAt
			}

			switch option.Currency {
			case wallet.CurrencyTypeBTC:
				newRecord.Valuation = newRecord.Amount.Mul(decimal.NewFromInt(30000))
			case wallet.CurrencyTypeETH:
				newRecord.Valuation = newRecord.Amount.Mul(decimal.NewFromInt(2000))
			case wallet.CurrencyTypeUSDC, wallet.CurrencyTypeUSDT:
				newRecord.Valuation = newRecord.Amount
			default:
				return errors.New("bad currency")
			}

			if err := seed.RetryCreateWhenDuplicate(db, &newRecord); err != nil {
				return err
			}

			createdAt = createdAt.Add(halfInterval + rand.Intn(halfInterval))
		}
	}

	return nil
}
