package userswithdrawalwhitelist

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Extra struct {
	// 備註
	Memo string `json:"memo"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Extra) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("users_withdrawal_whitelist.Extra: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Extra) Value() (driver.Value, error) {
	return json.Marshal(e)
}
