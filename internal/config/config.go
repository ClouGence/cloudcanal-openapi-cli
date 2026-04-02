package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
)

const (
	defaultHTTPTimeoutSeconds         = 10
	defaultHTTPReadRetryBackoffMillis = 250
	DefaultProfileName                = "dev"
)

var ErrLegacyFormat = errors.New("legacy config format")

type AppConfig struct {
	APIBaseURL                 string `json:"apiBaseUrl"`
	AccessKey                  string `json:"accessKey"`
	SecretKey                  string `json:"secretKey"`
	HTTPTimeoutSeconds         int    `json:"httpTimeoutSeconds,omitempty"`
	HTTPReadMaxRetries         int    `json:"httpReadMaxRetries,omitempty"`
	HTTPReadRetryBackoffMillis int    `json:"httpReadRetryBackoffMillis,omitempty"`
}

type State struct {
	Language       string               `json:"language,omitempty"`
	CurrentProfile string               `json:"currentProfile,omitempty"`
	Profiles       map[string]AppConfig `json:"profiles,omitempty"`
}

type ProfileSummary struct {
	Name       string `json:"name"`
	APIBaseURL string `json:"apiBaseUrl"`
	Current    bool   `json:"current"`
}

func (c AppConfig) Validate() error {
	return c.ValidateForLanguage(i18n.CurrentLanguage())
}

func (c AppConfig) ValidateForLanguage(language string) error {
	normalizedLanguage := normalizeLanguage(language)
	if strings.TrimSpace(c.APIBaseURL) == "" {
		return errors.New(i18n.TFor(normalizedLanguage, "config.apiBaseUrlRequired"))
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New(i18n.TFor(normalizedLanguage, "config.accessKeyRequired"))
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New(i18n.TFor(normalizedLanguage, "config.secretKeyRequired"))
	}
	if c.HTTPTimeoutSeconds < 0 {
		return errors.New(i18n.TFor(normalizedLanguage, "config.httpTimeoutSecondsInvalid"))
	}
	if c.HTTPReadMaxRetries < 0 {
		return errors.New(i18n.TFor(normalizedLanguage, "config.httpReadMaxRetriesInvalid"))
	}
	if c.HTTPReadRetryBackoffMillis < 0 {
		return errors.New(i18n.TFor(normalizedLanguage, "config.httpReadRetryBackoffMillisInvalid"))
	}

	parsed, err := url.Parse(strings.TrimSpace(c.APIBaseURL))
	if err != nil {
		return errors.New(i18n.TFor(normalizedLanguage, "config.apiBaseUrlInvalid"))
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New(i18n.TFor(normalizedLanguage, "config.apiBaseUrlScheme"))
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return errors.New(i18n.TFor(normalizedLanguage, "config.apiBaseUrlHost"))
	}
	return nil
}

func (c AppConfig) NormalizedBaseURL() string {
	value := strings.TrimSpace(c.APIBaseURL)
	return strings.TrimRight(value, "/")
}

func (c AppConfig) WithDefaults() AppConfig {
	return c
}

func (c AppConfig) HTTPTimeout() time.Duration {
	return time.Duration(c.HTTPTimeoutSecondsValue()) * time.Second
}

func (c AppConfig) HTTPTimeoutSecondsValue() int {
	if c.HTTPTimeoutSeconds <= 0 {
		return defaultHTTPTimeoutSeconds
	}
	return c.HTTPTimeoutSeconds
}

func (c AppConfig) HTTPReadMaxRetriesValue() int {
	if c.HTTPReadMaxRetries <= 0 {
		return 0
	}
	return c.HTTPReadMaxRetries
}

func (c AppConfig) HTTPReadRetryBackoff() time.Duration {
	return time.Duration(c.HTTPReadRetryBackoffMillisValue()) * time.Millisecond
}

func (c AppConfig) HTTPReadRetryBackoffMillisValue() int {
	if c.HTTPReadRetryBackoffMillis <= 0 {
		return defaultHTTPReadRetryBackoffMillis
	}
	return c.HTTPReadRetryBackoffMillis
}

func (s State) NormalizedLanguage() string {
	return normalizeLanguage(s.Language)
}

func (s State) ActiveProfileName() string {
	return NormalizeProfileName(s.CurrentProfile)
}

func (s State) Validate() error {
	language := s.NormalizedLanguage()
	if normalized := i18n.NormalizeLanguage(s.Language); normalized == "" && strings.TrimSpace(s.Language) != "" {
		return errors.New(i18n.TFor(language, "config.languageUnsupported"))
	}
	if len(s.Profiles) == 0 {
		return errors.New(i18n.TFor(language, "config.noProfilesConfigured"))
	}
	current := s.ActiveProfileName()
	if current == "" {
		return errors.New(i18n.TFor(language, "config.currentProfileRequired"))
	}
	if _, ok := s.Profiles[current]; !ok {
		return errors.New(i18n.TFor(language, "config.currentProfileMissing", current))
	}
	for name, cfg := range s.Profiles {
		if NormalizeProfileName(name) == "" {
			return errors.New(i18n.TFor(language, "config.profileNameRequired"))
		}
		if err := cfg.ValidateForLanguage(language); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func (s State) ActiveProfile() (string, AppConfig, error) {
	if err := s.Validate(); err != nil {
		return "", AppConfig{}, err
	}
	name := s.ActiveProfileName()
	return name, s.Profiles[name].WithDefaults(), nil
}

func (s State) Summaries() []ProfileSummary {
	names := make([]string, 0, len(s.Profiles))
	for name := range s.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	summaries := make([]ProfileSummary, 0, len(names))
	current := s.ActiveProfileName()
	for _, name := range names {
		summaries = append(summaries, ProfileSummary{
			Name:       name,
			APIBaseURL: s.Profiles[name].APIBaseURL,
			Current:    name == current,
		})
	}
	return summaries
}

type Service struct {
	path string
}

func NewService(path string) *Service {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath()
	}
	return &Service{path: path}
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cloudcanal-cli/config.json"
	}
	return filepath.Join(home, ".cloudcanal-cli", "config.json")
}

func (s *Service) Path() string {
	return s.path
}

func (s *Service) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}

func (s *Service) Load() (State, error) {
	return s.loadFromPath(s.path)
}

func (s *Service) LoadLanguage() string {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return i18n.DefaultLanguage()
	}

	var payload struct {
		Language string `json:"language"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return i18n.DefaultLanguage()
	}
	return normalizeLanguage(payload.Language)
}

func (s *Service) loadFromPath(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return State{}, err
	}
	if isLegacyConfig(data) {
		return State{}, fmt.Errorf("%w: %s", ErrLegacyFormat, i18n.T("config.legacyFormat"))
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, errors.New(i18n.T("config.invalidJSON"))
	}
	state = normalizeState(state)
	if err := state.Validate(); err != nil {
		return State{}, err
	}
	return state, nil
}

func (s *Service) Save(state State) error {
	state = normalizeState(state)
	if err := state.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, content, 0o600)
}

func NormalizeProfileName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func ValidateProfileName(name string) error {
	if NormalizeProfileName(name) == "" {
		return errors.New(i18n.T("config.profileNameRequired"))
	}
	return nil
}

func normalizeLanguage(language string) string {
	normalized := i18n.NormalizeLanguage(language)
	if normalized == "" {
		return i18n.DefaultLanguage()
	}
	return normalized
}

func normalizeState(state State) State {
	normalized := State{
		Language:       normalizeLanguage(state.Language),
		CurrentProfile: NormalizeProfileName(state.CurrentProfile),
		Profiles:       make(map[string]AppConfig, len(state.Profiles)),
	}
	for name, cfg := range state.Profiles {
		normalized.Profiles[NormalizeProfileName(name)] = cfg.WithDefaults()
	}
	return normalized
}

func isLegacyConfig(data []byte) bool {
	var payload struct {
		APIBaseURL json.RawMessage `json:"apiBaseUrl"`
		AccessKey  json.RawMessage `json:"accessKey"`
		SecretKey  json.RawMessage `json:"secretKey"`
		Profiles   json.RawMessage `json:"profiles"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return false
	}
	if len(payload.Profiles) != 0 {
		return false
	}
	return len(payload.APIBaseURL) != 0 || len(payload.AccessKey) != 0 || len(payload.SecretKey) != 0
}
