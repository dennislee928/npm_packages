package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689661893{})
}

// Migration1689661893 is a kind of IMigration, You can define schema here.
type Migration1689661893 struct {
	ID                                   int64          `gorm:"autoIncrement:false"`
	Account                              string         `gorm:"not null;uniqueIndex"`
	NationalID                           sql.NullString `gorm:"uniqueIndex"`
	PassportNumber                       sql.NullString
	CountriesCode                        sql.NullString
	DualNationalityCode                  sql.NullString
	IndustrialClassificationsID          sql.NullInt64
	InviterID                            sql.NullInt64
	BankAccountsID                       sql.NullInt64
	IDVerificationsID                    sql.NullInt64
	DueDiligencesID                      sql.NullInt64
	Password                             string    `gorm:"not null"`
	Type                                 int32     `gorm:"not null"`
	FirstName                            string    `gorm:"not null;default:''"`
	LastName                             string    `gorm:"not null;default:''"`
	Gender                               int32     `gorm:"not null;default:0"`
	BirthDate                            time.Time `gorm:"type:date;not null;default:'0001-01-01'"`
	Phone                                string    `gorm:"not null;default:''"`
	Address                              string    `gorm:"not null;default:''"`
	AnnualIncome                         string    `gorm:"not null;default:''"`
	FundsSources                         string    `gorm:"not null;default:''"`
	JuridicalPersonNature                string    `gorm:"not null;default:''"`
	JuridicalPersonCryptocurrencySources string    `gorm:"not null;default:''"`
	AuthorizedPersonName                 string    `gorm:"not null;default:''"`
	AuthorizedPersonNationalID           string    `gorm:"not null;default:''"`
	AuthorizedPersonPhone                string    `gorm:"not null;default:''"`
	PurposeOfUse                         string    `gorm:"not null;default:''"`
	InvestmentExperience                 string    `gorm:"not null;default:''"`
	Level                                int32     `gorm:"not null;default:0"`
	Comment                              string    `gorm:"type:text;not null;default:''"`
	Extra                                string    `gorm:"type:json;not null;default:'{}'"`
	Status                               int32     `gorm:"not null;default:0"`

	NameCheck               int32     `gorm:"default:0"`
	NameCheckPdfName        string    `gorm:"default:''"`
	NameCheckPdfData        string    `gorm:"default:''"`
	InternalRisksTotal      int64     `gorm:";default:0"`
	ComplianceReview        int32     `gorm:";default:0"`
	ComplianceReviewComment string    `gorm:";default:''"`
	FinalReview             int32     `gorm:";default:0"`
	FinalReviewNotice       string    `gorm:";default:''"`
	FinalReviewComment      string    `gorm:";default:''"`
	FinalReviewTime         time.Time `gorm:";default:'0001-01-01 00:00:00'"`

	CreatedAt time.Time `gorm:"index;not null;default:UTC_TIMESTAMP()"`
	DeletedAt *time.Time

	CountriesCodeFK           Migration1689659690 `gorm:"foreignKey:CountriesCode;references:Code;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	DualNationalityCodeFK     Migration1689659690 `gorm:"foreignKey:DualNationalityCode;references:Code;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	IndustrialClassifications Migration1689660236 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689661893) TableName() string {
	return "users"
}

// Version for IMigration
func (*Migration1689661893) Version() int64 {
	return 1689661893
}

func (m *Migration1689661893) GetMigrationWithInviterConstraint() interface{} {
	return &struct {
		Migration1689661893

		Inviter Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	}{}
}

// Up for IMigration
func (m *Migration1689661893) Up(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m); err != nil {
		return err
	} else if err = db.Migrator().CreateConstraint(m.GetMigrationWithInviterConstraint(), "Inviter"); err != nil {
		return err
	}

	return nil
}

// Down for IMigration
func (m *Migration1689661893) Down(db *gorm.DB) error {
	if err := db.Migrator().DropConstraint(m.GetMigrationWithInviterConstraint(), "Inviter"); err != nil {
		return err
	} else if err = db.Migrator().DropTable(m); err != nil {
		return err
	}

	return nil
}
