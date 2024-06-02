package brevo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"slices"
	"strings"

	"github.com/mrwormhole/emailer"
)

const endpoint = "https://api.brevo.com/v3/smtp/email"

type EmailClient struct {
	key    string
	client *http.Client
}

func New(key string, client *http.Client) (*EmailClient, error) {
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("brevo API key is blank")
	}

	return &EmailClient{
		key:    key,
		client: client,
	}, nil
}

type Detail struct {
	Email string `json:"email"`
}

type Payload struct {
	Sender      Detail   `json:"sender"`
	To          []Detail `json:"to"`
	BCC         []Detail `json:"bcc"`
	CC          []Detail `json:"cc"`
	Subject     string   `json:"subject"`
	HTMLContent string   `json:"htmlContent,omitempty"`
	TextContent string   `json:"textContent,omitempty"`
}

type CodeMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *EmailClient) Send(ctx context.Context, email emailer.Email) error {
	var p Payload
	p.Sender.Email = email.From
	for _, e := range email.To {
		p.To = append(p.To, Detail{Email: e})
	}
	for _, e := range email.BCC {
		p.BCC = append(p.BCC, Detail{Email: e})
	}
	for _, e := range email.CC {
		p.CC = append(p.CC, Detail{Email: e})
	}
	p.Subject = email.Subject
	p.HTMLContent = email.HTMLContent
	p.TextContent = email.TextContent

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(p); err != nil {
		return fmt.Errorf("json.NewEncoder().Encode(%v): %v", p, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}
	req.Header.Add("api-key", c.key)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do(%v): %v", req, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if slices.Contains([]int{http.StatusAccepted, http.StatusCreated, http.StatusOK}, resp.StatusCode) {
		return nil
	}

	var cm CodeMessage
	if err := json.NewDecoder(resp.Body).Decode(&cm); err != nil {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("json.NewDecoder(%v).Decode(): %v", dump, err)
	}

	return fmt.Errorf("unsuccessful response with status code(%d): %v", resp.StatusCode, cm)
}
