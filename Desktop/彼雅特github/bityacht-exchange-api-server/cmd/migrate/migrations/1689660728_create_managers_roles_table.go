package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689660728{})
}

// Migration1689660728 is a kind of IMigration, You can define schema here.
type Migration1689660728 struct {
	ID        int64
	Name      string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	DeletedAt *time.Time
}

// TableName for gorm
func (*Migration1689660728) TableName() string {
	return "managers_roles"
}

// Version for IMigration
func (*Migration1689660728) Version() int64 {
	return 1689660728
}

// Up for IMigration
func (m *Migration1689660728) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689660728) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
