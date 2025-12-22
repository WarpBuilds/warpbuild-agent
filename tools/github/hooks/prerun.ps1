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

Write-Host "`nPrehook for WarpBuild runner instance '$env:RUNNER_NAME' completed successfully."