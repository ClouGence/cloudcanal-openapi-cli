package repl

import (
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/util"
	"fmt"
	"strconv"
	"strings"
)

func (s *Shell) handleJobs(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println("Usage: jobs list | jobs show <jobId> | jobs schema <jobId> | jobs start <jobId> | jobs stop <jobId> | jobs delete <jobId> | jobs replay <jobId> [--auto-start] [--reset-to-created]")
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseJobListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printJobs(options)
	case "show":
		if len(tokens) != 3 {
			s.io.Println("Usage: jobs show <jobId>")
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		return s.printJob(jobID)
	case "schema":
		if len(tokens) != 3 {
			s.io.Println("Usage: jobs schema <jobId>")
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		return s.printJobSchema(jobID)
	case "start", "stop", "delete":
		if len(tokens) != 3 {
			s.io.Println("Usage: jobs " + strings.ToLower(tokens[1]) + " <jobId>")
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		switch strings.ToLower(tokens[1]) {
		case "start":
			if err := s.runtime.DataJobs().StartJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d started successfully", jobID))
		case "stop":
			if err := s.runtime.DataJobs().StopJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d stopped successfully", jobID))
		default:
			if err := s.runtime.DataJobs().DeleteJob(jobID); err != nil {
				return err
			}
			s.io.Println(fmt.Sprintf("Job %d deleted successfully", jobID))
		}
		return nil
	case "replay":
		if len(tokens) < 3 {
			s.io.Println("Usage: jobs replay <jobId> [--auto-start] [--reset-to-created]")
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		options, err := parseReplayOptions(tokens[3:])
		if err != nil {
			return err
		}
		if err := s.runtime.DataJobs().ReplayJob(jobID, options); err != nil {
			return err
		}
		s.io.Println(fmt.Sprintf("Job %d replay requested successfully", jobID))
		return nil
	default:
		s.io.Println("Usage: jobs list | jobs show <jobId> | jobs schema <jobId> | jobs start <jobId> | jobs stop <jobId> | jobs delete <jobId> | jobs replay <jobId> [--auto-start] [--reset-to-created]")
		return nil
	}
}

func (s *Shell) printJobs(options datajob.ListOptions) error {
	jobs, err := s.runtime.DataJobs().ListJobs(options)
	if err != nil {
		return err
	}

	headers := []string{"ID", "Name", "Type", "State", "Source", "Target"}
	rows := make([][]string, 0, len(jobs))
	for _, job := range jobs {
		rows = append(rows, []string{
			strconv.FormatInt(job.DataJobID, 10),
			orDash(job.DataJobName),
			orDash(job.DataJobType),
			orDash(job.DataTaskState),
			instanceDesc(job.SourceDS),
			instanceDesc(job.TargetDS),
		})
	}

	s.io.Println(util.FormatTable(headers, rows))
	s.io.Println(fmt.Sprintf("%d jobs", len(jobs)))
	return nil
}

func (s *Shell) printJob(jobID int64) error {
	job, err := s.runtime.DataJobs().GetJob(jobID)
	if err != nil {
		return err
	}

	s.io.Println("Job details:")
	s.io.Println("  ID: " + strconv.FormatInt(job.DataJobID, 10))
	s.io.Println("  Name: " + orDash(job.DataJobName))
	s.io.Println("  Description: " + orDash(job.DataJobDesc))
	s.io.Println("  Type: " + orDash(job.DataJobType))
	s.io.Println("  State: " + orDash(job.DataTaskState))
	s.io.Println("  Current Task Status: " + orDash(job.CurrTaskStatus))
	s.io.Println("  Lifecycle: " + orDash(job.LifeCycleState))
	s.io.Println("  User: " + orDash(job.UserName))
	s.io.Println("  Console Job ID: " + formatOptionalInt64(job.ConsoleJobID))
	s.io.Println("  Console Task State: " + orDash(job.ConsoleTaskState))
	s.io.Println("  Source: " + sourceSummary(job.SourceDS))
	s.io.Println("  Target: " + sourceSummary(job.TargetDS))
	s.io.Println("  Source Schema: " + orDash(job.SourceSchema))
	s.io.Println("  Target Schema: " + orDash(job.TargetSchema))
	s.io.Println("  Tasks: " + strconv.Itoa(len(job.DataTasks)))
	s.io.Println("  Has Exception: " + formatBool(job.HaveException))

	if len(job.DataTasks) > 0 {
		headers := []string{"Task ID", "Name", "Type", "Status", "Worker IP"}
		rows := make([][]string, 0, len(job.DataTasks))
		for _, task := range job.DataTasks {
			rows = append(rows, []string{
				strconv.FormatInt(task.DataTaskID, 10),
				orDash(task.DataTaskName),
				orDash(task.DataTaskType),
				orDash(task.DataTaskStatus),
				orDash(task.WorkerIP),
			})
		}
		s.io.Println("")
		s.io.Println(util.FormatTable(headers, rows))
	}
	return nil
}

func (s *Shell) printJobSchema(jobID int64) error {
	schema, err := s.runtime.DataJobs().GetJobSchema(jobID)
	if err != nil {
		return err
	}

	s.io.Println("Job schema:")
	s.io.Println("  Job ID: " + strconv.FormatInt(jobID, 10))
	s.io.Println("  Source Schema: " + orDash(schema.SourceSchema))
	s.io.Println("  Target Schema: " + orDash(schema.TargetSchema))
	s.io.Println("  Default Topic: " + orDash(schema.DefaultTopic))
	s.io.Println("  Default Topic Partition: " + formatOptionalInt64(int64(schema.DefaultTopicPartition)))
	s.io.Println("  Schema Whitelist Level: " + orDash(schema.SchemaWhiteListLevel))
	s.io.Println("  Source Schema Less Format: " + orDash(schema.SrcSchemaLessFormat))
	s.io.Println("  Target Schema Less Format: " + orDash(schema.DstSchemaLessFormat))
	if strings.TrimSpace(schema.MappingConfig) != "" {
		s.io.Println("")
		s.io.Println("Mapping Config:")
		s.io.Println(schema.MappingConfig)
	}
	return nil
}

func parseJobListOptions(tokens []string) (datajob.ListOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return datajob.ListOptions{}, err
	}

	result := datajob.ListOptions{}
	result.DataJobName, _ = popOption(options, "name", "data-job-name")
	result.DataJobType, _ = popOption(options, "type", "data-job-type")
	result.Desc, _ = popOption(options, "desc", "description")
	result.SourceInstanceID, err = parsePositiveInt64Option(options, "sourceInstanceId", "source-id", "source-instance-id")
	if err != nil {
		return datajob.ListOptions{}, err
	}
	result.TargetInstanceID, err = parsePositiveInt64Option(options, "targetInstanceId", "target-id", "target-instance-id")
	if err != nil {
		return datajob.ListOptions{}, err
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return datajob.ListOptions{}, err
	}
	return result, nil
}

func parseReplayOptions(tokens []string) (datajob.ReplayOptions, error) {
	options, err := parseFlagArgs(tokens)
	if err != nil {
		return datajob.ReplayOptions{}, err
	}

	var replayOptions datajob.ReplayOptions
	autoStart, err := parseBoolOption(options, "autoStart", "auto-start")
	if err != nil {
		return datajob.ReplayOptions{}, err
	}
	if autoStart != nil {
		replayOptions.AutoStart = *autoStart
	}
	resetToCreated, err := parseBoolOption(options, "resetToCreated", "reset-to-created")
	if err != nil {
		return datajob.ReplayOptions{}, err
	}
	if resetToCreated != nil {
		replayOptions.ResetToCreated = *resetToCreated
	}
	if err := ensureNoUnknownOptions(options); err != nil {
		return datajob.ReplayOptions{}, err
	}
	return replayOptions, nil
}

func instanceDesc(source *datajob.Source) string {
	if source == nil {
		return "-"
	}
	return orDash(source.InstanceDesc)
}

func sourceSummary(source *datajob.Source) string {
	if source == nil {
		return "-"
	}

	extras := make([]string, 0, 3)
	if strings.TrimSpace(source.DataSourceType) != "" {
		extras = append(extras, source.DataSourceType)
	}
	if strings.TrimSpace(source.HostType) != "" {
		extras = append(extras, source.HostType)
	}
	if strings.TrimSpace(source.Region) != "" {
		extras = append(extras, source.Region)
	}

	label := orDash(source.InstanceDesc)
	if len(extras) == 0 {
		return label
	}
	return label + " (" + strings.Join(extras, ", ") + ")"
}
