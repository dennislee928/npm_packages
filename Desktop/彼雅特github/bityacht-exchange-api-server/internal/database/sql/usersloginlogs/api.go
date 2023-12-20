package usersloginlogs

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

type Log struct {
	UserAgent string            `json:"userAgent"`
	Browser   string            `json:"browser"`
	Device    string            `json:"device"`
	Location  string            `json:"location" binding:"required"`
	IP        string            `json:"ip" binding:"required"`
	CreatedAt modelpkg.DateTime `json:"createdAt" binding:"required" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetLogCSVHeaders() []string {
	return []string{"位置", "瀏覽器", "裝置", "IP位置", "登入時間"}
}

func (l Log) ToCSV() []string {
	return []string{
		l.Location,
		l.Browser,
		l.Device,
		l.IP,
		l.CreatedAt.ToString(true),
	}
}
