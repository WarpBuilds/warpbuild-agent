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

Write-Host "[tailscale-addon] Ensuring Tailscale service is running..."
net start Tailscale 2>&1 | ForEach-Object { Write-Host "[tailscale-addon] $_" }

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
        Write-Host "[tailscale-addon]   extra args: $env:WARPBUILD_ADDON_TS_ARGS"
        $tsArgs += $env:WARPBUILD_ADDON_TS_ARGS -split '\s+'
    }

    $ErrorActionPreference = "Continue"
    $upOutput = & $TS_BIN @tsArgs 2>&1 | Out-String
    $upExit = $LASTEXITCODE
    $ErrorActionPreference = "Stop"

    if ($upExit -eq 0) {
        Write-Host "[tailscale-addon] Tailscale is connected"
        & $TS_BIN status 2>&1 | ForEach-Object { Write-Host "[tailscale-addon]   $_" }
        Write-Host "[tailscale-addon] Tailscale setup complete"
        exit 0
    }

    Write-Host "[tailscale-addon] tailscale up failed (exit $upExit), retrying in 2s... ($loginElapsed/$loginTimeout)"
    if ($upOutput.Trim()) {
        Write-Host "[tailscale-addon]   $($upOutput.Trim())"
    }
    Start-Sleep -Seconds 2
    $loginElapsed += 2
}

Write-Host "[tailscale-addon] ERROR: tailscale up failed after ${loginTimeout}s"
exit 1
