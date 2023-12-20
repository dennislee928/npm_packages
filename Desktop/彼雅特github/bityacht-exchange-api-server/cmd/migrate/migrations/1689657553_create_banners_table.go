package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689657553{})
}

// Migration1689657553 is a kind of IMigration, You can define schema here.
type Migration1689657553 struct {
	ID        int64
	Priority  int64     `gorm:"not null;default:0"`
	WebImage  string    `gorm:"not null;default:''"`
	AppImage  string    `gorm:"not null;default:''"`
	Title     string    `gorm:"not null;default:''"`
	SubTitle  string    `gorm:"not null;default:''"`
	Button    string    `gorm:"not null;default:''"`
	ButtonUrl string    `gorm:"not null;default:''"`
	Status    int32     `gorm:"not null;default:0"`
	StartAt   time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	EndAt     time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName for gorm
func (*Migration1689657553) TableName() string {
	return "banners"
}

// Version for IMigration
func (*Migration1689657553) Version() int64 {
	return 1689657553
}

// Up for IMigration
func (m *Migration1689657553) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689657553) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
