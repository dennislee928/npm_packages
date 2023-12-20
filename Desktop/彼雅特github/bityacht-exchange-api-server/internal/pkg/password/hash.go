package passwordpkg

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt the password
func Encrypt(input string) (string, *errpkg.Error) {
	if hashPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeEncryption, Err: err}
	} else {
		return string(hashPassword), nil
	}
}

// Validate the password
func Validate(encryptedPassword string, password string) *errpkg.Error {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeUnauthorized}
		}
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeEncryption, Err: err}
	}

	return nil
}
