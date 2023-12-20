package managers

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Status int32

const (
	StatusDisable Status = iota
	StatusEnable
)

var _ sql.Scanner = (*Extra)(nil)
var _ driver.Valuer = (*Extra)(nil)

type Extra struct {
	NeedChangePassword bool `json:"ncpw,omitempty"` // abbreviation for reduce the size in DB
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Extra) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("managers.Extra: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Extra) Value() (driver.Value, error) {
	return json.Marshal(e)
}
