package seeders

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"bityacht-exchange-api-server/internal/pkg/rand"
)

func init() {
	seed.Register(&SeederUsers{})
}

// SeederUsers is a kind of ISeeder, You can set it as same as model.
type SeederUsers users.Model

// SeederName for ISeeder
func (*SeederUsers) SeederName() string {
	return "Users"
}

// TableName for gorm
func (*SeederUsers) TableName() string {
	return users.TableName
}

// Default for ISeeder
func (*SeederUsers) Default(db *gorm.DB) error {
	return seed.ErrNotImplement
}

var fakeUserIDs []int64

func getFakeUserIDs(db *gorm.DB) ([]int64, error) {
	if len(fakeUserIDs) != 0 {
		return fakeUserIDs, nil
	} else if err := db.Select("`id`").
		Table((*SeederUsers)(nil).TableName()).
		Where("`account` LIKE ?", "FakeUser%@test.com.tw").
		Scopes(modelpkg.WithNotDeleted()).
		Pluck("id", &fakeUserIDs).Error; err != nil {
		return nil, err
	}

	return fakeUserIDs, nil
}

// Fake for ISeeder
func (s *SeederUsers) Fake(db *gorm.DB) error {
	password, err := passwordpkg.Encrypt("password")
	if err != nil {
		return err
	}

	_, rawErr := getFakeUserIDs(db)
	if rawErr != nil {
		return rawErr
	}

	var foreignCountryCodes []string
	if rawErr := db.Select("`code`").Table((*SeederCountries)(nil).TableName()).Where("`code` != ?", countryCodeTaiwan).Pluck("code", &foreignCountryCodes).Error; err != nil {
		return rawErr
	} else if len(foreignCountryCodes) == 0 {
		return errors.New("foreign country codes is nil")
	}

	var icsIDs []int64
	if rawErr := db.Select("`id`").Table((*SeederIndustrialClassifications)(nil).TableName()).Pluck("id", &icsIDs).Error; err != nil {
		return rawErr
	} else if len(icsIDs) == 0 {
		return errors.New("industrial classification ids is nil")
	}

	var (
		annualIncomeOptions         = []string{"0~30萬", "30~60萬", "60~100萬", "100~150萬", "150~200萬", "200萬以上"}
		fundsSourcesOptions         = []string{"薪資所得", "投資所得", "繼承", "贈與", "借貸"}
		purposeOfUseOptions         = []string{"短期投資", "長期持有", "資產配置"}
		investmentExperienceOptions = []string{"初次接觸", "1年以下", "1~3年", "3~5年", "5年以上"}
	)

	now := time.Now()
	for i := 0; i < seed.UserCount; i++ {
		usersIndex := i + seed.Offset
		user := users.Model{
			Account:                     fmt.Sprintf("FakeUser%d@test.com.tw", usersIndex),
			CountriesCode:               sql.NullString{String: countryCodeTaiwan, Valid: true},
			Status:                      usersmodifylogs.SLStatusEnable,
			IndustrialClassificationsID: sql.NullInt64{Int64: icsIDs[rand.Intn(len(icsIDs))], Valid: true},
			Password:                    password,
			Type:                        rand.Intn[users.Type](2) + 1,
			BirthDate:                   now.AddDate(-18, -rand.Intn(47), -rand.Intn(365)),
			Phone:                       "+8869" + rand.NumberString(8),
			Address:                     "fake address " + rand.LetterAndNumberString(20),
			Level:                       1,
		}
		if len(fakeUserIDs) != 0 {
			user.InviterID = sql.NullInt64{Int64: fakeUserIDs[rand.Intn(len(fakeUserIDs))], Valid: true}
		}
		if rand.Float64() < 0.5 {
			user.Extra.MobileBarcode = "/" + strings.ToUpper(rand.LetterAndNumberString(7))
		}

		switch user.Type {
		case users.TypeNaturalPerson:
			user.FirstName = fmt.Sprintf("User%d", usersIndex)
			user.LastName = "Fake"
			user.NationalID = sql.NullString{String: "FN" + rand.NumberString(8), Valid: true}
			user.AnnualIncome = annualIncomeOptions[rand.Intn(len(annualIncomeOptions))]
			user.FundsSources = fundsSourcesOptions[rand.Intn(len(fundsSourcesOptions))]
			user.PurposeOfUse = purposeOfUseOptions[rand.Intn(len(purposeOfUseOptions))]
			user.InvestmentExperience = investmentExperienceOptions[rand.Intn(len(investmentExperienceOptions))]

			if randNum := rand.Float64(); randNum > 0.5 {
				user.CountriesCode.String = foreignCountryCodes[rand.Intn(len(foreignCountryCodes))]
				user.PassportNumber = sql.NullString{String: "FAKE" + rand.NumberString(5), Valid: true}
			} else if randNum < 0.2 {
				user.DualNationalityCode = sql.NullString{String: foreignCountryCodes[rand.Intn(len(foreignCountryCodes))], Valid: true}
			}
		case users.TypeJuridicalPerson:
			user.FirstName = fmt.Sprintf("Company%d", usersIndex)
			user.NationalID = sql.NullString{String: "FJ" + rand.NumberString(6), Valid: true}
			user.JuridicalPersonNature = "fake nature: " + rand.LetterAndNumberString(10)
			user.FundsSources = "fake funds sources: " + rand.LetterAndNumberString(10)
			user.JuridicalPersonCryptocurrencySources = "fake cryptocurrency sources: " + rand.LetterAndNumberString(10)
			user.AuthorizedPersonName = "Fake Person " + rand.LetterAndNumberString(10)
			user.AuthorizedPersonNationalID = "FPN" + rand.NumberString(7)
			user.AuthorizedPersonPhone = "+8869" + rand.NumberString(8)
			user.Comment = "fake comment: " + rand.LetterAndNumberString(10)
		}

		if err := user.Create(db, ""); err != nil {
			if err.Code == errpkg.CodeAccountDuplicated {
				return errors.New("account duplicated")
			}
			return err
		}

		fakeUserIDs = append(fakeUserIDs, user.ID)
	}

	return nil
}
