package managers

import (
	"bityacht-exchange-api-server/internal/cache/redis"
	"bityacht-exchange-api-server/internal/cache/redis/token"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"context"
	"encoding"
	"encoding/json"
	"fmt"
)

var _ encoding.BinaryMarshaler = (*Model)(nil)
var _ encoding.BinaryUnmarshaler = (*Model)(nil)

type Model struct {
	token.Model
}

func (m Model) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Model) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func getKeyByID(id int64) string {
	return fmt.Sprintf(redis.ManagersTokenKeyFormat, id)
}

func Login(ctx context.Context, id int64, record Model) *errpkg.Error {
	return token.Login(ctx, getKeyByID(id), &record)
}

func Refresh(ctx context.Context, id int64, refreshToken string, record *Model) *errpkg.Error {
	return token.Refresh(ctx, getKeyByID(id), refreshToken, record)
}

func Logout(ctx context.Context, id int64, jti string) *errpkg.Error {
	return token.Logout(ctx, getKeyByID(id), jti)
}

func ForceLogout(ctx context.Context, id int64) *errpkg.Error {
	return token.ForceLogout(ctx, getKeyByID(id))
}

func Validate(ctx context.Context, id int64, jti string) *errpkg.Error {
	return token.Validate(ctx, getKeyByID(id), jti, &Model{})
}
