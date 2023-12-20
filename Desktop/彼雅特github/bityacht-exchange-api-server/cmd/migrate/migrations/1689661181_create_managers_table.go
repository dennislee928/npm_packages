package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689661181{})
}

// Migration1689661181 is a kind of IMigration, You can define schema here.
type Migration1689661181 struct {
	ID              int64
	Account         string    `gorm:"not null;uniqueIndex"`
	ManagersRolesID int64     `gorm:"not null"`
	Password        string    `gorm:"not null"`
	Name            string    `gorm:"not null;default:''"`
	Extra           string    `gorm:"type:json;not null;default:'{}'"`
	Status          int32     `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	DeletedAt       *time.Time

	ManagersRoles Migration1689660728 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689661181) TableName() string {
	return "managers"
}

// Version for IMigration
func (*Migration1689661181) Version() int64 {
	return 1689661181
}

// Up for IMigration
func (m *Migration1689661181) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689661181) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
