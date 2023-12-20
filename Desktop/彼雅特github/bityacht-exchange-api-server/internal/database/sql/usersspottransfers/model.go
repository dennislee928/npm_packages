package usersspottransfers

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	dbsql "database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TableName of spots_transfers table
const TableName = "users_spot_transfers"

// Model of spots_transfers table
type Model struct {
	TransfersID      string          `gorm:"primaryKey"`
	Type             Type            `gorm:"not null"`
	UsersID          int64           `gorm:"not null;index:idx_user_time"`
	CurrenciesSymbol string          `gorm:"not null"`
	Mainnet          string          `gorm:"not null"`
	FromAddress      string          `gorm:"not null;default:''"`
	ToAddress        string          `gorm:"not null;default:''"`
	Status           Status          `gorm:"not null;default:0"`
	Action           Action          `gorm:"not null;default:0"`
	TxID             string          `gorm:"not null;default:''"`
	Serial           dbsql.NullInt64 `gorm:"uniqueIndex"`
	Amount           decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Valuation        decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"` // in USDT
	HandlingCharge   decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Extra            Extra           `gorm:"type:json;not null;default:'{}'"`
	FinishedAt       *time.Time
	CreatedAt        time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.TransfersID = modelpkg.GetOrderID(modelpkg.TypeSpot, modelpkg.Action(m.Action))
	return nil
}

func (m *Model) updateValuationIfZero() {
	if m.Extra.ToUSDTPrice.GreaterThan(decimal.Zero) {
		return
	}

	currencyInfo, err := spottrend.GetCurrencyInfo(m.CurrenciesSymbol)
	if err != nil {
		return
	}

	m.Extra.ToTWDPrice = currencyInfo.OriToTWDPrice
	m.Extra.ToUSDTPrice = currencyInfo.OriToUSDTPrice
	m.Valuation = m.Amount.Mul(currencyInfo.OriToUSDTPrice)
}

func (m *Model) Create(tx *gorm.DB) *errpkg.Error {
	if tx == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil gorm DB")}
	}
	m.updateValuationIfZero()

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

func (m Model) Done() bool {
	return m.Status != StatusProcessing && m.Status != StatusReviewing
}

func Create(record *Model) *errpkg.Error {
	if record == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil record")}
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := record.Create(tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetAccWithdrawValuationByUser(db *gorm.DB, usersID int64) (AccWithdrawValuation, *errpkg.Error) {
	if db == nil {
		db = sql.DB()
	}

	now := time.Now().In(modelpkg.DefaultTimeLoc)
	todayStartAt := modelpkg.NewDate(now).UTC()
	thisMonthStartAt := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, modelpkg.DefaultTimeLoc).UTC()

	var output AccWithdrawValuation

	if err := db.Select("SUM(IF(`created_at` >= ?, `valuation`, 0)) AS `acc_withdraw_in_day`, SUM(`valuation`) AS `acc_withdraw_in_month`", todayStartAt).
		Table(TableName).
		Where("`users_id` = ? AND `action` = ? AND `status` IN ? AND `created_at` >= ?", usersID, ActionWithdraw, []Status{StatusProcessing, StatusReviewing, StatusFinished}, thisMonthStartAt).
		Find(&output).Error; err != nil {
		return output, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return output, nil
}

func GetList(req GetListRequest, paginator *modelpkg.Paginator, searcher modelpkg.Searcher) ([]Transfer, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]Transfer, 0)
	query := sql.DB().Table(TableName).Scopes(modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))
	if req.UsersID > 0 {
		query = query.Where("`users_id` = ?", req.UsersID)
	}
	if req.Status > 0 {
		query = query.Where("`status` = ?", req.Status)
	}
	if req.Coin != "" {
		query = query.Where("`currencies_symbol` = ?", req.Coin)
	}
	query = searcher.AddToQuery(query, []string{"transfers_id", "users_id"}).Session(&gorm.Session{})

	if err := query.Scopes(modelpkg.WithPaginator(paginator)).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetSpotTransfersForUser(usersID int64, req GetSpotTransferForUserRequest, paginator *modelpkg.Paginator) ([]SpotTransferForUser, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]SpotTransferForUser, 0)

	query := sql.DB().Table(TableName).Where("`users_id` = ?", usersID).Scopes(modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))

	if req.CurrenciesSymbol != "" {
		query = query.Where("`currencies_symbol` = ?", req.CurrenciesSymbol)
	}
	if req.Status > 0 {
		query = query.Where("`status` = ?", req.Status)
	}
	if req.Action > 0 {
		query = query.Where("`action` = ?", req.Action)
	}

	query = query.Session(&gorm.Session{})
	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetExport(params GetExportRequest) ([]csv.Record, *errpkg.Error) {
	var records []Transfer
	query := sql.DB().Table(TableName).Scopes(modelpkg.WithStartDateAndEnd("created_at", true, params.StartAt, params.EndAt))
	if params.UsersID > 0 {
		query = query.Where("`users_id` = ?", params.UsersID)
	}
	if len(params.StatusList) > 0 {
		query = query.Where("`status` IN (?)", params.StatusList)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}

	return output, nil
}

func GetAegisExport(mainnet wallet.Mainnet, startAt modelpkg.Date, endAt modelpkg.Date) ([]AegisTransfer, *errpkg.Error) {
	var records []AegisTransfer
	query := sql.DB().Table(TableName).
		Where("`type` = ? AND `mainnet` = ? AND `action` = ? AND `status` = ?", TypeAegisManual, mainnet.String(), ActionWithdraw, StatusReviewing).
		Scopes(modelpkg.WithStartDateAndEnd("created_at", true, startAt, endAt))

	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func AegisImport(records []AegisTransfer) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			updateMap := make(map[string]any, 3)

			if record.TxID != "" {
				updateMap["tx_id"] = record.TxID
			}
			if record.Status > 0 {
				updateMap["status"] = record.Status
			}
			if record.FinishedAtValid {
				updateMap["finished_at"] = record.FinishedAt
			}

			if err := tx.Table(TableName).Where("`transfers_id` = ?", record.TransfersID).Updates(updateMap).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

// TakeOldAndUpsert takes old record and upserts new record.
// If old record is not found, it will create a new record and retrun nil.
func TakeOldAndUpsert(m *Model) (*Model, *errpkg.Error) {
	var old *Model
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if m.TransfersID != "" {
			if err := tx.Scopes(modelpkg.WithLock(true)).Where("`transfers_id` = ?", m.TransfersID).Take(&old).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				} else {
					return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
				}
			}
		} else {
			if err := tx.Scopes(modelpkg.WithLock(true)).Where("`serial` = ?", m.Serial).Take(&old).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					old = nil // reset old
				} else {
					return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
				}
			}
		}

		if old != nil {
			m.TransfersID = old.TransfersID
			m.HandlingCharge = old.HandlingCharge
			m.Extra.IP = old.Extra.IP
			m.Extra.Location = old.Extra.Location
			m.CreatedAt = old.CreatedAt

			if old.Extra.ToUSDTPrice.GreaterThan(decimal.Zero) {
				m.Valuation = old.Valuation
				m.Extra.ToUSDTPrice = old.Extra.ToUSDTPrice
				m.Extra.ToTWDPrice = old.Extra.ToTWDPrice
			} else {
				m.updateValuationIfZero()
			}
		} else {
			m.CreatedAt = time.Now()
			m.updateValuationIfZero()
		}

		if old == nil {
			if wrapErr := m.Create(tx); wrapErr != nil {
				return wrapErr
			}
		} else {
			if err := tx.Save(m).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		return nil
	}); err != nil {
		return nil, err.(*errpkg.Error)
	}

	return old, nil
}
