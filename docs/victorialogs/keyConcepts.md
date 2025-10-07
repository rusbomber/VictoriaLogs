---
weight: 2
title: Key Concepts
menu:
  docs:
    identifier: vl-key-concepts
    parent: victorialogs
    weight: 2
    title: Key Concepts
tags:
  - logs
aliases:
- /victorialogs/keyConcepts.html
- /victorialogs/keyConcepts/
---
## Data model

[VictoriaLogs](https://docs.victoriametrics.com/victorialogs/) works with both structured and unstructured logs.
Every log entry must contain at least the [log message field](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field). An arbitrary number of additional `key=value` fields can be added to the log entry.
A single log entry can be expressed as a single-level [JSON](https://www.json.org/json-en.html) object with string keys and string values.
For example:

```json
{
  "job": "my-app",
  "instance": "host123:4567",
  "level": "error",
  "client_ip": "1.2.3.4",
  "trace_id": "1234-56789-abcdef",
  "_msg": "failed to serve the client request"
}
```

Empty values are treated the same as missing values. For example, the following log entries are equivalent,
since they have only one identical non-empty field - [`_msg`](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field):

```json
{
  "_msg": "foo bar",
  "some_field": "",
  "another_field": ""
}
```

```json
{
  "_msg": "foo bar",
  "third_field": ""
}
```

```json
{
  "_msg": "foo bar"
}
```

VictoriaLogs automatically transforms multi-level JSON (aka nested JSON) into single-level JSON
during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/) according to the following rules:

- Nested dictionaries are flattened by concatenating dictionary keys with the `.` character. For example, the following multi-level JSON
  is transformed into the following single-level JSON:

  ```json
  {
    "host": {
      "name": "foobar",
      "os": {
        "version": "1.2.3"
      }
    }
  }
  ```

  ```json
  {
    "host.name": "foobar",
    "host.os.version": "1.2.3"
  }
  ```

- Arrays, numbers and boolean values are converted into strings. This simplifies [full-text search](https://docs.victoriametrics.com/victorialogs/logsql/) over such values.
  For example, the following JSON with an array, a number and a boolean value is converted into the following JSON with string values:

  ```json
  {
    "tags": ["foo", "bar"],
    "offset": 12345,
    "is_error": false
  }
  ```

  ```json
  {
    "tags": "[\"foo\", \"bar\"]",
    "offset": "12345",
    "is_error": "false"
  }
  ```

Both the field name and the field value may contain arbitrary characters. Such characters must be encoded
during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/)
according to [JSON string encoding](https://www.rfc-editor.org/rfc/rfc7159.html#section-7).
Unicode characters must be encoded with [UTF-8](https://en.wikipedia.org/wiki/UTF-8) encoding:

```json
{
  "field with whitespace": "value\nwith\nnewlines",
  "Поле": "价值"
}
```

VictoriaLogs automatically indexes all the fields in all the [ingested](https://docs.victoriametrics.com/victorialogs/data-ingestion/) logs.
This enables [full-text search](https://docs.victoriametrics.com/victorialogs/logsql/) across all the fields.

VictoriaLogs supports the following special fields in addition to arbitrary [other fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#other-fields):

- [`_msg` field](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field)
- [`_time` field](https://docs.victoriametrics.com/victorialogs/keyconcepts/#time-field)
- [`_stream` and `_stream_id` fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields)

### Message field

Every ingested [log entry](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) must contain at least a `_msg` field with the actual log message. For example, this is the minimal
log entry for VictoriaLogs:

```json
{
  "_msg": "some log message"
}
```

If the actual log message uses a field name other than `_msg`, then it can be specified via `_msg_field` HTTP query arg or via `VL-Msg-Field` HTTP header
during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/).
For example, if the log message is located in the `event.original` field, then specify `_msg_field=event.original` query arg.
See [these docs](https://docs.victoriametrics.com/victorialogs/data-ingestion/#http-parameters) for details.

If the `_msg` field remains empty after an attempt to get it from `_msg_field`, then VictoriaLogs automatically sets it to the value specified
via `-defaultMsgValue` command-line flag.

```mermaid
flowchart LR
  A["Incoming log JSON"] --> B{"_msg_field set?"}
  B -- yes --> C["Use first non-empty (by list order) as _msg"]
  B -- no  --> D{"_msg present?"}
  D -- yes --> E["Use _msg"]
  D -- no  --> F["Use default message"]
```

### Time field

The ingested [log entries](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) may contain `_time` field with the timestamp of the ingested log entry.
This field must be in one of the following formats:

- [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) or [RFC3339](https://www.rfc-editor.org/rfc/rfc3339).
  For example, `2023-06-20T15:32:10Z` or `2023-06-20 15:32:10.123456789+02:00`.
  If timezone information is missing (for example, `2023-06-20 15:32:10`),
  then the time is parsed in the local timezone of the host where VictoriaLogs runs.

- Unix timestamp in seconds, milliseconds, microseconds or nanoseconds. For example, `1686026893` (seconds), `1686026893735` (milliseconds),
  `1686026893735321` (microseconds) or `1686026893735321098` (nanoseconds).

For example, the following [log entry](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) contains a valid timestamp with millisecond precision in the `_time` field:

```json
{
  "_msg": "some log message",
  "_time": "2023-04-12T06:38:11.095Z"
}
```

If the timestamp uses a field name other than `_time`, then it is possible to specify the real timestamp
field via `_time_field` HTTP query arg or via `VL-Time-Field` HTTP header during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/).
For example, if the timestamp is located in the `event.created` field, then specify `_time_field=event.created` query arg.
See [these docs](https://docs.victoriametrics.com/victorialogs/data-ingestion/#http-parameters) for details.

If `_time` field is missing, or if it equals `0`, or if it equals `-`, then the data ingestion time is used as log entry timestamp.

The `_time` field is used by [time filter](https://docs.victoriametrics.com/victorialogs/logsql/#time-filter) for quickly narrowing down
the search to the selected time range.

```mermaid
flowchart LR
  A["Incoming log JSON"] --> G{"_time_field set?"}
  G -- yes --> H["Use first parseable (by list order) as _time"]
  G -- no  --> I{"_time present & parseable?"}
  I -- yes --> J["Use _time"]
  I -- no  --> K["Use ingestion time"]
```

### Stream fields

Some [structured logging](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) fields may uniquely identify the application instance, which generates logs.
This may be either a single field such as `instance="host123:456"` or a set of fields such as
`{datacenter="...", env="...", job="...", instance="..."}` or
`{kubernetes.namespace="...", kubernetes.node.name="...", kubernetes.pod.name="...", kubernetes.container.name="..."}`.

Log entries received from a single application instance form a **log stream** in VictoriaLogs.
VictoriaLogs optimizes storing and [querying](https://docs.victoriametrics.com/victorialogs/logsql/#stream-filter) of individual log streams.
This provides the following benefits:

- Reduced disk space usage, since a log stream from a single application instance is usually compressed better
  than a mixed log stream from multiple distinct applications.

- Increased query performance, since VictoriaLogs needs to scan smaller amounts of data
  when [searching by stream fields](https://docs.victoriametrics.com/victorialogs/logsql/#stream-filter).

Every ingested log entry is associated with a log stream. Every log stream consists of the following special fields:

- `_stream_id` - this is a unique identifier for the log stream. All the logs for the particular stream can be selected
  via [`_stream_id:...` filter](https://docs.victoriametrics.com/victorialogs/logsql/#_stream_id-filter).

- `_stream` - this field contains stream labels in the format similar to [labels in Prometheus metrics](https://docs.victoriametrics.com/victoriametrics/keyconcepts/#labels):

  ```
  {field1="value1", ..., fieldN="valueN"}
  ```

  For example, if `host` and `app` fields are associated with the stream, then the `_stream` field will have `{host="host-123",app="my-app"}` value
  for the log entry with `host="host-123"` and `app="my-app"` fields. The `_stream` field can be searched
  with [stream filters](https://docs.victoriametrics.com/victorialogs/logsql/#stream-filter).

By default, the value of the `_stream` field is `{}`, since VictoriaLogs cannot automatically determine
which fields uniquely identify each log stream. This may lead to not-so-optimal resource usage and query performance.
Therefore we recommend specifying stream-level fields via `_stream_fields` query arg
during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/).
```mermaid
flowchart LR
  A["Incoming log JSON"] --> L{"_stream_fields set?"}
  L -- Yes --> X{"All listed fields have values?"}
  X -- Yes --> M["Build _stream from all listed fields"]
  X -- No  --> P["Build _stream from available fields only<br>(still creates a new stream)"]
  L -- No  --> N["Use empty stream {}"]
```
For example, if logs from Kubernetes containers have the following fields:

```json
{
  "kubernetes.namespace": "some-namespace",
  "kubernetes.node.name": "some-node",
  "kubernetes.pod.name": "some-pod",
  "kubernetes.container.name": "some-container",
  "_msg": "some log message"
}
```

then specify `_stream_fields=kubernetes.namespace,kubernetes.node.name,kubernetes.pod.name,kubernetes.container.name`
query arg during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/) in order to properly store
per-container logs into distinct streams.

#### How to determine which fields must be associated with log streams?

[Log streams](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields) must contain [fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model), which uniquely identify the application instance, which generates logs.
For example, `container`, `instance` and `host` are good candidates for stream fields.

Additional fields may be added to log streams if they **remain constant during the application instance lifetime**.
For example, `namespace`, `node`, `pod` and `job` are good candidates for additional stream fields. Adding such fields to log streams
makes sense if you are going to use these fields during search and want to speed up search with [stream filters](https://docs.victoriametrics.com/victorialogs/logsql/#stream-filter).

There is **no need to add all the constant fields to log streams**, since this may increase resource usage during data ingestion and querying.

**Never add non-constant fields to streams if these fields may change with every log entry of the same stream**.
For example, `ip`, `user_id` and `trace_id` **must never be associated with log streams**, since this may lead to [high cardinality issues](https://docs.victoriametrics.com/victorialogs/keyconcepts/#high-cardinality).

#### High cardinality

Some fields in the [ingested logs](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) may contain a large number of unique values across log entries.
For example, fields with names such as `ip`, `user_id` or `trace_id` tend to contain a large number of unique values.
VictoriaLogs works perfectly with such fields unless they are associated with [log streams](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields).

**Never** associate high-cardinality fields with [log streams](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields), since this may lead to the following issues:

- Performance degradation during [data ingestion](https://docs.victoriametrics.com/victorialogs/data-ingestion/)
  and [querying](https://docs.victoriametrics.com/victorialogs/querying/)
- Increased memory usage
- Increased CPU usage
- Increased disk space usage
- Increased disk read/write I/O

VictoriaLogs exposes `vl_streams_created_total` [metric](https://docs.victoriametrics.com/victorialogs/metrics/#vl_streams_created_total),
which shows the number of created streams since the last VictoriaLogs restart. If this metric grows at a rapid rate
over a long period of time, then there is a high chance of high-cardinality issues mentioned above.
VictoriaLogs can log all the newly registered streams when `-logNewStreams` command-line flag is passed to it.
This can help narrow down and eliminate high-cardinality fields from [log streams](https://docs.victoriametrics.com/victorialogs/keyconcepts/#stream-fields).

### Other fields

Every ingested log entry may contain an arbitrary number of [fields](https://docs.victoriametrics.com/victorialogs/keyconcepts/#data-model) in addition to [`_msg`](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field) and [`_time`](https://docs.victoriametrics.com/victorialogs/keyconcepts/#time-field).
For example, `level`, `ip`, `user_id`, `trace_id`, etc. Such fields can be used for simplifying and optimizing [search queries](https://docs.victoriametrics.com/victorialogs/logsql/).
It is usually faster to search over a dedicated `trace_id` field instead of searching for the `trace_id` inside a long [log message](https://docs.victoriametrics.com/victorialogs/keyconcepts/#message-field).
E.g. the `trace_id:="XXXX-YYYY-ZZZZ"` query usually works faster than the `_msg:"trace_id=XXXX-YYYY-ZZZZ"` query.

See [LogsQL docs](https://docs.victoriametrics.com/victorialogs/logsql/) for more details.
