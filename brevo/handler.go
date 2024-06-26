package brevo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrwormhole/emailer"
)

// EmailHandler is opinionated/reusable HTTP handler for brevo provider
func EmailHandler(sender emailer.Sender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var e emailer.Email
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
		_, _ = w.Write([]byte("Email successfully sent"))
	}
}
