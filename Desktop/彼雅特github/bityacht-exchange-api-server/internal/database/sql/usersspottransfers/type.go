package usersspottransfers

import (
	"bityacht-exchange-api-server/internal/pkg/mmdb"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/shopspring/decimal"
)

type Type int32

const (
	TypeCybavoAPI = iota + 1
	TypeAegisManual
)

type Status int32

const (
	StatusProcessing Status = iota + 1
	StatusFinished
	StatusFailed
	StatusCanceled
	StatusReviewing
)

func (s Status) Chinese() string {
	switch s {
	case StatusProcessing:
		return "處理中"
	case StatusFinished:
		return "已完成"
	case StatusFailed:
		return "已失敗"
	case StatusCanceled:
		return "已取消"
	case StatusReviewing:
		return "審核中"
	}

	return "未知錯誤"
}

func (s *Status) FromChinese(chinese string) {
	switch chinese {
	case "處理中":
		*s = StatusProcessing
	case "已完成":
		*s = StatusFinished
	case "已失敗":
		*s = StatusFailed
	case "已取消":
		*s = StatusCanceled
	case "審核中":
		*s = StatusReviewing
	default:
		*s = 0
	}
}

type Action int32

const (
	ActionDeposit Action = iota + 1
	ActionWithdraw
)

func (a Action) Chinese() string {
	switch a {
	case ActionDeposit:
		return "接收"
	case ActionWithdraw:
		return "發送"
	}

	return "未知錯誤"
}

type Extra struct {
	RawMessage  json.RawMessage `json:"rm,omitempty"`
	ToUSDTPrice decimal.Decimal `json:"tup"`
	ToTWDPrice  decimal.Decimal `json:"ttp"`
	IP          string          `json:"ip"`
	Location    mmdb.CityResult `json:"l"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Extra) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("usersspottransfers.Extra: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Extra) Value() (driver.Value, error) {
	return json.Marshal(e)
}
