package email

import (
	"testing"

	"bityacht-exchange-api-server/configs"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/spf13/viper"
)

func TestSendEmail(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()
	configs.Config.Email.Enable = true

	email := NewEmail()
	email.To = []string{getAddressOfAccount(configs.Config.Email.Account)}
	email.Subject = "Testing email_test.go"
	// email.Text = []byte("Testing email_test.go [Text]")
	email.HTML = []byte("<h1>Testing email_test.go [HTML]</h1>")

	if err := SendMail(email); err != nil {
		t.Error(err)
	}
}

func TestSendEmailWithTemplate(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()
	configs.Config.Email.Enable = true

	type args struct {
		ExecuteFunc func() (string, []byte, *errpkg.Error)
	}
	tests := []struct {
		name string
		args args
	}{
		{"account change disable", args{emailtemplates.ExecAccountChangeDisable}},
		{"account change forzen", args{emailtemplates.ExecAccountChangeForzen}},
		{"idv failed", args{emailtemplates.ExecIDVFailed}},
		{"idv pass", args{emailtemplates.ExecIDVPass}},
		{"idv risk", args{emailtemplates.ExecIDVRisk}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err *errpkg.Error

			email := NewEmail(WithLogo())
			email.To = []string{getAddressOfAccount(configs.Config.Email.Account)}

			if email.Subject, email.HTML, err = tt.args.ExecuteFunc(); err != nil {
				t.Error(err)
			} else if err = SendMail(email); err != nil {
				t.Error(err)
			}
		})
	}
}
