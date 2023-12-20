package modelpkg

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	DateFormat     = time.DateOnly
	JSONDateFormat = "2006/01/02"
)

var (
	_ json.Unmarshaler = (*Date)(nil)
	_ json.Marshaler   = (*Date)(nil)
	_ sql.Scanner      = (*Date)(nil)
	_ driver.Valuer    = (*Date)(nil)
)

// @Description Format: "YYYY/MM/DD"
type Date struct {
	time.Time
}

func NewDate(t time.Time) Date {
	if t.Compare(zeroTimeInTimeLoc) <= 0 {
		return Date{Time: zeroTime}
	}

	t = t.In(DefaultTimeLoc)
	return Date{
		Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, DefaultTimeLoc),
	}
}

func (d *Date) Parse(layout string, value string) *errpkg.Error {
	if t, err := time.ParseInLocation(layout, value, DefaultTimeLoc); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGormScan, Err: err}
	} else if t.Compare(zeroTimeInTimeLoc) <= 0 {
		d.Time = zeroTime
	} else {
		d.Time = t
	}

	return nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}

	if t, err := time.ParseInLocation(JSONDateFormat, str, DefaultTimeLoc); err != nil {
		return err
	} else if t.Compare(zeroTimeInTimeLoc) <= 0 {
		d.Time = zeroTime
	} else {
		d.Time = t
	}

	return nil
}

func (d Date) ToString(omitZeroTime bool) string {
	if d.Time.IsZero() {
		if omitZeroTime {
			return ""
		}

		return zeroTime.Format(JSONDateFormat)
	}

	return d.Time.In(DefaultTimeLoc).Format(JSONDateFormat)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.ToString(true) + `"`), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *Date) Scan(value interface{}) error {
	var strDate string

	switch value := value.(type) {
	case time.Time:
		if value.IsZero() {
			d.Time = zeroTime
		} else {
			d.Time = time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, DefaultTimeLoc)
		}
		return nil
	case string:
		strDate = value
	case []byte:
		strDate = string(value)
	default:
		return errors.New("modelpkg.Date: bad type")
	}

	return d.Parse(DateFormat, strDate)
}

// Value return json value, implement driver.Valuer interface
func (d Date) Value() (driver.Value, error) {
	return d.Time.In(DefaultTimeLoc).Format(DateFormat), nil
}

func (d Date) ToDBTime() time.Time {
	if d.IsZero() {
		return zeroTime
	}

	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
