package app

import (
	"cloudcanal-openapi-cli/internal/cluster"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/consolejob"
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/datasource"
	"cloudcanal-openapi-cli/internal/jobconfig"
	"cloudcanal-openapi-cli/internal/openapi"
	"cloudcanal-openapi-cli/internal/worker"
)

type RuntimeContext interface {
	Config() config.AppConfig
	DataJobs() datajob.Operations
	DataSources() datasource.Operations
	Clusters() cluster.Operations
	Workers() worker.Operations
	ConsoleJobs() consolejob.Operations
	JobConfigs() jobconfig.Operations
	Reinitialize(io console.IO) (bool, error)
}

type Runtime struct {
	configService *config.Service
	config        config.AppConfig
	dataJobs      datajob.Operations
	dataSources   datasource.Operations
	clusters      cluster.Operations
	workers       worker.Operations
	consoleJobs   consolejob.Operations
	jobConfigs    jobconfig.Operations
}

func NewRuntime(configService *config.Service) *Runtime {
	return &Runtime{configService: configService}
}

func (r *Runtime) InitializeIfNeeded(io console.IO) (bool, error) {
	if !r.configService.Exists() {
		return r.Reinitialize(io)
	}

	cfg, err := r.configService.Load()
	if err != nil {
		io.Println("Existing configuration is invalid: " + err.Error())
		return r.Reinitialize(io)
	}
	if err := r.activate(cfg); err != nil {
		io.Println("Existing configuration is invalid: " + err.Error())
		return r.Reinitialize(io)
	}
	return true, nil
}

func (r *Runtime) Reinitialize(io console.IO) (bool, error) {
	wizard := config.NewWizard(io, r.configService, r.validateConfig, r.config)
	cfg, err := wizard.Run()
	if err != nil {
		return false, err
	}
	if cfg == nil {
		io.Println("Initialization cancelled.")
		return false, nil
	}
	if err := r.activate(*cfg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Runtime) Config() config.AppConfig {
	return r.config
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

func (r *Runtime) validateConfig(cfg config.AppConfig) error {
	client, err := openapi.NewClient(cfg)
	if err != nil {
		return err
	}
	service := datajob.NewService(client)
	_, err = service.ListJobs(datajob.ListOptions{})
	return err
}

func (r *Runtime) activate(cfg config.AppConfig) error {
	client, err := openapi.NewClient(cfg)
	if err != nil {
		return err
	}
	r.config = cfg
	r.dataJobs = datajob.NewService(client)
	r.dataSources = datasource.NewService(client)
	r.clusters = cluster.NewService(client)
	r.workers = worker.NewService(client)
	r.consoleJobs = consolejob.NewService(client)
	r.jobConfigs = jobconfig.NewService(client)
	return nil
}
