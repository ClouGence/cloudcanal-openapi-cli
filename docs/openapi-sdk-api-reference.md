# CloudCanal OpenAPI SDK 对照 CLI 接口文档

本文档按 `cloudcanal-openapi-sdk` 当前支持的 API 生成，并给出 `cloudcanal-openapi-cli` 中对应的命令实现。

说明：

- SDK 当前共覆盖 7 个模块、25 个 OpenAPI path。
- CLI 已按 SDK 一比一补齐这些接口。
- `datasources show <dataSourceId>` 是 CLI 便捷命令，内部仍复用 `datasource/listds`，不是 SDK 的独立 API。
- 复杂请求统一支持 `--body '{...}'` 或 `--body-file FILE.json`。`datasource add` 额外支持 `--security-file`、`--secret-file`。

## Cluster

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListClusterRequest` | `/cloudcanal/console/api/v1/openapi/cluster/listclusters` | `cloudcanal clusters list [--name NAME] [--desc DESC] [--cloud CLOUD] [--region REGION]` | `cloudOrIdcName`, `clusterDescLike`, `clusterNameLike`, `region` | `id`, `clusterName`, `region`, `cloudOrIdcName`, `workerCount`, `runningCount`, `abnormalCount`, `ownerName` |

## ConsoleJob

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `QueryConsoleJobRequest` | `/cloudcanal/console/api/v1/openapi/consolejob/queryconsolejob` | `cloudcanal consolejobs show <consoleJobId>` | `consoleJobId` | `id`, `label`, `taskState`, `jobToken`, `resourceType`, `resourceId`, `taskVOList` |

## Constant

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListSpecsRequest` | `/cloudcanal/console/api/v1/openapi/constant/listspecs` | `cloudcanal job-config specs --type TYPE [--initial-sync=true\|false] [--short-term-sync=true\|false]` | `dataJobType`, `initialSync`, `shortTermSync` | `id`, `specKind`, `specKindCn`, `spec`, `fullMemoryMb`, `increMemoryMb`, `checkMemoryMb` |
| `TransformJobTypeRequest` | `/cloudcanal/console/api/v1/openapi/constant/transformjobtype` | `cloudcanal job-config transform-job-type --source-type TYPE --target-type TYPE` | `sourceType`, `targetType` | SDK 未定义固定响应类型，CLI 保留原始 `data` JSON |

## DataJob

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListJobsRequest` | `/cloudcanal/console/api/v1/openapi/datajob/list` | `cloudcanal jobs list [--name NAME] [--type TYPE] [--desc DESC] [--source-id ID] [--target-id ID]` | `dataJobName`, `dataJobType`, `desc`, `sourceInstanceId`, `targetInstanceId` | `dataJobId`, `dataJobName`, `dataTaskState`, `sourceDsVO`, `targetDsVO` |
| `QueryJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/queryjob` | `cloudcanal jobs show <jobId>` | `jobId` | `ApiDataJobDO` 详情 |
| `QueryJobSchemaRequest` | `/cloudcanal/console/api/v1/openapi/datajob/queryjobschemabyid` | `cloudcanal jobs schema <jobId>` | `jobId` | `sourceSchema`, `targetSchema`, `mappingConfig`, `defaultTopic`, `defaultTopicPartition` |
| `AddJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/create` | `cloudcanal jobs create --body-file FILE.json` | 见下方 `jobs create` 请求模板 | `data` 为创建结果，CLI 同时输出 `jobId` |
| `StartJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/start` | `cloudcanal jobs start <jobId>` | `jobId` | 通用成功/失败响应 |
| `StopJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/stop` | `cloudcanal jobs stop <jobId>` | `jobId` | 通用成功/失败响应 |
| `DeleteJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/delete` | `cloudcanal jobs delete <jobId>` | `jobId` | 通用成功/失败响应 |
| `ReplayJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/replay` | `cloudcanal jobs replay <jobId> [--auto-start] [--reset-to-created]` | `jobId`, `autoStart`, `resetToCreated` | 通用成功/失败响应 |
| `AttachIncreJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/attachincretask` | `cloudcanal jobs attach-incre-task <jobId>` | `jobId` | 通用成功/失败响应 |
| `DetachIncreJobRequest` | `/cloudcanal/console/api/v1/openapi/datajob/detachincretask` | `cloudcanal jobs detach-incre-task <jobId>` | `jobId` | 通用成功/失败响应 |
| `UpdateIncrePosRequest` | `/cloudcanal/console/api/v1/openapi/datajob/updateincrepos` | `cloudcanal jobs update-incre-pos --body-file FILE.json` | 见下方 `jobs update-incre-pos` 请求模板 | `data` 为更新结果 |

### `jobs create` 请求模板

```json
{
  "clusterId": 1,
  "srcDsId": 391,
  "dstDsId": 16,
  "srcHostType": "PRIVATE",
  "dstHostType": "PUBLIC",
  "jobType": "SYNC",
  "dataJobDesc": "api created",
  "specId": 18,
  "autoStart": true,
  "structMigration": true,
  "initialSync": true,
  "shortTermSync": false,
  "shortTermNum": 3,
  "filterDDL": true,
  "srcSchema": "{\"schema\":\"src\"}",
  "dstSchema": "{\"schema\":\"dst\"}",
  "mappingDef": "[{\"method\":\"DB_DB\"}]",
  "schemaWhiteListLevel": "TABLE",
  "checkOnce": false,
  "checkPeriod": false,
  "fullPeriod": false
}
```

SDK 里还有更多高级字段，CLI 已全部支持，例如：

- `srcCaseSensitiveType`, `dstCaseSensitiveType`
- `srcDsCharset`, `tarDsCharset`
- `keyConflictStrategy`
- `checkPeriodCronExpr`, `fullPeriodCronExpr`
- `dstMqDefaultTopic`, `dstMqDefaultTopicPartitions`
- `dstMqDdlTopic`, `dstMqDdlTopicPartitions`
- `srcSchemaLessFormat`, `dstSchemaLessFormat`
- `originDecodeMsgFormat`
- `dstCkTableEngine`, `dstSrOrDorisTableModel`
- `kafkaConsumerGroupId`, `kuduNumReplicas`
- `srcRocketMqGroupId`
- `srcRabbitMqVhost`, `srcRabbitExchange`
- `dstRabbitMqVhost`, `dstRabbitExchange`
- `obTenant`, `dbHeartbeatEnable`

### `jobs update-incre-pos` 请求模板

```json
{
  "taskId": 320,
  "posType": "MYSQL_LOG_FILE_POS",
  "journalFile": "binlog.000491",
  "filePosition": 794891
}
```

不同数据库位点字段按 SDK 原样支持：

- MySQL 类：`journalFile`, `filePosition`, `gtidPosition`, `positionTimestamp`, `serverId`
- PostgreSQL/SQL Server：`lsn`
- Oracle：`scn`, `scnIndex`
- 通用：`commonPosStr`
- HANA：`dataId`, `transactionId`

## DataSource

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListDsRequest` | `/cloudcanal/console/api/v1/openapi/datasource/listds` | `cloudcanal datasources list [--id ID] [--type TYPE] [--deploy-type TYPE] [--host-type TYPE] [--lifecycle STATE]` | `dataSourceId`, `deployType`, `hostType`, `lifeCycleState`, `type` | `id`, `instanceId`, `dataSourceType`, `hostType`, `deployType`, `lifeCycleState` |
| `AddDsRequest` | `/cloudcanal/console/api/v1/openapi/datasource/addds` | `cloudcanal datasources add --body-file FILE.json [--security-file FILE] [--secret-file FILE]` | `dataSourceAddData` + 可选文件 part | `data` 为创建结果 |
| `DeleteDsRequest` | `/cloudcanal/console/api/v1/openapi/datasource/deleteds` | `cloudcanal datasources delete <dataSourceId>` | `dataSourceId` | 通用成功/失败响应 |

### `datasources add` 请求体

`datasources add` 支持两种 body 形式：

1. 直接传 `ApiDsAddData` JSON
2. 传完整包装对象：

```json
{
  "dataSourceAddData": {
    "type": "MYSQL",
    "host": "127.0.0.1:3306",
    "privateHost": "127.0.0.1:3306",
    "hostType": "PRIVATE",
    "deployType": "ALIYUN",
    "region": "cn-hangzhou",
    "instanceDesc": "mysql source",
    "account": "root",
    "password": "secret",
    "securityType": "USER_PASSWD"
  },
  "securityFilePath": "/path/to/security.pem",
  "secretFilePath": "/path/to/secret.key"
}
```

CLI 也支持把文件路径拆到 flag：

```bash
cloudcanal datasources add \
  --body-file add-datasource.json \
  --security-file /path/to/security.pem \
  --secret-file /path/to/secret.key
```

`ApiDsAddData` 的高级字段也都支持，例如：

- `instanceId`, `version`, `dbName`
- `accessKey`, `secretKey`, `clientTrustStorePassword`
- `clusterIds`, `lifeCycleState`, `driver`, `connectType`
- `dsKvConfigs`, `extraData`, `parentDsId`

## Worker

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListWorkerRequest` | `/cloudcanal/console/api/v1/openapi/worker/listworkers` | `cloudcanal workers list --cluster-id ID [--source-id ID] [--target-id ID]` | `clusterId`, `sourceInstanceId`, `targetInstanceId` | `id`, `workerName`, `workerState`, `workerType`, `privateIp`, `healthLevel`, `workerLoad`，以及 SDK 中更多扩展字段 |
| `StartWorkerRequest` | `/cloudcanal/console/api/v1/openapi/worker/startWorker` | `cloudcanal workers start <workerId>` | `workerId` | 通用成功/失败响应 |
| `StopWorkerRequest` | `/cloudcanal/console/api/v1/openapi/worker/stopWorker` | `cloudcanal workers stop <workerId>` | `workerId` | 通用成功/失败响应 |
| `DeleteWorkerRequest` | `/cloudcanal/console/api/v1/openapi/worker/deleteWorker` | `cloudcanal workers delete <workerId>` | `workerId` | 通用成功/失败响应 |
| `ModifyMemOverSoldRequest` | `/cloudcanal/console/api/v1/openapi/worker/modifyMemOverSoldPercent` | `cloudcanal workers modify-mem-oversold <workerId> --percent N` | `workerId`, `memOverSoldPercent` | 通用成功/失败响应 |
| `UpdateWorkerAlertRequest` | `/cloudcanal/console/api/v1/openapi/worker/updateWorkerAlertConfig` | `cloudcanal workers update-alert <workerId> --phone=true\|false --email=true\|false --im=true\|false --sms=true\|false` | `workerId`, `phone`, `email`, `im`, `sms` | 通用成功/失败响应 |

## Schema

| SDK API | Path | CLI 命令 | 关键请求字段 | 返回摘要 |
| --- | --- | --- | --- | --- |
| `ListTransObjsByMetaRequest` | `/cloudcanal/console/api/v1/openapi/schema/listTransObjsByMeta` | `cloudcanal schemas list-trans-objs-by-meta [--src-db NAME] [--src-schema NAME] [--src-trans-obj NAME] [--dst-db NAME] [--dst-schema NAME] [--dst-tran-obj NAME]` | `srcDb`, `srcSchema`, `srcTransObj`, `dstDb`, `dstSchema`, `dstTranObj` | `dataJobId`, `dataJobName`, `srcFullTransferObjName`, `dstFullTransferObjName`, `srcDsType`, `dstDsType` |

## 通用返回规则

- 所有 OpenAPI 响应都遵循 SDK 的 `CcResponse` 约定。
- `code == "1"` 表示成功。
- `code != "1"` 时，CLI 会直接输出 `msg` 作为错误信息。
