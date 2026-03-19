package repl

import (
	"cloudcanal-openapi-cli/internal/datajob"
	"cloudcanal-openapi-cli/internal/util"
	"strconv"
	"strings"
)

func (s *Shell) handleJobs(tokens []string) error {
	if len(tokens) < 2 {
		s.io.Println(s.usageJobsGroup())
		return nil
	}

	switch strings.ToLower(tokens[1]) {
	case "list":
		options, err := parseJobListOptions(tokens[2:])
		if err != nil {
			return err
		}
		return s.printJobs(options)
	case "create":
		options, err := parseFlagArgs(tokens[2:])
		if err != nil {
			return err
		}
		var request datajob.CreateJobRequest
		if err := decodeBodyOptions(options, &request); err != nil {
			return err
		}
		if err := ensureNoUnknownOptions(options); err != nil {
			return err
		}
		result, err := s.runtime.DataJobs().CreateJob(request)
		if err != nil {
			return err
		}
		return s.printJobCreateResult(result)
	case "show":
		if len(tokens) != 3 {
			s.io.Println(s.usageJobAction("show"))
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		return s.printJob(jobID)
	case "schema":
		if len(tokens) != 3 {
			s.io.Println(s.usageJobAction("schema"))
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		return s.printJobSchema(jobID)
	case "start", "stop", "delete":
		if len(tokens) != 3 {
			s.io.Println(s.usageJobAction(strings.ToLower(tokens[1])))
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
			return s.printActionResult("job.started", "job", "started", jobID)
		case "stop":
			if err := s.runtime.DataJobs().StopJob(jobID); err != nil {
				return err
			}
			return s.printActionResult("job.stopped", "job", "stopped", jobID)
		default:
			if err := s.runtime.DataJobs().DeleteJob(jobID); err != nil {
				return err
			}
			return s.printActionResult("job.deleted", "job", "deleted", jobID)
		}
	case "replay":
		if len(tokens) < 3 {
			s.io.Println(s.usageJobReplay())
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
		return s.printActionResult("job.replayed", "job", "replayed", jobID)
	case "attach-incre-task", "detach-incre-task":
		if len(tokens) != 3 {
			s.io.Println(s.usageJobAction(strings.ToLower(tokens[1])))
			return nil
		}
		jobID, err := parsePositiveInt64(tokens[2], "jobId")
		if err != nil {
			return err
		}
		if strings.EqualFold(tokens[1], "attach-incre-task") {
			if err := s.runtime.DataJobs().AttachIncreJob(jobID); err != nil {
				return err
			}
			return s.printActionResult("job.increAttached", "job", "attach-incre-task", jobID)
		}
		if err := s.runtime.DataJobs().DetachIncreJob(jobID); err != nil {
			return err
		}
		return s.printActionResult("job.increDetached", "job", "detach-incre-task", jobID)
	case "update-incre-pos":
		options, err := parseFlagArgs(tokens[2:])
		if err != nil {
			return err
		}
		var request datajob.UpdateIncrePosRequest
		if err := decodeBodyOptions(options, &request); err != nil {
			return err
		}
		if err := ensureNoUnknownOptions(options); err != nil {
			return err
		}
		result, err := s.runtime.DataJobs().UpdateIncrePos(request)
		if err != nil {
			return err
		}
		return s.printUpdateIncrePosResult(request.TaskID, result)
	default:
		s.printUnknownSubcommand("jobs", tokens[1], jobsSubcommands, s.usageJobsGroup())
		return nil
	}
}

func (s *Shell) printJobs(options datajob.ListOptions) error {
	jobs, err := s.runtime.DataJobs().ListJobs(options)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(jobs)
	}

	headers := []string{s.label("id"), s.label("name"), s.label("type"), s.label("state"), s.label("source"), s.label("target")}
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
	s.io.Println(s.countLabel("jobs", len(jobs)))
	return nil
}

func (s *Shell) printJob(jobID int64) error {
	job, err := s.runtime.DataJobs().GetJob(jobID)
	if err != nil {
		return err
	}
	if s.isJSONOutput() {
		return s.printJSON(job)
	}

	s.io.Println(s.sectionTitle("job.details"))
	s.io.Println(s.line(s.label("id"), strconv.FormatInt(job.DataJobID, 10)))
	s.io.Println(s.line(s.label("name"), orDash(job.DataJobName)))
	s.io.Println(s.line(s.label("description"), orDash(job.DataJobDesc)))
	s.io.Println(s.line(s.label("type"), orDash(job.DataJobType)))
	s.io.Println(s.line(s.label("state"), orDash(job.DataTaskState)))
	s.io.Println(s.line(s.label("currentTaskStatus"), orDash(job.CurrTaskStatus)))
	s.io.Println(s.line(s.label("lifecycle"), orDash(job.LifeCycleState)))
	s.io.Println(s.line(s.label("user"), orDash(job.UserName)))
	s.io.Println(s.line(s.label("consoleJobId"), formatOptionalInt64(job.ConsoleJobID)))
	s.io.Println(s.line(s.label("consoleTaskState"), orDash(job.ConsoleTaskState)))
	s.io.Println(s.line(s.label("source"), sourceSummary(job.SourceDS)))
	s.io.Println(s.line(s.label("target"), sourceSummary(job.TargetDS)))
	s.io.Println(s.line(s.label("sourceSchema"), orDash(job.SourceSchema)))
	s.io.Println(s.line(s.label("targetSchema"), orDash(job.TargetSchema)))
	s.io.Println(s.line(s.label("tasks"), strconv.Itoa(len(job.DataTasks))))
	s.io.Println(s.line(s.label("hasException"), formatBool(job.HaveException)))

	if len(job.DataTasks) > 0 {
		headers := []string{s.label("taskId"), s.label("name"), s.label("type"), s.label("state"), s.label("workerIP")}
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
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"jobId":  jobID,
			"schema": schema,
		})
	}

	s.io.Println(s.sectionTitle("job.schema"))
	s.io.Println(s.line(s.label("id"), strconv.FormatInt(jobID, 10)))
	s.io.Println(s.line(s.label("sourceSchema"), orDash(schema.SourceSchema)))
	s.io.Println(s.line(s.label("targetSchema"), orDash(schema.TargetSchema)))
	s.io.Println(s.line(s.label("defaultTopic"), orDash(schema.DefaultTopic)))
	s.io.Println(s.line(s.label("defaultTopicPartition"), formatOptionalInt64(int64(schema.DefaultTopicPartition))))
	s.io.Println(s.line(s.label("schemaWhitelistLevel"), orDash(schema.SchemaWhiteListLevel)))
	s.io.Println(s.line(s.label("srcSchemaLessFormat"), orDash(schema.SrcSchemaLessFormat)))
	s.io.Println(s.line(s.label("dstSchemaLessFormat"), orDash(schema.DstSchemaLessFormat)))
	if strings.TrimSpace(schema.MappingConfig) != "" {
		s.io.Println("")
		s.io.Println(s.sectionTitle("job.mappingConfig"))
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
	return orDash(util.MaskSensitiveText(source.InstanceDesc))
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

	label := orDash(util.MaskSensitiveText(source.InstanceDesc))
	if len(extras) == 0 {
		return label
	}
	return label + " (" + strings.Join(extras, ", ") + ")"
}

func (s *Shell) printJobCreateResult(result datajob.CreateJobResult) error {
	message := "Job created successfully"
	if s.isChinese() {
		message = "任务创建成功"
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"resource": "job",
			"action":   "created",
			"jobId":    result.JobID,
			"data":     result.Data,
			"message":  message,
		})
	}
	s.io.Println(message)
	if result.JobID != "" {
		s.io.Println(s.line(s.label("jobId"), result.JobID))
	}
	return nil
}

func (s *Shell) printUpdateIncrePosResult(taskID int64, result datajob.UpdateIncrePosResult) error {
	message := "Increment position updated successfully"
	if s.isChinese() {
		message = "增量位点更新成功"
	}
	if s.isJSONOutput() {
		return s.printJSON(map[string]any{
			"resource": "task-position",
			"action":   "updated",
			"taskId":   taskID,
			"data":     result.Data,
			"message":  message,
		})
	}
	s.io.Println(message)
	s.io.Println(s.line(s.label("taskId"), strconv.FormatInt(taskID, 10)))
	if result.Data != "" {
		s.io.Println(s.line(s.label("result"), result.Data))
	}
	return nil
}
