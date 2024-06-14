package brevo

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	_, err := New("", nil)
	want := errors.New("brevo API key is blank")
	if !cmp.Equal(want.Error(), err.Error()) {
		t.Errorf("New(): got=%q want=%q", err, want)
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type FaultyRoundTripFunc func(req *http.Request) *http.Response

func (f FaultyRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), errors.New("boom baby boom")
}

func NewTestClient[T RoundTripFunc | FaultyRoundTripFunc](tripper T) (*EmailClient, error) {
	return New("key", &http.Client{Transport: http.RoundTripper(tripper)})
}
