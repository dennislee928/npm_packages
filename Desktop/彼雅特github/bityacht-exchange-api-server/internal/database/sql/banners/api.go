package banners

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"mime/multipart"
)

type Banner struct {
	ID        int64             `json:"id"`
	WebImage  string            `json:"webImage"`
	AppImage  string            `json:"appImage"`
	Priority  int64             `json:"priority"`
	Title     string            `json:"title"`
	SubTitle  string            `json:"subTitle"`
	Button    string            `json:"button"`
	ButtonUrl string            `json:"buttonUrl"`
	Status    int32             `json:"status"`
	StartAt   modelpkg.DateTime `json:"startAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
	EndAt     modelpkg.DateTime `json:"endAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

type CreateRequest struct {
	WebImage *multipart.FileHeader `form:"webImage" swaggerignore:"true"`
	AppImage *multipart.FileHeader `form:"appImage" swaggerignore:"true"`

	Title     string            `form:"title"`
	SubTitle  string            `form:"subTitle"`
	Button    string            `form:"button"`
	ButtonUrl string            `form:"buttonUrl"`
	Status    Status            `form:"status"`
	StartAt   modelpkg.DateTime `form:"startAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
	EndAt     modelpkg.DateTime `form:"endAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

type UpdateRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	WebImage *multipart.FileHeader `form:"webImage" swaggerignore:"true"`
	AppImage *multipart.FileHeader `form:"appImage" swaggerignore:"true"`

	Title     string            `form:"title"`
	SubTitle  string            `form:"subTitle"`
	Button    string            `form:"button"`
	ButtonUrl string            `form:"buttonUrl"`
	Status    Status            `form:"status"`
	StartAt   modelpkg.DateTime `form:"startAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
	EndAt     modelpkg.DateTime `form:"endAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

type PriorityUpdateRequset struct {
	Rows []struct {
		ID       int64 `json:"id"`
		Priority int64 `json:"priority"`
	} `json:"rows"`
}

type DeleteRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}
