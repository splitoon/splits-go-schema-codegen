package graphql

import "strings"

func initManualPart(manualParts []string) func() string {
	index := 0
	return func() string {
		if manualParts == nil {
			return ""
		}
		if index >= len(manualParts) {
			return ""
		}
		index++
		return strings.Replace(manualParts[index-1], "%", "%%", -1)
	}
}
