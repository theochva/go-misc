package strext

import "strings"

// TrimArray - trims whitespace from the array of strings.
//
// If removeEmpty=true, then any trimmed values are not included in the result array.
//
// If removeEmpty=false, then empty values may be included in the result array.
func TrimArray(values []string, removeEmpty bool) []string {
	if removeEmpty {
		result := []string{}

		for _, value := range values {
			if trimmedValue := strings.TrimSpace(value); trimmedValue != "" {
				result = append(result, trimmedValue)
			}
		}

		return result
	}

	result := make([]string, 0, len(values))

	for _, value := range values {
		result = append(result, strings.TrimSpace(value))
	}
	return result
}

// TrimLines - for each line in the string, it trims whitespace. The new
// "trimmed" string is returned. Empty lines remain in the resulting string.
func TrimLines(value string) string {
	return strings.Join(TrimArray(strings.Split(value, "\n"), false), "\n")
}
