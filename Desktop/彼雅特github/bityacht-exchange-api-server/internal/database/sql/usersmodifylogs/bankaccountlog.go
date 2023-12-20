package usersmodifylogs

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

// BALStatus is Status of BankAccountLog
type BALStatus int32

// Status List of BankAccountLog
const (
	BALStatusUnknown BALStatus = iota
	BALStatusPending
	BALStatusAccepted
	BALStatusRejected
)

func (s BALStatus) Chinese() string {
	switch s {
	case BALStatusUnknown:
		return "未綁定"
	case BALStatusPending:
		return "審核中"
	case BALStatusAccepted:
		return "已綁定"
	case BALStatusRejected:
		return "未通過"
	}

	return badStatusChinese
}

type BankAccountLog struct {
	// 日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 狀態
	// * 0: 未綁定
	// * 1: 審核中
	// * 2: 已綁定
	// * 3: 未通過
	Status BALStatus `json:"status"`

	// 備註
	Comment string `json:"comment"`

	// 操作者
	ManagersName string `json:"managersName"`
}

func GetBankAccountLogHeaders() []string {
	return []string{"日期", "狀態", "備註", "管理者"}
}

func (bal BankAccountLog) ToCSV() []string {
	return []string{
		bal.CreatedAt.ToString(true),
		bal.Status.Chinese(),
		bal.Comment,
		bal.ManagersName,
	}
}
