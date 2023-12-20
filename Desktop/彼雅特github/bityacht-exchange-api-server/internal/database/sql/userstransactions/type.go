package userstransactions

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Status int32

const (
	StatusFilled Status = iota + 1
	StatusKilled
)

func (s Status) Chinese() string {
	switch s {
	case StatusFilled:
		return "已完成"
	case StatusKilled:
		return "已失敗"
	}

	return "未知錯誤"
}

type Side int32

const (
	SideBuy Side = iota + 1
	SideSell
)

func (s Side) String() string {
	switch s {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	default:
		return "unknown"
	}
}

func (s Side) Chinese() string {
	switch s {
	case SideBuy:
		return "買"
	case SideSell:
		return "賣"
	default:
		return "未知錯誤"
	}
}

var _ sql.Scanner = (*Extra)(nil)
var _ driver.Valuer = (*Extra)(nil)

type Extra struct {
	BinanceResp string `json:"binanceResp,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Extra) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("userstransactions.Extra: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Extra) Value() (driver.Value, error) {
	return json.Marshal(e)
}
