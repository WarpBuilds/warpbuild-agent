$ErrorActionPreference = "Stop"

$TS_BIN = "C:\Program Files\Tailscale\tailscale.exe"
$INSTALLER_URL = "https://pkgs.tailscale.com/stable/tailscale-setup-1.96.3-amd64.msi"
$INSTALLER_PATH = "$env:TEMP\tailscale-setup.msi"

Write-Host "[tailscale-addon] Starting Tailscale setup for runner $env:WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID (Windows)"

if (-not (Test-Path $TS_BIN)) {
    Write-Host "[tailscale-addon] Tailscale not found, installing via MSI..."
    Invoke-WebRequest -Uri $INSTALLER_URL -OutFile $INSTALLER_PATH -UseBasicParsing
    Start-Process msiexec -ArgumentList "/i `"$INSTALLER_PATH`" /quiet /norestart" -Wait
    Remove-Item $INSTALLER_PATH -ErrorAction SilentlyContinue
    if (-not (Test-Path $TS_BIN)) {
        Write-Host "[tailscale-addon] ERROR: Tailscale installation failed"
        exit 1
    }
}

& $TS_BIN version

Write-Host "[tailscale-addon] Starting Tailscale service..."
Start-Service Tailscale -ErrorAction SilentlyContinue

Write-Host "[tailscale-addon] Performing OIDC login..."
$loginTimeout = 30
$loginElapsed = 0
while ($loginElapsed -lt $loginTimeout) {
    $tsArgs = @(
        "up",
        "--hostname=warpbuild-$env:WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID",
        "--client-id=$env:WARPBUILD_ADDON_TS_CLIENT_ID",
        "--id-token=$env:WARPBUILD_ADDON_TS_OIDC_TOKEN"
    )
    if ($env:WARPBUILD_ADDON_TS_ARGS) {
        $tsArgs += $env:WARPBUILD_ADDON_TS_ARGS -split '\s+'
    }

    $proc = Start-Process -FilePath $TS_BIN -ArgumentList $tsArgs -Wait -PassThru -NoNewWindow 2>&1
    if ($proc.ExitCode -eq 0) {
        Write-Host "[tailscale-addon] Tailscale is connected"
        & $TS_BIN status
        Write-Host "[tailscale-addon] Tailscale setup complete"
        exit 0
    }

    Write-Host "[tailscale-addon] Retrying in 2s... ($loginElapsed/$loginTimeout)"
    Start-Sleep -Seconds 2
    $loginElapsed += 2
}

Write-Host "[tailscale-addon] ERROR: tailscale up failed after ${loginTimeout}s"
exit 1
