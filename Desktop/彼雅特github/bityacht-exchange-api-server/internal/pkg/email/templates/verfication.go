package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

// #nosec G101
const (
	verificationLogin2FATitle    = "[BitYacht] 登入電子信箱驗證碼"
	verificationLogin2FATemplate = "login-2fa.html.tmpl"

	verificationRegisterTitle    = "[BitYacht] 會員註冊信箱驗證確認信"
	verificationRegisterTemplate = "register.html.tmpl"

	verificationResetPasswordTitle    = "[BitYacht] 重設密碼驗證信"
	verificationResetPasswordTemplate = "reset-password.html.tmpl"

	verificationUpdateWithdraw2FATitle    = "[BitYacht] 變更雙重驗證確認信"
	verificationUpdateWithdraw2FATemplate = "update-withdraw-2fa.html.tmpl"

	verificationWithdraw2FATitle    = "[BitYacht] 提領驗證碼"
	verificationWithdraw2FATemplate = "withdraw-2fa.html.tmpl"
)

type VerificationPayload struct {
	Code         string
	CodeLifeTime string
}

type UpdateWithdraw2FAPayload struct {
	VerificationPayload

	Action string
}

func ExecVerificationLogin2FA(payload VerificationPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(verificationLogin2FATemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return verificationLogin2FATitle, content, nil
}

func ExecVerificationRegister(payload VerificationPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(verificationRegisterTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return verificationRegisterTitle, content, nil
}

func ExecVerificationResetPassword(payload VerificationPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(verificationResetPasswordTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return verificationResetPasswordTitle, content, nil
}

func ExecVerificationUpdateWithdraw2FA(payload UpdateWithdraw2FAPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(verificationUpdateWithdraw2FATemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return verificationUpdateWithdraw2FATitle, content, nil
}

func ExecVerificationWithdraw2FA(payload VerificationPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(verificationWithdraw2FATemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return verificationWithdraw2FATitle, content, nil
}
