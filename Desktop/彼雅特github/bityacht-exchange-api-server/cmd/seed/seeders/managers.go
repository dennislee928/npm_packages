package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederManagers{})
}

// SeederManagers is a kind of ISeeder, You can set it as same as model.
type SeederManagers migrations.Migration1689661181

// SeederName for ISeeder
func (*SeederManagers) SeederName() string {
	return "Managers"
}

// TableName for gorm
func (*SeederManagers) TableName() string {
	return (*migrations.Migration1689661181)(nil).TableName()
}

// Default for ISeeder
func (*SeederManagers) Default(db *gorm.DB) error {
	records := []SeederManagers{
		{
			Account:         "admin",
			ManagersRolesID: 1,
			Password:        "admin",
			Name:            "Admin",
			Extra:           `{}`,
			Status:          1,
		},
	}

	for i := 0; i < len(records); i++ {
		record := &records[i]
		encryptedPassword, err := passwordpkg.Encrypt(record.Password)
		if err != nil {
			return err
		}

		record.Password = encryptedPassword
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederManagers) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
