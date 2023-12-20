package tobereview

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"errors"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

const lowRiskThreshold = 7
const mediumRiskThreshold = 12
const longReviewPeriod = -3  // years
const shortReviewPeriod = -1 // years

func UpdateFinalReviewToBeReview(job gocron.Job) {
	jobLogger := logger.Logger.With().Str("service", "to be review check schedule").Logger()
	now := time.Now()

	query := sql.DB().Table(users.TableName).
		Where("`final_review` = ?", usersmodifylogs.RLStatusApproved).
		Where("`internal_risks_total` <= ? AND `final_review_time` < ?", lowRiskThreshold, now.AddDate(longReviewPeriod, 0, 0)).
		Or("`internal_risks_total` > ? AND `internal_risks_total` <= ? AND `final_review_time` < ?", lowRiskThreshold, mediumRiskThreshold, now.AddDate(shortReviewPeriod, 0, 0))

	if err := query.Update("final_review", usersmodifylogs.RLStatusToBeReview).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		jobLogger.Err(err).Msg("sql update failed")
	}
}
