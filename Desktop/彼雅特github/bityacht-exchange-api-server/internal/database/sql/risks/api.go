package risks

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

type CreateRequest struct {
	Factor    string `json:"factor" binding:"required"`
	SubFactor string `json:"subFactor" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	Score     int64  `json:"score" binding:"required"`
}

type UpdateRequest struct {
	ID        int64  `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	Factor    string `json:"factor"`
	SubFactor string `json:"subFactor"`
	Detail    string `json:"detail"`
	Score     int64  `json:"score" binding:"gt=0"`
}

type DeleteRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
}

type Risk struct {
	ID        int64             `json:"id"`
	Factor    string            `json:"factor"`
	SubFactor string            `json:"subFactor"`
	Detail    string            `json:"detail"`
	Score     int64             `json:"score"`
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}
type UpdateRisksRequest struct {
	ID       int64   `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	RisksIDs []int64 `json:"risksIDs"`
}
