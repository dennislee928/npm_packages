package memory

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/cache/memory/banners"
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"context"
	"time"
)

func Init(ctx context.Context) {
	go func() {
		updateSpotTrendTimer := time.NewTimer(0)
		updateSqlCacherTimer := time.NewTimer(0)
		banners.Init(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-updateSpotTrendTimer.C:
				if err := spottrend.Update(ctx); err != nil {
					logger.Logger.Err(err.Err).Msg("memory spottrend Update error")
					updateSpotTrendTimer.Reset(time.Second)
					continue
				}

				updateSpotTrendTimer.Reset(configs.Config.Exchange.UpdateTrendInterval)
			case <-updateSqlCacherTimer.C:
				if err := sqlcache.Update(); err != nil {
					logger.Logger.Err(err.Err).Msg("sqlcache Update error")
					updateSqlCacherTimer.Reset(time.Second)
					continue
				}

				updateSqlCacherTimer.Reset(30 * time.Second)
			}
		}
	}()
}
