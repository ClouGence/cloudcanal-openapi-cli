package updatecheck

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultLatestReleaseURL = "https://github.com/ClouGence/cloudcanal-openapi-cli/releases/latest"
	DefaultUpgradeCommand   = "curl -fsSL https://raw.githubusercontent.com/ClouGence/cloudcanal-openapi-cli/main/scripts/install.sh | bash"
	defaultTimeout          = 2 * time.Second
	defaultUserAgent        = "cloudcanal-cli-update-checker"
)

type Notice struct {
	CurrentVersion string
	LatestVersion  string
	UpgradeCommand string
}

type Checker struct {
	LatestReleaseURL string
	UpgradeCommand   string
	HTTPClient       *http.Client
}

type version struct {
	major int
	minor int
	patch int
}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) Check(currentVersion string) (Notice, error) {
	current, ok := parseVersion(currentVersion)
	if !ok {
		return Notice{}, nil
	}

	latestRaw, err := c.latestVersion()
	if err != nil {
		return Notice{}, err
	}
	latest, ok := parseVersion(latestRaw)
	if !ok || current.compare(latest) >= 0 {
		return Notice{}, nil
	}

	return Notice{
		CurrentVersion: displayVersion(currentVersion),
		LatestVersion:  displayVersion(latestRaw),
		UpgradeCommand: c.upgradeCommand(),
	}, nil
}

func (c *Checker) latestVersion() (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.latestReleaseURL(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" && resp.Request != nil {
		location = resp.Request.URL.String()
	}
	if location == "" {
		return "", fmt.Errorf("latest release redirect missing")
	}

	latest, ok := releaseTagFromLocation(location)
	if !ok {
		return "", fmt.Errorf("latest release tag missing")
	}
	return latest, nil
}

func (c *Checker) latestReleaseURL() string {
	if strings.TrimSpace(c.LatestReleaseURL) != "" {
		return strings.TrimSpace(c.LatestReleaseURL)
	}
	return DefaultLatestReleaseURL
}

func (c *Checker) upgradeCommand() string {
	if strings.TrimSpace(c.UpgradeCommand) != "" {
		return strings.TrimSpace(c.UpgradeCommand)
	}
	return DefaultUpgradeCommand
}

func (c *Checker) httpClient() *http.Client {
	if c.HTTPClient == nil {
		return &http.Client{
			Timeout:       defaultTimeout,
			CheckRedirect: stopRedirects,
		}
	}

	client := *c.HTTPClient
	if client.Timeout <= 0 {
		client.Timeout = defaultTimeout
	}
	if client.CheckRedirect == nil {
		client.CheckRedirect = stopRedirects
	}
	return &client
}

func stopRedirects(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func releaseTagFromLocation(location string) (string, bool) {
	parsed, err := url.Parse(strings.TrimSpace(location))
	if err != nil {
		return "", false
	}

	tag := strings.TrimSpace(path.Base(strings.TrimRight(parsed.Path, "/")))
	if tag == "" || strings.EqualFold(tag, "latest") || strings.EqualFold(tag, "tag") {
		return "", false
	}
	return tag, true
}

func displayVersion(raw string) string {
	parsed, ok := parseVersion(raw)
	if !ok {
		return strings.TrimSpace(raw)
	}
	return "v" + parsed.String()
}

func parseVersion(raw string) (version, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return version{}, false
	}
	value = strings.TrimPrefix(strings.TrimPrefix(value, "v"), "V")
	if cut := strings.IndexAny(value, "-+"); cut >= 0 {
		value = value[:cut]
	}

	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return version{}, false
	}

	numbers := [3]int{}
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return version{}, false
		}
		number, err := strconv.Atoi(part)
		if err != nil || number < 0 {
			return version{}, false
		}
		numbers[i] = number
	}

	return version{major: numbers[0], minor: numbers[1], patch: numbers[2]}, true
}

func (v version) compare(other version) int {
	if v.major != other.major {
		return compareInts(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInts(v.minor, other.minor)
	}
	return compareInts(v.patch, other.patch)
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func compareInts(left int, right int) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}
