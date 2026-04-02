package updatecheck

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckerReturnsNoticeWhenLatestReleaseIsNewer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases/latest" {
			t.Fatalf("request path = %q, want /releases/latest", r.URL.Path)
		}
		w.Header().Set("Location", "/releases/tag/v0.1.3")
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()

	checker := &Checker{
		LatestReleaseURL: server.URL + "/releases/latest",
		UpgradeCommand:   "curl upgrade",
		HTTPClient:       server.Client(),
	}

	notice, err := checker.Check("0.1.2")
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if notice.CurrentVersion != "v0.1.2" {
		t.Fatalf("CurrentVersion = %q, want v0.1.2", notice.CurrentVersion)
	}
	if notice.LatestVersion != "v0.1.3" {
		t.Fatalf("LatestVersion = %q, want v0.1.3", notice.LatestVersion)
	}
	if notice.UpgradeCommand != "curl upgrade" {
		t.Fatalf("UpgradeCommand = %q, want curl upgrade", notice.UpgradeCommand)
	}
}

func TestCheckerSkipsInvalidOrDevelopmentVersions(t *testing.T) {
	checker := NewChecker()

	for _, current := range []string{"", "dev", "main", "unknown"} {
		notice, err := checker.Check(current)
		if err != nil {
			t.Fatalf("Check(%q) error = %v", current, err)
		}
		if notice != (Notice{}) {
			t.Fatalf("Check(%q) = %#v, want zero notice", current, notice)
		}
	}
}

func TestCheckerSkipsWhenCurrentVersionIsLatestOrNewer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/releases/tag/v0.1.2")
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()

	checker := &Checker{
		LatestReleaseURL: server.URL + "/releases/latest",
		HTTPClient:       server.Client(),
	}

	for _, current := range []string{"0.1.2", "0.1.3"} {
		notice, err := checker.Check(current)
		if err != nil {
			t.Fatalf("Check(%q) error = %v", current, err)
		}
		if notice != (Notice{}) {
			t.Fatalf("Check(%q) = %#v, want zero notice", current, notice)
		}
	}
}

func TestReleaseTagFromLocation(t *testing.T) {
	tests := []struct {
		location string
		want     string
		ok       bool
	}{
		{location: "https://github.com/ClouGence/cloudcanal-openapi-cli/releases/tag/v0.1.2", want: "v0.1.2", ok: true},
		{location: "/ClouGence/cloudcanal-openapi-cli/releases/tag/v1.2.3", want: "v1.2.3", ok: true},
		{location: "https://github.com/ClouGence/cloudcanal-openapi-cli/releases/latest", ok: false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s", tc.location), func(t *testing.T) {
			got, ok := releaseTagFromLocation(tc.location)
			if ok != tc.ok {
				t.Fatalf("releaseTagFromLocation(%q) ok = %v, want %v", tc.location, ok, tc.ok)
			}
			if got != tc.want {
				t.Fatalf("releaseTagFromLocation(%q) = %q, want %q", tc.location, got, tc.want)
			}
		})
	}
}
