package emailer

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	BCC         []string `json:"bcc"`
	CC          []string `json:"cc"`
	Subject     string   `json:"subject"`
	HTMLContent string   `json:"htmlContent"`
	TextContent string   `json:"textContent"`
}

func (e Email) ValidationMsg() string {
	if strings.TrimSpace(e.From) == "" {
		return "from field must not be blank"
	}
	if !emailRegex.MatchString(e.From) {
		return fmt.Sprintf("%q is not a valid email", e.From)
	}
	if len(e.To) == 0 {
		return "to field must not be blank"
	}
	for _, s := range e.To {
		if !emailRegex.MatchString(s) {
			return fmt.Sprintf("%q is not a valid email", s)
		}
	}
	if strings.TrimSpace(e.Subject) == "" {
		return "subject field must not be blank"
	}
	if strings.TrimSpace(e.HTMLContent) == "" && strings.TrimSpace(e.TextContent) == "" {
		return "either the htmlContent or textContent field must be filled"
	}
	for _, s := range e.BCC {
		if !emailRegex.MatchString(s) {
			return fmt.Sprintf("%q is not a valid email", s)
		}
	}
	for _, s := range e.CC {
		if !emailRegex.MatchString(s) {
			return fmt.Sprintf("%q is not a valid email", s)
		}
	}
	return ""
}

type Sender interface {
	Send(ctx context.Context, e Email) error
}
