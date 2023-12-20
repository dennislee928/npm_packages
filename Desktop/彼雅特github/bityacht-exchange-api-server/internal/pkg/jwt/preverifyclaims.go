package jwt

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PreverifyUsage string

const (
	PreverifyUsageVerifyPhoneForIDV PreverifyUsage = "phoneIDV"
	PreverifyUsageForgotPassword    PreverifyUsage = "forgotPW"
	PreverifyUsageUpdateWithdraw2FA PreverifyUsage = "updateWithdraw2FA"
)

type PreverifyClaims struct {
	jwt.RegisteredClaims
	Type Type `json:"claimsType"`
	PreverifyPayload
}

func (pc PreverifyClaims) ID() int64 {
	return pc.PreverifyPayload.ID
}

type PreverifyPayload struct {
	AccountType Type           // Manager or User
	ID          int64          `json:"id"`
	Usage       PreverifyUsage `json:"usage"`
	Phone       string         `json:"phone,omitempty"`
	GASecret    string         `json:"gaSecret,omitempty"`
}

func IssuePreverifyToken(payload PreverifyPayload) (string, PreverifyClaims, *errpkg.Error) {
	now := time.Now()
	claims := PreverifyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Issuer: "",
			// Subject: "",
			// Audience: "",
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
			// NotBefore: ,
			IssuedAt: jwt.NewNumericDate(now),
			ID:       uuid.NewString(),
		},
		Type:             TypePreverify,
		PreverifyPayload: payload,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(configs.Config.JWT.Key)
	if err != nil {
		return "", PreverifyClaims{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJWTIssueToken, Err: err}
	}

	return token, claims, nil
}

func ValidatePreverify(accountType Type, id int64, usage PreverifyUsage, accessToken string) (PreverifyClaims, *errpkg.Error) {
	var claims PreverifyClaims

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
	} else if claims.Type != TypePreverify {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad type")}
	} else if claims.AccountType != accountType {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad accountType")}
	} else if claims.ID() != 0 && id != 0 && claims.ID() != id {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad id")}
	} else if claims.Usage != usage {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad usage")}
	}

	return claims, nil
}
