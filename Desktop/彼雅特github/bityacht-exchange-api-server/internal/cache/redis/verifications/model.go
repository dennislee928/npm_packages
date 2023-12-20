package verifications

import (
	"bityacht-exchange-api-server/internal/cache/redis"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/googleauthenticator"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	UsageVerifyEmail       = "email"
	UsageForgotPassword    = "forgotPW"
	UsageLogin2FA          = "login2FA"
	UsageVerifyPhoneForIDV = "phoneIDV"
	UsageWithdraw          = "withdraw"
	UsageUpdateWithdraw2FA = "updateWithdraw2FA"
)

func getKey(jwtType jwt.Type, id int64, usage string) string {
	return fmt.Sprintf(redis.VerificationKeyFormat, int8(jwtType), id, usage)
}

func IssueVerificationCode(ctx context.Context, jwtType jwt.Type, id int64, usage string, code string, expiration time.Duration) *errpkg.Error {
	if err := redis.Client().Set(ctx, getKey(jwtType, id, usage), code, expiration).Err(); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return nil
}

func Verify(ctx context.Context, jwtType jwt.Type, id int64, usage string, code string) *errpkg.Error {
	result, err := verifyScript.Run(ctx, redis.Client(), []string{getKey(jwtType, id, usage)}, code).Int()
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

type Withdraw struct {
	ID        int64               `redis:"id"`
	EmailCode string              `redis:"emailCode"`
	SMSCode   string              `redis:"smsCode"`
	GASecret  string              `redis:"gaCode"`
	Currency  wallet.CurrencyType `redis:"currency"`
	Mainnet   wallet.Mainnet      `redis:"mainnet"`
	Address   string              `redis:"address"`
	Amount    string              `redis:"amount"`
}

func IssueWithdrawVerification(ctx context.Context, data Withdraw, expired time.Duration) (string, *errpkg.Error) {
	key := uuid.NewString()

	pipe := redis.Client().TxPipeline()
	if err := pipe.HSet(ctx, key, data).Err(); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	if err := pipe.Expire(ctx, key, expired).Err(); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	return key, nil
}

func GetWithdrawVerification(ctx context.Context, reqUnixTime, id int64, key, emailCode, smsCode, gaCode string) (*Withdraw, *errpkg.Error) {
	result := redis.Client().HGetAll(ctx, key)
	if err := result.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeVerificationCodeExpired}
		}

		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	var data Withdraw
	if err := result.Scan(&data); err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRedis, Err: err}
	}

	if data.ID != id || data.EmailCode != emailCode || data.SMSCode != smsCode {
		return nil, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadVerificationCode}
	}
	if data.GASecret != "" {
		if err := googleauthenticator.VerifyTOTP(data.GASecret, reqUnixTime, gaCode); err != nil {
			return nil, err
		}
	}

	if err := redis.Client().Del(ctx, key).Err(); err != nil {
		logger.Logger.Warn().Err(err).Msg("failed to delete verification key")
	}

	return &data, nil
}
