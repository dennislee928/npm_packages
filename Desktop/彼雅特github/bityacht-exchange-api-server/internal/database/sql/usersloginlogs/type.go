package usersloginlogs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var _ sql.Scanner = (*Headers)(nil)
var _ driver.Valuer = (*Headers)(nil)

type Headers map[string][]string

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Headers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("usersloginlogs.Headers: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Headers) Value() (driver.Value, error) {
	return json.Marshal(e)
}
