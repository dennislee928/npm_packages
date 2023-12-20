package usersmodifylogs

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

// SLStatus is Status of StatusLog
type SLStatus int32

// Status List of StatusLog
const (
	SLStatusUnverified SLStatus = iota
	SLStatusEnable
	SLStatusDisable
	SLStatusForzen
)

func (s SLStatus) Chinese() string {
	switch s {
	case SLStatusUnverified:
		return "未啟用"
	case SLStatusEnable:
		return "已啟用"
	case SLStatusDisable:
		return "已停用"
	case SLStatusForzen:
		return "凍結中"
	}

	return badStatusChinese
}

type StatusLog struct {
	// 日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 狀態
	// * 0: 未啟用
	// * 1: 已啟用
	// * 2: 已停用
	// * 3: 凍結中
	Status SLStatus `json:"status"`

	// 備註
	Comment string `json:"comment"`

	// 管理者
	ManagersName string `json:"managersName"`
}

func GetStatusLogCSVHeaders() []string {
	return []string{"日期", "狀態", "備註", "管理者"}
}

func (sl StatusLog) ToCSV() []string {
	return []string{
		sl.CreatedAt.ToString(true),
		sl.Status.Chinese(),
		sl.Comment,
		sl.ManagersName,
	}
}
