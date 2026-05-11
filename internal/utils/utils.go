package utils

import (
	"strings"
	"time"
)

type SenderType int

const (
    SelfType SenderType = iota
    ContactType
)

type RawMessage struct {
    Sender   SenderType  `json:"sender"`
    Message  string      `json:"message"`
    Time     time.Time   `json:"time"`
}

func FormatOutput(raw string) string {
	paragraphs := strings.Split(raw, "\n\n")
	for i, p := range paragraphs {
		paragraphs[i] = strings.Join(strings.Split(p, "\n"), " ")
	}
	return strings.Join(paragraphs, "\n\n")
}

