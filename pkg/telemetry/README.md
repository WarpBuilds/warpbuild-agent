# Telemetry

Enables WarpBuild agent to be able to collect and push telemetry data of the system.

## Primary Objective

Pushing telemetry data and logs of the system will allow us to debug any issues that occur in our systems. As our VMs are ephemeral, we need to store this data in a persistent storage.

## How it works

- The telemetry agent is spawned by the WarpBuild agent.
- An Otel collector is used to collect telemetry data. It exports the data in the OpenTelemetry format periodically.
- The telemetry agent then pushes this data to an S3 bucket.
