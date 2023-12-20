package passwordpkg

import (
	"testing"
)

func TestStrengthValidate(t *testing.T) {
	inValidPasswords := []string{"abc", "123", "abc123", "abcd1234", "aA12345", "Aa12345", "acxZ123456789", "12345678", "a1234567", "A1234567"}
	validPasswords := []string{"Az123456", "Qs123456", "123456789z1A", "Q1234567890a", "1q234567P", "!@#$A12s"}

	for _, password := range inValidPasswords {
		if err := StrengthValidate(password); err == nil {
			t.Errorf("bad IsValid function, password: %q cannot be valid.", password)
		}
	}

	for _, password := range validPasswords {
		if err := StrengthValidate(password); err != nil {
			t.Errorf("bad IsValid function, password: %q should be valid.", password)
		}
	}
}
