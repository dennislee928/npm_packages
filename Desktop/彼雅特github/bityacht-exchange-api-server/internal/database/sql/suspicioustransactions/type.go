package suspicioustransactions

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
)

type RiskReviewResult int32

const (
	RiskReviewResultPending = iota + 1
	RiskReviewResultNonSuspicious
	RiskReviewResultSuspicious
)

func (rrr RiskReviewResult) Validate() *errpkg.Error {
	if rrr < RiskReviewResultPending || rrr > RiskReviewResultSuspicious {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad risk review result")}
	}

	return nil
}

type DedicatedReviewResult int32

const (
	DedicatedReviewResultPending = iota + 1
	DedicatedReviewResultPass
	DedicatedReviewResultReject
)

func (drr DedicatedReviewResult) Validate() *errpkg.Error {
	if drr < DedicatedReviewResultPending || drr > DedicatedReviewResultReject {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad dedicated review result")}
	}

	return nil
}

func (drr DedicatedReviewResult) Chinese() string {
	switch drr {
	case DedicatedReviewResultPending:
		return "待審核"
	case DedicatedReviewResultPass:
		return "通過"
	case DedicatedReviewResultReject:
		return "駁回"
	}

	return "未知錯誤"
}

var _ sql.Scanner = (*Files)(nil)
var _ driver.Valuer = (*Files)(nil)

type Files []string

// Scan scan value into Jsonb, implements sql.Scanner interface
func (f *Files) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("suspicious_transactions.Files: bad type")
	}

	return json.Unmarshal(bytes, f)
}

// Value return json value, implement driver.Valuer interface
func (f Files) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f Files) Add(filename string) Files {
	var (
		exist    bool
		newFiles = make([]string, 0, len(f)+1)
	)

	for _, v := range f {
		newFiles = append(newFiles, v)
		if v == filename {
			exist = true
		}
	}

	if !exist {
		newFiles = append(newFiles, filename)
	}

	return newFiles
}

func (f Files) Remove(filename string) Files {
	newFiles := make([]string, 0, len(f))

	for _, v := range f {
		if v != filename {
			newFiles = append(newFiles, v)
		}
	}

	return newFiles
}
