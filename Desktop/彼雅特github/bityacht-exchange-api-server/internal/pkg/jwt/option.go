package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsOption interface {
	apply(*jwt.RegisteredClaims)
}

type jwtIDOption struct {
	JWTID string
}

func (o jwtIDOption) apply(claims *jwt.RegisteredClaims) {
	claims.ID = o.JWTID
}

func WithJWTID(jwtID string) ClaimsOption {
	return jwtIDOption{JWTID: jwtID}
}
