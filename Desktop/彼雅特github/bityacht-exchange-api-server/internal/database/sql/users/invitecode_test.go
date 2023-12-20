package users

import (
	"testing"
)

func TestInviteCode(t *testing.T) {
	for id := int64(10000000); id <= 99999999; id++ {
		code := GetInviteCodeByID(id)

		if parsedID, err := ParseInviteCode(code); err != nil {
			t.Error(id, code, err)
			break
		} else if parsedID != id {
			t.Errorf("[id] %q != [parsed id] %q, code = %q\n", id, parsedID, code)
			break
		}
	}
}

// Ref: https://math.tools/calculator/base/10-36	https://go.dev/play/p/apdQUL-qmSq
func TestGetInviteCodeByID(t *testing.T) {
	testCases := []struct {
		ID   int64
		Code string
	}{
		{99999999, "1NJCHRJ"},
		{10000000, "5YC1S8="},
	}

	for _, testCase := range testCases {
		code := GetInviteCodeByID(testCase.ID)

		if code != testCase.Code {
			t.Errorf("ID: %d, GetInviteCodeByID: %q != [Answer] %q\n", testCase.ID, code, testCase.Code)
		}
	}
}
