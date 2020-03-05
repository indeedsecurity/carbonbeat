# Carbonbeat

Carbonbeat currently supports shipping notifications from the Carbon Black Defense notifications API.

## Getting Started with Carbonbeat
You'll need to provide your API credentials in `carbonbeat.yml`. CB Defense notifications api requires a `SIEM` type API key.
As of carbonbeat 2.0 you need to provide both a `SIEM` type key for CB Defense notifications and an `API` type key for audit logging.
Like any other beat, customize `carbonbeat.full.yml` to your liking, rename to `carbonbeat.yml` and you're ready to go.
You can customize the outputs per the [beats outputs documentation](https://www.elastic.co/guide/en/beats/filebeat/current/configuring-output.html).

There is a multistage Dockerfile included. It does not include the config so you need to mount it when you run the container.

## Output example

Carbonbeat ships events in JSON format to its outputs. Here is an example of an event indexed into Elasticsearch and displayed by Kibana:

![example notification event](docs/carbonbeat-event.png "example notification event")
