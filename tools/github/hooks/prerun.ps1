Write-Host "Running prehook for WarpBuild runner instance '$env:RUNNER_NAME'..."
Write-Host "`nLogging environment variables..."

# Environment variable logging
Write-Host "GITHUB_RUN_ID=$env:GITHUB_RUN_ID"
Write-Host "GITHUB_RUN_ATTEMPT=$env:GITHUB_RUN_ATTEMPT"
Write-Host "GITHUB_JOB=$env:GITHUB_JOB"
Write-Host "GITHUB_REPOSITORY=$env:GITHUB_REPOSITORY"
Write-Host "GITHUB_BASE_REF=$env:GITHUB_BASE_REF"
Write-Host "GITHUB_HEAD_REF=$env:GITHUB_HEAD_REF"
Write-Host "GITHUB_REF=$env:GITHUB_REF"
Write-Host "GITHUB_REF_TYPE=$env:GITHUB_REF_TYPE"
Write-Host "RUNNER_NAME=$env:RUNNER_NAME"
Write-Host "RUNNER_OS=$env:RUNNER_OS"
Write-Host "WARPBUILD_RUNNER_SET_ID=$env:WARPBUILD_RUNNER_SET_ID"

if ($env:WARPBUILD_SNAPSHOT_KEY) {
    Write-Host "WARPBUILD_SNAPSHOT_KEY=$env:WARPBUILD_SNAPSHOT_KEY"
}

if (-not $env:WARPBUILD_RUNNER_VERIFICATION_TOKEN) {
    Write-Host "WARPBUILD_RUNNER_VERIFICATION_TOKEN is not set."
    exit 1
}

if (-not $env:RUNNER_NAME) {
    Write-Host "RUNNER_NAME is not set."
    exit 1
}

# Create request body
$requestBody = @{
    vcs_workflow_run_id = $env:GITHUB_RUN_ID
    vcs_workflow_run_attempt = $env:GITHUB_RUN_ATTEMPT
    repo_entity = $env:GITHUB_REPOSITORY
    repo_base_ref = $env:GITHUB_BASE_REF
    repo_head_ref = $env:GITHUB_HEAD_REF
    repo_ref = $env:GITHUB_REF
    repo_ref_type = $env:GITHUB_REF_TYPE
} | ConvertTo-Json

Write-Host "`nMaking a request to WarpBuild..."

# Bypass SSL certificate validation
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = { $true }

try {
    $headers = @{
        'Content-Type' = 'application/json'
        'Authorization' = "Bearer $env:WARPBUILD_RUNNER_VERIFICATION_TOKEN"
    }

    # PowerShell equivalent of wget with retries
    $maxRetries = 5
    $retryCount = 0
    $success = $false

    while (-not $success -and $retryCount -lt $maxRetries) {
        try {
            $response = Invoke-WebRequest -Uri "$env:WARPBUILD_HOST_URL/api/v1/runners_instance/$env:RUNNER_NAME/pre_hook" `
                -Method Post `
                -Headers $headers `
                -Body $requestBody `
                -UseBasicParsing `
                -ErrorAction Stop
            $success = $true
        }
        catch {
            $retryCount++
            if ($retryCount -eq $maxRetries) {
                throw
            }
            Start-Sleep -Seconds 2
        }
    }
}
catch {
    Write-Host "Failed to send request to warpbuild. Logging response. Exiting..."
    if ($response) {
        $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 100
    }
    else {
        Write-Host $_.Exception.Message
    }
    exit 1
}

# Execute addon setup scripts returned by the backend
$toolsDir = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path

if ($response -and $response.Content) {
    try {
        $data = $response.Content | ConvertFrom-Json
        $scripts = $data.setup_scripts
        if ($scripts -and $scripts.Count -gt 0) {
            Write-Host "`nExecuting $($scripts.Count) addon setup script(s)..."
            for ($i = 0; $i -lt $scripts.Count; $i++) {
                $scriptName = if ($scripts[$i].name) { $scripts[$i].name } else { "addon-$i" }
                Write-Host "`n[addon:$scriptName] Starting..."

                $addonScriptPath = $scripts[$i].script_path
                $fullPath = Join-Path $toolsDir $addonScriptPath
                if (-not $addonScriptPath -or -not (Test-Path $fullPath)) {
                    Write-Host "[addon:$scriptName] FAILED: script not found at $fullPath"
                    exit 1
                }

                # Set env vars for addon subprocess, clean up after
                $envMap = $scripts[$i].env
                $envKeys = @()
                if ($envMap) {
                    $envMap.PSObject.Properties | ForEach-Object {
                        $envKeys += $_.Name
                        [System.Environment]::SetEnvironmentVariable($_.Name, $_.Value, "Process")
                    }
                }

                if ($fullPath -match '\.ps1$') {
                    & powershell -ExecutionPolicy Bypass -File $fullPath
                } else {
                    & cmd.exe /c $fullPath
                }
                $addonExit = $LASTEXITCODE

                # Clean up addon env vars
                foreach ($key in $envKeys) {
                    [System.Environment]::SetEnvironmentVariable($key, $null, "Process")
                }

                if ($addonExit -ne 0) {
                    Write-Host "[addon:$scriptName] FAILED with exit code $addonExit"
                    exit 1
                }
                Write-Host "[addon:$scriptName] Completed successfully."
            }
        }
    }
    catch {
        Write-Host "[addon] FAILED to parse addon setup scripts: $($_.Exception.Message)"
        exit 1
    }
}

# BYOC substitute to ACTIONS_RUNNER_HOOK_JOB_STARTED
# See: https://docs.github.com/en/actions/how-tos/manage-runners/self-hosted-runners/run-scripts
$byocPreHook = $env:WARPBUILD_ACTIONS_RUNNER_HOOK_JOB_STARTED
if ($byocPreHook) {
    Write-Host "Found user-defined pre-hook script (WARPBUILD_ACTIONS_RUNNER_HOOK_JOB_STARTED): $byocPreHook"

    if (-not [System.IO.Path]::IsPathRooted($byocPreHook)) {
        Write-Host "User-defined pre-hook script path must be absolute: $byocPreHook"
        exit 1
    }

    if (-not (Test-Path -Path $byocPreHook -PathType Leaf)) {
        Write-Host "User-defined pre-hook script not found at: $byocPreHook"
        exit 1
    }

    Write-Host "Executing user-defined pre-hook"
    $hookExt = [System.IO.Path]::GetExtension($byocPreHook).ToLowerInvariant()
    if ($hookExt -eq '.ps1') {
        $hookCommand = ". '$byocPreHook'"
        if (Get-Command pwsh -ErrorAction SilentlyContinue) {
            & pwsh -command $hookCommand
        } else {
            & powershell -command $hookCommand
        }
    } elseif ($hookExt -eq '.sh') {
        if (Get-Command bash -ErrorAction SilentlyContinue) {
            & bash --noprofile --norc -e -o pipefail $byocPreHook
        } else {
            Write-Host "Cannot run .sh pre-hook: bash is not installed."
            exit 1
        }
    } else {
        Write-Host "User-defined pre-hook script has an unsupported extension. Supported: .sh, .ps1"
        exit 1
    }
    $hookExitCode = $LASTEXITCODE
    if ($hookExitCode -ne 0) {
        Write-Host "User-defined pre-hook exited with non-zero status: $hookExitCode"
        exit $hookExitCode
    }
    Write-Host "User-defined pre-hook completed successfully."
}

Write-Host "`nPrehook for WarpBuild runner instance '$env:RUNNER_NAME' completed successfully."