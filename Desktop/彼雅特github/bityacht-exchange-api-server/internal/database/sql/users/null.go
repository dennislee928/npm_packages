package users

import modelpkg "bityacht-exchange-api-server/internal/pkg/model"

func (m Model) GetCountriesCode() string {
	return modelpkg.GetSqlNullString(m.CountriesCode, "")
}

func (m Model) GetNationalID() string {
	return modelpkg.GetSqlNullString(m.NationalID, "")
}

func (m Model) GetIndustrialClassificationsID() int64 {
	return modelpkg.GetSqlNullInt64(m.IndustrialClassificationsID, 0)
}

func (m Model) GetBankAccountsID() int64 {
	return modelpkg.GetSqlNullInt64(m.BankAccountsID, 0)
}

func (m Model) GetIDVerificationsID() int64 {
	return modelpkg.GetSqlNullInt64(m.IDVerificationsID, 0)
}

func (m Model) GetDueDiligencesID() int64 {
	return modelpkg.GetSqlNullInt64(m.DueDiligencesID, 0)
}
