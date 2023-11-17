#!/bin/bash

go build -o bin/warpbuild-agentd cmd/agentd/main.go
cp bin/warpbuild-agentd /usr/local/bin/warpbuild-agentd
cp tools/systemd/warpbuild-agentd.service /etc/systemd/system/warpbuild-agentd.service
systemctl daemon-reload
systemctl enable warpbuild-agentd
systemctl start warpbuild-agentd

echo "Using agent id: $AGENT_ID"
cat <<EOF > /var/lib/warpbuild-agentd/settings.json
{
  "agent": {
    "id": "$AGENT_ID"
  },
  "runner": {
    "provider": "github",
    "github": {
      "runner_dir": "/runner",
      "script": "run.sh",
      "stdout_file": "/var/log/warpbuild-agentd/runner.github.stdout.log",  
      "stderr_file": "/var/log/warpbuild-agentd/runner.github.stderr.log"
    }
  }
}
EOF