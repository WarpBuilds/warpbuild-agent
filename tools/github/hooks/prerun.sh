#!/bin/bash

# this only used for script debug (uncomment if needed)
# echo "Script ID: 3"
echo "Running prehook for WarpBuild runner instance '$RUNNER_NAME'..."
echo -e "\nLogging environment variables..."

echo "GITHUB_RUN_ID=$GITHUB_RUN_ID"
echo "GITHUB_RUN_ATTEMPT=$GITHUB_RUN_ATTEMPT"
echo "GITHUB_JOB=$GITHUB_JOB"
echo "GITHUB_REPOSITORY=$GITHUB_REPOSITORY"
echo "GITHUB_BASE_REF=$GITHUB_BASE_REF"
echo "GITHUB_HEAD_REF=$GITHUB_HEAD_REF"
echo "GITHUB_REF_TYPE=$GITHUB_REF_TYPE"
echo "GITHUB_REF=$GITHUB_REF"
echo "GITHUB_REF_TYPE=$GITHUB_REF_TYPE"
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
  "repo_entity": "$GITHUB_REPOSITORY",
  "repo_base_ref": "$GITHUB_BASE_REF",
  "repo_head_ref": "$GITHUB_HEAD_REF",
  "repo_ref_type": "$GITHUB_REF_TYPE",
  "repo_ref": "$GITHUB_REF",
  "repo_ref_type": "$GITHUB_REF_TYPE",
}
EOF

echo -e "\nMaking a request to WarpBuild..."

# Use wget with retries, retry interval, no certificate check, and exit on failure
wget --tries=5 --waitretry=2 --retry-connrefused \
  --retry-on-host-error --retry-on-http-error=502 \
  --content-on-error \
  --no-check-certificate --continue --no-verbose \
  --header="Content-Type: application/json" \
  --header="X-Warpbuild-Scope-Token: $WARPBUILD_SCOPE_TOKEN" \
  -O warpbuild_response.json --post-file=warpbuild_body.json \
  "$WARPBUILD_HOST_URL/api/v1/job" || exit_code=$? || true

if [ -n "$exit_code" ]; then
    echo "Failed to send request to warpbuild. Logging response. Exiting..."
    # check if jq is installed if so then pretty print the json response
    if ! command -v jq &> /dev/null; then
        cat warpbuild_response.json
    else
        jq . warpbuild_response.json
    fi
    exit 1
fi

rm warpbuild_body.json

echo -e "\nPrehook for WarpBuild runner instance '$RUNNER_NAME' completed succesfully."