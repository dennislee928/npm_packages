package seeders

import (
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/cmd/seed/seeders/img"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederIDVerifications{})
}

// SeederIDVerifications is a kind of ISeeder, You can set it as same as model.
type SeederIDVerifications idverifications.Model

// SeederName for ISeeder
func (*SeederIDVerifications) SeederName() string {
	return "IDVerifications"
}

// TableName for gorm
func (*SeederIDVerifications) TableName() string {
	return idverifications.TableName
}

// Default for ISeeder
func (*SeederIDVerifications) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

// Fake for ISeeder
func (*SeederIDVerifications) Fake(db *gorm.DB) error {
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
			} else if userRecord.Type != users.TypeNaturalPerson || userRecord.CountriesCode.String == "" {
				return nil
			}

			record := SeederIDVerifications{
				UsersID:   userID,
				TaskID:    "fake" + rand.NumberString(6),
				AuditTime: sql.NullTime{Time: time.Now().AddDate(0, 0, -rand.Intn(365)).Add(-time.Duration(rand.Intn(24*60)) * time.Minute), Valid: true},
			}
			if rand.Float64() < 0.75 {
				record.State = idverifications.StateAccept
				record.AuditStatus = rand.Intn[idverifications.AuditStatus](2) + idverifications.AuditStatusAccepted
			} else {
				record.State = idverifications.StateReject
				record.AuditStatus = idverifications.AuditStatusRejected
			}

			switch userRecord.CountriesCode.String {
			case countryCodeTaiwan:
				record.Type = idverifications.TypeKryptoGO
				record.IDImage = img.Identification()
				record.IDBackImage = img.IdentificationBack()
			default:
				record.Type = idverifications.TypeManual
				record.IDImage = img.ResidentCertificate()
				record.IDBackImage = img.ResidentCertificateBack()
				record.PassportImage = img.Passport()

				if record.AuditStatus == idverifications.AuditStatusAccepted {
					record.ResultImage = img.ResultPass()
				} else {
					record.ResultImage = img.ResultReject()
				}
			}

			if err := tx.Create(&record).Error; err != nil {
				return err
			} else if err = tx.Model(&userRecord).Update("id_verifications_id", record.ID).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
