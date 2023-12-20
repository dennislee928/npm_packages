package levellimits

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

func TestGetAllLevelLimit(t *testing.T) {
	viper.AddConfigPath("../../../../configs")
	logger.Init()
	configs.Init()

	t.Run("getAllLevelLimit", func(t *testing.T) {
		levelLimits, err := GetAllLevelLimit()
		if err != nil {
			t.Errorf("GetAllLevelLimit() got err = %v", err)
		}

		levelLimitMap := make(map[int64]LevelLimit)
		for _, levelLimit := range levelLimits {
			levelLimitMap[int64(levelLimit.Type)<<32|int64(levelLimit.Level)] = levelLimit
		}

		nd := decimal.New
		twdTests := []struct {
			usersType           int32
			usersLevel          int32
			depositAmount       decimal.Decimal
			withdrawAmount      decimal.Decimal
			accOfDay            decimal.Decimal
			accOfMonth          decimal.Decimal
			wantDepositErrCode  errpkg.Code
			wantWithdrawErrCode errpkg.Code
		}{
			{1, 0, nd(1234, -2), nd(4567, -2), nd(0, 0), nd(0, 0), 4033, 4033},
			{1, 0, nd(0, 0), nd(0, 0), nd(0, 0), nd(0, 0), 0, 0},
			{1, 1, nd(1234, -2), nd(4567, -2), nd(0, 0), nd(0, 0), 4033, 4033},
			{1, 1, nd(0, 0), nd(0, 0), nd(0, 0), nd(0, 0), 0, 0},
			{1, 2, nd(1234, 0), nd(15, 0), nd(0, 0), nd(0, 0), 0, 0},
			{1, 2, nd(0, 0), nd(13, 0), nd(15, 4), nd(0, 0), 0, 4033},
			{1, 2, nd(1, 0), nd(1, 5), nd(15, 4), nd(0, 0), 4033, 0},
			{1, 2, nd(-1, 0), nd(100001, 0), nd(15, 4), nd(0, 0), 4033, 4033},
		}

		for _, tt := range twdTests {
			t.Run(fmt.Sprintf("%d-%d", tt.usersType, tt.usersLevel), func(t *testing.T) {
				levelLimit, ok := levelLimitMap[int64(tt.usersType)<<32|int64(tt.usersLevel)]
				if !ok {
					t.Error("level limit not found")
					return
				}

				if err := levelLimit.TWD.Deposit.Validate(tt.depositAmount, tt.accOfDay, tt.accOfMonth); (err == nil) != (tt.wantDepositErrCode == 0) {
					t.Logf("%+v\n", levelLimit.TWD)
					t.Errorf("Deposit.Validate() err = %v, want %v", err, tt.wantDepositErrCode)
				} else if err = levelLimit.TWD.Withdraw.Validate(tt.withdrawAmount, tt.accOfDay, tt.accOfMonth); (err == nil) != (tt.wantWithdrawErrCode == 0) {
					t.Logf("%+v\n", levelLimit.TWD)
					t.Errorf("Withdraw.Validate() err = %v, want %v", err, tt.wantWithdrawErrCode)
				}
			})
		}
	})
}
