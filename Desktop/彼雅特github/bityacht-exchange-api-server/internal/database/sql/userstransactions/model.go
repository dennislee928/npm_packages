package userstransactions

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/receipts"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"fmt"

	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/receipt"
	dbsql "database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TableName of users_transactions table
const TableName = "users_transactions"

// Model of users_transactions table
type Model struct {
	TransactionsID        string          `gorm:"not null;primaryKey"`
	UsersID               int64           `gorm:"not null;default:0"`
	BaseSymbol            string          `gorm:"not null;default:''"`
	QuoteSymbol           string          `gorm:"not null;default:''"`
	Status                Status          `gorm:"not null;default:0"`
	Side                  Side            `gorm:"not null;default:0"`
	Quantity              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Price                 decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Amount                decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	TwdExchangeRate       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	TwdTotalValue         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	HandlingCharge        decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceID             int64           `gorm:"not null;default:0"`
	BinanceQuantity       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinancePrice          decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceAmount         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BinanceHandlingCharge decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Extra                 Extra           `gorm:"type:json;not null;default:'{}'"`
	CreatedAt             time.Time       `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.TransactionsID = modelpkg.GetOrderID(modelpkg.TypeSpot, modelpkg.Action(m.Side))
	return nil
}

func (m *Model) Create(tx *gorm.DB) *errpkg.Error {
	if tx == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil gorm DB")}
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

func Create(record *Model, commissionRecord *userscommissions.Model, paySymbol string, payAmount decimal.Decimal, earnSymbol string, earnAmount decimal.Decimal) (*receipts.Model, *errpkg.Error) {
	var receiptRecord *receipts.Model

	if twdHandlingCharge := record.HandlingCharge.Mul(record.TwdExchangeRate).Round(0); twdHandlingCharge.GreaterThanOrEqual(decimal.NewFromInt(1)) {
		receiptRecord = &receipts.Model{
			Status:        receipts.StatusPending,
			UserID:        record.UsersID,
			InvoiceAmount: twdHandlingCharge.IntPart(),
		}
		receiptRecord.SalesAmount, receiptRecord.Tax = receipt.CalcSalesAndTaxFromTotal(twdHandlingCharge, decimal.New(5, -2))

		userRecord, err := users.GetByID(record.UsersID)
		if err != nil {
			return nil, err
		}
		receiptRecord.Barcode = userRecord.Extra.MobileBarcode
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := record.Create(tx); err != nil {
			return err
		} else if record.Status != StatusFilled {
			return nil
		}

		if receiptRecord != nil {
			receiptRecord.ID = record.TransactionsID

			if err := tx.Create(&receiptRecord).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		if commissionRecord != nil {
			commissionRecord.TransactionsID = dbsql.NullString{String: record.TransactionsID, Valid: true}
			if err := tx.Create(commissionRecord).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		// StatusFilled -> Modify the balance
		payWallet, err := userswallets.GetOrCreateWithLock(tx, record.UsersID, paySymbol)
		if err != nil {
			return err
		} else if ok := payWallet.Withdraw(payAmount); !ok {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadAmount, Err: errors.New("tx.create failed to withdraw from pay wallet")}
		}

		earnWallet, err := userswallets.GetOrCreateWithLock(tx, record.UsersID, earnSymbol)
		if err != nil {
			return err
		}
		earnWallet.Deposit(earnAmount)

		if err := tx.Save(&payWallet).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err = tx.Save(&earnWallet).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return nil, err.(*errpkg.Error)
	}

	return receiptRecord, nil
}

func GetList(req GetListRequest, paginator *modelpkg.Paginator, searcher modelpkg.Searcher) ([]TransactionForManager, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]TransactionForManager, 0)
	query := sql.DB().Select("`t`.*, `r`.`invoice_id` AS `invoice_id`, `r`.`invoice_amount` AS `invoice_amount`, `r`.`status` AS `invoice_status`").
		Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `r` ON `r`.`id` = `t`.`transactions_id`", receipts.TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt))

	if req.UsersID > 0 {
		query = query.Where("`t`.`users_id` = ?", req.UsersID)
	}
	if req.Status > 0 {
		query = query.Where("`t`.`status` = ?", req.Status)
	}
	if req.Side > 0 {
		query = query.Where("`t`.`side` = ?", req.Side)
	}

	query = searcher.AddToQuery(query, []string{"id", "users_id"}).Session(&gorm.Session{})
	if err := query.Scopes(modelpkg.WithPaginator(paginator)).Order("`t`.`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func GetExport(params GetExportRequest) ([]csv.Record, *errpkg.Error) {
	var records []TransactionForManager

	query := sql.DB().Select("`t`.*, `r`.`invoice_id` AS `invoice_id`, `r`.`invoice_amount` AS `invoice_amount`, `r`.`status` AS `invoice_status`").
		Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `r` ON `r`.`id` = `t`.`transactions_id`", receipts.TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, params.StartAt, params.EndAt))

	if params.UsersID > 0 {
		query = query.Where("`t`.`users_id` = ?", params.UsersID)
	}
	if len(params.StatusList) > 0 {
		query = query.Where("`t`.`status` IN (?)", params.StatusList)
	}
	if err := query.Order("`t`.`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}
	return output, nil
}

func GetTransactionForUserListByUser(usersID int64, req GetTransactionForUserListRequest, paginator *modelpkg.Paginator) ([]TransactionForUser, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]TransactionForUser, 0)
	query := sql.DB().Select("`t`.*, `r`.`invoice_amount` AS `twd_handling_charge`").
		Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `r` ON `r`.`id` = `t`.`transactions_id`", receipts.TableName)).
		Where("`t`.`users_id` = ?", usersID).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt))

	if req.Status > 0 {
		query = query.Where("`t`.`status` = ?", req.Status)
	}
	if req.Side > 0 {
		query = query.Where("`t`.`side` = ?", req.Side)
	}

	query = query.Session(&gorm.Session{})
	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}
