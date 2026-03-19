package config

import (
	"cloudcanal-openapi-cli/internal/i18n"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultHTTPTimeoutSeconds         = 10
	defaultHTTPReadRetryBackoffMillis = 250
)

type AppConfig struct {
	APIBaseURL                 string `json:"apiBaseUrl"`
	AccessKey                  string `json:"accessKey"`
	SecretKey                  string `json:"secretKey"`
	Language                   string `json:"language,omitempty"`
	HTTPTimeoutSeconds         int    `json:"httpTimeoutSeconds,omitempty"`
	HTTPReadMaxRetries         int    `json:"httpReadMaxRetries,omitempty"`
	HTTPReadRetryBackoffMillis int    `json:"httpReadRetryBackoffMillis,omitempty"`
}

func (c AppConfig) Validate() error {
	language := c.NormalizedLanguage()
	if strings.TrimSpace(c.APIBaseURL) == "" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlRequired"))
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New(i18n.TFor(language, "config.accessKeyRequired"))
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New(i18n.TFor(language, "config.secretKeyRequired"))
	}
	if normalized := i18n.NormalizeLanguage(c.Language); normalized == "" && strings.TrimSpace(c.Language) != "" {
		return errors.New(i18n.T("config.languageUnsupported"))
	}
	if c.HTTPTimeoutSeconds < 0 {
		return errors.New(i18n.TFor(language, "config.httpTimeoutSecondsInvalid"))
	}
	if c.HTTPReadMaxRetries < 0 {
		return errors.New(i18n.TFor(language, "config.httpReadMaxRetriesInvalid"))
	}
	if c.HTTPReadRetryBackoffMillis < 0 {
		return errors.New(i18n.TFor(language, "config.httpReadRetryBackoffMillisInvalid"))
	}

	parsed, err := url.Parse(strings.TrimSpace(c.APIBaseURL))
	if err != nil {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlInvalid"))
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlScheme"))
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return errors.New(i18n.TFor(language, "config.apiBaseUrlHost"))
	}
	return nil
}

func (c AppConfig) NormalizedBaseURL() string {
	value := strings.TrimSpace(c.APIBaseURL)
	return strings.TrimRight(value, "/")
}

func (c AppConfig) NormalizedLanguage() string {
	normalized := i18n.NormalizeLanguage(c.Language)
	if normalized == "" {
		return i18n.DefaultLanguage()
	}
	return normalized
}

func (c AppConfig) WithDefaults() AppConfig {
	c.Language = c.NormalizedLanguage()
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
		return ".cloudcanal/config.json"
	}
	return filepath.Join(home, ".cloudcanal", "config.json")
}

func (s *Service) Path() string {
	return s.path
}

func (s *Service) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}

func (s *Service) Load() (AppConfig, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, errors.New(i18n.T("config.invalidJSON"))
	}
	cfg = cfg.WithDefaults()
	if err := cfg.Validate(); err != nil {
		return AppConfig{}, err
	}
	return cfg, nil
}

func (s *Service) Save(cfg AppConfig) error {
	cfg = cfg.WithDefaults()
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, content, 0o600)
}
