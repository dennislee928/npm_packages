package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

const (
	accountChangeDisableTitle    = "[BitYacht] 會員帳號已停權"
	accountChangeDisableTemplate = "disable.html.tmpl"

	accountChangeForzenTitle    = "[BitYacht] 會員帳號已凍結"
	accountChangeForzenTemplate = "forzen.html.tmpl"
)

func ExecAccountChangeDisable() (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(accountChangeDisableTemplate, nil)
	if err != nil {
		return "", nil, err
	}

	return accountChangeDisableTitle, content, nil
}

func ExecAccountChangeForzen() (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(accountChangeForzenTemplate, nil)
	if err != nil {
		return "", nil, err
	}

	return accountChangeForzenTitle, content, nil
}
