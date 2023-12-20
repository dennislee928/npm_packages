package modelpkg

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var _ json.Marshaler = (*TWCellPhone)(nil)
var _ json.Unmarshaler = (*TWCellPhone)(nil)

// Ref: https://en.wikipedia.org/wiki/List_of_country_calling_codes
const taiwanCallingCode = "+886"

var twCellPhoneLocalFormat = regexp.MustCompile(`^09[0-9]{8}$`)

type TWCellPhone string

func (twcp TWCellPhone) String() string {
	return string(twcp)
}

func (twcp TWCellPhone) GetLocalString() string {
	return "0" + strings.TrimPrefix(string(twcp), taiwanCallingCode)
}

func (twcp *TWCellPhone) UnmarshalJSON(data []byte) error {
	number := strings.Trim(string(data), `"`)

	if number == "" {
		return nil
	} else if !twCellPhoneLocalFormat.MatchString(number) {
		return errors.New("bad phone format")
	}

	*twcp = TWCellPhone(taiwanCallingCode + number[1:])
	return nil
}

func (twcp TWCellPhone) MarshalJSON() ([]byte, error) {
	if twcp == "" {
		return []byte(`""`), nil
	}

	return []byte(`"` + twcp.GetLocalString() + `"`), nil
}
