package app

import (
	"errors"

	"github.com/ClouGence/cloudcanal-openapi-cli/internal/cluster"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/config"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/console"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/consolejob"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/datajob"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/datasource"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/i18n"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/jobconfig"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/openapi"
	ccschema "github.com/ClouGence/cloudcanal-openapi-cli/internal/schema"
	"github.com/ClouGence/cloudcanal-openapi-cli/internal/worker"
)

type RuntimeContext interface {
	Config() config.AppConfig
	CurrentProfile() string
	Language() string
	ProfileSummaries() []config.ProfileSummary
	DataJobs() datajob.Operations
	DataSources() datasource.Operations
	Clusters() cluster.Operations
	Workers() worker.Operations
	ConsoleJobs() consolejob.Operations
	JobConfigs() jobconfig.Operations
	Schemas() ccschema.Operations
	Reinitialize(io console.IO) (bool, error)
	AddProfile(name string, io console.IO) (bool, error)
	UseProfile(name string) error
	RemoveProfile(name string) error
	SetLanguage(language string) error
}

type Runtime struct {
	configService  *config.Service
	state          config.State
	currentProfile string
	config         config.AppConfig
	dataJobs       datajob.Operations
	dataSources    datasource.Operations
	clusters       cluster.Operations
	workers        worker.Operations
	consoleJobs    consolejob.Operations
	jobConfigs     jobconfig.Operations
	schemas        ccschema.Operations
}

func NewRuntime(configService *config.Service) *Runtime {
	return &Runtime{configService: configService}
}

func (r *Runtime) InitializeIfNeeded(io console.IO) (bool, error) {
	if !r.configService.Exists() {
		return r.Reinitialize(io)
	}

	state, err := r.configService.Load()
	if err != nil {
		if errors.Is(err, config.ErrLegacyFormat) {
			io.Println(i18n.T("runtime.legacyConfig"))
		} else {
			io.Println(i18n.T("runtime.invalidConfig", err.Error()))
		}
		return r.Reinitialize(io)
	}
	if err := r.activateState(state); err != nil {
		io.Println(i18n.T("runtime.invalidConfig", err.Error()))
		return r.Reinitialize(io)
	}
	return true, nil
}

func (r *Runtime) Reinitialize(io console.IO) (bool, error) {
	state := r.state
	if state.Language == "" {
		state.Language = r.Language()
	}
	_ = i18n.SetLanguage(state.NormalizedLanguage())

	profileName := r.currentProfile
	if profileName == "" {
		profileName = config.DefaultProfileName
	}

	initial := r.config
	if existing, ok := state.Profiles[profileName]; ok {
		initial = existing
	}

	wizard := config.NewWizard(io, r.validateConfig, profileName, initial)
	cfg, err := wizard.Run()
	if err != nil {
		return false, err
	}
	if cfg == nil {
		io.Println(i18n.T("runtime.initCancelled"))
		return false, nil
	}

	state.Language = state.NormalizedLanguage()
	if state.Profiles == nil {
		state.Profiles = make(map[string]config.AppConfig)
	}
	state.CurrentProfile = profileName
	state.Profiles[profileName] = *cfg
	if err := r.saveAndActivate(state); err != nil {
		return false, err
	}
	io.Println(i18n.T("wizard.savedTo", r.configService.Path()))
	return true, nil
}

func (r *Runtime) AddProfile(name string, io console.IO) (bool, error) {
	if err := config.ValidateProfileName(name); err != nil {
		return false, err
	}

	state := r.state
	profileName := config.NormalizeProfileName(name)
	if state.Profiles == nil {
		state.Profiles = make(map[string]config.AppConfig)
	}
	if _, exists := state.Profiles[profileName]; exists {
		return false, errors.New(i18n.T("config.profileExists", profileName))
	}

	_ = i18n.SetLanguage(state.NormalizedLanguage())
	wizard := config.NewWizard(io, r.validateConfig, profileName, config.AppConfig{})
	cfg, err := wizard.Run()
	if err != nil {
		return false, err
	}
	if cfg == nil {
		io.Println(i18n.T("runtime.initCancelled"))
		return false, nil
	}

	state.Profiles[profileName] = *cfg
	if state.CurrentProfile == "" {
		state.CurrentProfile = profileName
	}
	if err := r.configService.Save(state); err != nil {
		return false, err
	}
	if state.ActiveProfileName() == profileName {
		if err := r.activateState(state); err != nil {
			return false, err
		}
	} else {
		r.state = state
	}
	io.Println(i18n.T("wizard.savedTo", r.configService.Path()))
	return true, nil
}

func (r *Runtime) UseProfile(name string) error {
	if err := config.ValidateProfileName(name); err != nil {
		return err
	}

	state := r.state
	profileName := config.NormalizeProfileName(name)
	if _, ok := state.Profiles[profileName]; !ok {
		return errors.New(i18n.T("config.profileNotFound", profileName))
	}

	next := state
	next.CurrentProfile = profileName
	if err := r.prepareState(next); err != nil {
		return err
	}
	if err := r.configService.Save(next); err != nil {
		return err
	}
	return r.activateState(next)
}

func (r *Runtime) RemoveProfile(name string) error {
	if err := config.ValidateProfileName(name); err != nil {
		return err
	}

	state := r.state
	profileName := config.NormalizeProfileName(name)
	if _, ok := state.Profiles[profileName]; !ok {
		return errors.New(i18n.T("config.profileNotFound", profileName))
	}
	if state.ActiveProfileName() == profileName {
		return errors.New(i18n.T("config.profileRemoveActive", profileName))
	}

	delete(state.Profiles, profileName)
	if err := r.configService.Save(state); err != nil {
		return err
	}
	r.state = state
	return nil
}

func (r *Runtime) Config() config.AppConfig {
	return r.config
}

func (r *Runtime) CurrentProfile() string {
	return r.currentProfile
}

func (r *Runtime) Language() string {
	return r.state.NormalizedLanguage()
}

func (r *Runtime) ProfileSummaries() []config.ProfileSummary {
	return r.state.Summaries()
}

func (r *Runtime) DataJobs() datajob.Operations {
	return r.dataJobs
}

func (r *Runtime) DataSources() datasource.Operations {
	return r.dataSources
}

func (r *Runtime) Clusters() cluster.Operations {
	return r.clusters
}

func (r *Runtime) Workers() worker.Operations {
	return r.workers
}

func (r *Runtime) ConsoleJobs() consolejob.Operations {
	return r.consoleJobs
}

func (r *Runtime) JobConfigs() jobconfig.Operations {
	return r.jobConfigs
}

func (r *Runtime) Schemas() ccschema.Operations {
	return r.schemas
}

func (r *Runtime) SetLanguage(language string) error {
	normalized := i18n.NormalizeLanguage(language)
	if normalized == "" {
		return errors.New(i18n.T("config.languageUnsupported"))
	}

	state := r.state
	if state.Profiles == nil {
		return errors.New(i18n.T("config.noProfilesConfigured"))
	}
	state.Language = normalized
	if err := r.configService.Save(state); err != nil {
		return err
	}
	r.state = state
	return i18n.SetLanguage(normalized)
}

func (r *Runtime) validateConfig(cfg config.AppConfig) error {
	_ = i18n.SetLanguage(r.Language())
	client, err := openapi.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.ProbeAuthentication()
}

func (r *Runtime) saveAndActivate(state config.State) error {
	if err := r.configService.Save(state); err != nil {
		return err
	}
	return r.activateState(state)
}

func (r *Runtime) prepareState(state config.State) error {
	_, cfg, err := state.ActiveProfile()
	if err != nil {
		return err
	}
	_, err = openapi.NewClient(cfg)
	return err
}

func (r *Runtime) activateState(state config.State) error {
	if err := r.prepareState(state); err != nil {
		return err
	}

	profileName, cfg, _ := state.ActiveProfile()
	_ = i18n.SetLanguage(state.NormalizedLanguage())
	client, err := openapi.NewClient(cfg)
	if err != nil {
		return err
	}

	r.state = state
	r.currentProfile = profileName
	r.config = cfg
	r.dataJobs = datajob.NewService(client)
	r.dataSources = datasource.NewService(client)
	r.clusters = cluster.NewService(client)
	r.workers = worker.NewService(client)
	r.consoleJobs = consolejob.NewService(client)
	r.jobConfigs = jobconfig.NewService(client)
	r.schemas = ccschema.NewService(client)
	return nil
}
