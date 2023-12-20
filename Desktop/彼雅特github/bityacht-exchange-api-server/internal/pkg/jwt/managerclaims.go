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

// ClaimsKeyInGin for authorization
const ClaimsKeyInGin = "claims"

type ManagerClaims struct {
	jwt.RegisteredClaims
	Type Type `json:"claimsType"`
	ManagerPayload
}

func (mc ManagerClaims) ID() int64 {
	return mc.ManagerPayload.ID
}

type ManagerPayload struct {
	ManagersRolesID int64  `json:"managersRolesID"`
	ID              int64  `json:"id,omitempty"`
	Name            string `json:"name"`
}

func IssueManagerToken(payload ManagerPayload) (token string, claims ManagerClaims, err *errpkg.Error) {
	now := time.Now()
	claims = ManagerClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Issuer: "",
			// Subject: "",
			// Audience: "",
			ExpiresAt: jwt.NewNumericDate(now.Add(configs.Config.JWT.AccessTokenLifetime)),
			// NotBefore: ,
			IssuedAt: jwt.NewNumericDate(now),
			ID:       uuid.NewString(),
		},
		Type:           TypeManager,
		ManagerPayload: payload,
	}

	var signErr error
	if token, signErr = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(configs.Config.JWT.Key); signErr != nil {
		return "", ManagerClaims{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJWTIssueToken, Err: signErr}
	}

	return
}

func ValidateManager(accessToken string) (ManagerClaims, *errpkg.Error) {
	var claims ManagerClaims

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
	} else if claims.Type != TypeManager {
		return claims, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTInvalid, Err: errors.New("bad type")}
	}

	return claims, nil
}
