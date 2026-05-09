package utils

import "strings"

func FormatOutput(raw string) string {
	paragraphs := strings.Split(raw, "\n\n")
	for i, p := range paragraphs {
		paragraphs[i] = strings.Join(strings.Split(p, "\n"), " ")
	}
	return strings.Join(paragraphs, "\n\n")
}

