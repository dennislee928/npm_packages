package users

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/bankaccounts"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"bityacht-exchange-api-server/internal/pkg/rand"
	dbsql "database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// TableName of users table
const TableName = "users"

// Model of users table
type Model struct {
	ID                                   int64            `gorm:"autoIncrement:false"`
	Account                              string           `gorm:"not null;uniqueIndex"`
	NationalID                           dbsql.NullString `gorm:"uniqueIndex"`
	PassportNumber                       dbsql.NullString
	CountriesCode                        dbsql.NullString
	DualNationalityCode                  dbsql.NullString
	IndustrialClassificationsID          dbsql.NullInt64
	InviterID                            dbsql.NullInt64
	BankAccountsID                       dbsql.NullInt64
	IDVerificationsID                    dbsql.NullInt64
	DueDiligencesID                      dbsql.NullInt64
	Password                             string    `gorm:"not null"`
	Type                                 Type      `gorm:"not null"`
	FirstName                            string    `gorm:"not null;default:''"`
	LastName                             string    `gorm:"not null;default:''"`
	Gender                               Gender    `gorm:"not null;default:0"`
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
	Extra                                Extra     `gorm:"type:json;not null;default:'{}'"`
	Status                               Status    `gorm:"not null;default:0"`

	NameCheck               usersmodifylogs.RLStatus `gorm:"not null;default:0"`
	NameCheckPdfName        string                   `gorm:"not null;default:''"`
	NameCheckPdfData        string                   `gorm:"type:longblob;default:''"`
	InternalRisksTotal      int64                    `gorm:"not null;default:0"`
	ComplianceReview        usersmodifylogs.RLStatus `gorm:"not null;default:0"`
	ComplianceReviewComment string                   `gorm:"not null;default:''"`
	FinalReview             usersmodifylogs.RLStatus `gorm:"not null;default:0"`
	FinalReviewNotice       string                   `gorm:"not null;default:''"`
	FinalReviewComment      string                   `gorm:"not null;default:''"`
	FinalReviewTime         time.Time                `gorm:"not null;default:'0001-01-01 00:00:00'"`

	CreatedAt time.Time `gorm:"index;not null;default:UTC_TIMESTAMP()"`
	DeletedAt *time.Time
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = rand.Intn[int64](89999999) + 10000000 // 10000000 ~ 99999999
	m.Extra.InviteCode = GetInviteCodeByID(m.ID)

	return nil
}

func (m *Model) Create(tx *gorm.DB, ip string) *errpkg.Error {
	if tx == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil gorm DB")}
	}

	if err := m.Extra.setRegisterIP(ip); err != nil {
		return err
	}

	if m.NameCheck == 0 {
		m.NameCheck = usersmodifylogs.RLStatusPending
	}
	if m.ComplianceReview == 0 {
		m.ComplianceReview = usersmodifylogs.RLStatusPending
	}
	if m.FinalReview == 0 {
		m.FinalReview = usersmodifylogs.RLStatusPending
	}

	// TODO: Optimize Performance
	for retry := 0; retry < 50; retry++ {
		err := tx.Create(m).Error
		if err == nil {
			return nil
		}

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			if strings.Contains(mysqlErr.Message, "key 'PRIMARY'") {
				continue
			}

			return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeAccountDuplicated}
		} else {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
	}

	return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeFailedToGenerateID}
}

func (m Model) AllowWithdraw() *errpkg.Error {
	if time.Since(m.Extra.LastChangePasswordAt) < 24*time.Hour {
		return &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeTemporaryForbidden}
	}

	return nil
}

func SetRegisterIP(id int64, ip string) *errpkg.Error {
	var record Model
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		newExtra := record.Extra
		if newExtra.RegisterIP != "" {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadAction, Err: errors.New("cannot change register ip")}
		}

		if err := newExtra.setRegisterIP(ip); err != nil {
			return err
		}

		if newExtra.RegisterIP == record.Extra.RegisterIP {
			return nil
		}

		if err := tx.Model(&record).Update("extra", newExtra).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

//! Deprecated (Message from TG 2023/11/28)
// type UpdateFromIDVCallbackArgs struct {
// 	FirstName  string
// 	LastName   string
// 	BirthDate  modelpkg.Date
// 	NationalID string
// 	Gender     Gender
// }

//! Deprecated (Message from TG 2023/11/28)
// func UpdateFromIDVCallback(id int64, idvsID int64, args UpdateFromIDVCallbackArgs) *errpkg.Error {
// 	var record Model

// 	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
// 			}
// 			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
// 		}

// 		if record.IDVerificationsID.Int64 != idvsID {
// 			return nil
// 		}

// 		updateMap := make(map[string]any)
// 		if args.FirstName != "" && record.FirstName != args.FirstName {
// 			updateMap["first_name"] = args.FirstName
// 		}
// 		if args.LastName != "" && record.LastName != args.LastName {
// 			updateMap["last_name"] = args.LastName
// 		}
// 		if !args.BirthDate.IsZero() && !record.BirthDate.Equal(args.BirthDate.Time) {
// 			updateMap["birth_date"] = args.BirthDate
// 		}
// 		if args.NationalID != "" && record.NationalID.String != args.NationalID {
// 			updateMap["national_id"] = args.NationalID
// 		}
// 		if args.Gender != GenderUnknown && record.Gender != args.Gender {
// 			updateMap["gender"] = args.Gender
// 		}

// 		if err := tx.Model(&record).Updates(updateMap).Error; err != nil {
// 			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
// 		}

// 		return nil
// 	}); err != nil {
// 		return err.(*errpkg.Error)
// 	}

// 	return nil
// }

func CheckPhoneExist(id int64, phone string) *errpkg.Error {
	var count int64
	if err := sql.DB().Table(TableName).Where("`id` != ? AND `phone` = ?", id, phone).Count(&count).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if count != 0 {
		return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodePhoneNumberDuplicated, Err: err}
	}

	return nil
}

func CheckNationalIDExist(id int64, nationalID string) *errpkg.Error {
	var count int64
	if err := sql.DB().Table(TableName).Where("`id` != ? AND `national_id` = ?", id, nationalID).Count(&count).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if count != 0 {
		return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeNationalIDDuplicated, Err: err}
	}

	return nil
}

func GetUserByID(id int64) (User, *errpkg.Error) {
	var record User

	if err := sql.DB().Where("`id` = ?", id).Scopes(modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetByID(id int64) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`id` = ?", id).Scopes(modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetByAccount(account string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`account` = ?", account).Scopes(modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func Login(account string, password string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`account` = ?", account).Scopes(modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeUnauthorized}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := passwordpkg.Validate(record.Password, password); err != nil {
		return record, err
	}

	return record, nil
}

func CreateNaturalPerson(account string, password string, inviterID int64, ip string) (Model, *errpkg.Error) {
	if err := passwordpkg.StrengthValidate(password); err != nil {
		return Model{}, err
	} else if password, err = passwordpkg.Encrypt(password); err != nil {
		return Model{}, err
	}

	record := Model{
		Account:  account,
		Password: password,
		Type:     TypeNaturalPerson,
		Level:    0,
		Status:   usersmodifylogs.SLStatusUnverified,
	}
	record.Extra.Login2FAType = record.Extra.GetLogin2FAType()
	record.Extra.Withdraw2FAType = record.Extra.GetWithdraw2FAType()

	if inviterID != 0 {
		if _, err := GetByID(inviterID); err != nil {
			return Model{}, err
		}
		record.InviterID = dbsql.NullInt64{Int64: inviterID, Valid: true}
	}

	if err := record.Create(sql.DB(), ip); err != nil {
		return Model{}, err
	}

	return record, nil
}

func EmailVerified(id int64) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if record.Status != usersmodifylogs.SLStatusUnverified {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeEmailAlreadyVerified}
		} else if err = tx.Model(&record).Update("status", int32(usersmodifylogs.SLStatusEnable)).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return record, err.(*errpkg.Error)
	}

	return record, nil
}

func GetList(req GetListRequest, paginator *modelpkg.Paginator, searcher *modelpkg.Searcher) ([]User, *errpkg.Error) {
	records := make([]User, 0)
	query := sql.DB().Table(TableName).
		Scopes(modelpkg.WithNotDeleted(), modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))

	if req.Status != nil {
		query = query.Where("`status` = ?", req.Status)
	}
	if req.Type != nil {
		query = query.Where("`type` = ?", req.Type)
	}
	query = searcher.AddToQuery(query, []string{"id", "account", "CONCAT(`last_name`, `first_name`)", "phone"}).Session(&gorm.Session{})
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}
	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetExport(req ExportRequest) ([]csv.Record, *errpkg.Error) {
	var records []User

	query := sql.DB().Table(TableName).
		Scopes(modelpkg.WithNotDeleted(), modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))

	if len(req.StatusList) > 0 {
		query = query.Where("`status` IN (?)", req.StatusList)
	}

	if err := query.Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for i, v := range records {
		output[i] = v
	}

	return output, nil
}

func UpdateStatus(managerID, id int64, status Status, comment string) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(TableName).Where("`id` = ?", id).Update("status", status).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		record := &usersmodifylogs.Model{
			ManagersID: dbsql.NullInt64{
				Int64: managerID,
				Valid: true,
			},
			UsersID: id,
			Type:    usersmodifylogs.TypeStatusLog,
			Status:  int32(status),
			Comment: comment,
		}
		if err := tx.Create(&record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func UpdateLevel(id int64, level int32) *errpkg.Error {
	var record Model

	if err := sql.DB().Where("`id` = ?", id).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}

		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	switch record.Type {
	case TypeNaturalPerson: // 2 ~ 5
		if level < 2 || level > 5 {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad level")}
		}
	case TypeJuridicalPerson: // 1 ~ 2
		if level < 1 || level > 2 {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad level")}
		}
	default:
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadRecord, Err: errors.New("bad user type")}
	}

	if err := sql.DB().Model(&record).Update("level", level).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return nil
}

type UpdateReviewParams struct {
	ReviewType usersmodifylogs.RLType
	ManagerID  int64
	UserID     int64
	Result     usersmodifylogs.RLStatus
	Notice     string
	Comment    string
}

func UpdateReview(params UpdateReviewParams) (Model, *errpkg.Error) {
	var (
		record    Model
		logRecord = &usersmodifylogs.Model{
			ManagersID: dbsql.NullInt64{Int64: params.ManagerID, Valid: true},
			UsersID:    params.UserID,
			Type:       usersmodifylogs.TypeReviewLog,
			SubType:    int32(params.ReviewType),
			Status:     int32(params.Result),
			Comment:    params.Comment,
		}
	)

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", params.UserID).Scopes(modelpkg.WithNotDeleted(), modelpkg.WithLock(true)).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		updateMap := make(map[string]interface{}, 4)
		switch params.ReviewType {
		case usersmodifylogs.RLTypeNameCheckUploadPDF:
			updateMap["name_check"] = params.Result
			updateMap["name_check_pdf_name"] = params.Notice
			updateMap["name_check_pdf_data"] = params.Comment
		case usersmodifylogs.RLTypeComplianceReview:
			updateMap["compliance_review"] = params.Result
			updateMap["compliance_review_comment"] = params.Comment
		case usersmodifylogs.RLTypeFinalReview:
			updateMap["final_review"] = params.Result
			updateMap["final_review_comment"] = params.Comment
			updateMap["final_review_time"] = time.Now().UTC()

			if record.Level > 0 && params.Result != usersmodifylogs.RLStatusApproved {
				updateMap["level"] = 0
			} else if record.Level == 0 && params.Result == usersmodifylogs.RLStatusApproved {
				updateMap["level"] = 1
			}
		default:
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("bad review type")}
		}

		if err := tx.Model(&record).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err := tx.Create(&logRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		return nil
	}); err != nil {
		return Model{}, err.(*errpkg.Error)
	}
	return record, nil
}

func getNew2FAType(old TwoFAType, modifier *TwoFAType) TwoFAType {
	if modifier == nil {
		return old
	}

	if *modifier > 0 {
		return old | *modifier
	} else {
		return old &^ -*modifier
	}
}

func UpdateExtra(id int64, barcode *string, login2FAType *TwoFAType, withdraw2FAType *TwoFAType, gaSecret *string) *errpkg.Error {
	if barcode == nil && login2FAType == nil && withdraw2FAType == nil && gaSecret == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("args all nil")}
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var record Model

		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeUnauthorized}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		updated := false
		if barcode != nil && record.Extra.MobileBarcode != *barcode {
			updated = true
			record.Extra.MobileBarcode = *barcode
		}
		if newLogin2FAType := getNew2FAType(record.Extra.Login2FAType, login2FAType); record.Extra.Login2FAType != newLogin2FAType {
			updated = true
			record.Extra.Login2FAType = newLogin2FAType
		}
		if newWithdraw2FAType := getNew2FAType(record.Extra.Withdraw2FAType, withdraw2FAType); record.Extra.Withdraw2FAType != newWithdraw2FAType {
			updated = true
			record.Extra.Withdraw2FAType = newWithdraw2FAType
		}
		if gaSecret != nil && record.Extra.GoogleAuthenticatorSecret != *gaSecret {
			updated = true
			record.Extra.GoogleAuthenticatorSecret = *gaSecret
		}

		//! Don't Remove it, for forcing open some 2FA types
		record.Extra.Login2FAType = record.Extra.GetLogin2FAType()
		record.Extra.Withdraw2FAType = record.Extra.GetWithdraw2FAType()

		if !updated {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
		} else if err := tx.Model(&record).Update("extra", record.Extra).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetUserInfoByID(id int64) (UserInfo, *errpkg.Error) {
	var totalInvited, totalSucceed int64
	record := struct {
		ID                   int64
		FinalReview          usersmodifylogs.RLStatus
		IDVerificationsID    int64
		BanksCode            string
		BranchsCode          string
		BankAccount          string
		BankAccountStatus    bankaccounts.Status
		BankAccountCreatedAt modelpkg.DateTime
		Extra                Extra
	}{}

	// TODO: Join bank_accounts
	if err := sql.DB().
		Select(
			"`t`.`id` AS `id`",
			"`t`.`final_review` AS `final_review`",
			"`t`.`extra` AS `extra`",
			"`t`.`id_verifications_id` AS `id_verifications_id`",
			"`b`.`banks_code` AS `banks_code`",
			"`b`.`branchs_code` AS `branchs_code`",
			"`b`.`account` AS `bank_account`",
			"`b`.`status` AS `bank_account_status`",
			"`b`.`created_at` AS `bank_account_created_at`",
		).Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `b` ON `t`.`bank_accounts_id` = `b`.`id`", bankaccounts.TableName)).
		Where("`t`.`id` = ? AND `t`.`deleted_at` IS NULL", id).
		Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserInfo{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return UserInfo{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	totalInvited, totalSucceed, err := GetInviteCount(id)
	if err != nil {
		return UserInfo{}, err
	}

	if accountLength := len(record.BankAccount); accountLength > 0 {
		endIndex := accountLength/2 + accountLength%2
		if endIndex < accountLength-5 {
			endIndex = accountLength - 5
		}

		account := []rune(record.BankAccount)
		for i := 0; i < endIndex; i++ {
			account[i] = '*'
		}

		record.BankAccount = string(account)
	}

	return UserInfo{
		BankAccountStatus:    record.BankAccountStatus,
		BanksCode:            record.BanksCode,
		BranchsCode:          record.BranchsCode,
		BankAccount:          record.BankAccount,
		BankAccountCreatedAt: record.BankAccountCreatedAt,
		InviteCode:           record.Extra.InviteCode,
		MobileBarcode:        record.Extra.MobileBarcode,
		// Login2FAType:         record.Extra.GetLogin2FAType(), //! Deprecated (Meeting at 2023/10/2)
		GoogleAuthenticator:  record.Extra.IsEnableWithdrawGA2FA(),
		TotalInvited:         totalInvited,
		TotalSucceed:         totalSucceed,
		LastChangePasswordAt: modelpkg.NewDateTime(record.Extra.LastChangePasswordAt),
		FinalReview:          record.FinalReview,
		IDVerificationsID:    record.IDVerificationsID,
	}, nil
}

func GetInviteCount(id int64) (totalInvited int64, totalSucceed int64, err *errpkg.Error) {
	if err := sql.DB().Select("COUNT(*), COUNT(IF(`level` >= 2, TRUE, NULL))").
		Table(TableName).
		Where("`inviter_id` = ?", id).
		Scopes(modelpkg.WithNotDeleted()).
		Row().Scan(&totalInvited, &totalSucceed); err != nil {
		return 0, 0, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return
}

func GetInviteeList(id int64, paginator *modelpkg.Paginator) ([]Invitee, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	query := sql.DB().Table(TableName).Where("`inviter_id` = ?", id).Scopes(modelpkg.WithNotDeleted()).Session(&gorm.Session{})
	records := make([]Invitee, 0)

	if err := query.Select(
		"`account`, `created_at`, IF(`level` >= 2, ?, ?) AS `status`", int32(InviteStatusFinished), int32(InviteStatusNotFinish),
	).Scopes(modelpkg.WithPaginator(paginator)).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err = query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

// if needValidate, it will validate the strength of new password & old password is correct or not.
func UpdatePassword(needValidate bool, id int64, oldPassword string, newPassword string) *errpkg.Error {
	if needValidate {
		if err := passwordpkg.StrengthValidate(newPassword); err != nil {
			return err
		}
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var record Model

		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeUnauthorized}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if needValidate {
			if err := passwordpkg.Validate(record.Password, oldPassword); err != nil {
				if err.HttpStatus == http.StatusInternalServerError {
					return err
				}

				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad old password")}
			}
		}

		var (
			updateMap = make(map[string]any, 2)
			err       *errpkg.Error
		)

		record.Extra.LastChangePasswordAt = time.Now()
		updateMap["extra"] = record.Extra
		if updateMap["password"], err = passwordpkg.Encrypt(newPassword); err != nil {
			return err
		} else if err := tx.Model(&record).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func UpsertBankAccount(id int64, bankAccountsRecord *bankaccounts.Model) *errpkg.Error {
	if bankAccountsRecord == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("bank accounts record is nil")}
	}
	bankAccountsRecord.UsersID = id
	bankAccountsRecord.Status = usersmodifylogs.BALStatusPending

	var (
		record    Model
		logRecord *usersmodifylogs.Model = &usersmodifylogs.Model{
			UsersID: id,
			Type:    usersmodifylogs.TypeBankAccountLog,
			Status:  int32(bankAccountsRecord.Status),
			Comment: fmt.Sprintf("使用者自行綁定 (%s%s)%s", bankAccountsRecord.BanksCode, bankAccountsRecord.BranchsCode, bankAccountsRecord.Account),
		}
	)
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true)).Take(&record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Create(bankAccountsRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Create(logRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Model(&record).Update("bank_accounts_id", bankAccountsRecord.ID).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func DeleteBankAccount(id int64) *errpkg.Error {
	var (
		record    Model
		logRecord *usersmodifylogs.Model = &usersmodifylogs.Model{
			UsersID: id,
			Type:    usersmodifylogs.TypeBankAccountLog,
			Status:  int32(usersmodifylogs.BALStatusUnknown),
			Comment: "使用者自行刪除",
		}
	)
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true)).Take(&record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if !record.BankAccountsID.Valid {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAction, Err: errors.New("bank accounts id is nil")}
		}

		if err := tx.Create(logRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Model(&record).Update("bank_accounts_id", dbsql.NullInt64{}).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

// Get BankAccount by Users ID
func GetBankAccountByID(id int64) (bankaccounts.BankAccount, *errpkg.Error) {
	record, err := GetByID(id)
	if err != nil {
		return bankaccounts.BankAccount{}, err
	}

	if !record.BankAccountsID.Valid || record.BankAccountsID.Int64 == 0 {
		return bankaccounts.BankAccount{}, nil
	}

	return bankaccounts.GetBankAccountByID(record.BankAccountsID.Int64)
}

// Update BankAccount Status by Users ID
func UpdateBankAccountStatusByID(id int64, managersID int64, bankAccountsID int64, status usersmodifylogs.BALStatus, comment string) *errpkg.Error {
	var (
		record            Model
		bankAccountRecord bankaccounts.Model
		logRecord         *usersmodifylogs.Model = &usersmodifylogs.Model{
			ManagersID: dbsql.NullInt64{Int64: managersID, Valid: true},
			UsersID:    id,
			Type:       usersmodifylogs.TypeBankAccountLog,
			Status:     int32(status),
			Comment:    comment,
		}
	)

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if !record.BankAccountsID.Valid || record.BankAccountsID.Int64 != bankAccountsID {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAction, Err: errors.New("bad bank accounts id")}
		}

		if err := tx.Where("`id` = ?", record.BankAccountsID.Int64).Scopes(modelpkg.WithLock(true)).Take(&bankAccountRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if bankAccountRecord.Status == status {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
		}

		var auditTime *dbsql.NullTime
		if status == usersmodifylogs.BALStatusAccepted || status == usersmodifylogs.BALStatusRejected {
			auditTime = &dbsql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}

		updateMap := map[string]any{
			"status":     int32(status),
			"audit_time": auditTime,
		}

		if err := tx.Model(&bankAccountRecord).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Create(logRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
