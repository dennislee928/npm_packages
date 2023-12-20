package banners

import (
	"bityacht-exchange-api-server/internal/database/sql/banners"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

var cache cacheStruct

type cacheStruct struct {
	Banners           []banners.Model
	ActivatingBanners []Banner
	EventTimer        *time.Timer
	BannersResp       json.RawMessage // Marshal from ActivatingBanners

	mux sync.RWMutex
}

type Banner struct {
	WebImage   string `json:"webImage"`   // Web 圖檔 URL
	AppImage   string `json:"appImage"`   // 手機版圖檔 URL
	Title      string `json:"title"`      // 大標
	SubTitle   string `json:"subTitle"`   // 副標
	ButtonText string `json:"buttonText"` // 按鈕文字
	ButtonUrl  string `json:"buttonUrl"`  // 按鈕連結
}

const updateInterval = 30 * time.Second

func Init(ctx context.Context) {
	go func() {
		updateFromDBTimer := time.NewTimer(updateInterval)
		cache.EventTimer = time.NewTimer(0)
		resetTimer(cache.EventTimer, 0)
		UpdateBanners(true, time.Now())

		for {
			select {
			case <-ctx.Done():
				return
			case now := <-cache.EventTimer.C:
				UpdateBanners(false, now)
			case now := <-updateFromDBTimer.C:
				UpdateBanners(true, now)
				resetTimer(updateFromDBTimer, updateInterval)
			}
		}
	}()
}

func GetBannersResp() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.BannersResp
}

func UpdateBanners(updateFromDB bool, now time.Time) *errpkg.Error {
	var bannersToCheck []banners.Model

	if updateFromDB {
		var err *errpkg.Error
		bannersToCheck, err = banners.GetEnableList()

		if err != nil {
			logger.Logger.Err(err.Err).Msg("banners cache get list for user error")
			return err
		}
	} else {
		cache.mux.RLock()
		bannersToCheck = cache.Banners
		cache.mux.RUnlock()
	}

	var nextEventTime time.Time
	newActivatingBanners := make([]Banner, 0, len(bannersToCheck))

	for _, record := range bannersToCheck {
		if !record.StartAt.IsZero() && !record.StartAt.Before(now) { // start >= now: Not Start
			if nextEventTime.IsZero() || record.StartAt.Before(nextEventTime) {
				nextEventTime = record.StartAt
			}
			continue
		}
		if !record.EndAt.IsZero() {
			if !record.EndAt.After(now) { // end <= now: End!
				continue
			} else if nextEventTime.IsZero() || record.EndAt.Before(nextEventTime) {
				nextEventTime = record.EndAt
			}
		}

		newActivatingBanners = append(newActivatingBanners, Banner{
			WebImage:   record.WebImage,
			AppImage:   record.AppImage,
			Title:      record.Title,
			SubTitle:   record.SubTitle,
			ButtonText: record.Button,
			ButtonUrl:  record.ButtonUrl,
		})
	}

	newBannersResp, err := json.Marshal(newActivatingBanners)
	if err != nil {
		logger.Logger.Err(err).Msg("banners cache marshal new banners error")
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
	}

	cache.mux.Lock()
	defer cache.mux.Unlock()

	if updateFromDB {
		cache.Banners = bannersToCheck
	}
	if !nextEventTime.IsZero() {
		resetTimer(cache.EventTimer, nextEventTime.Sub(now))
	}
	cache.ActivatingBanners = slices.Clip(newActivatingBanners)
	cache.BannersResp = newBannersResp

	return nil
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C: // Clean up the timer.
		default:
		}
	}

	if duration > 0 {
		timer.Reset(duration)
	}
}
