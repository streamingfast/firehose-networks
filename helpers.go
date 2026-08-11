package networks

import (
	"slices"
	"strings"
)

func ptr[T any](v T) *T { return &v }

// mergeEndpoints returns the concatenation of preferred and existing endpoints, in that
// order and without duplicates. Both input slices are left untouched.
func mergeEndpoints(preferred, existing []string) []string {
	if len(preferred) == 0 {
		return existing
	}

	merged := make([]string, 0, len(preferred)+len(existing))
	for _, endpoint := range slices.Concat(preferred, existing) {
		if !slices.Contains(merged, endpoint) {
			merged = append(merged, endpoint)
		}
	}

	return merged
}

// some chains have noisy '0x' prefixes, some don't, normalize it without 0x
func nox(s string) string {
	return strings.TrimPrefix(s, "0x")
}
