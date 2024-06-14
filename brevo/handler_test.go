package brevo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mrwormhole/emailer"
)

func TestEmailHandler_BrokenRequest(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/email", bytes.NewBuffer([]byte{'h', 'e', 'l', 'l', 'o'}))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EmailHandler(nil))
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusBadRequest, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Failed to decode request: invalid character 'h' looking for beginning of value\n", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}

func TestEmailHandler_FailedValidation(t *testing.T) {
	email := emailer.Email{
		From: "a@a.com",
	}
	raw, err := json.Marshal(email)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", email, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/email", bytes.NewBuffer(raw))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EmailHandler(nil))
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusBadRequest, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Failed to validate: to field must not be blank\n", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}

func TestEmailHandler_Success(t *testing.T) {
	tripper := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
		}
	}
	client, err := NewTestClient[RoundTripFunc](tripper)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	email := emailer.Email{
		From:        "a@a.com",
		To:          []string{"b@b.com"},
		BCC:         []string{"bcc@bcc.com"},
		CC:          []string{"cc@cc.com"},
		Subject:     "sub",
		HTMLContent: "html",
	}
	raw, err := json.Marshal(email)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", email, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(raw))
	rr := httptest.NewRecorder()
	var handler http.HandlerFunc = EmailHandler(client)
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusOK, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Email successfully sent", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}

func TestEmailHandler_FaultyClient(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(100))
	tripper := func(req *http.Request) *http.Response {
		return nil
	}
	client, err := NewTestClient[FaultyRoundTripFunc](tripper)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	email := emailer.Email{
		From:        "a@a.com",
		To:          []string{"b@b.com"},
		Subject:     "sub",
		HTMLContent: "html",
	}
	raw, err := json.Marshal(email)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", email, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(raw))
	rr := httptest.NewRecorder()
	var handler http.HandlerFunc = EmailHandler(client)
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusInternalServerError, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Failed to send email\n", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}

func TestEmailHandler_TeapotClient(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(100))
	tripper := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusTeapot,
		}
	}
	client, err := NewTestClient[RoundTripFunc](tripper)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	email := emailer.Email{
		From:        "a@a.com",
		To:          []string{"b@b.com"},
		Subject:     "sub",
		HTMLContent: "html",
	}
	raw, err := json.Marshal(email)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", email, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(raw))
	rr := httptest.NewRecorder()
	var handler http.HandlerFunc = EmailHandler(client)
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusInternalServerError, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Failed to send email\n", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}

func TestEmailHandler_ProviderCodeMessage(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(100))
	tripper := func(req *http.Request) *http.Response {
		cm := CodeMessage{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: "brevo don't like that",
		}
		raw, err := json.Marshal(cm)
		if err != nil {
			t.Fatalf("json.Marshal(%v): %v", cm, err)
		}

		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBuffer(raw)),
		}
	}
	client, err := NewTestClient[RoundTripFunc](tripper)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	email := emailer.Email{
		From:        "a@a.com",
		To:          []string{"b@b.com"},
		Subject:     "sub",
		HTMLContent: "html",
	}
	raw, err := json.Marshal(email)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", email, err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(raw))
	rr := httptest.NewRecorder()
	var handler http.HandlerFunc = EmailHandler(client)
	handler.ServeHTTP(rr, req)

	if diff := cmp.Diff(http.StatusInternalServerError, rr.Code); diff != "" {
		t.Errorf("EmailHandler(): HTTP code diff=\n %v", diff)
	}

	if diff := cmp.Diff("Failed to send email\n", rr.Body.String()); diff != "" {
		t.Errorf("EmailHandler(): HTTP body diff=\n %v", diff)
	}
}
