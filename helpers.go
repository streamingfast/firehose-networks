package networks

import "strings"

func ptr[T any](v T) *T { return &v }

// some chains have noisy '0x' prefixes, some don't, normalize it without 0x
func nox(s string) string {
	return strings.TrimPrefix(s, "0x")
}
