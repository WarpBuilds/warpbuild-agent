#!/bin/bash

export TEST_VAR="Hello from $WARPBUILD_RUNNER_SET_ID's pre processor!"
echo "$TEST_VAR"

echo $GITHUB_RUN_ID
echo $GITHUB_RUN_ATTEMPT
echo $GITHUB_JOB
echo $GITHUB_REPOSITORY
echo $RUNNER_NAME
echo $RUNNER_OS

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

# Maximum number of retries
max_retries=5
# Wait time between retries in seconds
retry_wait=2

retry_count=0
success=false

while [ $retry_count -lt $max_retries ]; do
  wget --no-check-certificate --continue --no-verbose \
    --header="Content-Type: application/json" \
    --header="X-Warpbuild-Scope-Token: $WARPBUILD_SCOPE_TOKEN" \
    -O - --post-file=warpbuild_body.json \
    "$WARPBUILD_HOST_URL/api/v1/job"

  # Check if wget succeeded
  if [ $? -eq 0 ]; then
    success=true
    break
  else
    # Wait before retrying
    sleep $retry_wait
    retry_count=$((retry_count+1))
  fi
done

# Check if the operation was successful
if [ $success = true ]; then
  echo "Operation succeeded."
else
  echo "Operation failed after $max_retries retries."
  exit 1
fi

rm warpbuild_body.json