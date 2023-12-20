package seeders

import (
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/receipts"
	"bityacht-exchange-api-server/internal/database/sql/transactionpairs"
	"bityacht-exchange-api-server/internal/database/sql/userstransactions"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/receipt"
	"errors"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederUsersTransactions{})
}

// SeederUsersTransactions is a kind of ISeeder, You can set it as same as model.
type SeederUsersTransactions userstransactions.Model

// SeederName for ISeeder
func (*SeederUsersTransactions) SeederName() string {
	return "UsersTransactions"
}

// TableName for gorm
func (*SeederUsersTransactions) TableName() string {
	return "users_transactions"
}

func (s *SeederUsersTransactions) BeforeCreate(tx *gorm.DB) error {
	s.TransactionsID = modelpkg.GetOrderID(modelpkg.TypeSpot, modelpkg.Action(s.Side))
	s.TransactionsID = s.TransactionsID[:len(s.TransactionsID)-1] + "F"
	return nil
}

// Default for ISeeder
func (*SeederUsersTransactions) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (*SeederUsersTransactions) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	var txPairs []transactionpairs.Model
	if err := db.Find(&txPairs).Error; err != nil {
		return err
	} else if len(txPairs) == 0 {
		return errors.New("tx pairs is nil")
	}

	aHundred := decimal.NewFromInt(100)
	startAt := time.Now().AddDate(0, 0, -7)
	for _, userID := range userIDs {
		recordCount := 25 + rand.Intn(75)
		createdAt := startAt
		halfInterval := 84 * time.Hour / time.Duration(recordCount) // 7*24/2

		for i := 0; i < recordCount; i++ {
			txPair := txPairs[rand.Intn(len(txPairs))]

			newRecord := SeederUsersTransactions{
				UsersID:         userID,
				BaseSymbol:      txPair.BaseCurrenciesSymbol,
				QuoteSymbol:     txPair.QuoteCurrenciesSymbol,
				Status:          rand.Intn[userstransactions.Status](2) + userstransactions.StatusFilled,
				Side:            rand.Intn[userstransactions.Side](2) + userstransactions.SideBuy,
				Price:           decimal.New(rand.Intn[int64](300000000)+100000, -4),   // 10.0000 ~ 30009.9999
				TwdExchangeRate: decimal.New(rand.Intn[int64](10000000000)+300000, -4), // 30.0000 ~ 1000029.9999
				BinanceID:       -rand.Intn[int64](999999999) - 1,
				CreatedAt:       createdAt,
			}

			switch newRecord.Side {
			case userstransactions.SideBuy: // Pay Quote [Amount], Earn Base [Quantity]
				newRecord.Amount = decimal.New(rand.Intn[int64](100000000000)+5000000000, -8).RoundDown(txPair.QuoteCurrenciesDecimalPrecision) // 50.00000000 ~ 1049.99999999
				exchange := newRecord.Amount.Div(newRecord.Price).RoundDown(txPair.BaseCurrenciesDecimalPrecision)
				newRecord.HandlingCharge = exchange.Mul(txPair.HandlingChargeRate).Div(aHundred).RoundUp(txPair.BaseCurrenciesDecimalPrecision)
				newRecord.Quantity = exchange.Sub(newRecord.HandlingCharge)
				newRecord.TwdTotalValue = newRecord.Quantity.Mul(newRecord.TwdExchangeRate)
				newRecord.BinancePrice = newRecord.Price.Div((txPair.SpreadsOfBuy.Add(aHundred)).Div(aHundred))
			case userstransactions.SideSell: // Pay Base [Quantity], Earn Quote [Amount]
				newRecord.Quantity = decimal.New(rand.Intn[int64](200000000)+10000, -8).RoundDown(txPair.BaseCurrenciesDecimalPrecision) // 0.00010000 ~ 2.00009999
				exchange := newRecord.Amount.Mul(newRecord.Price).RoundDown(txPair.QuoteCurrenciesDecimalPrecision)
				newRecord.HandlingCharge = exchange.Mul(txPair.HandlingChargeRate).Div(decimal.NewFromInt(100)).RoundUp(txPair.QuoteCurrenciesDecimalPrecision)
				newRecord.Amount = exchange.Sub(newRecord.HandlingCharge)
				newRecord.TwdTotalValue = newRecord.Amount.Mul(newRecord.TwdExchangeRate)
				newRecord.BinancePrice = newRecord.Price.Div((aHundred.Sub(txPair.SpreadsOfSell)).Div(aHundred))
			default:
				return errors.New("bad side")
			}

			newRecord.BinanceQuantity = newRecord.Quantity.Add(decimal.New(rand.Intn[int64](1000)+100, -txPair.BaseCurrenciesDecimalPrecision))
			newRecord.BinanceAmount = newRecord.Amount.Add(decimal.New(rand.Intn[int64](1000)+100, -txPair.QuoteCurrenciesDecimalPrecision))
			newRecord.BinanceHandlingCharge = decimal.New(rand.Intn[int64](10000), -8)

			if newRecord.Status != userstransactions.StatusFilled {
				newRecord.HandlingCharge = decimal.Zero
			}

			if twdHandlingCharge := newRecord.HandlingCharge.Mul(newRecord.TwdTotalValue).Round(0); twdHandlingCharge.GreaterThan(decimal.Zero) {
				if err := db.Transaction(func(tx *gorm.DB) error {
					if err := seed.RetryCreateWhenDuplicate(tx, &newRecord); err != nil {
						return err
					}

					receiptRecord := receipts.Model{
						ID:              newRecord.TransactionsID,
						Status:          rand.Intn[receipts.Status](2) + receipts.StatusIssued,
						UserID:          newRecord.UsersID,
						InvoiceAmount:   twdHandlingCharge.IntPart(),
						InvoiceID:       "F" + rand.NumberString(7),
						InvoiceIssuedAt: newRecord.CreatedAt.Add(24*time.Hour + rand.Intn[time.Duration](24*time.Hour)),
						CreatedAt:       newRecord.CreatedAt,
					}
					receiptRecord.SalesAmount, receiptRecord.Tax = receipt.CalcSalesAndTaxFromTotal(twdHandlingCharge, decimal.New(5, -2))
					if rand.Float64() < 0.5 {
						receiptRecord.Barcode = "/" + strings.ToUpper(rand.LetterAndNumberString(7))
					}

					if err := tx.Create(&receiptRecord).Error; err != nil {
						return err
					}

					return nil
				}); err != nil {
					return err
				}
			} else if err := seed.RetryCreateWhenDuplicate(db, &newRecord); err != nil {
				return err
			}

			createdAt = createdAt.Add(halfInterval + rand.Intn(halfInterval))
		}
	}

	return nil
}
