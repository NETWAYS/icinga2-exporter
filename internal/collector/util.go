package collector

import (
	"regexp"
	"strings"
)

var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9_:]`)

func ensureValidMetricName(s string) string {
	return invalidChars.ReplaceAllString(s, "_")
}

func isPerfdataMetric(s string) bool {
	if strings.HasPrefix(s, "data") {
		return true
	}

	if strings.HasPrefix(s, "work") {
		return true
	}

	return false
}
