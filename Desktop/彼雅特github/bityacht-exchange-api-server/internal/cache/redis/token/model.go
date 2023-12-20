package token

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/cache/redis"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
)

var _ encoding.BinaryMarshaler = (*Model)(nil)
var _ encoding.BinaryUnmarshaler = (*Model)(nil)

// Basic Model for Managers And Users
type Model struct {
	Token                 string    `json:"Token"` // JWT ID
	RefreshToken          string    `json:"RT"`
	RefreshTokenExpiredAt time.Time `json:"RTEXP"`
	LoginAt               time.Time `json:"LA"`
}

func (m Model) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Model) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Model) Setup(token string, loginAt time.Time) {
	m.Token = token
	m.RefreshToken = uuid.NewString()
	m.RefreshTokenExpiredAt = time.Now().Add(configs.Config.JWT.RefreshTokenLifetime)
	m.LoginAt = loginAt
}

func (m Model) GetToken() string {
	return m.Token
}

type AttemptLogin struct {
	Count  int64 `json:"c"`
	Locked bool  `json:"l"`
}

func (a AttemptLogin) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AttemptLogin) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

type IRecord interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	GetToken() string
}

func CheckLoginLock(ctx context.Context, key string) *errpkg.Error {
	var record AttemptLogin

	if err := redis.Client().Get(ctx, key).Scan(&record); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	} else if record.Locked {
		return &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeTemporaryForbidden}
	}

	return nil
}

func UpdateAttemptLogin(ctx context.Context, key string, loggedIn bool) *errpkg.Error {
	const (
		attemptLoginExp = 30 * 60 // 30 min
		lockLoginExp    = 60 * 60 // 60 min
	)

	_, err := updateAttemptLoginScript.Run(ctx, redis.Client(), []string{key}, loggedIn, attemptLoginExp, lockLoginExp).Int()
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return nil
}

func Login(ctx context.Context, key string, record IRecord) *errpkg.Error {
	if err := redis.Client().Set(ctx, key, record, configs.Config.JWT.RefreshTokenLifetime).Err(); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return nil
}

func getKeyForPreverify(payload jwt.PreverifyPayload) string {
	return fmt.Sprintf(redis.PreverificationKeyFormat, payload.AccountType, payload.ID, payload.Usage)
}

func SetForPreverify(ctx context.Context, claims jwt.PreverifyClaims) *errpkg.Error {
	if claims.RegisteredClaims.ExpiresAt == nil || claims.RegisteredClaims.IssuedAt == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJWTBadPayload, Err: errors.New("bad exp or iat")}
	}

	expiration := claims.RegisteredClaims.ExpiresAt.Sub(claims.RegisteredClaims.IssuedAt.Time)
	if err := redis.Client().Set(ctx, getKeyForPreverify(claims.PreverifyPayload), Model{Token: claims.RegisteredClaims.ID}, expiration).Err(); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return nil
}

func CheckForPreverify(ctx context.Context, claims jwt.PreverifyClaims) *errpkg.Error {
	result, err := logoutScript.Run(ctx, redis.Client(), []string{getKeyForPreverify(claims.PreverifyPayload)}, claims.RegisteredClaims.ID).Int()
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	switch result {
	case 0:
		return nil
	case -1:
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeVerificationCodeExpired}
	case -2:
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadVerificationCode}
	}

	return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedisBadScript, Err: fmt.Errorf("bad result %+v", result)}
}

func ForceLogout(ctx context.Context, key string) *errpkg.Error {
	if err := redis.Client().Del(ctx, key).Err(); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return nil
}

func Refresh(ctx context.Context, key string, refreshToken string, record IRecord) *errpkg.Error {
	result, err := refreshTokenScript.Run(ctx, redis.Client(), []string{key}, refreshToken, record).Result()
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	switch result := result.(type) {
	case int:
		switch result {
		case -1:
			return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeBadAuthorizationToken, Err: errors.New("not found")}
		case -2:
			return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeBadAuthorizationToken, Err: errors.New("bad token")}
		}
	case string:
		if err := json.Unmarshal([]byte(result), &record); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedisBadScript, Err: fmt.Errorf("unmarshal result failed err: %+v, result: %q", err, result)}
		}

		return nil
	}

	return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedisBadScript, Err: fmt.Errorf("bad result %+v, type %+v", result, reflect.TypeOf(result))}
}

func Logout(ctx context.Context, key string, jti string) *errpkg.Error {
	result, err := logoutScript.Run(ctx, redis.Client(), []string{key}, jti).Int()
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	switch result {
	case 0:
		return nil
	case -1:
		return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeBadAuthorizationToken, Err: errors.New("not found")}
	case -2:
		return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeBadAuthorizationToken, Err: errors.New("bad token")}
	}

	return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedisBadScript, Err: fmt.Errorf("bad result %+v", result)}
}

func Validate(ctx context.Context, key, jti string, record IRecord) *errpkg.Error {
	if err := redis.Client().Get(ctx, key).Scan(record); err != nil {
		if errors.Is(err, redis.Nil) {
			return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeBadAuthorizationToken}
		}
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	} else if record.GetToken() != jti {
		return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeJWTRevoked}
	}

	return nil
}
