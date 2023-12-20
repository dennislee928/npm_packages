package loginlogs

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

type ExportRequest struct {
	ID      int64         `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt   modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}
