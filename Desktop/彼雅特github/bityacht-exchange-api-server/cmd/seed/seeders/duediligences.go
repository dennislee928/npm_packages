package seeders

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/duediligences"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/rand"
)

func init() {
	seed.Register(&SeederDueDiligences{})
}

// SeederUsers is a kind of ISeeder, You can set it as same as model.
type SeederDueDiligences duediligences.Model

// SeederName for ISeeder
func (*SeederDueDiligences) SeederName() string {
	return "DueDiligences"
}

// TableName for gorm
func (*SeederDueDiligences) TableName() string {
	return duediligences.TableName
}

// Default for ISeeder
func (*SeederDueDiligences) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (s *SeederDueDiligences) Fake(db *gorm.DB) error {
	var userIDs []int64

	if seed.UsersID > 0 {
		userIDs = []int64{seed.UsersID}
	} else {
		var err error
		if userIDs, err = getFakeUserIDs(db); err != nil {
			return err
		}
	}

	for _, userID := range userIDs {
		if err := db.Transaction(func(tx *gorm.DB) error {
			var userRecord users.Model

			if err := tx.Where("`id` = ?", userID).Take(&userRecord).Error; err != nil {
				return err
			} else if userRecord.Type != users.TypeNaturalPerson {
				return nil
			}

			record := SeederDueDiligences{
				UsersID:         userID,
				Type:            duediligences.TypeCreateByIDV,
				TaskID:          "fake_" + rand.NumberString(6),
				SanctionMatched: rand.Intn[duediligences.Bool](2) + 1,
				PotentialRisk:   rand.Intn[int64](100),
				AuditAccepted:   rand.Intn[duediligences.Bool](2) + 1,
				Comment:         "fake comment: " + rand.LetterAndNumberString(10),
				AuditTime:       sql.NullTime{Time: time.Now().AddDate(0, 0, -rand.Intn(365)).Add(-time.Duration(rand.Intn(24*60)) * time.Minute), Valid: true},
			}
			if userRecord.CountriesCode.String == countryCodeTaiwan {
				record.Type = duediligences.TypeManualSet
			}

			if err := tx.Create(&record).Error; err != nil {
				return err
			} else if err = tx.Model(&userRecord).Update("due_diligences_id", record.ID).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
