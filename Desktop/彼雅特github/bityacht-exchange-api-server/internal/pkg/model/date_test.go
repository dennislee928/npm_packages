package modelpkg

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"testing"
)

func TestDate_ToString(t *testing.T) {
	tests := []struct {
		name         string
		jsonInput    string
		omitZeroTime bool
		want         string
		wantErrCode  errpkg.Code
	}{
		{"", "0000/01/01", true, "", 0},            // Zero Time
		{"", "0000/01/01", false, "0001/01/01", 0}, // Zero Time
		{"", "0001/01/01", true, "", 0},            // Zero Time
		{"", "0001/01/01", false, "0001/01/01", 0}, // Zero Time
		{"", "2023/06/30", true, "2023/06/30", 0},
		{"", "2023/06/30", false, "2023/06/30", 0},
		{"", "9999/12/31", true, "9999/12/31", 0},
		{"", "9999/12/31", false, "9999/12/31", 0},
		{"", "9999-12-31", false, "", errpkg.CodeJSONUnmarshal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date

			if err := d.UnmarshalJSON([]byte(tt.jsonInput)); (err == nil && tt.wantErrCode != 0) || (err != nil && tt.wantErrCode != errpkg.CodeJSONUnmarshal) {
				t.Errorf("Date.UnmarshalJSON() = %v, want err code %v", err, errpkg.CodeJSONUnmarshal)
			} else if tt.wantErrCode == 0 {
				if got := d.ToString(tt.omitZeroTime); got != tt.want {
					t.Errorf("Date.ToString() = %v, want %v", got, tt.want)
				} else {
					t.Log(d)
				}
			}
		})
	}
}
