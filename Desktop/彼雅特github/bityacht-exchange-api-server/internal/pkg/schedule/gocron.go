package schedule

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/schedule/tobereview"
	"bityacht-exchange-api-server/internal/pkg/schedule/usersvaluation"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler

func init() {
	s = gocron.NewScheduler(modelpkg.DefaultTimeLoc)
	s.TagsUnique()

	if _, err := s.Tag("UsersValuation: Daily").Every(1).Day().At("00:00").DoWithJobDetails(usersvaluation.CalcUsersValuation); err != nil {
		panic(err)
	}

	if _, err := s.Tag("ToBeReview: Daily").Every(1).Day().At("00:00").DoWithJobDetails(tobereview.UpdateFinalReviewToBeReview); err != nil {
		panic(err)
	}

	s.StartAsync()
}
