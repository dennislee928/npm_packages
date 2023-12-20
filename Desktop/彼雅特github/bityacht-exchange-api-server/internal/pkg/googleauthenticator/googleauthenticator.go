package googleauthenticator

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bytes"
	"crypto/hmac"
	"crypto/sha1" // #nosec G505
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	keyUriType = "totp" // Valid types are hotp and totp, to distinguish whether the key will be used for counter-based HOTP or for TOTP.
	issuer     = "BitYacht"
	label      = issuer
)

// Ref: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func GenerateSecret(userEmail string) (string, []byte, *errpkg.Error) {
	secret := rand.Base32String(32)

	qrCodeInfo := fmt.Sprintf("otpauth://%s/%s:%s?secret=%s&issuer=%s", keyUriType, label, userEmail, secret, issuer)
	qrCode, err := qrcode.Encode(qrCodeInfo, qrcode.Medium, 256)
	if err != nil {
		return "", nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGenQRCode, Err: err}
	}

	return secret, qrCode, nil
}

// Ref RFC 4226: https://datatracker.ietf.org/doc/html/rfc4226#section-5.3
func generateHOTP(key string, counter int64) (string, *errpkg.Error) {
	byteKey, err := base32.StdEncoding.DecodeString(key)
	if err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadBase32String, Err: err}
	}

	byteCounter := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCounter, uint64(counter))

	// We can describe the operations in 3 distinct steps:
	// Step 1: Generate an HMAC-SHA-1 value Let HS = HMAC-SHA-1(K,C)  // HS is a 20-byte string

	hashFunc := hmac.New(sha1.New, byteKey)
	if _, err := hashFunc.Write(byteCounter); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGenerateHOTP, Err: err}
	}
	hs := hashFunc.Sum(nil)

	// Step 2: Generate a 4-byte string (Dynamic Truncation)
	// Let Sbits = DT(HS)   //  DT, defined below, returns a 31-bit string

	// DT(String) // String = String[0]...String[19]
	// Let OffsetBits be the low-order 4 bits of String[19]
	// Offset = StToNum(OffsetBits) // 0 <= OffSet <= 15
	// Let P = String[OffSet]...String[OffSet+3]
	// Return the Last 31 bits of P

	offset := hs[19] & 0b00001111
	hs[offset] = hs[offset] & 0b01111111
	sBitsReader := bytes.NewReader(hs[offset : offset+4])

	// Step 3: Compute an HOTP value
	// Let Snum  = StToNum(Sbits)   // Convert S to a number in 0...2^{31}-1
	// Return D = Snum mod 10^Digit //  D is a number in the range 0...10^{Digit}-1

	var sNum int32
	if err := binary.Read(sBitsReader, binary.BigEndian, &sNum); err != nil {
		return "", &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGenerateHOTP, Err: err}
	}
	sNum %= 1000000 //! 6 Digits, if 8 -> 10^8

	return fmt.Sprintf("%06d", sNum), nil
}

func VerifyTOTP(secret string, unixTime int64, otp string) *errpkg.Error {
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format#period
	// Default: 30
	const period = 30

	if unixTime == 0 {
		unixTime = time.Now().Unix()
	}

	otpAns, err := generateHOTP(secret, unixTime/period)
	if err != nil {
		return err
	}

	if otp != otpAns {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadVerificationCode}
	}

	return nil
}
