package idv

import (
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/countries"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/pkg/datauri"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"net/http"
)

type IssuePhoneVerificationCodeRequest struct {
	CreateIDVerificationRequest

	Phone modelpkg.TWCellPhone `json:"phone" binding:"required"` // 行動電話, 格式： ^09[0-9]{8}$
}

type CheckPhoneRequest struct {
	Phone modelpkg.TWCellPhone `json:"phone" binding:"required"` // 行動電話, 格式： ^09[0-9]{8}$
}

type VerifyPhoneRequest struct {
	Phone            modelpkg.TWCellPhone `json:"phone" binding:"required"`            // 行動電話, 格式： ^09[0-9]{8}$
	VerificationCode string               `json:"verificationCode" binding:"required"` // 行動電話驗證碼
}

type VerifyPhoneResponse struct {
	PhoneToken string `json:"phoneToken" binding:"required"` // 行動電話 - 驗證 Token
}

type UpdateIDVImageRequest struct {
	// 居留證正面 (Data URI Scheme with Base64)
	IDImage string `json:"idImage"`

	// 居留證背面 (Data URI Scheme with Base64)
	IDBackImage string `json:"idBackImage"`

	// 護照證件照 (Data URI Scheme with Base64)
	PassportImage string `json:"passportImage"`

	// 手持自拍照 (Data URI Scheme with Base64)
	IDAndFaceImage string `json:"idAndFaceImage"`
}

func (uidvir UpdateIDVImageRequest) Validate() *errpkg.Error {
	for _, img := range []string{
		uidvir.IDImage,
		uidvir.IDBackImage,
		uidvir.PassportImage,
		uidvir.IDAndFaceImage,
	} {
		if err := datauri.ValidateImage(img); err != nil {
			return err
		}
	}

	return nil
}

func (uidvir UpdateIDVImageRequest) ToModel(usersID int64) idverifications.Model {
	return idverifications.Model{
		UsersID:        usersID,
		Type:           idverifications.TypeManual,
		IDImage:        []byte(uidvir.IDImage),
		IDBackImage:    []byte(uidvir.IDBackImage),
		PassportImage:  []byte(uidvir.PassportImage),
		IDAndFaceImage: []byte(uidvir.IDAndFaceImage),
		State:          idverifications.StateAccept,
		AuditStatus:    idverifications.AuditStatusPending,
	}
}

type CreateIDVerificationRequest struct {
	UpdateIDVImageRequest

	// 身分證號碼 or 居留證號碼
	NationalID string `json:"nationalID" binding:"required"`

	// 護照號碼
	PassportNumber string `json:"passportNumber"`

	// 真實名字
	FirstName string `json:"firstName" binding:"required"`

	// 真實姓氏
	LastName string `json:"lastName" binding:"required"`

	// 出生年月日
	BirthDate modelpkg.Date `json:"birthDate" binding:"required" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// // 性別
	// // Gender users.Gender `json:"gender" binding:"required,gte=1,lte=2"`

	// 行動電話 - 驗證 Token (實際創建時為必填)
	PhoneToken string `json:"phoneToken"`

	// 國籍
	CountriesCode string `json:"countriesCode" binding:"required"`

	// 雙重國籍
	DualNationalityCode string `json:"dualNationalityCode"`

	// 行業別
	IndustrialClassificationsID int64 `json:"industrialClassificationsID" binding:"required"`

	// 居住地址
	Address string `json:"address" binding:"required"`

	// 年收入
	AnnualIncome string `json:"annualIncome" binding:"required"`

	// 資金來源
	FundsSources string `json:"fundsSources" binding:"required"`

	// 使用目的
	PurposeOfUse string `json:"purposeOfUse" binding:"required"`

	// 投資經驗
	InvestmentExperience string `json:"investmentExperience" binding:"required"`
}

func (cidvr CreateIDVerificationRequest) Validate() (countries.Country, *errpkg.Error) {
	errData := make(map[string]errpkg.Code)

	isForeigner := !countriesCodeIsTW(cidvr.CountriesCode)
	if err := validateNationalID(cidvr.NationalID, isForeigner); err != nil {
		errData["nationalID"] = err.Code
	}

	country, err := sqlcache.GetCountry(cidvr.CountriesCode)
	if err != nil {
		errData["countriesCode"] = err.Code
	}

	if cidvr.DualNationalityCode != "" {
		if countriesCodeIsTW(cidvr.DualNationalityCode) {
			errData["dualNationalityCode"] = errpkg.CodeBadBody
		} else if _, err = sqlcache.GetCountry(cidvr.DualNationalityCode); err != nil {
			errData["dualNationalityCode"] = err.Code
		}
	}

	if isForeigner {
		if err = datauri.ValidateImage(cidvr.IDImage); err != nil {
			errData["idImage"] = err.Code
		}
		if err = datauri.ValidateImage(cidvr.IDBackImage); err != nil {
			errData["idBackImage"] = err.Code
		}
		if err = datauri.ValidateImage(cidvr.PassportImage); err != nil {
			errData["passportImage"] = err.Code
		}
		if err = datauri.ValidateImage(cidvr.IDAndFaceImage); err != nil {
			errData["idAndFaceImage"] = err.Code
		}
		if len(cidvr.PassportNumber) == 0 {
			errData["passportNumber"] = errpkg.CodeBadBody
		}
	}

	if _, err := sqlcache.GetIC(cidvr.IndustrialClassificationsID); err != nil {
		errData["industrialClassificationsID"] = err.Code
	}

	if len(errData) != 0 {
		return country, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Data: errData}
	}

	return country, nil
}

func (cidvr CreateIDVerificationRequest) ToModel(usersID int64, isForeigner bool) idverifications.Model {
	if isForeigner {
		return cidvr.UpdateIDVImageRequest.ToModel(usersID)
	}

	return idverifications.Model{
		UsersID: usersID,
		Type:    idverifications.TypeKryptoGO,
	}
}

func (cidvr CreateIDVerificationRequest) ToUpdateMap(isForeigner bool, phone string) map[string]any {
	updateMap := map[string]any{
		"first_name":                    cidvr.FirstName,
		"last_name":                     cidvr.LastName,
		"birth_date":                    cidvr.BirthDate,
		"phone":                         phone,
		"countries_code":                cidvr.CountriesCode,
		"industrial_classifications_id": cidvr.IndustrialClassificationsID,
		"address":                       cidvr.Address,
		"annual_income":                 cidvr.AnnualIncome,
		"funds_sources":                 cidvr.FundsSources,
		"purpose_of_use":                cidvr.PurposeOfUse,
		"investment_experience":         cidvr.InvestmentExperience,
		// // "gender":                        cidvr.Gender,
	}

	if isForeigner {
		updateMap["passport_number"] = cidvr.PassportNumber
	}
	if cidvr.DualNationalityCode != "" {
		updateMap["dual_nationality_code"] = cidvr.DualNationalityCode
	}

	return updateMap
}

type CreateIDVerificationResponse struct {
	IDVerificationURL string `json:"idVerificationUrl"`
}
