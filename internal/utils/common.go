package utils

import "strings"

// ReplacePlaceholders Replaces placeholders in a template string with actual values.
//
// It takes a template string and a map of placeholders as input, where each key in the map is a placeholder and its corresponding value is the actual value to replace the placeholder.
// Returns the modified template string with all placeholders replaced.
func ReplacePlaceholders(template string, placeholders map[string]string) string {
	for placeholder, value := range placeholders {
		template = strings.Replace(template, ":"+placeholder, value, -1)
	}

	return template
}
