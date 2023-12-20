package jwt

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql/users"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Type Type `json:"claimsType"`
	UserPayload
}

func (u UserClaims) ID() int64 {
	return u.UserPayload.ID
}

type UserPayload struct {
	ID            int64        `json:"id"`
	Account       string       `json:"account"`
	CountriesCode string       `json:"countriesCode"`
	Type          users.Type   `json:"type"`
	FirstName     string       `json:"firstName"`
	LastName      string       `json:"lastName"`
	Level         int32        `json:"level"`
	Status        users.Status `json:"status"`
	InviterID     int64        `json:"inviterID,omitempty"`
	LoginAt       time.Time    `json:"loginAt"`
}

func NewUserPayload(record users.Model, loginAt time.Time) UserPayload {
	return UserPayload{
		ID:            record.ID,
		Account:       record.Account,
		CountriesCode: record.GetCountriesCode(),
		Type:          record.Type,
		FirstName:     record.FirstName,
		LastName:      record.LastName,
		Level:         record.Level,
		Status:        record.Status,
		InviterID:     record.InviterID.Int64,
		LoginAt:       loginAt,
	}
}

func IssueUserToken(payload UserPayload, opts ...ClaimsOption) (string, UserClaims, *errpkg.Error) {
	now := time.Now()
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Issuer: "",
			// Subject: "",
			// Audience: "",
			ExpiresAt: jwt.NewNumericDate(now.Add(configs.Config.JWT.AccessTokenLifetime)),
			// NotBefore: ,
			IssuedAt: jwt.NewNumericDate(now),
			ID:       uuid.NewString(),
		},
		Type:        TypeUser,
		UserPayload: payload,
	}

	for _, opt := range opts {
		opt.apply(&claims.RegisteredClaims)
	}

	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(configs.Config.JWT.Key); err != nil {
		return "", UserClaims{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJWTIssueToken, Err: err}
	} else {
		return token, claims, nil
	}
}

func ValidateUser(accessToken string) (UserClaims, *errpkg.Error) {
	var claims UserClaims

	if _, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTBadSigning, Err: fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])}
		}

		return configs.Config.JWT.Key, nil
	}); err != nil {
		switch err := err.(type) {
		case *errpkg.Error: // For keyFunc
			return claims, err
		default:
			return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: err}
		}
	} else if claims.Type != TypeUser {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad type")}
	}

	return claims, nil
}
