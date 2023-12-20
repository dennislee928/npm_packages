package passwordpkg

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"net/http"
)

func StrengthValidate(password string) *errpkg.Error {
	if passwordLen := len(password); passwordLen < 8 || passwordLen > 16 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadPasswordStrength, Err: errors.New("bad length")}
	}

	var hasLowerLetter, hasCapitalLetter bool
	for _, char := range password {
		if !hasLowerLetter && char >= 'a' && char <= 'z' {
			hasLowerLetter = true
		}
		if !hasCapitalLetter && char >= 'A' && char <= 'Z' {
			hasCapitalLetter = true
		}

		if hasLowerLetter && hasCapitalLetter {
			return nil
		}
	}

	return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadPasswordStrength, Err: errors.New("at least one lower and one capital letter")}
}
