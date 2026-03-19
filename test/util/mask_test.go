package util_test

import (
	"cloudcanal-openapi-cli/internal/util"
	"testing"
)

func TestMaskSensitiveTextMasksCredentialsAndKeys(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "dsn password",
			input: "localhost:3306(root/123456)",
			want:  "localhost:3306(root/******)",
		},
		{
			name:  "assignment keys",
			input: "oss(access-key=test-access-1234 secret-key=test-secret-5678)",
			want:  "oss(access-key=test********1234 secret-key=test********5678)",
		},
		{
			name:  "single token key",
			input: "clougence(token-test-abcdef12)",
			want:  "clougence(toke***********ef12)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := util.MaskSensitiveText(tc.input); got != tc.want {
				t.Fatalf("MaskSensitiveText(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
