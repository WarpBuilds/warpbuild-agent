#!/bin/bash

echo 'Running cloudinit script...'

cd ~
mkdir runner
cd runner

export TARGETARCH=${TARGETARCH:-amd64}
export TARGETOS=${TARGETOS:-linux}
export RUNNER_VERSION=${RUNNER_VERSION:-2.311.0}
export RUNNER_CONTAINER_HOOKS_VERSION=${RUNNER_CONTAINER_HOOKS_VERSION:-0.3.2}
export RUNNER_ARCH=$TARGETARCH \
	&& if [ "$RUNNER_ARCH" = "amd64" ]; then export RUNNER_ARCH=x64 ; fi \
	&& curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-${TARGETOS}-${RUNNER_ARCH}-${RUNNER_VERSION}.tar.gz \
	&& tar xzf ./runner.tar.gz \
	&& rm runner.tar.gz

curl -f -L -o runner-container-hooks.zip https://github.com/actions/runner-container-hooks/releases/download/v${RUNNER_CONTAINER_HOOKS_VERSION}/actions-runner-hooks-k8s-${RUNNER_CONTAINER_HOOKS_VERSION}.zip \
	&& unzip ./runner-container-hooks.zip -d ./k8s \
	&& rm runner-container-hooks.zip

export DEBIAN_FRONTEND=noninteractive
export RUNNER_MANUALLY_TRAP_SIG=1
export ACTIONS_RUNNER_PRINT_LOG_TO_STDOUT=1

# gh auth login
# install go 1.21 using snap
# sudo snap install go --classic --channel=1.21/stable


git clone https://github.com/WarpBuilds/warpbuild-agent.git
cd warpbuild-agent
git checkout feat/agent-v0
make build-agentd
sudo cp bin/warpbuild-agentd /usr/local/bin/warpbuild-agentd

sudo cp tools/systemd/warpbuild-agentd.service /etc/systemd/system/warpbuild-agentd.service
sudo systemctl daemon-reload
sudo systemctl enable warpbuild-agentd
sudo systemctl start warpbuild-agentd

sudo cp tools/github/hooks/prerun.sh ~/runner/prerun.sh

# ? set agent id and polling secret in settings.json
# ? you might need to replace the host url as well
#
echo "Using agent id: $AGENT_ID"

cat <<EOF > /var/lib/warpbuild-agentd/settings.json
{
  "agent": {
    "id": "$(echo $AGENT_ID)",
    "polling_secret": "$(echo $POLLING_SECRET)",
    "host_url": "https://api.dev.warpbuild.dev/api/v1",
    "exit_file_location": "/var/log/warpbuild-agentd/exit.json"
  },
  "runner": {
    "provider": "github",
    "github": {
      "runner_dir": "/home/prashant/runner",
      "script": "./run.sh",
      "stdout_file": "/var/log/warpbuild-agentd/runner.github.stdout.log",  
      "stderr_file": "/var/log/warpbuild-agentd/runner.github.stderr.log"
    }
  }
}
EOF
