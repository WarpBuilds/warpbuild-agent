$BASE_DIR = "C:\warpbuilds"
$MACHINE_USER = "runneradmin"
$MACHINE_PASSWORD = "WARPPassw0rd!2024"

# Create the base directory if it doesn't exist
if (-not (Test-Path -Path $BASE_DIR)) {
    New-Item -Path $BASE_DIR -ItemType Directory -Force
}

Set-Location -Path $BASE_DIR

$AGENT_DIR = "$BASE_DIR\warpbuild-agent"
if (-not (Test-Path -Path $AGENT_DIR)) {
    New-Item -ItemType Directory -Path $AGENT_DIR
}

# Get the latest agent tag
# $latestAgentTag = (Invoke-RestMethod -Uri "https://api.github.com/repos/WarpBuilds/warpbuild-agent/releases/latest").tag_name
$latestAgentTag = "feat/windows-telemetry"  # Uncomment this line to use a specific version

# Download and extract the agent
$agentUrl = "https://pub-b4f7dbd911ef411ca27c8befa94bb744.r2.dev/warpbuild-agentd/$latestAgentTag/warpbuild-agentd_Windows_x86_64.zip"
Write-Host "Downloading warpbuild-agent using aria2 from r2..."
aria2c -s16 -x16 $agentUrl -d "$AGENT_DIR" -o "warpbuild-agent.zip"
Expand-Archive -Path "$AGENT_DIR\warpbuild-agent.zip" -DestinationPath "$AGENT_DIR" -Force
Remove-Item -Path "$AGENT_DIR\warpbuild-agent.zip"

# Copy the agent executable to the system path
Copy-Item -Path "$AGENT_DIR\warpbuild-agentd.exe" -Destination "C:\Windows\System32\warpbuild-agentd.exe"

$restarterZipPath = "warpbuild-agentd-restarter_Windows_x86_64.zip"
$restarterUrl = "https://pub-b4f7dbd911ef411ca27c8befa94bb744.r2.dev/warpbuild-agentd/$latestAgentTag/$restarterZipPath"
Write-Host "Downloading warpbuild-agentd-restarter using aria2..."
aria2c -s16 -x16 $restarterUrl -d "$AGENT_DIR" -o "$restarterZipPath"
Write-Host "Downloaded $restarterZipPath"
Expand-Archive -Path "$AGENT_DIR\$restarterZipPath" -DestinationPath "$AGENT_DIR" -Force
Remove-Item -Path "$AGENT_DIR\$restarterZipPath"
Write-Host "Extracted warpbuild-agentd-restarter.exe"
Copy-Item -Path "$AGENT_DIR\warpbuild-agentd-restarter.exe" -Destination "C:\Windows\System32\warpbuild-agentd-restarter.exe"
Write-Host "Copied warpbuild-agentd-restarter.exe to C:\Windows\System32"

# Get the domain of the current user
$domain = (Get-WmiObject Win32_ComputerSystem).Domain
Write-Host "Current user domain: $domain"

$SETTINGS_FILE = "$BASE_DIR\settings.json"
$AGENTD_STDOUT_FILE = "$BASE_DIR\warpbuild-agentd.stdout.log"
$AGENTD_STDERR_FILE = "$BASE_DIR\warpbuild-agentd.stderr.log"

$TELEMETRY_STDOUT_FILE = "$BASE_DIR\warpbuild-telemetryd.stdout.log"
$TELEMETRY_STDERR_FILE = "$BASE_DIR\warpbuild-telemetryd.stderr.log"

$PROXY_STDOUT_FILE = "$BASE_DIR\warpbuild-proxyd.stdout.log"
$PROXY_STDERR_FILE = "$BASE_DIR\warpbuild-proxyd.stderr.log"

$RESTARTER_STDOUT_FILE = "$BASE_DIR\warpbuild-agentd-restarter.stdout.log"
$RESTARTER_STDERR_FILE = "$BASE_DIR\warpbuild-agentd-restarter.stderr.log"

# Create and configure the services
$services = @(
    @{
        Name = "warpbuild-agentd"
        DisplayName = "WarpBuild Agent"
        Description = "WarpBuild Agent Service"
        BinaryPath = "C:\Windows\System32\warpbuild-agentd.exe --settings $SETTINGS_FILE --stdout $AGENTD_STDOUT_FILE --stderr $AGENTD_STDERR_FILE"
        UserName = ".\$MACHINE_USER"
        Password = $MACHINE_PASSWORD
        StartupType = "Manual"
        Dependencies = @()
        Environment = @{}
    },
    @{
        Name = "warpbuild-agentd-restarter"
        DisplayName = "WarpBuild Agent Restarter"
        Description = "WarpBuild Agent Restarter Service"
        BinaryPath = "C:\Windows\System32\warpbuild-agentd-restarter.exe --restart-interval 100ms --stdout $RESTARTER_STDOUT_FILE --stderr $RESTARTER_STDERR_FILE"
        StartupType = "Manual"
        Dependencies = @()
        Environment = @{}
    },
    @{
        Name = "warpbuild-telemetryd"
        DisplayName = "WarpBuild Telemetry"
        Description = "WarpBuild Telemetry Service"
        BinaryPath = "C:\Windows\System32\warpbuild-agentd.exe --settings $SETTINGS_FILE --launch-telemetry=true --stdout $TELEMETRY_STDOUT_FILE --stderr $TELEMETRY_STDERR_FILE"
        UserName = ".\$MACHINE_USER"
        Password = $MACHINE_PASSWORD
        StartupType = "Automatic"
        IsDelayed = $true
        Dependencies = @()
        Environment = @{}
    },
    @{
        Name = "warpbuild-proxyd"
        DisplayName = "WarpBuild Proxy"
        Description = "WarpBuild Proxy Service"
        BinaryPath = "C:\Windows\System32\warpbuild-agentd.exe --settings $SETTINGS_FILE --launch-proxy-server=true --stdout $PROXY_STDOUT_FILE --stderr $PROXY_STDERR_FILE"
        StartupType = "Automatic"
        Dependencies = @()
        Environment = @{}
    }
)

foreach ($service in $services) {
    $existingService = Get-Service -Name $service.Name -ErrorAction SilentlyContinue
    if ($existingService) {
        Write-Host "Service $($service.Name) already exists. Removing it."
        Stop-Service -Name $service.Name -Force
        sc.exe delete $service.Name
    }

    # $username = "$service.UserName"
    # $password = "$service.Password"
    # $securepassword = ConvertTo-SecureString $password -AsPlainText -Force
    # $cred = New-Object System.Management.Automation.PSCredential ($username, $securepassword)

    Write-Host "Creating service: $($service.Name)"
    Write-Host "Running Command: New-Service -Name $($service.Name) -BinaryPathName $($service.BinaryPath) -DisplayName $($service.DisplayName) -Description $($service.Description) -StartupType $($service.StartupType)"

    New-Service -Name $service.Name -BinaryPathName $service.BinaryPath  -DisplayName $service.DisplayName -Description $service.Description -StartupType $service.StartupType
    # $command = "sc.exe config $($service.Name) obj= `"$($service.UserName)`" password= `"$($service.Password)`""
    # Write-Host "Executing: $command"
    # Invoke-Expression $command

    # Set service dependencies
    if ($service.Dependencies.Count -gt 0) {
        $dependencies = $service.Dependencies -join "/"
        sc.exe config $service.Name depend= $dependencies
    }

    # Set environment variables
    $envString = ""
    foreach ($key in $service.Environment.Keys) {
        $envString += "$key=$($service.Environment[$key])`0"
    }
    if ($envString -ne "") {
        sc.exe config $service.Name env= $envString
    }

    # if this is the telemetry service, flip on delayed auto‑start
    if ($service.IsDelayed) {
        Write-Host "Enabling Delayed‑AutoStart for $($service.Name)…"
        sc.exe config $service.Name start= delayed-auto
    }

    # Define variables for the service configuration
    $serviceName = $service.Name       # Replace with the name of your service
    $windowsLogonUser = $service.UserName  # Replace DOMAIN\Username with the user account for the service
    $windowsLogonPassword = $service.Password # Replace with the password for the user account

    # Validate user input
    if ([string]::IsNullOrWhiteSpace($serviceName) -or 
        [string]::IsNullOrWhiteSpace($windowsLogonUser) -or 
        [string]::IsNullOrWhiteSpace($windowsLogonPassword)) {
        Write-Host "Please ensure all variables (serviceName, windowsLogonUser, windowsLogonPassword) are populated to override the default values."
        continue
    }

    # Configure the service to use the specified user account and password
    try {
        $command = "sc.exe config $serviceName obj= `"$windowsLogonUser`" password= `"$windowsLogonPassword`""
        Write-Host "Executing: $command"
        Invoke-Expression $command

        # Check for success
        $service = Get-Service -Name $serviceName -ErrorAction Stop
        if ($service) {
            Write-Host "Service '$serviceName' has been configured to run as $windowsLogonUser."
            
            # Restart the service to apply changes
            # Restart-Service -Name $serviceName -Force
            # Write-Host "Service '$serviceName' has been restarted successfully."

            # Get-Service -Name $serviceName
        }
    } catch {
        Write-Host "An error occurred: $_"
    }

}