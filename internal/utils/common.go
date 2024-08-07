package utils

import "strings"

func ReplacePlaceholders(template string, placeholders map[string]string) string {
	for placeholder, value := range placeholders {
		template = strings.Replace(template, ":"+placeholder, value, -1)
	}

	return template
}
