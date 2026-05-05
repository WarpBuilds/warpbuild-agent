#!/bin/bash
set -euo pipefail

echo "[tailscale-addon] Starting Tailscale setup for runner ${WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID} (macOS)"

export PATH="/opt/homebrew/bin:$PATH"

if command -v tailscale &>/dev/null; then
  echo "[tailscale-addon] Tailscale CLI found at $(which tailscale)"
else
  echo "[tailscale-addon] Tailscale CLI not found, installing via brew..."
  brew install tailscale
fi

tailscale version

echo "[tailscale-addon] Starting tailscaled daemon (ephemeral)..."
sudo tailscaled --state=mem: &

echo "[tailscale-addon] Performing OIDC login..."
login_timeout=30
login_elapsed=0
while [ $login_elapsed -lt $login_timeout ]; do
  if sudo tailscale up \
    --hostname="warpbuild-${WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID}" \
    --client-id="${WARPBUILD_ADDON_TS_CLIENT_ID}" \
    --id-token="${WARPBUILD_ADDON_TS_OIDC_TOKEN}" \
    ${WARPBUILD_ADDON_TS_ARGS:-} 2>&1; then
    echo "[tailscale-addon] Tailscale is connected"
    sudo tailscale status --peers=false || true
    echo "[tailscale-addon] Tailscale setup complete"
    exit 0
  fi
  echo "[tailscale-addon] Retrying in 2s... ($((login_elapsed + 1))/${login_timeout})"
  sleep 2
  login_elapsed=$((login_elapsed + 2))
done

echo "[tailscale-addon] ERROR: tailscale up failed after ${login_timeout}s"
exit 1
