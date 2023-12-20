package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

const (
	idvFailedTitle    = "[BitYacht] 身分驗證審核未通過"
	idvFailedTemplate = "failed.html.tmpl"

	idvOnGoingTitle    = "[BitYacht] 身分驗證資料審核中"
	idvOnGoingTemplate = "on-going.html.tmpl"

	idvPassTitle    = "[BitYacht] 身分驗證資料已通過" // #nosec G101
	idvPassTemplate = "pass.html.tmpl"       // #nosec G101

	idvRiskTitle    = "[BitYacht] 身分驗證審核未通過"
	idvRiskTemplate = "risk.html.tmpl"
)

func ExecIDVFailed() (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(idvFailedTemplate, nil)
	if err != nil {
		return "", nil, err
	}

	return idvFailedTitle, content, nil
}

type IDVOnGoingPayload struct {
	Time string
}

func ExecIDVOnGoing(payload IDVOnGoingPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(idvOnGoingTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return idvOnGoingTitle, content, nil
}

func ExecIDVPass() (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(idvPassTemplate, nil)
	if err != nil {
		return "", nil, err
	}

	return idvPassTitle, content, nil
}

func ExecIDVRisk() (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(idvRiskTemplate, nil)
	if err != nil {
		return "", nil, err
	}

	return idvRiskTitle, content, nil
}
