#!/bin/bash

export TEST_VAR="Hello from $WARPBUILD_RUNNER_SET_ID's pre processor!"
echo "$TEST_VAR"

echo $GITHUB_RUN_ID
echo $GITHUB_RUN_ATTEMPT
echo $GITHUB_JOB
echo $GITHUB_REPOSITORY
echo $WARPBUILD_RUNNER_NAME
echo $WARPBUILD_RUNNER_OS

if [ -z "$WARPBUILD_SCOPE_TOKEN" ]; then
    echo "WARPBUILD_SCOPE_TOKEN is not set."
    exit 1
fi

cat <<EOF > warpbuild_body.json
{
  "runner_id": "$WARPBUILD_RUNNER_SET_ID",
  "runner_name": "$WARPBUILD_RUNNER_NAME",
  "orchestrator_job_id": "$GITHUB_JOB_ID",
  "orchestrator_job_group_id": "$GITHUB_RUN_ID",
  "orchestrator_job_group_attempt": "$GITHUB_RUN_ATTEMPT",
  "repo_entity": "$GITHUB_REPOSITORY"
}
EOF

curl -X POST --silent --show-error --fail \
     -H "Content-Type: application/json" \
     -H "X-Warpbuild-Scope-Token: $WARPBUILD_SCOPE_TOKEN" \
     -d @warpbuild_body.json \
     $WARPBUILD_HOST_URL/api/v1/job

rm warpbuild_body.json