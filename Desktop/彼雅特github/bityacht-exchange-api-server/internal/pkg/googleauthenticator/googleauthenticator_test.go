package googleauthenticator

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"encoding/base64"
	"testing"
)

func TestGenerateSecret(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantErrCode errpkg.Code
	}{
		{"gen key", "test@test.bityacht.go", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := GenerateSecret(tt.email)
			if !got2.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("GenerateSecret() got2 = %v, wantErrCode %v", got2, tt.wantErrCode)
			} else {
				// Prview: https://jaredwinick.github.io/base64-image-viewer/
				// Base64 Image String: "data:image/png;base64," + QR Code
				t.Logf("Key: %q\nQR Code: %q", got, base64.StdEncoding.EncodeToString(got1))
			}
		})
	}
}

func TestVerifyTOTP(t *testing.T) {
	// 	Key: "NUXHGDKQAV7MSGKATTFY3MRWSMYDINP7"
	// Prview QR Code: https://jaredwinick.github.io/base64-image-viewer/
	// Base64 Image String: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACvElEQVR42uyZwY3zIBSEx/KBIyXQSWjMih25MdwJJXDkYDG/5jnJbvYvwFgK0kqb+Dss+L2ZeSy+67suuBxJbtz40C95JMDmNn0oFwIy4JJv/kHmOSO45JrbCjD2BJAuYWDzO6Ywit/qUGaunQH6o/XdEkjEenzoD7j5B1cDSLYLAjr8Ern4lcwAk4r8T0WdDVjraRdlCmveg0v19n9vngvYGsoNQxnzFBAAsv6VupMBlxEdS5SI5Skwu2Q1wfLzLroA6kDywQXSWUc2/WC+EsAaXULkgllIRmTCDTN3dAQ44uaTXzy5cqRejbb1fhddABXyC6lWnsKuqk6AX9ATkHWyvmHgmu18I+Cab37BdQDphVkE/BFxVOR1KHhVdR+AOTE8y5gR9lCjhZzhVTBdAMHE4Ia7XwLzaCpXZcSvo+4BQI3aBYAJh1uYpTX/NrXzAUemoxrKnHXYiK4hfsSD/gHlcnsMjNyxB0QMTD+N1wUQKqjxgVzCrj7U+ND8o/QE2PIsdy8XG7POHvF3b54PPOPiwEcZuZqg4YYPU7sAgFAjIhvufrU8+RyCF7/2Azi6dHixV57cofpBLHfpST+AdK7AxgclMVqgEF+6AfTYJb+VmeQiIaZrAN676AFQEpP3KtNmHLcHCjUsP8G7fwCBySVvcqzWOwRENvg6hw4Am3EiN3sTaj3SJTPntXQDKInZncbdL5gwat7VYL78jovnAyQiKYMgM7Mp2oehXACQ8qpgfIPGIA1BNfJj5j0fOO5AVDBK5tAubHSbS0fA6zZJoTcom1vM2cr8x5rPBQ7ViuWZzGXE+uLu154AmxeOaUFaW6M7xse3iF0FkBA3v4cFNGtWFnsXTD+ABGFSesi2i/acKLoB7D8gm6Q4rNQgpqrm0hVwtB4wk1lBB9FtTL9DzgWA7/qurta/AAAA///OLrID+UxD+QAAAABJRU5ErkJggg=="

	type args struct {
		secret   string
		unixTime int64
		otp      string
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		{"ok", args{"NUXHGDKQAV7MSGKATTFY3MRWSMYDINP7", 1696493410, "821198"}, 0},
		{"bad code", args{"NUXHGDKQAV7MSGKATTFY3MRWSMYDINP7", 1696493410, ""}, errpkg.CodeBadVerificationCode},
		{"bad secret", args{"NUXHGDKQAV7MSGKATTFY3MRWSMYDINp7", 1696493410, "821198"}, errpkg.CodeBadBase32String},
		{"ok2", args{"NUXHGDKQAV7MSGKATTFY3MRWSMYDINP7", 1696552633, "804617"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyTOTP(tt.args.secret, tt.args.unixTime, tt.args.otp); !got.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("VerifyTOTP() = %v, want %v", got, tt.wantErrCode)
			}
		})
	}
}
