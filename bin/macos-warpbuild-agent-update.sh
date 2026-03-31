#!/bin/bash -e -o pipefail
################################################################################
##  File:  install-agent.sh
##  Desc:  Installs and configures warpbuild agent
################################################################################

# source ~/utils/utils.sh

echo "Creating warpbuild directories ..."
mkdir -p ~/.warpbuild/agent

BASE_DIR=$(echo ~/.warpbuild/agent)

# Source switcher: "github" or "r2"
DOWNLOAD_SOURCE="r2"

# When using R2, specify the branch name here
R2_BRANCH_NAME="feat/macos-init-commands"

if [ "$DOWNLOAD_SOURCE" = "r2" ]; then
  # R2 Configuration - branch builds are uploaded to warpbuild-agentd/{BRANCH_NAME}/
  # This matches the release-testing.yaml workflow which syncs dist to:
  # bucket: warpbuild-packages-dev (public URL: pub-b4f7dbd911ef411ca27c8befa94bb744.r2.dev)
  # path: warpbuild-agentd/${GITHUB_REF_NAME}/

  R2_DEV_PUBLIC_URL="https://pub-b4f7dbd911ef411ca27c8befa94bb744.r2.dev"
  VERSION=$R2_BRANCH_NAME
  DOWNLOAD_URL="$R2_DEV_PUBLIC_URL/warpbuild-agentd/$VERSION/warpbuild-agentd_Darwin_arm64.tar.gz"

  echo "Using R2 dev bucket with branch: $VERSION"
  echo "Download URL: $DOWNLOAD_URL"
else
  # GitHub releases
  VERSION=$(curl -sL https://api.github.com/repos/WarpBuilds/warpbuild-agent/releases/latest | jq -r '.tag_name')
  DOWNLOAD_URL="https://github.com/WarpBuilds/warpbuild-agent/releases/download/$VERSION/warpbuild-agentd_Darwin_arm64.tar.gz"
fi

USERNAME=$(whoami)

echo "Downloading warpbuild-agentd $VERSION from $DOWNLOAD_SOURCE..."
curl -fL -o $BASE_DIR/warpbuild-agentd.tar.gz "$DOWNLOAD_URL"
mkdir -p $BASE_DIR/warpbuild-agentd
tar -xzf $BASE_DIR/warpbuild-agentd.tar.gz -C $BASE_DIR/warpbuild-agentd
rm $BASE_DIR/warpbuild-agentd.tar.gz
sudo mkdir -p /usr/local/bin
sudo mv $BASE_DIR/warpbuild-agentd/warpbuild-agentd /usr/local/bin/warpbuild-agentd
sudo chmod +x /usr/local/bin/warpbuild-agentd

echo "Configuring agent ..."

echo "Adding launcher"
cat << EOF > $BASE_DIR/warpbuild-agentd-launcher.sh
#!/bin/bash

. $HOME/.bashrc

/usr/local/bin/warpbuild-agentd --settings=$BASE_DIR/settings.json --with-sysinit
EOF

echo "Setting launcher permissions ..."
chmod +x $BASE_DIR/warpbuild-agentd-launcher.sh

echo "Adding settings ..."
mv $BASE_DIR/warpbuild-agentd-launcher.sh /usr/local/bin/warpbuild-agentd-launcher.sh

echo "Adding launchd config ..."

cat << EOF > $BASE_DIR/com.warpbuild.warpbuild-agentd-launcher.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>warpbuild-agentd-launcher</string>
    <key>Nice</key>
    <integer>-10</integer>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/warpbuild-agentd-launcher.sh</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
      <key>LANG</key>
      <string>en_US.UTF-8</string>
    </dict>
    <key>UserName</key>
    <string>$USERNAME</string>
    <key>RunAtLoad</key>
    <true/>
    <key>ThrottleInterval</key>
    <integer>15</integer>
    <key>KeepAlive</key>
    <dict>
      <key>SuccessfulExit</key>
      <false/>
    </dict>
    <key>StandardErrorPath</key>
    <string>$BASE_DIR/log/warpbuild-agentd-launcher.error.log</string>
    <key>StandardOutPath</key>
    <string>$BASE_DIR/log/warpbuild-agentd-launcher.out.log</string>
  </dict>
</plist>
EOF

echo "Adding telemetry launcher"
cat << EOF > $BASE_DIR/warpbuild-telemetryd-launcher.sh
#!/bin/bash

. $HOME/.bashrc

/usr/local/bin/warpbuild-agentd --settings=$BASE_DIR/settings.json --launch-telemetry
EOF

echo "Setting launcher permissions ..."
chmod +x $BASE_DIR/warpbuild-telemetryd-launcher.sh

echo "Adding settings ..."
mv $BASE_DIR/warpbuild-telemetryd-launcher.sh /usr/local/bin/warpbuild-telemetryd-launcher.sh

cat << EOF > $BASE_DIR/com.warpbuild.warpbuild-telemetryd-launcher.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>warpbuild-telemetryd-launcher</string>
    <key>Nice</key>
    <integer>-10</integer>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/warpbuild-telemetryd-launcher.sh</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
      <key>LANG</key>
      <string>en_US.UTF-8</string>
    </dict>
    <key>UserName</key>
    <string>$USERNAME</string>
    <key>RunAtLoad</key>
    <true/>
    <key>ThrottleInterval</key>
    <integer>15</integer>
    <key>KeepAlive</key>
    <dict>
      <key>SuccessfulExit</key>
      <false/>
    </dict>
    <key>StandardErrorPath</key>
    <string>$BASE_DIR/log/warpbuild-telemetryd-launcher.error.log</string>
    <key>StandardOutPath</key>
    <string>$BASE_DIR/log/warpbuild-telemetryd-launcher.out.log</string>
  </dict>
</plist>
EOF

echo "Adding proxy launcher"
cat << EOF > $BASE_DIR/warpbuild-proxyd-launcher.sh
#!/bin/bash

. $HOME/.bashrc

/usr/local/bin/warpbuild-agentd --settings=$BASE_DIR/settings.json --launch-proxy-server
EOF

echo "Setting launcher permissions ..."
chmod +x $BASE_DIR/warpbuild-proxyd-launcher.sh

echo "Adding settings ..."
mv $BASE_DIR/warpbuild-proxyd-launcher.sh /usr/local/bin/warpbuild-proxyd-launcher.sh

cat << EOF > $BASE_DIR/com.warpbuild.warpbuild-proxyd-launcher.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>warpbuild-proxyd-launcher</string>
    <key>Nice</key>
    <integer>-10</integer>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/warpbuild-proxyd-launcher.sh</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
      <key>LANG</key>
      <string>en_US.UTF-8</string>
    </dict>
    <key>UserName</key>
    <string>$USERNAME</string>
    <key>RunAtLoad</key>
    <true/>
    <key>ThrottleInterval</key>
    <integer>15</integer>
    <key>KeepAlive</key>
    <dict>
      <key>SuccessfulExit</key>
      <false/>
    </dict>
    <key>StandardErrorPath</key>
    <string>$BASE_DIR/log/warpbuild-proxyd-launcher.error.log</string>
    <key>StandardOutPath</key>
    <string>$BASE_DIR/log/warpbuild-proxyd-launcher.out.log</string>
  </dict>
</plist>
EOF

echo "Adding agent to launchd ..."
LAUNCH_AGENTS_DIR=~/Library/LaunchAgents
mkdir -p $LAUNCH_AGENTS_DIR
mv $BASE_DIR/com.warpbuild.warpbuild-agentd-launcher.plist $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-agentd-launcher.plist
mv $BASE_DIR/com.warpbuild.warpbuild-telemetryd-launcher.plist $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-telemetryd-launcher.plist
mv $BASE_DIR/com.warpbuild.warpbuild-proxyd-launcher.plist $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-proxyd-launcher.plist
echo "Loading agent ..."
launchctl unload -w $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-agentd-launcher.plist || true
sleep 5
launchctl load -w $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-agentd-launcher.plist || true
echo "Loading telemetry agent ..."
launchctl load -w $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-telemetryd-launcher.plist || true
echo "Loading proxy agent ..."
launchctl load -w $LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-proxyd-launcher.plist || true

echo "Give permission to runner user to read /var/log/system.log"
sudo chmod +a "runner allow read" /var/log/system.log

echo "Verifying plist files creation..."

if [ -f "$LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-agentd-launcher.plist" ]; then
  echo "Agent plist file exists."
else
  echo "Error: Agent plist file does not exist."
  exit 1
fi

if [ -f "$LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-telemetryd-launcher.plist" ]; then
  echo "Telemetry plist file exists."
else
  echo "Error: Telemetry plist file does not exist."
  exit 1
fi

if [ -f "$LAUNCH_AGENTS_DIR/com.warpbuild.warpbuild-proxyd-launcher.plist" ]; then
  echo "Proxy plist file exists."
else
  echo "Error: Proxy plist file does not exist."
  exit 1
fi

echo "All plist files verified successfully."


echo "Agent setup complete."
