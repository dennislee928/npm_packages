package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689902251{})
}

// Migration1689902251 is a kind of IMigration, You can define schema here.
type Migration1689902251 struct {
	ID      int64
	Ptype   string `gorm:"column:ptype"`
	Subject string `gorm:"column:v0"`
	Object  string `gorm:"column:v1"`
	Action  string `gorm:"column:v2"`
	V3      string `gorm:"column:v3"`
	V4      string `gorm:"column:v4"`
	V5      string `gorm:"column:v5"`
}

// TableName for gorm
func (*Migration1689902251) TableName() string {
	return "managers_roles_policies"
}

// Version for IMigration
func (*Migration1689902251) Version() int64 {
	return 1689902251
}

// Up for IMigration
func (m *Migration1689902251) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689902251) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
