package suspicioustransactionspkg

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"time"

	"net/http"

	"gorm.io/gorm"
)

type scopeFunc = func(db *gorm.DB) *gorm.DB

func withDuration(duration time.Duration) scopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		if duration <= 0 {
			return db
		}

		return db.Where("`created_at` >= ?", time.Now().UTC().Add(-duration))
	}
}

func withActions(actions ...usersspottransfers.Action) scopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		if len(actions) == 0 {
			return db
		}

		return db.Where("`action` IN ?", actions)
	}
}

// func withSides(sides ...userstransactions.Side) scopeFunc {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if len(sides) == 0 {
// 			return db
// 		}

// 		return db.Where("`side` IN ?", sides)
// 	}
// }

// func withLimit(limit int) scopeFunc {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if limit <= 0 {
// 			return db
// 		}

// 		return db.Limit(limit)
// 	}
// }

func queryTransferByUser(usersID int64, scopFunctions ...scopeFunc) ([]usersspottransfers.Transfer, *errpkg.Error) {
	var records []usersspottransfers.Transfer

	query := sql.DB().Table(usersspottransfers.TableName).
		Where("`users_id` = ?", usersID).
		Order("`created_at` ASC")

	if len(scopFunctions) > 0 {
		query = query.Scopes(scopFunctions...)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func queryWithdrawTransferByToAddress(address string, scopFunctions ...scopeFunc) ([]usersspottransfers.Transfer, *errpkg.Error) {
	var (
		records []usersspottransfers.Transfer
		query   = sql.DB().Table(usersspottransfers.TableName).
			Where("`action` = ? AND `to_address` = ?", usersspottransfers.ActionWithdraw, address).
			Order("`created_at` ASC")
	)

	if len(scopFunctions) > 0 {
		query = query.Scopes(scopFunctions...)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}
