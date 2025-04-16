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

if (-not $env:WARPBUILD_SCOPE_TOKEN) {
    Write-Host "WARPBUILD_SCOPE_TOKEN is not set."
    exit 1
}

if (-not $env:WARPBUILD_RUNNER_SET_ID) {
    Write-Host "WARPBUILD_RUNNER_SET_ID is not set."
    exit 1
}

# Create request body
$requestBody = @{
    runner_id = $env:WARPBUILD_RUNNER_SET_ID
    runner_name = $env:RUNNER_NAME
    orchestrator_job_id = $env:GITHUB_JOB_ID
    orchestrator_job_group_id = $env:GITHUB_RUN_ID
    orchestrator_job_group_attempt = $env:GITHUB_RUN_ATTEMPT
    repo_entity = $env:GITHUB_REPOSITORY
    repo_base_ref = $env:GITHUB_BASE_REF
    repo_head_ref = $env:GITHUB_HEAD_REF
    repo_ref = $env:GITHUB_REF
    repo_ref_type = $env:GITHUB_REF_TYPE
} | ConvertTo-Json -Depth 10

Write-Host "`nMaking a request to WarpBuild..."

# Bypass SSL certificate validation (only if needed)
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = { $true }

$headers = @{
    'Content-Type' = 'application/json'
    'X-Warpbuild-Scope-Token' = $env:WARPBUILD_SCOPE_TOKEN
}

# Retry logic
$maxRetries = 5
$retryCount = 0
$success = $false
$response = $null

while (-not $success -and $retryCount -lt $maxRetries) {
    try {
        $response = Invoke-RestMethod -Uri "$env:WARPBUILD_HOST_URL/api/v1/job" `
            -Method Post `
            -Headers $headers `
            -Body $requestBody `
            -ContentType 'application/json' `
            -ErrorAction Stop

        $success = $true
    }
    catch {
        $retryCount++
        Write-Host "Request failed (attempt $retryCount/$maxRetries): $($_.Exception.Message)"
        if ($retryCount -lt $maxRetries) {
            Start-Sleep -Seconds 2
        } else {
            Write-Host "`nFailed to send request to warpbuild. Exiting..."
            exit 1
        }
    }
}

Write-Host "`nPrehook for WarpBuild runner instance '$env:RUNNER_NAME' completed successfully."
