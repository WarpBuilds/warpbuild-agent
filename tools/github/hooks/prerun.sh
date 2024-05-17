#!/bin/bash

echo "Setting up WarpBuild runner for job: $GITHUB_JOB_ID"
echo "Job will be run on $RUNNER_NAME with image version: $(cat /.warp_image_version)"

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
  -O - --post-file=warpbuild_body.json \
  "$WARPBUILD_HOST_URL/api/v1/job" || exit 1

rm warpbuild_body.json