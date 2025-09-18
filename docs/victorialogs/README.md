---
build:
  list: never
  publishResources: false
  render: never
sitemap:
  disable: true
---
VictoriaLogs is [open source](https://github.com/VictoriaMetrics/VictoriaLogs/) user-friendly database for logs
from [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics/).

## Features

VictoriaLogs provides the following features:

- It is resource-efficient and fast. It uses up to 30x less RAM and up to 15x less disk space than other solutions such as Elasticsearch and Grafana Loki.
  See [these benchmarks](https://docs.victoriametrics.com/victorialogs/#benchmarks) and [this article](https://itnext.io/how-do-open-source-solutions-for-logs-work-elasticsearch-loki-and-victorialogs-9f7097ecbc2f) for details.
  See also [the post from a happy user, who replaced 27-node Elasticsearch with a single-node VictoriaLogs](https://aus.social/@phs/114583927679254536),
  [this post from happy users, who replaced Loki with VictoriaLogs](https://www.truefoundry.com/blog/victorialogs-vs-loki)
  and [this post from a happy user who replaced grep with VictoriaLogs](https://chronicles.mad-scientist.club/tales/grepping-logs-remains-terrible/).
- VictoriaLogs' capacity and performance scales linearly with the available resources (CPU, RAM, disk IO, disk space).
  It runs smoothly on Raspberry PI and on servers with hundreds of CPU cores and terabytes of RAM.
  It can scale horizontally to hundreds of nodes in [cluster mode](https://docs.victoriametrics.com/victorialogs/cluster/).
- It can accept logs from popular log collectors. See [these docs](https://docs.victoriametrics.com/victorialogs/data-ingestion/).
- It is much easier to set up and operate compared to Elasticsearch and Grafana Loki, since it is a single zero-config executable.
  See [these docs](https://docs.victoriametrics.com/victorialogs/quickstart/).
- It provides easy yet powerful query language, which supports fast full-text search, fast advanced analytics and fast data extraction and transformation at query time.
  See [LogsQL docs](https://docs.victoriametrics.com/victorialogs/logsql/).
- It provides [built-in web UI](https://docs.victoriametrics.com/victorialogs/querying/#web-ui) for logs' exploration.
- It provides [Grafana plugin](https://docs.victoriametrics.com/victorialogs/victorialogs-datasource/) for building arbitrary dashboards in Grafana.
- It provides [interactive command-line tool for querying VictoriaLogs](https://docs.victoriametrics.com/victorialogs/querying/vlogscli/).
- It can be seamlessly combined with good old Unix tools for log analysis such as `grep`, `less`, `sort`, `jq`, etc.
  See [these docs](https://docs.victoriametrics.com/victorialogs/querying/#command-line) for details.
- It support [log fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) with high cardinality (e.g. high number of unique values) such as `trace_id`, `user_id` and `ip`.
- It is optimized for logs with hundreds of fields (aka [`wide events`](https://jeremymorrell.dev/blog/a-practitioners-guide-to-wide-events/)).
- It supports multitenancy - see [these docs](https://docs.victoriametrics.com/victorialogs/#multitenancy).
- It supports out-of-order logs' ingestion aka backfilling.
- It supports live tailing for newly ingested logs. See [these docs](https://docs.victoriametrics.com/victorialogs/querying/#live-tailing).
- It supports selecting surrounding logs in front and after the selected logs. See [these docs](https://docs.victoriametrics.com/victorialogs/logsql/#stream_context-pipe).
- It supports alerting - see [these docs](https://docs.victoriametrics.com/victorialogs/vmalert/).
- It fits well [RUM](https://en.wikipedia.org/wiki/Real_user_monitoring) and [SIEM](https://en.wikipedia.org/wiki/Security_information_and_event_management) use cases.

See also [articles about VictoriaLogs](https://docs.victoriametrics.com/victorialogs/articles/).

If you have questions about VictoriaLogs, then read [this FAQ](https://docs.victoriametrics.com/victorialogs/faq/).
Also feel free asking any questions at [VictoriaMetrics community Slack chat](https://victoriametrics.slack.com/),
you can join it via [Slack Inviter](https://slack.victoriametrics.com/).

See [quick start docs](https://docs.victoriametrics.com/victorialogs/quickstart/) for start working with VictoriaLogs.

If you want playing with VictoriaLogs web UI and [LogsQL](https://docs.victoriametrics.com/victorialogs/logsql/) query language,
then go to [VictoriaLogs demo playground](https://play-vmlogs.victoriametrics.com/) and
to [Grafana plugin playground for VictoriaLogs](https://play-grafana.victoriametrics.com/d/be5zidev72m80f/k8s-logs-via-victorialogs).

## Tuning

- No need in tuning for VictoriaLogs - it uses reasonable defaults for command-line flags, which are automatically adjusted for the available CPU and RAM resources.
- No need in tuning for Operating System - VictoriaLogs is optimized for default OS settings.
  The only option is increasing the limit on [the number of open files in the OS](https://medium.com/@muhammadtriwibowo/set-permanently-ulimit-n-open-files-in-ubuntu-4d61064429a).
- The recommended filesystem is `ext4`, the recommended persistent storage is [persistent HDD-based disk on GCP](https://cloud.google.com/compute/docs/disks/#pdspecs),
  since it is protected from hardware failures via internal replication and it can be [resized on the fly](https://cloud.google.com/compute/docs/disks/add-persistent-disk#resize_pd).
  If you plan to store more than 1TB of data on `ext4` partition or plan extending it to more than 16TB,
  then the following options are recommended to pass to `mkfs.ext4`:

  ```sh
  mkfs.ext4 ... -O 64bit,huge_file,extent -T huge
  ```

## Monitoring

VictoriaLogs exposes internal metrics in Prometheus exposition format at `http://localhost:9428/metrics` page.
It is recommended to set up monitoring of these metrics via VictoriaMetrics
(see [these docs](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#how-to-scrape-prometheus-exporters-such-as-node-exporter)),
vmagent (see [these docs](https://docs.victoriametrics.com/victoriametrics/vmagent/#how-to-collect-metrics-in-prometheus-format)) or via Prometheus.

See [metrics reference](https://docs.victoriametrics.com/victorialogs/metrics/) for a comprehensive list of all available metrics with detailed descriptions.

We recommend installing Grafana dashboard for [VictoriaLogs single-node](https://grafana.com/grafana/dashboards/22084) or [cluster](https://grafana.com/grafana/dashboards/23274).

We recommend setting up [alerts](https://github.com/VictoriaMetrics/VictoriaLogs/blob/master/deployment/docker/rules/alerts-vlogs.yml)
via [vmalert](https://docs.victoriametrics.com/victoriametrics/vmalert/) or via Prometheus.

VictoriaLogs emits its own logs to stdout. It is recommended to investigate these logs during troubleshooting.

## Upgrading

It is safe upgrading VictoriaLogs to new versions unless [release notes](https://docs.victoriametrics.com/victorialogs/changelog/) say otherwise.
It is safe to skip multiple versions during the upgrade unless [release notes](https://docs.victoriametrics.com/victorialogs/changelog/) say otherwise.
It is recommended to perform regular upgrades to the latest version, since it may contain important bug fixes, performance optimizations or new features.

It is also safe to downgrade to older versions unless [release notes](https://docs.victoriametrics.com/victorialogs/changelog/) say otherwise.

The following steps must be performed during the upgrade / downgrade procedure:

- Send `SIGINT` signal to VictoriaLogs process in order to gracefully stop it.
  See [how to send signals to processes](https://stackoverflow.com/questions/33239959/send-signal-to-process-from-command-line).
- Wait until the process stops. This can take a few seconds.
- Start the upgraded VictoriaLogs.

## Retention

By default, VictoriaLogs stores log entries with timestamps in the time range `[now-7d, now]`, while dropping logs outside the given time range.
E.g. it uses the retention of 7 days. The retention can be configured with `-retentionPeriod` command-line flag.
This flag accepts values starting from `1d` (one day) up to `100y` (100 years). See [these docs](https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations)
for the supported duration formats.

For example, the following command starts VictoriaLogs with the retention of 8 weeks:

```sh
/path/to/victoria-logs -retentionPeriod=8w
```

See also [retention by disk space usage](https://docs.victoriametrics.com/victorialogs/#retention-by-disk-space-usage).

VictoriaLogs stores the [ingested](https://docs.victoriametrics.com/victorialogs/data-ingestion/) logs in per-day partition directories.
It automatically drops partition directories outside the configured retention.

VictoriaLogs automatically drops logs at [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/) stage
if they have timestamps outside the configured retention. A sample of dropped logs is logged with `WARN` message in order to simplify troubleshooting.
The `vl_rows_dropped_total` [metric](https://docs.victoriametrics.com/victorialogs/metrics/#vl_rows_dropped_total) is incremented each time an ingested log entry is dropped because of timestamp outside the retention.
It is recommended to set up the following alerting rule at [vmalert](https://docs.victoriametrics.com/victoriametrics/vmalert/) in order to be notified
when logs with wrong timestamps are ingested into VictoriaLogs:

```metricsql
rate(vl_rows_dropped_total[5m]) > 0
```

By default, VictoriaLogs doesn't accept log entries with timestamps bigger than `now+2d`, e.g. 2 days in the future.
If you need accepting logs with bigger timestamps, then specify the desired "future retention" via `-futureRetention` command-line flag.
This flag accepts values starting from `1d`. See [these docs](https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations)
for the supported duration formats.

For example, the following command starts VictoriaLogs, which accepts logs with timestamps up to a year in the future:

```sh
/path/to/victoria-logs -futureRetention=1y
```

## Retention by disk space usage

VictoriaLogs can be configured to automatically drop older per-day partitions based on disk space usage using one of two approaches:

### Absolute disk space limit

Use the `-retention.maxDiskSpaceUsageBytes` command-line flag to set a fixed threshold. VictoriaLogs will drop old per-day partitions
if the total size of data at [`-storageDataPath` directory](https://docs.victoriametrics.com/victorialogs/#storage) becomes bigger than the specified limit.
For example, the following command starts VictoriaLogs, which drops old per-day partitions if the total [storage](https://docs.victoriametrics.com/victorialogs/#storage) size becomes bigger than `100GiB`:

```sh
/path/to/victoria-logs -retention.maxDiskSpaceUsageBytes=100GiB
```

### Percentage-based disk space limit

Use the `-retention.maxDiskUsagePercent` command-line flag to set a dynamic threshold based on the filesystem's total capacity.
VictoriaLogs will drop old per-day partitions if the filesystem containing the [`-storageDataPath` directory](https://docs.victoriametrics.com/victorialogs/#storage) exceeds the specified percentage usage.
For example, the following command starts VictoriaLogs, which drops old per-day partitions if the filesystem usage exceeds 80%:

```sh
/path/to/victoria-logs -retention.maxDiskUsagePercent=80
```

This approach is particularly useful in environments where the total disk capacity may vary (e.g., cloud environments with resizable volumes)
or when you want to maintain a consistent percentage of free space regardless of the total disk size.

**Important:** The `-retention.maxDiskSpaceUsageBytes` and `-retention.maxDiskUsagePercent` flags are mutually exclusive.
VictoriaLogs will refuse to start if both flags are set simultaneously.

VictoriaLogs usually compresses logs by 10x or more times. This means that VictoriaLogs can store more than a terabyte of uncompressed
logs when it runs with `-retention.maxDiskSpaceUsageBytes=100GiB` or when using percentage-based retention on a large filesystem.

VictoriaLogs keeps at least two last days of data in order to guarantee that the logs for the last day can be returned in queries.
This means that the total disk space usage may exceed the configured threshold if the size of the last two days of data
exceeds the limit.

The [`-retentionPeriod`](https://docs.victoriametrics.com/victorialogs/#retention) is applied independently to the disk space usage limits. This means that
VictoriaLogs automatically drops logs older than 7 days by default if only a disk space usage flag is set.
Set the `-retentionPeriod` to some big value (e.g. `100y` - 100 years) if logs shouldn't be dropped because of time-based retention.
For example:

```sh
/path/to/victoria-logs -retention.maxDiskSpaceUsageBytes=10TiB -retentionPeriod=100y
```

or

```sh
/path/to/victoria-logs -retention.maxDiskUsagePercent=85 -retentionPeriod=100y
```

## Storage

VictoriaLogs stores all its data in a single directory - `victoria-logs-data`. The path to the directory can be changed via `-storageDataPath` command-line flag.
For example, the following command starts VictoriaLogs, which stores the data at `/var/lib/victoria-logs`:

```sh
/path/to/victoria-logs -storageDataPath=/var/lib/victoria-logs
```

VictoriaLogs automatically creates the `-storageDataPath` directory on the first run if it is missing. VictoriaLogs stores logs
per every day into a spearate subdirectory (aka per-day partition). See [partitions lifecycle](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) for details.

VictoriaLogs switches to cluster mode if `-storageNode` command-line flag is specified:

- It stops storing the ingested logs locally in cluster mode. It spreads them evenly among `vlstorage` nodes specified via the `-storageNode` command-line flag.
- It stops querying the locally stored logs in cluster mode. It queries `vlstorage` nodes specified via `-storageNode` command-line flag.

See [cluster mode docs](https://docs.victoriametrics.com/victorialogs/cluster/) for details.

## Partitions lifecycle

The ingested logs are stored in per-day subdirectories (partitions) at the `<-storageDataPath>/partitions/` directory. The per-day subdirectories have `YYYYMMDD` names.
For example, the directory with the name `20250418` contains logs with [`_time` field](https://docs.victoriametrics.com/victorialogs/keyconcepts/#time-field) values
at April 18, 2025 UTC. This allows flexible data management.

For example, old per-day data is automatically and quickly deleted according to the provided [retention policy](https://docs.victoriametrics.com/victorialogs/#retention) by removing the corresponding per-day subdirectory (partition).

VictoriaLogs supports the following HTTP API endpoints at `victoria-logs:9428` address for managing partitions:

- `/internal/partition/attach?name=YYYYMMDD` - attaches the partition directory with the given name `YYYYMMDD` to VictoriaLogs,
  so it becomes visible for querying and can be used for data ingestion.
  The directory must be placed inside `<-storageDataPath>/partitions` and it must contain valid data for the given `YYYYMMDD` day.
- `/internal/partition/detach?name=YYYYMMDD` - detaches the partition directory with the given name `YYYYMMDD` from VictoriaLogs,
  so it is no longer visible for querying and cannot be used for data ingestion.
  The `/internal/partition/detach` endpoint waits until all the concurrently executed queries stop reading the data from the detached partition
  before returning. This allows safe on-disk manipulions of the detached partitions by external tools after returning from the `/internal/partition/detach` endpoint.
  Detached partitions are automatically attached after VictoriaLogs restart if the corresponding subdirectories at `<-storageDataPath>/partitions/` aren't removed.
- `/internal/partition/list` - returns JSON-encoded list of currently active partitions, which can be passed to `/internal/partition/detach` endpoint via `name` query arg.
- `/internal/partition/snapshot/create?name=YYYYMMDD` - creates a [snapshot](https://medium.com/@valyala/how-victoriametrics-makes-instant-snapshots-for-multi-terabyte-time-series-data-e1f3fb0e0282)
  for the partition for the given day `YYYYMMDD`. The endpoint returns a JSON string with the absolute filesystem path to the created snapshot. It is safe to make backups from
  the created snapshots according to [these instructions](https://docs.victoriametrics.com/victorialogs/#backup-and-restore). It is safe removing the created snapshots with `rm -rf` command.
  It is recommended removing unneeded snapshots on a regular basis in order to free up storage space occupied by these snapshots.
- `/internal/partition/snapshot/list` - returns JSON-encoded list of absolute paths to per-day partition snapshots created via `/internal/partition/snapshot/create`.

These endpoints can be protected from unauthorized access via `-partitionManageAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).

These endpoints can be used also for setting up automated multi-tier storage schemes where recently ingested logs are stored to VictoriaLogs instances
with fast NVMe (SSD) disks, while historical logs are gradully migrated to VictoriaLogs instances with slower, but bigger and less expensive HDD disks.
This scheme can be implemented with the following simple cron job, which must run once per day:

1. To make a snapshot for the older day stored at NVMe via `/internal/partition/snapshot/create?name=YYYYMMDD` endpoint.
1. To copy the snapshot to the `<-storageDataPath>/partitions/YYYYMMDD` directory at VictoriaMetrics with HDD via [`rsync`](https://en.wikipedia.org/wiki/Rsync).
1. To detach the copied partition from the VictoriaLogs with NVMe via `/internal/partition/detach?name=YYYYMMDD` endpoint.
1. To attach the copied partition to the VictoriaLogs with HDD via `/internal/partition/attach?name=YYYYMMDD` endpoint.
1. To delete the copied partition directory from the VictoriaLogs with NVMe via `rm -rf <-storageDataPath>/partitions/YYYYMMDD` command.

All the VictoriaLogs instances with NVMe and HDD disks can be queried simultaneously via `vlselect` component of [VictoriaLogs cluster](https://docs.victoriametrics.com/victorialogs/cluster/),
since [single-node VictoriaLogs instances can be a part of cluster](https://docs.victoriametrics.com/victorialogs/cluster/#single-node-and-cluster-mode-duality).

## Forced merge

VictoriaLogs performs data compactions in background in order to keep good performance characteristics when accepting new data.
These compactions (merges) are performed independently on per-day partitions.
This means that compactions are stopped for per-day partitions if no new data is ingested into these partitions.
Sometimes it is necessary to trigger compactions for old partitions. In this case forced compaction may be initiated on the specified per-day partition
by sending request to `/internal/force_merge?partition_prefix=YYYYMMDD`,
where `YYYYMMDD` is per-day partition name. For example, `http://victoria-logs:9428/internal/force_merge?partition_prefix=20240921` would initiate forced
merge for September 21, 2024 partition. The call to `/internal/force_merge` returns immediately, while the corresponding forced merge continues running in background.

Forced merges may require additional CPU, disk IO and storage space resources. It is unnecessary to run forced merge under normal conditions,
since VictoriaLogs automatically performs optimal merges in background when new data is ingested into it.

The `/internal/force_merge` endpoint can be protected from unauthorized access via `-forceMergeAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).

## Forced flush

VictoriaLogs puts the recently [ingested logs](https://docs.victoriametrics.com/victorialogs/data-ingestion/) into in-memory buffers,
which aren't available for [querying](https://docs.victoriametrics.com/victorialogs/querying/) for up to a second.
If you need querying logs immediately after their ingestion, then the `/internal/force_flush` HTTP endpoint must be requested
before querying. This endpoint converts in-memory buffers with the recently ingested logs into searchable data blocks.

It isn't recommended requesting the `/internal/force_flush` HTTP endpoint on a regular basis, since this increases CPU usage
and slows down data ingestion. It is expected that the `/internal/force_flush` is requested in automated tests, which need querying
the recently ingested data.

The `/internal/force_flush` endpoint can be protected from unauthorized access via `-forceFlushAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).

## High Availability

### High Availability (HA) Setup with VictoriaLogs Single-Node Instances

The setup consists of the following components:

- **Log Collector**: The log collector should support sending the same collected data to multiple destinations (aka replication).
It is recommended to use [vlagent](https://docs.victoriametrics.com/victorialogs/vlagent/). Otherp popular log collectors also provide this ability:
- [How to setup replication at FluentBit](https://docs.fluentbit.io/manual/concepts/data-pipeline/router)
- [How to setup replication at Logstash](https://www.elastic.co/guide/en/logstash/current/output-plugins.html)
- [How to setup replication at Fluentd](https://docs.fluentd.org/output/copy)
- [How to setup replication at Vector](https://vector.dev/docs/reference/configuration/sinks/)

- **VictoriaLogs Single-Node Instances**: send copies of the collected logs to multiple instances of VictoriaLogs in distinct availability zones to achieve HA.

- **[vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/#load-balancing)**: query logs via `vmauth` - it balances incoming queries among available VictoriaLogs instances,
  and automatically re-routes requests to healthy backends if some of the instances are temporarily unavailable.

![VictoriaLogs Single-Node Instance High-Availability schema](ha-victorialogs-single-node.webp)

Here are the working examples of HA configuration for VictoriaLogs using Docker Compose:

- [Fluent Bit + VictoriaLogs Single-Node + vmauth](https://github.com/VictoriaMetrics/VictoriaLogs/tree/master/deployment/docker/victorialogs/fluentbit/jsonline-ha)
- [Logstash + VictoriaLogs Single-Node + vmauth](https://github.com/VictoriaMetrics/VictoriaLogs/tree/master/deployment/docker/victorialogs/logstash/jsonline-ha)
- [Vector + VictoriaLogs Single-Node + vmauth](https://github.com/VictoriaMetrics/VictoriaLogs/tree/master/deployment/docker/victorialogs/vector/jsonline-ha)

## Backup and restore

VictoriaLogs stores data into independent per-day partitions. Every partition is stored in a separate directory - `<-storageDataPath>/partitions/YYYYMMDD`.

The following steps must be performed to make a backup of the given `YYYYMMDD` partition:

1. To create a snapshot for the given per-day partition via `/internal/partition/snapshot/create?name=YYYYMMDD` HTTP endpoint (see [partitions lifecycle](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) docs).
   This endpoint returns an absolute filesystem path to the created snapshot - `<path-to-snapshot>`.

1. To backup the created snapshot with [`rsync`](https://en.wikipedia.org/wiki/Rsync):

   ```sh
   rsync -avh --progress --delete <path-to-snapshot> <username>@<host>:<path-to-backup>/YYYYMMDD
   ```

   The `--delete` option is required in the command above in order to ensures that the backup contains the full copy of the original data without superfluous files.

1. To remove the snapshot with `rm -rf <path-to-snapshot>` command. It is important to remove unneeded snapshots in order to free up storage space.


The following steps must be performed for restoring the partition data from backup:

1. To stop VictoriaLogs instance or to detach the `YYYYMMDD` partition, which is going to be restored from backup,
   from the running VictoriaLogs via `/internal/partition/detach?name=YYYYMMDD` HTTP endpoint according to [these docs](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle).

1. To copy the partition from backup with `rsync`:

   ```sh
   rsync -avh --progress --delete <username>@<host>:<path-to-backup>/YYYYMMDD <-storageDataPath>/partitions/
   ```

   The `--delete` option is required in the command above in order to ensure that the partition contains the full copy of the backup without superfluous files.

1. To start VictoriaLogs instance or to attach the restored partition to the running VictoriaLogs instance via `/internal/partition/attach?name=YYYYMMDD` HTTP endpoint
   according to [these docs](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle).

It is also possible to use **the disk snapshot** feature provided by the operating system or cloud provider in order to perform a backup.

## Multitenancy

VictoriaLogs supports multitenancy. A tenant is identified by `(AccountID, ProjectID)` pair, where `AccountID` and `ProjectID` are arbitrary 32-bit unsigned integers.
The `AccountID` and `ProjectID` fields can be set during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/)
and [querying](https://docs.victoriametrics.com/victorialogs/querying/) via `AccountID` and `ProjectID` request headers.

If `AccountID` and/or `ProjectID` request headers aren't set, then the default `0` value is used.

VictoriaLogs has very low overhead for per-tenant management, so it is OK to have thousands of tenants in a single VictoriaLogs instance.

VictoriaLogs doesn't perform per-tenant authorization. Use [vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/) or similar tools for per-tenant authorization.
See [Security and Load balancing docs](https://docs.victoriametrics.com/victorialogs/security-and-lb/) for details.

### Multitenancy access control

Enforce access control for tenants by using [vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/). Access control can be configured for each tenant by setting up the following rules:

```yaml
users:
  - username: "foo"
    password: "bar"
    url_map:
      - src_paths:
        - "/select/.*"
        - "/insert/.*"
        headers:
          - "AccountID: 1"
          - "ProjectID: 0"
        url_prefix:
          - "http://localhost:9428/"

  - username: "baz"
    password: "bar"
    url_map:
      - src_paths: ["/select/.*"]
        headers:
          - "AccountID: 2"
          - "ProjectID: 0"
        url_prefix:
          - "http://localhost:9428/"
```

This configuration allows `foo` to use the `/select/.*` and `/insert/.*` endpoints with `AccountID: 1` and `ProjectID: 0`, while `baz` can only use the `/select/.*` endpoint with `AccountID: 2` and `ProjectID: 0`.

See also [Security and Load balancing docs](https://docs.victoriametrics.com/victorialogs/security-and-lb/).

## Security

It is expected that VictoriaLogs runs in a protected environment, which is unreachable from the Internet without proper authorization.
It is recommended providing access to VictoriaLogs [data ingestion APIs](https://docs.victoriametrics.com/victorialogs/data-ingestion/)
and [querying APIs](https://docs.victoriametrics.com/victorialogs/querying/#http-api) via [vmauth](https://docs.victoriametrics.com/victoriametrics/vmauth/)
or similar authorization proxies. See [Security and Load balancing docs](https://docs.victoriametrics.com/victorialogs/security-and-lb/) for details.

It is recommended protecting internal HTTP endpoints from unauthorized access:

- [`/internal/force_flush`](https://docs.victoriametrics.com/victorialogs/#forced-flush) - via `-forceFlushAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).
- [`/internal/force_merge`](https://docs.victoriametrics.com/victorialogs/#forced-merge) - via `-forceMergeAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).
- [`/internal/partition/*`](https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle) - via `-partitionManageAuthKey` [command-line flag](https://docs.victoriametrics.com/victorialogs/#list-of-command-line-flags).

### mTLS

[Enterprise version of VictoriaLogs](https://docs.victoriametrics.com/victoriametrics/enterprise/) supports verification of client TLS certificates
for TCP connections at the address specified via `-httpListenAddr` command-line flag (by default, this is `9428` TCP port).
This is known as [mTLS authentication](https://en.wikipedia.org/wiki/Mutual_authentication#mTLS).

Pass `-mtls` command-line flag to VictoriaLogs in order to enable mTLS authentication for incoming requests.

By default the system-wide [root CA certificates](https://en.wikipedia.org/wiki/Root_certificate) are used for verifying client TLS certificates.
The `-mtlsCAFile` command-line flag can be used for pointing to custom root CA certificates.

[Enterprise version of VictoriaLogs](https://docs.victoriametrics.com/victoriametrics/enterprise/) can be downloaded and evaluated for free
from [the releases page](https://github.com/VictoriaMetrics/VictoriaLogs/releases/latest). See [how to request a free trial license](https://victoriametrics.com/products/enterprise/trial/).

### Automatic issuing of TLS certificates

All the [VictoriaLogs Enterprise](https://docs.victoriametrics.com/victoriametrics/enterprise/) components support automatic issuing of TLS certificates
for public HTTPS server running at `-httpListenAddr` via [Let's Encrypt service](https://letsencrypt.org/).
The following command-line flags must be set in order to enable automatic issuing of TLS certificates:

- `-httpListenAddr` must be set for listening TCP port `443`. For example, `-httpListenAddr=:443`. This port must be accessible by the [Let's Encrypt service](https://letsencrypt.org/).
- `-tls` must be set in order to accept HTTPS requests at `-httpListenAddr`. Note that `-tlcCertFile` and `-tlsKeyFile` aren't needed when automatic TLS certificate issuing is enabled.
- `-tlsAutocertHosts` must be set to comma-separated list of hosts, which can be reached via `-httpListenAddr`. TLS certificates are automatically issued for these hosts.
- `-tlsAutocertEmail` must be set to contact email for the issued TLS certificates.
- `-tlsAutocertCacheDir` may be set to the directory path for persisting the issued TLS certificates between VictoriaMetrics restarts. If this flag isn't set,
  then TLS certificates are re-issued on every restart.

This functionality can be evaluated for free according to [these docs](https://docs.victoriametrics.com/victoriametrics/enterprise/).

See also [security recommendations](https://docs.victoriametrics.com/victorialogs/#security).

## Benchmarks

See the following benchmark results:

- [JSONBench: the comparison of VictoriaLogs with Elasticsearch, MongoDB, DuckDB and PostgreSQL](https://jsonbench.com/#eyJzeXN0ZW0iOnsiQ2xpY2tIb3VzZSAobHo0KSI6ZmFsc2UsIkNsaWNrSG91c2UgKHpzdGQpIjpmYWxzZSwiRHVja0RCIjp0cnVlLCJFbGFzdGljc2VhcmNoIChubyBzb3VyY2UsIGJlc3QgY29tcHJlc3Npb24pIjpmYWxzZSwiRWxhc3RpY3NlYXJjaCAobm8gc291cmNlLCBkZWZhdWx0KSI6ZmFsc2UsIkVsYXN0aWNzZWFyY2ggKGJlc3QgY29tcHJlc3Npb24pIjpmYWxzZSwiRWxhc3RpY3NlYXJjaCAoZGVmYXVsdCkiOnRydWUsIkVsYXN0aWNzZWFyY2giOmZhbHNlLCJNb25nb0RCIChzbmFwcHksIGNvdmVyZWQgaW5kZXgpIjpmYWxzZSwiTW9uZ29EQiAoenN0ZCwgY292ZXJlZCBpbmRleCkiOmZhbHNlLCJNb25nb0RCIChzbmFwcHkpIjpmYWxzZSwiTW9uZ29EQiAoenN0ZCkiOnRydWUsIlBvc3RncmVTUUwgKGx6NCkiOnRydWUsIlBvc3RncmVTUUwgKHBnbHopIjpmYWxzZSwiVmljdG9yaWFMb2dzIjp0cnVlLCJFbGFzdGljc2VhcmNoIChubyBzb3VyY2UsIHpzdGQpIjp0cnVlLCJFbGFzdGljc2VhcmNoIChubyBzb3VyY2UsIGx6NCkiOnRydWUsIkVsYXN0aWNzZWFyY2ggKHpzdGQpIjp0cnVlLCJFbGFzdGljc2VhcmNoIChsejQpIjp0cnVlLCJQb3N0Z3JlU1FMIjp0cnVlfSwic2NhbGUiOjEwMDAwMDAwMDAsIm1ldHJpYyI6ImhvdCIsInF1ZXJpZXMiOlt0cnVlLHRydWUsdHJ1ZSx0cnVlLHRydWVdLCJyZXRhaW5fc3RydWN0dXJlIjp7InllcyI6dHJ1ZSwibm8iOnRydWV9fQ==). The benchmark can be reproduced by running `main.sh` file inside `victorialogs` directory of the [JSONBench repository](https://github.com/ClickHouse/JSONBench).
- [ClickBench: the comparison of VictoriaLogs with Elasticsearch, MongoDB, TimescaleDB, PostgreSQL, MySQL and SQLite](<https://benchmark.clickhouse.com/#system=+ltrc|ehed|noB|yL|gS(|gQ|Lt|m%E2%98%81|%20nu|coog&type=-&machine=-&cluster_size=-&opensource=-&tuned=-&metric=combined&queries=->). The benchmark can be reproduced by running `benchmark.sh` file inside `victorialogs` directory of the [ClickBench repository](https://github.com/ClickHouse/ClickBench/).

Here is a [benchmark suite](https://github.com/VictoriaMetrics/VictoriaLogs/tree/master/deployment/logs-benchmark) for comparing data ingestion performance
and resource usage between VictoriaLogs and Elasticsearch or Loki.

It is recommended [setting up VictoriaLogs](https://docs.victoriametrics.com/victorialogs/quickstart/) in production alongside the existing
log management systems and comparing resource usage + query performance between VictoriaLogs and your system such as Elasticsearch or Grafana Loki.

Please share benchmark results and ideas on how to improve benchmarks / VictoriaLogs
via [VictoriaMetrics community channels](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#community-and-contributions).

## Profiling

VictoriaLogs provides handlers for collecting the following [Go profiles](https://blog.golang.org/profiling-go-programs):

- Memory profile. It can be collected with the following command (replace `0.0.0.0` with hostname if needed):

```sh
curl http://0.0.0.0:9428/debug/pprof/heap > mem.pprof
```

- CPU profile. It can be collected with the following command (replace `0.0.0.0` with hostname if needed):

```sh
curl http://0.0.0.0:9428/debug/pprof/profile > cpu.pprof
```

The command for collecting CPU profile waits for 30 seconds before returning.

The collected profiles may be analyzed with [go tool pprof](https://github.com/google/pprof).
It is safe sharing the collected profiles from security point of view, since they do not contain sensitive information.

## Environment variables

All VictoriaLogs components support configuring command-line flags via environment variables.
You can define flags using environment variables, and you can also
reference environment variables as values, allowing you to reuse or dynamically inject configuration values at application startup.
See [these docs](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#environment-variables) for details.

## List of command-line flags

Pass `-help` to VictoriaLogs in order to see the list of supported command-line flags with their description:

```
  -blockcache.missesBeforeCaching int
        The number of cache misses before putting the block into cache. Higher values may reduce indexdb/dataBlocks cache size at the cost of higher CPU and disk read usage (default 2)
  -datadog.ignoreFields array
        Comma-separated list of fields to ignore for logs ingested via DataDog protocol. See https://docs.victoriametrics.com/victorialogs/data-ingestion/datadog-agent/#dropping-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -datadog.maxRequestSize size
        The maximum size in bytes of a single DataDog request
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 67108864)
  -datadog.streamFields array
        Comma-separated list of fields to use as log stream fields for logs ingested via DataDog protocol. See https://docs.victoriametrics.com/victorialogs/data-ingestion/datadog-agent/#stream-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -defaultMsgValue string
        Default value for _msg field if the ingested log entry doesn't contain it; see https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field (default "missing _msg field; see https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field")
  -defaultParallelReaders int
        Default number of parallel data readers to use for executing every query; higher number of readers may help increasing query performance on high-latency storage such as NFS or S3 at the cost of higher RAM usage; see https://docs.victoriametrics.com/victorialogs/logsql/#parallel_readers-query-option (default 32)
  -elasticsearch.version string
        Elasticsearch version to report to client (default "8.9.0")
  -enableTCP6
        Whether to enable IPv6 for listening and dialing. By default, only IPv4 TCP and UDP are used
  -envflag.enable
        Whether to enable reading flags from environment variables in addition to the command line. Command line flag values have priority over values from environment vars. Flags are read only from the command line if this flag isn't set. See https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#environment-variables for more details
  -envflag.prefix string
        Prefix for environment variables if -envflag.enable is set
  -eula
        Deprecated, please use -license or -licenseFile flags instead. By specifying this flag, you confirm that you have an enterprise license and accept the ESA https://victoriametrics.com/legal/esa/ . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
  -filestream.disableFadvise
        Whether to disable fadvise() syscall when reading large data files. The fadvise() syscall prevents from eviction of recently accessed data from OS page cache during background merges and backups. In some rare cases it is better to disable the syscall if it uses too much CPU
  -flagsAuthKey value
        Auth key for /flags endpoint. It must be passed via authKey query arg. It overrides -httpAuth.*
        Flag value can be read from the given file when using -flagsAuthKey=file:///abs/path/to/file or -flagsAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -flagsAuthKey=http://host/path or -flagsAuthKey=https://host/path
  -forceFlushAuthKey value
        authKey, which must be passed in query string to /internal/force_flush . It overrides -httpAuth.* . See https://docs.victoriametrics.com/victorialogs/#forced-flush
        Flag value can be read from the given file when using -forceFlushAuthKey=file:///abs/path/to/file or -forceFlushAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -forceFlushAuthKey=http://host/path or -forceFlushAuthKey=https://host/path
  -forceMergeAuthKey value
        authKey, which must be passed in query string to /internal/force_merge . It overrides -httpAuth.* . See https://docs.victoriametrics.com/victorialogs/#forced-merge
        Flag value can be read from the given file when using -forceMergeAuthKey=file:///abs/path/to/file or -forceMergeAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -forceMergeAuthKey=http://host/path or -forceMergeAuthKey=https://host/path
  -fs.disableMmap
        Whether to use pread() instead of mmap() for reading data files. By default, mmap() is used for 64-bit arches and pread() is used for 32-bit arches, since they cannot read data files bigger than 2^32 bytes in memory. mmap() is usually faster for reading small data chunks than pread()
  -futureRetention value
        Log entries with timestamps bigger than now+futureRetention are rejected during data ingestion; see https://docs.victoriametrics.com/victorialogs/#retention
        The following optional suffixes are supported: s (second), h (hour), d (day), w (week), y (year). If suffix isn't set, then the duration is counted in months (default 2d)
  -http.connTimeout duration
        Incoming connections to -httpListenAddr are closed after the configured timeout. This may help evenly spreading load among a cluster of services behind TCP-level load balancer. Zero value disables closing of incoming connections (default 2m0s)
  -http.disableCORS
        Disable CORS for all origins (*)
  -http.disableKeepAlive
        Whether to disable HTTP keep-alive for incoming connections at -httpListenAddr
  -http.disableResponseCompression
        Disable compression of HTTP responses to save CPU resources. By default, compression is enabled to save network bandwidth
  -http.header.csp string
        Value for 'Content-Security-Policy' header, recommended: "default-src 'self'"
  -http.header.frameOptions string
        Value for 'X-Frame-Options' header
  -http.header.hsts string
        Value for 'Strict-Transport-Security' header, recommended: 'max-age=31536000; includeSubDomains'
  -http.idleConnTimeout duration
        Timeout for incoming idle http connections (default 1m0s)
  -http.maxGracefulShutdownDuration duration
        The maximum duration for a graceful shutdown of the HTTP server. A highly loaded server may require increased value for a graceful shutdown (default 7s)
  -http.pathPrefix string
        An optional prefix to add to all the paths handled by http server. For example, if '-http.pathPrefix=/foo/bar' is set, then all the http requests will be handled on '/foo/bar/*' paths. This may be useful for proxied requests. See https://www.robustperception.io/using-external-urls-and-proxies-with-prometheus
  -http.shutdownDelay duration
        Optional delay before http server shutdown. During this delay, the server returns non-OK responses from /health page, so load balancers can route new requests to other servers
  -httpAuth.password value
        Password for HTTP server's Basic Auth. The authentication is disabled if -httpAuth.username is empty
        Flag value can be read from the given file when using -httpAuth.password=file:///abs/path/to/file or -httpAuth.password=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -httpAuth.password=http://host/path or -httpAuth.password=https://host/path
  -httpAuth.username string
        Username for HTTP server's Basic Auth. The authentication is disabled if empty. See also -httpAuth.password
  -httpListenAddr array
        TCP address to listen for incoming http requests. See also -httpListenAddr.useProxyProtocol
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -httpListenAddr.useProxyProtocol array
        Whether to use proxy protocol for connections accepted at the given -httpListenAddr . See https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt . With enabled proxy protocol http server cannot serve regular /metrics endpoint. Use -pushmetrics.url for metrics pushing
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -inmemoryDataFlushInterval duration
        The interval for guaranteed saving of in-memory data to disk. The saved data survives unclean shutdowns such as OOM crash, hardware reset, SIGKILL, etc. Bigger intervals may help increase the lifetime of flash storage with limited write cycles (e.g. Raspberry PI). Smaller intervals increase disk IO load. Minimum supported value is 1s (default 5s)
  -insert.concurrency int
        The average number of concurrent data ingestion requests, which can be sent to every -storageNode (default 2)
  -insert.disable
        Whether to disable /insert/* HTTP endpoints
  -insert.disableCompression
        Whether to disable compression when sending the ingested data to -storageNode nodes. Disabled compression reduces CPU usage at the cost of higher network usage
  -insert.maxFieldsPerLine int
        The maximum number of log fields per line, which can be read by /insert/* handlers; see https://docs.victoriametrics.com/victorialogs/faq/#how-many-fields-a-single-log-entry-may-contain (default 1000)
  -insert.maxLineSizeBytes size
        The maximum size of a single line that can be read by /insert/* handlers. Regardless of this flag, entries above the 2 MB limit are ignored, see https://docs.victoriametrics.com/victorialogs/faq/#what-length-a-log-record-is-expected-to-have
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 262144)
  -insert.maxQueueDuration duration
        The maximum duration to wait in the queue when -maxConcurrentInserts concurrent insert requests are executed (default 1m0s)
  -internStringCacheExpireDuration duration
        The expiry duration for caches for interned strings. See https://en.wikipedia.org/wiki/String_interning . See also -internStringMaxLen and -internStringDisableCache (default 6m0s)
  -internStringDisableCache
        Whether to disable caches for interned strings. This may reduce memory usage at the cost of higher CPU usage. See https://en.wikipedia.org/wiki/String_interning . See also -internStringCacheExpireDuration and -internStringMaxLen
  -internStringMaxLen int
        The maximum length for strings to intern. A lower limit may save memory at the cost of higher CPU usage. See https://en.wikipedia.org/wiki/String_interning . See also -internStringDisableCache and -internStringCacheExpireDuration (default 500)
  -internalinsert.disable
        Whether to disable /internal/insert HTTP endpoint. See https://docs.victoriametrics.com/victorialogs/cluster/#security
  -internalinsert.maxRequestSize size
        The maximum size in bytes of a single request, which can be accepted at /internal/insert HTTP endpoint
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 67108864)
  -internalselect.disable
        Whether to disable /internal/select/* HTTP endpoints
  -internalselect.maxConcurrentRequests int
        The limit on the number of concurrent requests to /internal/select/* endpoints; other requests are put into the wait queue; see https://docs.victoriametrics.com/victorialogs/cluster/ (default 100)
  -journald.ignoreFields array
        Comma-separated list of fields to ignore for logs ingested over journald protocol. See https://docs.victoriametrics.com/victorialogs/data-ingestion/journald/#dropping-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -journald.includeEntryMetadata
        Include Journald fields with double underscore prefixes
  -journald.streamFields array
        Comma-separated list of fields to use as log stream fields for logs ingested over journald protocol. See https://docs.victoriametrics.com/victorialogs/data-ingestion/journald/#stream-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -journald.tenantID string
        TenantID for logs ingested via the Journald endpoint. See https://docs.victoriametrics.com/victorialogs/data-ingestion/journald/#multitenancy (default "0:0")
  -journald.timeField string
        Field to use as a log timestamp for logs ingested via journald protocol. See https://docs.victoriametrics.com/victorialogs/data-ingestion/journald/#time-field (default "__REALTIME_TIMESTAMP")
  -license string
        License key for VictoriaMetrics Enterprise. See https://victoriametrics.com/products/enterprise/ . Trial Enterprise license can be obtained from https://victoriametrics.com/products/enterprise/trial/ . This flag is available only in Enterprise binaries. The license key can be also passed via file specified by -licenseFile command-line flag
  -license.forceOffline
        Whether to enable offline verification for VictoriaMetrics Enterprise license key, which has been passed either via -license or via -licenseFile command-line flag. The issued license key must support offline verification feature. Contact info@victoriametrics.com if you need offline license verification. This flag is available only in Enterprise binaries
  -licenseFile string
        Path to file with license key for VictoriaMetrics Enterprise. See https://victoriametrics.com/products/enterprise/ . Trial Enterprise license can be obtained from https://victoriametrics.com/products/enterprise/trial/ . This flag is available only in Enterprise binaries. The license key can be also passed inline via -license command-line flag
  -licenseFile.reloadInterval duration
        Interval for reloading the license file specified via -licenseFile. See https://victoriametrics.com/products/enterprise/ . This flag is available only in Enterprise binaries (default 1h0m0s)
  -logIngestedRows
        Whether to log all the ingested log entries; this can be useful for debugging of data ingestion; see https://docs.victoriametrics.com/victorialogs/data-ingestion/ ; see also -logNewStreams
  -logNewStreams
        Whether to log creation of new streams; this can be useful for debugging of high cardinality issues with log streams; see https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields ; see also -logIngestedRows
  -loggerDisableTimestamps
        Whether to disable writing timestamps in logs
  -loggerErrorsPerSecondLimit int
        Per-second limit on the number of ERROR messages. If more than the given number of errors are emitted per second, the remaining errors are suppressed. Zero values disable the rate limit
  -loggerFormat string
        Format for logs. Possible values: default, json (default "default")
  -loggerJSONFields string
        Allows renaming fields in JSON formatted logs. Example: "ts:timestamp,msg:message" renames "ts" to "timestamp" and "msg" to "message". Supported fields: ts, level, caller, msg
  -loggerLevel string
        Minimum level of errors to log. Possible values: INFO, WARN, ERROR, FATAL, PANIC (default "INFO")
  -loggerMaxArgLen int
        The maximum length of a single logged argument. Longer arguments are replaced with 'arg_start..arg_end', where 'arg_start' and 'arg_end' is prefix and suffix of the arg with the length not exceeding -loggerMaxArgLen / 2 (default 5000)
  -loggerOutput string
        Output for the logs. Supported values: stderr, stdout (default "stderr")
  -loggerTimezone string
        Timezone to use for timestamps in logs. Timezone must be a valid IANA Time Zone. For example: America/New_York, Europe/Berlin, Etc/GMT+3 or Local (default "UTC")
  -loggerWarnsPerSecondLimit int
        Per-second limit on the number of WARN messages. If more than the given number of warns are emitted per second, then the remaining warns are suppressed. Zero values disable the rate limit
  -loki.disableMessageParsing
        Whether to disable automatic parsing of JSON-encoded log fields inside Loki log message into distinct log fields
  -loki.maxRequestSize size
        The maximum size in bytes of a single Loki request
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 67108864)
  -maxConcurrentInserts int
        The maximum number of concurrent insert requests. Set higher value when clients send data over slow networks. Default value depends on the number of available CPU cores. It should work fine in most cases since it minimizes resource usage. See also -insert.maxQueueDuration (default 32)
  -memory.allowedBytes size
        Allowed size of system memory VictoriaMetrics caches may occupy. This option overrides -memory.allowedPercent if set to a non-zero value. Too low a value may increase the cache miss rate usually resulting in higher CPU and disk IO usage. Too high a value may evict too much data from the OS page cache resulting in higher disk IO usage
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 0)
  -memory.allowedPercent float
        Allowed percent of system memory VictoriaMetrics caches may occupy. See also -memory.allowedBytes. Too low a value may increase cache miss rate usually resulting in higher CPU and disk IO usage. Too high a value may evict too much data from the OS page cache which will result in higher disk IO usage (default 60)
  -metrics.exposeMetadata
        Whether to expose TYPE and HELP metadata at the /metrics page, which is exposed at -httpListenAddr . The metadata may be needed when the /metrics page is consumed by systems, which require this information. For example, Managed Prometheus in Google Cloud - https://cloud.google.com/stackdriver/docs/managed-prometheus/troubleshooting#missing-metric-type
  -metricsAuthKey value
        Auth key for /metrics endpoint. It must be passed via authKey query arg. It overrides -httpAuth.*
        Flag value can be read from the given file when using -metricsAuthKey=file:///abs/path/to/file or -metricsAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -metricsAuthKey=http://host/path or -metricsAuthKey=https://host/path
  -mtls array
        Whether to require valid client certificate for https requests to the corresponding -httpListenAddr . This flag works only if -tls flag is set. See also -mtlsCAFile . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -mtlsCAFile array
        Optional path to TLS Root CA for verifying client certificates at the corresponding -httpListenAddr when -mtls is enabled. By default the host system TLS Root CA is used for client certificate verification. This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -opentelemetry.maxRequestSize size
        The maximum size in bytes of a single OpenTelemetry request
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 67108864)
  -partitionManageAuthKey value
        authKey, which must be passed in query string to /internal/partition/* . It overrides -httpAuth.* . See https://docs.victoriametrics.com/victorialogs/#partitions-lifecycle
        Flag value can be read from the given file when using -partitionManageAuthKey=file:///abs/path/to/file or -partitionManageAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -partitionManageAuthKey=http://host/path or -partitionManageAuthKey=https://host/path
  -pprofAuthKey value
        Auth key for /debug/pprof/* endpoints. It must be passed via authKey query arg. It overrides -httpAuth.*
        Flag value can be read from the given file when using -pprofAuthKey=file:///abs/path/to/file or -pprofAuthKey=file://./relative/path/to/file.
        Flag value can be read from the given http/https url when using -pprofAuthKey=http://host/path or -pprofAuthKey=https://host/path
  -pushmetrics.disableCompression
        Whether to disable request body compression when pushing metrics to every -pushmetrics.url
  -pushmetrics.extraLabel array
        Optional labels to add to metrics pushed to every -pushmetrics.url . For example, -pushmetrics.extraLabel='instance="foo"' adds instance="foo" label to all the metrics pushed to every -pushmetrics.url
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -pushmetrics.header array
        Optional HTTP request header to send to every -pushmetrics.url . For example, -pushmetrics.header='Authorization: Basic foobar' adds 'Authorization: Basic foobar' header to every request to every -pushmetrics.url
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -pushmetrics.interval duration
        Interval for pushing metrics to every -pushmetrics.url (default 10s)
  -pushmetrics.url array
        Optional URL to push metrics exposed at /metrics page. See https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/#push-metrics . By default, metrics exposed at /metrics page aren't pushed to any remote storage
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -retention.maxDiskSpaceUsageBytes size
        The maximum disk space usage at -storageDataPath before older per-day partitions are automatically dropped; see https://docs.victoriametrics.com/victorialogs/#retention-by-disk-space-usage ; see also -retentionPeriod
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 0)
  -retention.maxDiskUsagePercent int
        The maximum allowed disk usage percentage (1-100) for the filesystem that contains -storageDataPath before older per-day partitions are automatically dropped; mutually exclusive with -retention.maxDiskSpaceUsageBytes; see https://docs.victoriametrics.com/victorialogs/#retention-by-disk-space-usage-percent
  -retentionPeriod value
        Log entries with timestamps older than now-retentionPeriod are automatically deleted; log entries with timestamps outside the retention are also rejected during data ingestion; the minimum supported retention is 1d (one day); see https://docs.victoriametrics.com/victorialogs/#retention ; see also -retention.maxDiskSpaceUsageBytes and -retention.maxDiskUsagePercent
        The following optional suffixes are supported: s (second), h (hour), d (day), w (week), y (year). If suffix isn't set, then the duration is counted in months (default 7d)
  -search.allowPartialResponse
        Whether to allow returning partial responses when some of vlstorage nodes from the -storageNode list are unavaialbe for querying. This flag works only for cluster setup of VictoriaLogs. See https://docs.victoriametrics.com/victorialogs/querying/#partial-responses
  -search.maxConcurrentRequests int
        The maximum number of concurrent search requests. It shouldn't be high, since a single request can saturate all the CPU cores, while many concurrently executed requests may require high amounts of memory. See also -search.maxQueueDuration (default 16)
  -search.maxQueryDuration duration
        The maximum duration for query execution. It can be overridden to a smaller value on a per-query basis via 'timeout' query arg (default 30s)
  -search.maxQueryTimeRange duration
        The maximum time range, which can be set in the query sent to querying APIs. Queries with bigger time ranges are rejected. See https://docs.victoriametrics.com/victorialogs/querying/#resource-usage-limits
  -search.maxQueueDuration duration
        The maximum time the search request waits for execution when -search.maxConcurrentRequests limit is reached; see also -search.maxQueryDuration (default 10s)
  -select.disable
        Whether to disable /select/* HTTP endpoints
  -select.disableCompression
        Whether to disable compression for select query responses received from -storageNode nodes. Disabled compression reduces CPU usage at the cost of higher network usage
  -storage.minFreeDiskSpaceBytes size
        The minimum free disk space at -storageDataPath after which the storage stops accepting new data
        Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 10000000)
  -storageDataPath string
        Path to directory where to store VictoriaLogs data; see https://docs.victoriametrics.com/victorialogs/#storage (default "victoria-logs-data")
  -storageNode array
        Comma-separated list of TCP addresses for storage nodes to route the ingested logs to and to send select queries to. If the list is empty, then the ingested logs are stored and queried locally from -storageDataPath
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.bearerToken array
        Optional bearer auth token to use for the corresponding -storageNode
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.bearerTokenFile array
        Optional path to bearer token file to use for the corresponding -storageNode. The token is re-read from the file every second
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.password array
        Optional basic auth password to use for the corresponding -storageNode
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.passwordFile array
        Optional path to basic auth password to use for the corresponding -storageNode. The file is re-read every second
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.tls array
        Whether to use TLS (HTTPS) protocol for communicating with the corresponding -storageNode. By default communication is performed via HTTP
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -storageNode.tlsCAFile array
        Optional path to TLS CA file to use for verifying connections to the corresponding -storageNode. By default, system CA is used
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.tlsCertFile array
        Optional path to client-side TLS certificate file to use when connecting to the corresponding -storageNode
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.tlsInsecureSkipVerify array
        Whether to skip tls verification when connecting to the corresponding -storageNode
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -storageNode.tlsKeyFile array
        Optional path to client-side TLS certificate key to use when connecting to the corresponding -storageNode
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.tlsServerName array
        Optional TLS server name to use for connections to the corresponding -storageNode. By default, the server name from -storageNode is used
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.username array
        Optional basic auth username to use for the corresponding -storageNode
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -storageNode.usernameFile array
        Optional path to basic auth username to use for the corresponding -storageNode. The file is re-read every second
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.compressMethod.tcp array
        Compression method for syslog messages received at the corresponding -syslog.listenAddr.tcp. Supported values: none, gzip, deflate. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#compression
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.compressMethod.udp array
        Compression method for syslog messages received at the corresponding -syslog.listenAddr.udp. Supported values: none, gzip, deflate. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#compression
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.compressMethod.unix array
        Compression method for syslog messages received at the corresponding -syslog.listenAddr.unix. Supported values: none, gzip, deflate. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#compression
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.decolorizeFields.tcp array
        Fields to remove ANSI color codes across logs ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#decolorizing-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.decolorizeFields.udp array
        Fields to remove ANSI color codes across logs ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#decolorizing-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.decolorizeFields.unix array
        Fields to remove ANSI color codes across logs ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#decolorizing-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.extraFields.tcp array
        Fields to add to logs ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#adding-extra-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.extraFields.udp array
        Fields to add to logs ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#adding-extra-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.extraFields.unix array
        Fields to add to logs ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#adding-extra-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.ignoreFields.tcp array
        Fields to ignore at logs ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#dropping-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.ignoreFields.udp array
        Fields to ignore at logs ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#dropping-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.ignoreFields.unix array
        Fields to ignore at logs ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#dropping-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.listenAddr.tcp array
        Comma-separated list of TCP addresses to listen to for Syslog messages. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.listenAddr.udp array
        Comma-separated list of UDP addresses to listen to for Syslog messages. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.listenAddr.unix array
        Comma-separated list of Unix socket filepaths to listen to for Syslog messages. Filepaths may be prepended with 'unixgram:'  for listening for SOCK_DGRAM sockets. By default SOCK_STREAM sockets are used. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.mtls array
        Whether to require valid client certificate for https requests to the corresponding -syslog.listenAddr.tcp. This flag works only if -syslog.tls flag is set for the corresponding -syslog.listenAddr.tcp. See also -syslog.mtlsCAFile. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#mtls . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.mtlsCAFile array
        Optional path to TLS Root CA for verifying client certificates at the corresponding -syslog.listenAddr.tcp when the corresponding -syslog.mtls is enabled. By default the host system TLS Root CA is used for client certificate verification. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#mtls . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.streamFields.tcp array
        Fields to use as log stream labels for logs ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#stream-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.streamFields.udp array
        Fields to use as log stream labels for logs ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#stream-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.streamFields.unix array
        Fields to use as log stream labels for logs ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#stream-fields
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tenantID.tcp array
        TenantID for logs ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#multitenancy
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tenantID.udp array
        TenantID for logs ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#multitenancy
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tenantID.unix array
        TenantID for logs ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#multitenancy
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.timezone string
        Timezone to use when parsing timestamps in RFC3164 syslog messages. Timezone must be a valid IANA Time Zone. For example: America/New_York, Europe/Berlin, Etc/GMT+3 . See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/ (default "Local")
  -syslog.tls array
        Whether to enable TLS for receiving syslog messages at the corresponding -syslog.listenAddr.tcp. The corresponding -syslog.tlsCertFile and -syslog.tlsKeyFile must be set if -syslog.tls is set. See also -syslog.mtls. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#security
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.tlsCertFile array
        Path to file with TLS certificate for the corresponding -syslog.listenAddr.tcp if the corresponding -syslog.tls is set. Prefer ECDSA certs instead of RSA certs as RSA certs are slower. The provided certificate file is automatically re-read every second, so it can be dynamically updated. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#security
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tlsCipherSuites array
        Optional list of TLS cipher suites for -syslog.listenAddr.tcp if -syslog.tls is set. See the list of supported cipher suites at https://pkg.go.dev/crypto/tls#pkg-constants . See also https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#security
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tlsKeyFile array
        Path to file with TLS key for the corresponding -syslog.listenAddr.tcp if the corresponding -syslog.tls is set. The provided key file is automatically re-read every second, so it can be dynamically updated. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#security
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -syslog.tlsMinVersion string
        The minimum TLS version to use for -syslog.listenAddr.tcp if -syslog.tls is set. Supported values: TLS10, TLS11, TLS12, TLS13. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#security (default "TLS13")
  -syslog.useLocalTimestamp.tcp array
        Whether to use local timestamp instead of the original timestamp for the ingested syslog messages at the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#log-timestamps
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.useLocalTimestamp.udp array
        Whether to use local timestamp instead of the original timestamp for the ingested syslog messages at the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#log-timestamps
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.useLocalTimestamp.unix array
        Whether to use local timestamp instead of the original timestamp for the ingested syslog messages at the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#log-timestamps
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.useRemoteIP.tcp array
        Whether to add remote ip address as 'remote_ip' log field for syslog messages ingested via the corresponding -syslog.listenAddr.tcp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#capturing-remote-ip-address
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.useRemoteIP.udp array
        Whether to add remote ip address as 'remote_ip' log field for syslog messages ingested via the corresponding -syslog.listenAddr.udp. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#capturing-remote-ip-address
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -syslog.useRemoteIP.unix array
        Whether to add remote ip address as 'remote_ip' log field for syslog messages ingested via the corresponding -syslog.listenAddr.unix. See https://docs.victoriametrics.com/victorialogs/data-ingestion/syslog/#capturing-remote-ip-address
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -tls array
        Whether to enable TLS for incoming HTTP requests at the given -httpListenAddr (aka https). -tlsCertFile and -tlsKeyFile must be set if -tls is set. See also -mtls
        Supports array of values separated by comma or specified via multiple flags.
        Empty values are set to false.
  -tlsAutocertCacheDir string
        Directory to store TLS certificates issued via Let's Encrypt. Certificates are lost on restarts if this flag isn't set. This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
  -tlsAutocertEmail string
        Contact email for the issued Let's Encrypt TLS certificates. See also -tlsAutocertHosts and -tlsAutocertCacheDir . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
  -tlsAutocertHosts array
        Optional hostnames for automatic issuing of Let's Encrypt TLS certificates. These hostnames must be reachable at -httpListenAddr . The -httpListenAddr must listen tcp port 443 . The -tlsAutocertHosts overrides -tlsCertFile and -tlsKeyFile . See also -tlsAutocertEmail and -tlsAutocertCacheDir . This flag is available only in Enterprise binaries. See https://docs.victoriametrics.com/victoriametrics/enterprise/
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -tlsCertFile array
        Path to file with TLS certificate for the corresponding -httpListenAddr if -tls is set. Prefer ECDSA certs instead of RSA certs as RSA certs are slower. The provided certificate file is automatically re-read every second, so it can be dynamically updated. See also -tlsAutocertHosts
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -tlsCipherSuites array
        Optional list of TLS cipher suites for incoming requests over HTTPS if -tls is set. See the list of supported cipher suites at https://pkg.go.dev/crypto/tls#pkg-constants
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -tlsKeyFile array
        Path to file with TLS key for the corresponding -httpListenAddr if -tls is set. The provided key file is automatically re-read every second, so it can be dynamically updated. See also -tlsAutocertHosts
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -tlsMinVersion array
        Optional minimum TLS version to use for the corresponding -httpListenAddr if -tls is set. Supported values: TLS10, TLS11, TLS12, TLS13
        Supports an array of values separated by comma or specified via multiple flags.
        Value can contain comma inside single-quoted or double-quoted string, {}, [] and () braces.
  -version
        Show VictoriaMetrics version
```
