package validation

import (
	"fmt"
	"unicode/utf8"
)

// CheckLimitNameLen checking length of name, limit in 1000 symbols
func CheckLimitNameLen(name string) error {
	if utf8.RuneCountInString(name) > 1000 {
		return fmt.Errorf("limit up to 1000 sign for name of key")
	}
	return nil
}
