package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/pkg/mmdb"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederUsersLoginLogs{})
}

// SeederUsersLoginLogs is a kind of ISeeder, You can set it as same as model.
type SeederUsersLoginLogs migrations.Migration1689733636

// SeederName for ISeeder
func (*SeederUsersLoginLogs) SeederName() string {
	return "UsersLoginLogs"
}

// TableName for gorm
func (*SeederUsersLoginLogs) TableName() string {
	return "users_login_logs"
}

// Default for ISeeder
func (*SeederUsersLoginLogs) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (*SeederUsersLoginLogs) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	startAt := time.Now().AddDate(0, 0, -7)
	for _, userID := range userIDs {
		records := make([]SeederUsersLoginLogs, 0, rand.Intn(75)+25)
		createdAt := startAt
		halfInterval := 84 * time.Hour / time.Duration(cap(records)) // 7*24/2

		for i := 0; i < cap(records); i++ {
			createdAt = createdAt.Add(halfInterval + rand.Intn(halfInterval))

			newRecord := SeederUsersLoginLogs{
				UsersID:          userID,
				UserAgent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (Fake " + rand.LetterAndNumberString(5) + ")",
				Browser:          "Fake Browser " + rand.LetterAndNumberString(5),
				Device:           "Fake Device " + rand.LetterAndNumberString(5),
				IP:               fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255)+1, rand.Intn(255)+1, rand.Intn(255)+1, rand.Intn(255)+1),
				IPRelteadHeaders: "{}",
				CreatedAt:        createdAt,
			}

			if city, err := mmdb.LookupCity(newRecord.IP); err == nil {
				newRecord.Location = city.String()
			}

			records = append(records, newRecord)
		}

		if err := db.Create(&records).Error; err != nil {
			return err
		}
	}

	return nil
}
