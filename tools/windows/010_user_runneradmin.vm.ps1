# Define user credentials
$MACHINE_USER = "runneradmin"
$MACHINE_PASSWORD = "WARPPassw0rd!2024"
$Password = ConvertTo-SecureString $MACHINE_PASSWORD -AsPlainText -Force

# Check if the user exists
if (-not (Get-LocalUser -Name $MACHINE_USER -ErrorAction SilentlyContinue)) {
    Write-Host "Creating user '$MACHINE_USER'"
    New-LocalUser -Name $MACHINE_USER -Password $Password -FullName "Runner Admin" -Description "Runner admin user for CI/CD"
} else {
    Write-Host "User '$MACHINE_USER' already exists"
}

# Ensure user is a member of 'Administrators' and 'Users' groups
foreach ($group in @("Administrators", "Users")) {
    $groupMembers = Get-LocalGroupMember -Group $group -ErrorAction SilentlyContinue
    if ($groupMembers -and $groupMembers.Name -contains $MACHINE_USER) {
        Write-Host "User '$MACHINE_USER' is already a member of '$group'"
    } else {
        Write-Host "Adding user '$MACHINE_USER' to group '$group'"
        Add-LocalGroupMember -Group $group -Member $MACHINE_USER
    }
}

# Set password to never expire
Set-LocalUser -Name $MACHINE_USER -PasswordNeverExpires $true

# Enable auto-login for the user
$RegistryPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon"
Set-ItemProperty -Path $RegistryPath -Name "AutoAdminLogon" -Value "1"
Set-ItemProperty -Path $RegistryPath -Name "DefaultUsername" -Value $MACHINE_USER
Set-ItemProperty -Path $RegistryPath -Name "DefaultPassword" -Value $MACHINE_PASSWORD

# Disable UAC
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "EnableLUA" -Value 0

# Enable PowerShell remoting
Enable-PSRemoting -Force

# Grant 'Log on as a service' right to the user
$UserName = $MACHINE_USER
$policyPath = "C:\temp\secpol.cfg"

# Ensure the temporary directory exists
if (-not (Test-Path -Path "C:\temp")) {
    New-Item -ItemType Directory -Path "C:\temp" | Out-Null
}

# Export the current security policy
secedit /export /cfg $policyPath

# Read the exported policy content
$content = Get-Content $policyPath

# Check if the user already has 'Log on as a service' right
if ($content -match "SeServiceLogonRight\s*=\s*(.*)") {
    $existingUsers = $matches[1].Split(",") | ForEach-Object { $_.Trim() }
    if ($existingUsers -contains $UserName) {
        Write-Host "User '$UserName' already has 'Log on as a service' permission."
    } else {
        # Add the user to the policy
        $updatedUsers = ($existingUsers + $UserName) -join ","
        $updatedContent = $content -replace "SeServiceLogonRight\s*=.*", "SeServiceLogonRight = $updatedUsers"
        $updatedContent | Set-Content $policyPath

        # Apply the updated policy
        secedit /configure /db secedit.sdb /cfg $policyPath /overwrite
        gpupdate /force

        Write-Host "Granted 'Log on as a service' permission to user '$UserName'."
    }
} else {
    # If the policy line doesn't exist, add it
    Add-Content -Path $policyPath -Value "SeServiceLogonRight = $UserName"
    secedit /configure /db secedit.sdb /cfg $policyPath /overwrite
    gpupdate /force

    Write-Host "Added 'Log on as a service' permission for user '$UserName'."
}
