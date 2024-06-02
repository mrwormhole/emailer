package brevo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrwormhole/emailer"
)

func EmailHandler(sender emailer.Sender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		var e emailer.Email
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
			return
		}

		if m := e.ValidationMsg(); m != "" {
			http.Error(w, fmt.Sprintf("Failed to validate: %v", m), http.StatusBadRequest)
			return
		}

		if err := sender.Send(ctx, e); err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, "failed to send email", slog.String("err", err.Error()))
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Email successfully sent"))
	}
}
