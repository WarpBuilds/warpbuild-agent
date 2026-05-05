#!/bin/bash
set -euo pipefail

TAILSCALE_VERSION="1.96.4"

echo "[tailscale-addon] Starting Tailscale setup for runner ${WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID}"

if ! command -v tailscale &>/dev/null; then
  echo "[tailscale-addon] Tailscale not found, installing..."
  if command -v apt-get &>/dev/null; then
    curl -fsSL https://pkgs.tailscale.com/stable/ubuntu/$(lsb_release -cs).noarmor.gpg | sudo tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null
    curl -fsSL https://pkgs.tailscale.com/stable/ubuntu/$(lsb_release -cs).tailscale-keyring.list | sudo tee /etc/apt/sources.list.d/tailscale.list
    sudo apt-get update -qq
    sudo apt-get install -y tailscale=${TAILSCALE_VERSION}
  else
    curl -fsSL https://tailscale.com/install.sh | sh
  fi
fi

tailscale version

echo "[tailscale-addon] Stopping any existing tailscaled instance..."
sudo systemctl stop tailscaled.service 2>/dev/null || true
sudo systemctl disable tailscaled.service 2>/dev/null || true
sleep 1

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
