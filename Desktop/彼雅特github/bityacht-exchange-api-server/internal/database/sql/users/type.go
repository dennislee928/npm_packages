package users

import (
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/mmdb"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Type int32

const (
	TypeNaturalPerson Type = iota + 1
	TypeJuridicalPerson
)

func (t Type) Chinese() string {
	switch t {
	case TypeNaturalPerson:
		return "自然人"
	case TypeJuridicalPerson:
		return "法人"
	}

	return "未知錯誤"
}

type Gender int32

const (
	GenderUnknown = iota
	GenderMale
	GenderFemale
	GenderX
)

func (g *Gender) ParseFromKryptoGO(gender string) {
	switch gender {
	case "M":
		*g = GenderMale
	case "F":
		*g = GenderFemale
	case "X":
		*g = GenderX
	default:
		*g = GenderUnknown
	}
}

type Status = usersmodifylogs.SLStatus

type TwoFAType uint32

const (
	TwoFATypeEmail TwoFAType = 1 << iota
	TwoFATypeSMS
	TwoFATypeGoogleAuthenticator
)

var _ sql.Scanner = (*Extra)(nil)
var _ driver.Valuer = (*Extra)(nil)

type Extra struct {
	MobileBarcode        string           `json:"MB,omitempty"`
	Login2FAType         TwoFAType        `json:"L2FAT,omitempty"`
	Withdraw2FAType      TwoFAType        `json:"W2FAT,omitempty"`
	LastChangePasswordAt time.Time        `json:"LCPA,omitempty"`
	InviteCode           string           `json:"IC"` //! if change the JSON key, Remember Update for Seeder also.
	RegisterIP           string           `json:"RIP,omitempty"`
	RegisterLocation     *mmdb.CityResult `json:"RL,omitempty"`

	GoogleAuthenticatorSecret string `json:"GAS,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (e *Extra) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("users.Extra: bad type")
	}

	return json.Unmarshal(bytes, e)
}

// Value return json value, implement driver.Valuer interface
func (e Extra) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e Extra) GetLogin2FAType() TwoFAType {
	return e.Login2FAType | TwoFATypeEmail
}

func (e Extra) GetWithdraw2FAType() TwoFAType {
	return e.Withdraw2FAType | TwoFATypeEmail | TwoFATypeSMS
}

func (e Extra) IsEnableWithdrawGA2FA() bool {
	return (e.Withdraw2FAType&TwoFATypeGoogleAuthenticator > 0) && e.GoogleAuthenticatorSecret != ""
}

func (e *Extra) setRegisterIP(ip string) *errpkg.Error {
	if ip == "" {
		return nil
	}

	location, err := mmdb.LookupCity(ip)
	if err != nil {
		return err
	}

	if location.Country.ISOCode != "" {
		e.RegisterIP = ip
		e.RegisterLocation = &location
	}

	return nil
}
