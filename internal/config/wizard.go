package config

import (
	"errors"
	"io"
	"strings"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/console"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
)

type Validator func(AppConfig) error

type Wizard struct {
	io          console.IO
	validator   Validator
	initial     AppConfig
	profileName string
}

func NewWizard(io console.IO, validator Validator, profileName string, initial AppConfig) *Wizard {
	return &Wizard{
		io:          io,
		validator:   validator,
		initial:     initial,
		profileName: NormalizeProfileName(profileName),
	}
}

func (w *Wizard) Run() (*AppConfig, error) {
	current := w.initial.WithDefaults()
	profileName := w.profileName
	if profileName == "" {
		profileName = DefaultProfileName
	}

	w.io.Println(i18n.T("wizard.title", profileName))
	w.io.Println(i18n.T("wizard.cancelHint"))
	w.io.Println(i18n.T("wizard.apiHostHint"))
	if w.hasInitialValue(w.initial) {
		w.io.Println(i18n.T("wizard.keepCurrent"))
	}

	for {
		apiBaseURL, cancelled, err := w.promptRequired("apiHost", current.APIBaseURL, validateAPIBaseURL)
		if err != nil {
			if isWizardCancelled(err) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		accessKey, cancelled, err := w.promptRequired("ak", current.AccessKey, validateAccessKey)
		if err != nil {
			if isWizardCancelled(err) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		secretKey, cancelled, err := w.promptSecret("sk", current.SecretKey)
		if err != nil {
			if isWizardCancelled(err) {
				return nil, nil
			}
			return nil, err
		}
		if cancelled {
			return nil, nil
		}

		current = AppConfig{
			APIBaseURL: apiBaseURL,
			AccessKey:  accessKey,
			SecretKey:  secretKey,
		}

		if err := current.Validate(); err != nil {
			w.io.Println(i18n.T("wizard.invalidConfig", err.Error()))
			continue
		}
		w.io.Println(i18n.T("wizard.checkingConnection"))
		if err := w.validator(current); err != nil {
			w.io.Println(i18n.T("wizard.validationFailed", err.Error()))
			w.io.Println(i18n.T("wizard.reuseValues"))
			continue
		}
		return &current, nil
	}
}

func (w *Wizard) promptRequired(label, current string, validate func(string) error) (string, bool, error) {
	for {
		value, cancelled, err := w.promptWithDefault(label, current)
		if err != nil || cancelled {
			return "", cancelled, err
		}
		if err := validate(value); err != nil {
			w.io.Println(i18n.T("wizard.invalidField", label, err.Error()))
			continue
		}
		return value, false, nil
	}
}

func (w *Wizard) promptSecret(label, current string) (string, bool, error) {
	for {
		prompt := label + ": "
		if strings.TrimSpace(current) != "" {
			prompt = label + " [hidden]: "
		}

		value, err := w.io.ReadSecret(prompt)
		if err != nil {
			return "", false, err
		}
		trimmed := strings.TrimSpace(value)
		if strings.EqualFold(trimmed, "exit") {
			return "", true, nil
		}
		if trimmed == "" {
			if strings.TrimSpace(current) == "" {
				w.io.Println(i18n.T("wizard.invalidField", label, i18n.T("config.secretKeyRequired")))
				continue
			}
			return current, false, nil
		}
		return trimmed, false, nil
	}
}

func (w *Wizard) promptWithDefault(label, current string) (string, bool, error) {
	prompt := label + ": "
	if strings.TrimSpace(current) != "" {
		prompt = label + " [" + current + "]: "
	}

	value, err := w.io.ReadLine(prompt)
	if err != nil {
		return "", false, err
	}
	trimmed := strings.TrimSpace(value)
	if strings.EqualFold(trimmed, "exit") {
		return "", true, nil
	}
	if trimmed == "" {
		return strings.TrimSpace(current), false, nil
	}
	return trimmed, false, nil
}

func (w *Wizard) hasInitialValue(cfg AppConfig) bool {
	return strings.TrimSpace(cfg.APIBaseURL) != "" ||
		strings.TrimSpace(cfg.AccessKey) != "" ||
		strings.TrimSpace(cfg.SecretKey) != ""
}

func validateAPIBaseURL(value string) error {
	return AppConfig{APIBaseURL: value, AccessKey: "ak", SecretKey: "sk"}.Validate()
}

func validateAccessKey(value string) error {
	return AppConfig{APIBaseURL: "https://cc.example.com", AccessKey: value, SecretKey: "sk"}.Validate()
}

func isWizardCancelled(err error) bool {
	return errors.Is(err, io.EOF) || console.IsPromptAborted(err)
}
