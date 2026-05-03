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

# Ensure a single clean daemon: stop service and kill any stray tailscaled processes
Write-Host "[tailscale-addon] Ensuring clean daemon state..."
Stop-Service Tailscale -ErrorAction SilentlyContinue
Start-Sleep -Seconds 1
Get-Process tailscaled -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 1

Write-Host "[tailscale-addon] Starting Tailscale service..."
Start-Service Tailscale
Start-Sleep -Seconds 2

$svc = Get-Service Tailscale -ErrorAction SilentlyContinue
Write-Host "[tailscale-addon] Service state: $($svc.Status)"
$procs = @(Get-Process tailscaled -ErrorAction SilentlyContinue)
Write-Host "[tailscale-addon] tailscaled process count: $($procs.Count)"

Write-Host "[tailscale-addon] Performing OIDC login..."
$loginTimeout = 30
$loginElapsed = 0
while ($loginElapsed -lt $loginTimeout) {
    Write-Host "[tailscale-addon] Attempt at ${loginElapsed}s - calling tailscale up..."
    Write-Host "[tailscale-addon]   --hostname=warpbuild-$env:WARPBUILD_ADDON_TS_RUNNER_INSTANCE_ID"
    Write-Host "[tailscale-addon]   --client-id=$env:WARPBUILD_ADDON_TS_CLIENT_ID"
    Write-Host "[tailscale-addon]   --id-token=<redacted:$($env:WARPBUILD_ADDON_TS_OIDC_TOKEN.Length) chars>"

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

    $upOutput = & $TS_BIN @tsArgs 2>&1 | Out-String
    $upExit = $LASTEXITCODE

    Write-Host "[tailscale-addon] tailscale up exit code: $upExit"
    if ($upOutput.Trim()) {
        Write-Host "[tailscale-addon] tailscale up output: $($upOutput.Trim())"
    }

    if ($upExit -ne 0) {
        Write-Host "[tailscale-addon] tailscale up failed, retrying in 2s... ($loginElapsed/$loginTimeout)"
        Start-Sleep -Seconds 2
        $loginElapsed += 2
        continue
    }

    # tailscale up returned 0 — now wait for the connection to actually establish
    Write-Host "[tailscale-addon] tailscale up succeeded, waiting for connection..."
    $connectTimeout = 30
    $connectElapsed = 0
    while ($connectElapsed -lt $connectTimeout) {
        $statusJson = & $TS_BIN status --json 2>&1 | Out-String
        try {
            $statusObj = $statusJson | ConvertFrom-Json
            $backendState = $statusObj.BackendState
            Write-Host "[tailscale-addon] Backend state: $backendState (${connectElapsed}s)"
            if ($backendState -eq "Running") {
                Write-Host "[tailscale-addon] Tailscale connected:"
                & $TS_BIN status 2>&1 | ForEach-Object { Write-Host "[tailscale-addon]   $_" }
                Write-Host "[tailscale-addon] Tailscale setup complete"
                exit 0
            }
        } catch {
            Write-Host "[tailscale-addon] Could not parse status JSON, waiting..."
        }
        Start-Sleep -Seconds 2
        $connectElapsed += 2
    }
    Write-Host "[tailscale-addon] ERROR: Tailscale did not reach Running state after ${connectTimeout}s"
    & $TS_BIN status 2>&1 | ForEach-Object { Write-Host "[tailscale-addon]   $_" }
    exit 1
}

Write-Host "[tailscale-addon] ERROR: tailscale up failed after ${loginTimeout}s"
exit 1
