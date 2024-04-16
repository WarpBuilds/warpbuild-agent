#!/bin/bash

echo "Running prehook for warpbuild runner instance '$WARPBUILD_RUNNER_SET_ID'..."
echo "Logging environment variables..."

echo "GITHUB_RUN_ID=$GITHUB_RUN_ID"
echo "GITHUB_RUN_ATTEMPT=$GITHUB_RUN_ATTEMPT"
echo "GITHUB_JOB=$GITHUB_JOB"
echo "GITHUB_REPOSITORY=$GITHUB_REPOSITORY"
echo "RUNNER_NAME=$RUNNER_NAME"
echo "RUNNER_OS=$RUNNER_OS"

if [ -z "$WARPBUILD_SCOPE_TOKEN" ]; then
    echo "WARPBUILD_SCOPE_TOKEN is not set."
    exit 1
fi

cat <<EOF > warpbuild_body.json
{
  "runner_id": "$WARPBUILD_RUNNER_SET_ID",
  "runner_name": "$RUNNER_NAME",
  "orchestrator_job_id": "$GITHUB_JOB_ID",
  "orchestrator_job_group_id": "$GITHUB_RUN_ID",
  "orchestrator_job_group_attempt": "$GITHUB_RUN_ATTEMPT",
  "repo_entity": "$GITHUB_REPOSITORY"
}
EOF

# Use wget with retries, retry interval, no certificate check, and exit on failure
wget --tries=5 --waitretry=2 --retry-connrefused \
  --retry-on-host-error --retry-on-http-error=502 \
  --no-check-certificate --continue --no-verbose \
  --header="Content-Type: application/json" \
  --header="X-Warpbuild-Scope-Token: $WARPBUILD_SCOPE_TOKEN" \
  -O warpbuild_response.json --post-file=warpbuild_body.json \
  "$WARPBUILD_HOST_URL/api/v1/job"

exit_code=$?

if [ $exit_code -ne 0 ]; then
    echo "Failed to send job request to warpbuild. Logging response..."
    cat warpbuild_response.json
    exit $exit_code
fi

rm warpbuild_body.json

echo "Prehook for warpbuild runner instance '$WARPBUILD_RUNNER_SET_ID' completed."