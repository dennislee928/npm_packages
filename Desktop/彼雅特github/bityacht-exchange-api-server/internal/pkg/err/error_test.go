package errpkg

import (
	"testing"
)

func TestError_CodeEqualTo(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code Code
		want bool
	}{
		{"nil 0 -> true", nil, 0, true},
		{"nil err -> false", nil, CodeRecordNotFound, false},
		{"err 0 -> false", &Error{Code: CodeCallMaxAPI}, 0, false},
		{"err err -> true", &Error{Code: CodeScheduleJobKeyDuplicated}, CodeScheduleJobKeyDuplicated, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.CodeEqualTo(tt.code); got != tt.want {
				t.Errorf("Error.CodeEqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
