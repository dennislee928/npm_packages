package seeders

import (
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/usersvaluationhistories"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederUsersValuationHistories{})
}

// SeederUsersValuationHistories is a kind of ISeeder, You can set it as same as model.
type SeederUsersValuationHistories struct {
	usersvaluationhistories.Model
}

// SeederName for ISeeder
func (*SeederUsersValuationHistories) SeederName() string {
	return "UsersValuationHistories"
}

// Default for ISeeder
func (*SeederUsersValuationHistories) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (*SeederUsersValuationHistories) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	startAt := time.Now().AddDate(0, 0, -364)
	for _, userID := range userIDs {
		records := make([]SeederUsersValuationHistories, 0, 364)

		date := startAt
		for i := 0; i < cap(records); i++ {
			records = append(records, SeederUsersValuationHistories{
				Model: usersvaluationhistories.Model{
					UsersID:   userID,
					Date:      time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
					Valuation: decimal.New(rand.Intn[int64](999999999999999), -8),
					CreatedAt: date,
				},
			})

			date = date.AddDate(0, 0, 1)
		}

		if err := db.Create(&records).Error; err != nil {
			return err
		}
	}

	return nil
}
