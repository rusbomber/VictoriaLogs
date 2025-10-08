#!/bin/bash

year=2025
month=08
day=20
parallel_workers=24

for hour in $(seq 0 23);
do
	curl -s https://data.gharchive.org/$year-$month-$day-$hour.json.gz \
		| curl -T - -X POST -H 'Content-Encoding: gzip' 'http://localhost:9428/insert/jsonline?_time_field=created_at&_stream_fields=type' &

	[[ $(jobs -p -r | wc -l) -ge $parallel_workers ]] && wait -n
done

wait
