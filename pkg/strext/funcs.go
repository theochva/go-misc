package strext

import "strings"

// Trim - trims whitespace from the array of strings.  Any strings that are empty after
// the trimming are omitted from the output
func Trim(values []string) (result []string) {
	for _, value := range values {
		if trimmedValue := strings.TrimSpace(value); trimmedValue != "" {
			result = append(result, trimmedValue)
		}
	}

	return result
}
