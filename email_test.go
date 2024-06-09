package emailer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidationMsg(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{
			name: "empty from",
			email: Email{
				From: "",
			},
			want: "from field must not be blank",
		},
		{
			name: "invalid from",
			email: Email{
				From: "a",
			},
			want: `"a" is not a valid email`,
		},
		{
			name: "empty to",
			email: Email{
				From: "a@a.com",
			},
			want: "to field must not be blank",
		},
		{
			name: "invalid to",
			email: Email{
				From: "a@a.com",
				To:   []string{"b"},
			},
			want: `"b" is not a valid email`,
		},
		{
			name: "empty subject",
			email: Email{
				From: "a@a.com",
				To:   []string{"b@b.com"},
			},
			want: "subject field must not be blank",
		},
		{
			name: "missing html or text content",
			email: Email{
				From:    "a@a.com",
				To:      []string{"b@b.com"},
				Subject: "subj",
			},
			want: "either the htmlContent or textContent field must be filled",
		},
		{
			name: "invalid BCC",
			email: Email{
				From:        "a@a.com",
				To:          []string{"b@b.com"},
				Subject:     "subj",
				TextContent: "text",
				BCC:         []string{"duck"},
			},
			want: `"duck" is not a valid email`,
		},
		{
			name: "invalid CC",
			email: Email{
				From:        "a@a.com",
				To:          []string{"b@b.com"},
				Subject:     "subj",
				HTMLContent: "html",
				CC:          []string{"ant"},
			},
			want: `"ant" is not a valid email`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.email.ValidationMsg()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ValidationMsg(): diff=\n %v", diff)
			}
		})
	}
}
