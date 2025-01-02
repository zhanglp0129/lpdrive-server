package dbutil

import "strings"

func IsDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate")
}
