package modelpkg

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func WithSoftID[T any](id T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(WithNotDeleted()).Where("`id` = ?", id)
	}
}

func WithNotDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("`deleted_at` IS NULL")
	}
}

func WithStartAtAndEnd(field string, startAt, endAt DateTime) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || startAt.IsZero() || endAt.IsZero() {
			return db
		}

		return db.Where(addBackQuote(field)+" BETWEEN ? AND ?", startAt.UTC(), endAt.UTC())
	}
}

const endTimeOfDate = 24*time.Hour - 1 // 23:59:59.9999....

func WithStartDateAndEnd(field string, queryWithTime bool, startAt, endAt Date) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || startAt.IsZero() || endAt.IsZero() {
			return db
		}
		field = addBackQuote(field)

		if queryWithTime {
			return db.Where(field+" BETWEEN ? AND ?", startAt.UTC(), endAt.Time.Add(endTimeOfDate).UTC())
		}

		return db.Where(field+" BETWEEN ? AND ?", startAt, endAt)
	}
}

func addBackQuote(field string) string {
	if strings.HasPrefix(field, "`") || strings.HasSuffix(field, ")") {
		return field
	}

	return "`" + field + "`"
}

func WithPaginator(paginator *Paginator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if paginator == nil {
			return db
		}

		return db.Limit(paginator.PageSize).Offset(paginator.Offset())
	}
}

func WithLock(noWait bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if noWait {
			return db.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"})
		}

		return db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
}
