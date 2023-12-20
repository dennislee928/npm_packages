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

var (
	DefaultTimeLoc    *time.Location
	zeroTime          time.Time
	zeroTimeInTimeLoc time.Time
)

func init() {
	var err error

	if DefaultTimeLoc, err = time.LoadLocation("Asia/Taipei"); err != nil {
		panic(err)
	}
	zeroTimeInTimeLoc = time.Date(1, 1, 1, 0, 0, 0, 0, DefaultTimeLoc)
}

const (
	DateTimeFormat     = time.DateTime
	JSONDateTimeFormat = "2006/01/02 15:04:05"
)

var (
	_ json.Unmarshaler = (*DateTime)(nil)
	_ json.Marshaler   = (*DateTime)(nil)
	_ sql.Scanner      = (*DateTime)(nil)
	_ driver.Valuer    = (*DateTime)(nil)
)

// @Description Format: "YYYY/MM/DD HH:MM:SS"
type DateTime struct {
	time.Time
}

func NewDateTime(t time.Time) DateTime {
	if t.Compare(zeroTimeInTimeLoc) <= 0 {
		return DateTime{Time: zeroTime}
	}

	return DateTime{Time: t.In(DefaultTimeLoc)}
}

func (d *DateTime) Parse(layout string, value string) *errpkg.Error {
	if t, err := time.ParseInLocation(layout, value, DefaultTimeLoc); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGormScan, Err: err}
	} else if t.Compare(zeroTimeInTimeLoc) <= 0 {
		d.Time = zeroTime
	} else {
		d.Time = t
	}

	return nil
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}

	if t, err := time.ParseInLocation(JSONDateTimeFormat, str, DefaultTimeLoc); err != nil {
		return err
	} else if t.Compare(zeroTimeInTimeLoc) <= 0 {
		d.Time = zeroTime
	} else {
		d.Time = t
	}

	return nil
}

func (d DateTime) ToString(omitZeroTime bool) string {
	if d.Time.IsZero() {
		if omitZeroTime {
			return ""
		}

		return zeroTime.Format(JSONDateTimeFormat)
	}

	return d.Time.In(DefaultTimeLoc).Format(JSONDateTimeFormat)
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.ToString(true) + `"`), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DateTime) Scan(value interface{}) error {
	var strUTCDateTime string

	switch value := value.(type) {
	case time.Time:
		if value.IsZero() {
			d.Time = zeroTime
		} else {
			d.Time = value.In(DefaultTimeLoc)
		}
		return nil
	case string:
		strUTCDateTime = value
	case []byte:
		strUTCDateTime = string(value)
	default:
		return errors.New("modelpkg.DateTime: bad type")
	}

	if t, err := time.ParseInLocation(DateTimeFormat, strUTCDateTime, time.UTC); err != nil {
		return err
	} else if t.IsZero() {
		d.Time = zeroTime
	} else {
		d.Time = t.In(DefaultTimeLoc)
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (d DateTime) Value() (driver.Value, error) {
	return d.Time.UTC().Format(DateTimeFormat), nil
}

func (d DateTime) ToDBTime() time.Time {
	return d.Time.UTC()
}

func ValidStartAndEndAt(startAt time.Time, endAt time.Time) *errpkg.Error {
	if !startAt.IsZero() && !endAt.IsZero() && startAt.After(endAt) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("startAt must before endAt")}
	}

	return nil
}
