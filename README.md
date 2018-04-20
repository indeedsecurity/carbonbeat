# Carbonbeat

Carbonbeat currently supports shipping notifications from the Carbon Black Defense notifications API.

## Getting Started with Carbonbeat

Like any other beat, customize `carbonbeat.full.yml` to your liking, rename to `carbonbeat.yml` and you're ready to go.
You'll need to provide your API credentials. CB Defense notifications api requires a `SIEM` type API key.

There is a multistage Dockerfile included. It does not include the config so you need to mount it when you run the container.
