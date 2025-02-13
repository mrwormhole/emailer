// Package emailer provides shared foundation stones for email providers
package emailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Config configures the email clients
type Config struct {
	Key string
	http.Client
}

// Email is generic email structure for all providers
type Email struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	BCC         []string `json:"bcc"`
	CC          []string `json:"cc"`
	Subject     string   `json:"subject"`
	HTMLContent string   `json:"htmlContent"`
	TextContent string   `json:"textContent"`
}

// ValidationMsg returns empty if all validations passed, else it will return failed validation message
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

// Sender is a behaviour for email senders
type Sender interface {
	Send(ctx context.Context, e Email) error
}

// HandlerFunc is opinionated/reusable HTTP handler for brevo provider
func HandlerFunc(sender Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e Email
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
			return
		}

		if m := e.ValidationMsg(); m != "" {
			http.Error(w, fmt.Sprintf("Failed to validate: %v", m), http.StatusBadRequest)
			return
		}

		if err := sender.Send(r.Context(), e); err != nil {
			slog.LogAttrs(r.Context(), slog.LevelError, fmt.Sprintf("%T.Send(%v)", sender, e), slog.String("err", err.Error()))
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "Email successfully sent")
	}
}
