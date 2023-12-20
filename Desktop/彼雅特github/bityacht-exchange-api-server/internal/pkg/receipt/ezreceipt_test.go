package receipt

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func testEZFixture() {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	configs.Config.Receipt.APIMode = int(APIModeTest)
	MustInit()
}

func makeTestingOrder(t *testing.T) (*CreateOrderRespItem, func()) {
	t.Helper()

	ctx := context.Background()
	orderNo := rand.LetterAndNumberString(30)

	order, err := EZ.CreateOrder(ctx, CreateOrderPayload{
		OrderNo: orderNo,
		Title:   "test",
		ProdList: []ProdItem{
			{
				Title:  "test_prod",
				Sales:  7,
				Qty:    1,
				IncTax: true,
			},
		},
		Confirm: true,
	})
	if err != nil {
		t.Errorf("ezreceipt.CreateOrder() error = %v", err)
		return nil, nil
	}

	tearDown := func() {
		if err := EZ.VoidOrder(ctx, order.OrderID); err != nil {
			t.Errorf("ezreceipt.VoidOrder() error = %v", err)
			return
		}
	}

	return order, tearDown
}

func Test_loginPwd(t *testing.T) {
	type fields struct {
		apiAcc string
		apiPwd string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"test",
			fields{
				apiAcc: "test",
				apiPwd: "5678",
			},
			"e6bae8c0bfbd44d29944eab05ec4e08a807313b0",
			false,
		},
		{
			"empty acc",
			fields{
				apiAcc: "",
				apiPwd: "5678",
			},
			"",
			true,
		},
		{
			"empty pwd",
			fields{
				apiAcc: "test",
				apiPwd: "",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loginPwd(tt.fields.apiAcc, tt.fields.apiPwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ezreceipt.loginPwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ezreceipt.loginPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ezrecepit_Login(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := EZ.Login(ctx); err != nil {
		t.Errorf("ezrecepit.Login() error = %v", err)
		return
	}

	if EZ.(*ezreceipt).token.token == "" {
		t.Error("ezrecepit.Login() token is empty")
		return
	}

	if EZ.(*ezreceipt).token.validTo.IsZero() {
		t.Error("ezrecepit.Login() expireAt is zero")
		return
	}
}

func Test_ezreceipt_InvNumberList(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := EZ.InvNumberList(ctx)
	if err != nil {
		t.Errorf("ezreceipt.InvNumberList() error = %v", err)
		return
	}

	t.Logf("ezreceipt.InvNumberList() data = %+v", data)
}

func TestMultipleReqs(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const goroutineCount = 30
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			if err := EZ.Login(ctx); err != nil {
				t.Errorf("ezreceipt.Login() error = %v", err)
				return
			}

			if _, err := EZ.InvNumberList(ctx); err != nil {
				t.Errorf("ezreceipt.InvNumberList() error = %v", err)
			}

			if err := EZ.Login(ctx); err != nil {
				t.Errorf("ezreceipt.Login() error = %v", err)
				return
			}
		}()
	}

	wg.Wait()
}

func Test_ezreceipt_CreateOrder(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prodItem := ProdItem{
		Title:  "test_prod",
		Sales:  7,
		Qty:    1,
		IncTax: true,
	}

	order, err := EZ.CreateOrder(ctx, CreateOrderPayload{
		Title:    "test",
		ProdList: []ProdItem{prodItem},
		Confirm:  true,
	})
	if err != nil {
		t.Errorf("ezreceipt.CreateOrder() error = %v", err)
		return
	}

	if len(order.ProdList) != 1 {
		t.Errorf("ezreceipt.CreateOrder() resp pord list size not match, expect 1, got %d", len(order.ProdList))
		return
	}

	if order.ProdList[0].Title != prodItem.Title {
		t.Errorf("ezreceipt.CreateOrder() resp pord title not match, expect test_prod, got %s", order.ProdList[0].Title)
		return
	}

	if order.ProdList[0].Qty != prodItem.Qty {
		t.Errorf("ezreceipt.CreateOrder() resp pord qty not match, expect 1, got %d", order.ProdList[0].Qty)
		return
	}

	t.Logf("ezreceipt.CreateOrder() data = %+v", order)

	// tear down
	if err := EZ.VoidOrder(ctx, order.OrderID); err != nil {
		t.Errorf("ezreceipt.VoidOrder() error = %v", err)
		return
	}
}

func Test_ezreceipt_VoidOrder(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderNo := rand.LetterAndNumberString(30)

	order, err := EZ.CreateOrder(ctx, CreateOrderPayload{
		OrderNo: orderNo,
		Title:   "test",
		ProdList: []ProdItem{
			{
				Title:  "test_prod",
				Sales:  7,
				Qty:    1,
				IncTax: true,
			},
		},
		Confirm: true,
	})
	if err != nil {
		t.Errorf("ezreceipt.CreateOrder() error = %v", err)
		return
	}

	if err := EZ.VoidOrder(ctx, order.OrderID); err != nil {
		t.Errorf("ezreceipt.VoidOrder() error = %v", err)
		return
	}
}

func Test_ezreceipt_SetCarrier(t *testing.T) {
	testEZFixture()

	order, tearDown := makeTestingOrder(t)
	defer tearDown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := EZ.SetCarrier(ctx, order.OrderID, SetCarrierPayload{
		CarrierType: 10,
	})
	if err != nil {
		t.Errorf("ezreceipt.SetCarrier() error = %v", err)
		return
	}

	t.Logf("ezreceipt.SetCarrier() data = %+v", data)
}

func Test_ezreceipt_CheckMobileCode(t *testing.T) {
	testEZFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mobileCode := "/" + strings.ToUpper(rand.LetterAndNumberString(7))

	if err := EZ.CheckMobileCode(ctx, mobileCode); err != nil {
		if err.Code == errpkg.CodeMobileBarcodeNotExist {
			t.Logf("%q is not exist", mobileCode)
		} else {
			t.Errorf("ezreceipt.CheckMobileCode() error = %v", err)
		}
	} else {
		t.Logf("%q is exist", mobileCode)
	}
}
