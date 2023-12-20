package invoice

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
	"regexp"
)

// Ref: https://www.ecpay.com.tw/CascadeFAQ/CascadeFAQ_Qa?nID=3531
var mobileBarcodePattern = regexp.MustCompile(`^\/[0-9A-Z\+\-\.]{7}$`)

func ValidateMobileBarcode(barcode string) *errpkg.Error {
	if !mobileBarcodePattern.MatchString(barcode) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadMobileBarcodeFormat}
	}

	return nil
}
