package users

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"math/bits"
	"net/http"
	"strconv"
	"strings"
)

const inviteCodeBase = 36
const inviteCodeLength = 7
const inviteCodePaddingChar = "="

func GetInviteCodeByID(id int64) string {
	// ID(in base 36) + Checksum
	code := strconv.FormatInt(id, inviteCodeBase) + strconv.FormatInt(int64(bits.OnesCount64(uint64(id))%inviteCodeBase), inviteCodeBase)
	code = strings.ToUpper(code)

	if codeLength := len(code); codeLength < inviteCodeLength {
		code += strings.Repeat(inviteCodePaddingChar, inviteCodeLength-codeLength)
	}

	return code
}

func ParseInviteCode(code string) (int64, *errpkg.Error) {
	if len(code) != inviteCodeLength {
		return 0, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadInviteCode, Err: errors.New("bad length")}
	}
	code = strings.TrimRight(code, inviteCodePaddingChar)

	codeLength := len(code)
	if codeLength < 2 {
		return 0, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadInviteCode, Err: errors.New("bad valid length")}
	}

	id, err := strconv.ParseInt(code[:codeLength-1], inviteCodeBase, 64)
	if err != nil {
		return 0, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadInviteCode, Err: err}
	}

	checksum, err := strconv.ParseInt(code[codeLength-1:], inviteCodeBase, 64)
	if err != nil {
		return 0, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadInviteCode, Err: err}
	}

	if int(checksum) != bits.OnesCount64(uint64(id))%inviteCodeBase {
		return 0, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadInviteCode, Err: errors.New("bad checksum")}
	}

	return id, nil
}
