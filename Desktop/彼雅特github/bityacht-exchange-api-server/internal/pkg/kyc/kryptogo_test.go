package kyc

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql/countries"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func Test_kryptoGO_InitIDV(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	logger.Init()
	Init()

	type args struct {
		country    countries.Country
		usersID    int64
		idvsID     int64
		ddsID      int64
		name       string
		birthDate  string
		nationalID string
	}
	tests := []struct {
		name string
		args args
	}{
		{"init idv", args{countries.Country{Code: "TWN", Locale: "zh-HK"}, 0, 0, 0, "某測試", "2006-01-02", "A123456789"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			birthDate, err := time.Parse(time.DateOnly, tt.args.birthDate)
			if err != nil {
				t.Errorf("bad birthDate in %q, err: %v", tt.name, err)
				return
			}

			rawResp := new(bytes.Buffer)
			if resp, err := KryptoGO.InitIDV(ctx, tt.args.country, tt.args.usersID, tt.args.idvsID, tt.args.ddsID, tt.args.name, modelpkg.Date{Time: birthDate}, tt.args.nationalID, rawResp); err != nil {
				t.Errorf("kryptoGO.InitIDV() got err = %v, want nil", err)
			} else {
				t.Log(resp, rawResp.String())
			}
		})
	}
}

func Test_kryptoGO_CreateTasks(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	t.Run("Test_kryptoGO_CreateTasks", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tasks := []CreateTasksParams{{
			Name:      "葉子",
			FromIDVID: 0}}

		if resp, err := KryptoGO.CreateTasks(ctx, tasks); err != nil {
			t.Errorf("kryptoGO.InitIDV() got err = %v, want nil", err)
		} else {
			t.Log(resp)
		}
	})
}
