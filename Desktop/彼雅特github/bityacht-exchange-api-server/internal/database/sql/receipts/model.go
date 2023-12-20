package receipts

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of receipts table
const TableName = "receipts"

// Model of receipts table
type Model struct {
	ID              string    // Use Order ID
	Status          Status    `gorm:"not null;default:1"`
	UserID          int64     `gorm:"not null;index"`
	InvoiceAmount   int64     `gorm:"not null;default:0"`
	InvoiceID       string    `gorm:"not null;index;default:''"`
	SalesAmount     int64     `gorm:"not null;default:0"`
	Tax             int64     `gorm:"not null;default:0"`
	InvoiceIssuedAt time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	Barcode         string    `gorm:"not null;default:''"`
	CreatedAt       time.Time `gorm:"not null;index;default:UTC_TIMESTAMP()"`
	DeletedAt       *time.Time
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func Create(m *Model) *errpkg.Error {
	if err := sql.DB().Create(m).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func GetList(paginator *modelpkg.Paginator, req GetListRequest, searcher modelpkg.Searcher) ([]ListItem, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]ListItem, 0)
	query := sql.DB().Table(TableName).Scopes(
		modelpkg.WithNotDeleted(),
		modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt),
	)
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	query = searcher.AddToQuery(query, []string{"id", "invoice_id"}).Session(&gorm.Session{})

	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetByID(id string) (*Model, *errpkg.Error) {
	var record Model
	if err := sql.DB().Table(TableName).Scopes(modelpkg.WithSoftID(id)).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return &record, nil
}

func GetDetailByID(id string) (*DetailItem, *errpkg.Error) {
	var record DetailItem
	if err := sql.DB().Table(TableName).Scopes(modelpkg.WithSoftID(id)).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return &record, nil
}

// SetReciptStatusByID set receipt status by id.
//
// NOTE: Use IssueReceipt() to issue receipt.
func SetReciptStatusByID(id string, status Status) *errpkg.Error {
	if status == StatusIssued {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("can not set receipt status to issued")}
	}

	if err := sql.DB().Transaction(
		func(tx *gorm.DB) error {
			query := tx.Table(TableName).Scopes(modelpkg.WithSoftID(id)).Session(&gorm.Session{})

			var record Model
			if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&record).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			if record.Status == StatusIssued {
				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("receipt status is issued")}
			}

			if record.Status == status {
				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("receipt status is already set")}
			}

			if err := query.Updates(map[string]interface{}{
				"status": status,
			}).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			return nil
		}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

type IssueReceiptParam struct {
	InvoiceID   string
	CreatedAt   time.Time
	SalesAmount int64
	Tax         int64
}

func IssueReceiptByID(id string, param IssueReceiptParam) *errpkg.Error {
	if err := sql.DB().Transaction(
		func(tx *gorm.DB) error {
			query := tx.Table(TableName).Scopes(modelpkg.WithSoftID(id)).Session(&gorm.Session{})

			var record Model
			if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&record).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			if !canIssueStatus(record.Status) {
				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("receipt status is ready to issue")}
			}

			if err := query.Updates(map[string]interface{}{
				"status":            StatusIssued,
				"invoice_id":        param.InvoiceID,
				"invoice_issued_at": param.CreatedAt.UTC(),
				"sales_amount":      param.SalesAmount,
				"tax":               param.Tax,
			}).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			return nil
		}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetExport(req ExportRequest) ([]csv.Record, *errpkg.Error) {
	var records []DetailItem

	query := sql.DB().Table(TableName).Scopes(
		modelpkg.WithNotDeleted(),
		modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt),
	)
	if len(req.Statuses) != 0 {
		query = query.Where("`status` IN ?", req.Statuses)
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

func GetIssuingReceiptOrderIDs() ([]string, *errpkg.Error) {
	var ids []string
	if err := sql.DB().Table(TableName).Scopes(
		modelpkg.WithNotDeleted(),
	).Where("status = ?", StatusIssuing).Pluck("id", &ids).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return ids, nil
}
