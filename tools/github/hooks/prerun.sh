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
echo "GITHUB_REF=$GITHUB_REF"
echo "GITHUB_REF_TYPE=$GITHUB_REF_TYPE"
echo "RUNNER_NAME=$RUNNER_NAME"
echo "RUNNER_OS=$RUNNER_OS"
echo "WARPBUILD_RUNNER_SET_ID=$WARPBUILD_RUNNER_SET_ID"

if [ -n "$WARPBUILD_SNAPSHOT_KEY" ]; then
    echo "WARPBUILD_SNAPSHOT_KEY=$WARPBUILD_SNAPSHOT_KEY"
fi

if [ -z "$WARPBUILD_RUNNER_VERIFICATION_TOKEN" ]; then
    echo "WARPBUILD_RUNNER_VERIFICATION_TOKEN is not set."
    exit 1
fi

if [ -z "$RUNNER_NAME" ]; then
    echo "RUNNER_NAME is not set."
    exit 1
fi

cat <<EOF > warpbuild_body.json
{
  "vcs_workflow_run_id": "$GITHUB_RUN_ID",
  "vcs_workflow_run_attempt": "$GITHUB_RUN_ATTEMPT",
  "repo_entity": "$GITHUB_REPOSITORY",
  "repo_base_ref": "$GITHUB_BASE_REF",
  "repo_head_ref": "$GITHUB_HEAD_REF",
  "repo_ref": "$GITHUB_REF",
  "repo_ref_type": "$GITHUB_REF_TYPE"
}
EOF


echo -e "\n Invoking WarpBuild pre_hook..."

max_parent_retries=10
retry_delay_seconds=2 # Define delay between parent retries (in seconds)

while [[ max_parent_retries -gt 0 ]]; do
  # Use wget with retries, retry interval, no certificate check, and exit on failure
  wget --tries=5 --waitretry=2 --retry-connrefused \
    --retry-on-host-error --retry-on-http-error=502 \
    --retry-on-http-error=504 --retry-on-http-error=401 \
    --content-on-error \
    --no-check-certificate --continue \
    --header="Content-Type: application/json" \
    --header="Authorization: Bearer $WARPBUILD_RUNNER_VERIFICATION_TOKEN" \
    -O warpbuild_response --post-file=warpbuild_body.json \
    "$WARPBUILD_HOST_URL/api/v1/runners_instance/$RUNNER_NAME/pre_hook" || exit_code=$? || true

  if [ -n "$exit_code" ]; then
      echo "Failed to send request to warpbuild. Logging response..."
      cat warpbuild_body.json
      cat warpbuild_response
      max_parent_retries=$(expr $max_parent_retries - 1)
      echo "Retries left: $max_parent_retries"
      unset exit_code
      rm warpbuild_response
      sleep $retry_delay_seconds
  else
      echo "Request completed successfully."
      break
  fi

done

if [[ max_parent_retries -eq 0 ]]; then
    echo "All retries exhausted. StartJob API call failed. Exiting..."
    exit 1
fi

rm warpbuild_body.json

# Execute addon setup scripts returned by the backend
TOOLS_DIR="$(cd "$(dirname "$0")/../.." && pwd)"

if [ -f warpbuild_response ]; then
  script_count=$(jq -r '.setup_scripts | length // empty' warpbuild_response 2>/dev/null || echo "0")
  [ -z "$script_count" ] && script_count=0

  if [ "$script_count" -gt 0 ]; then
    echo -e "\nExecuting $script_count addon setup script(s)..."
    for i in $(seq 0 $((script_count - 1))); do
      script_name=$(jq -r ".setup_scripts[$i].name // \"addon-$i\"" warpbuild_response 2>/dev/null || echo "addon-$i")
      script_path=$(jq -r ".setup_scripts[$i].script_path // empty" warpbuild_response 2>/dev/null || true)

      echo -e "\n[addon:$script_name] Starting..."

      addon_script="${TOOLS_DIR}/${script_path}"
      if [ -z "$script_path" ] || [ ! -f "$addon_script" ]; then
        echo "[addon:$script_name] FAILED: script not found at $addon_script"
        exit 1
      fi

      # Run addon in subshell so env vars don't leak to parent
      chmod +x "$addon_script"
      (
        env_keys=$(jq -r ".setup_scripts[$i].env // {} | keys[]" warpbuild_response 2>/dev/null || true)
        for key in $env_keys; do
          val=$(jq -r ".setup_scripts[$i].env[\"$key\"]" warpbuild_response 2>/dev/null || true)
          export "$key=$val"
        done
        bash "$addon_script"
      )
      addon_exit=$?

      if [ $addon_exit -ne 0 ]; then
        echo "[addon:$script_name] FAILED with exit code $addon_exit"
        exit 1
      fi
      echo "[addon:$script_name] Completed successfully."
    done
  fi
fi

rm -f warpbuild_response

echo -e "\nPrehook for WarpBuild runner instance '$RUNNER_NAME' completed successfully."
