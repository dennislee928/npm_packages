package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1698116559{})
}

// Migration1698116559 is a kind of IMigration, You can define schema here.
type Migration1698116559 struct {
	ID                       int64
	UsersID                  int64     `gorm:"not null"`
	Type                     int32     `gorm:"not null"`
	OrderID                  string    `gorm:"not null;default:''"`
	Informations             string    `gorm:"type:json;not null;default:'{}'"`
	InformationReviewFiles   string    `gorm:"type:json;not null;default:'null'"`
	InformationReviewComment string    `gorm:"not null;default:''"`
	RiskReviewResult         int32     `gorm:"not null;default:0"`
	RiskReviewFiles          string    `gorm:"type:json;not null;default:'null'"`
	DedicatedReviewResult    int32     `gorm:"not null;default:0"`
	DedicatedReviewComment   string    `gorm:"not null;default:''"`
	CreatedAt                time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
	ReportMJIBAt             sql.NullTime

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1698116559) TableName() string {
	return "suspicious_transactions"
}

// Version for IMigration
func (*Migration1698116559) Version() int64 {
	return 1698116559
}

// Up for IMigration
func (m *Migration1698116559) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1698116559) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
