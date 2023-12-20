package schedulelogs

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Status int32

const (
	StatusRunning Status = iota + 1
	StatusAbort
	StatusFinished
)

var _ sql.Scanner = (*Result)(nil)
var _ driver.Valuer = (*Result)(nil)

type Result struct {
	Error        string                            `json:"err,omitempty"`
	CurrencyInfo map[string]spottrend.CurrencyInfo `json:"currencyInfo"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (r *Result) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("schedule_logs.Result: bad type")
	}

	return json.Unmarshal(bytes, r)
}

// Value return json value, implement driver.Valuer interface
func (r Result) Value() (driver.Value, error) {
	return json.Marshal(r)
}
