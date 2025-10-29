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

$UserName = $MACHINE_USER

Write-Host "UserName for service elevate access: $UserName"

# Define the path for the temporary policy file
$policyPath = "C:\temp\secpol.cfg"

# Ensure the temporary directory exists
if (-not (Test-Path -Path "C:\temp")) {
    New-Item -ItemType Directory -Path "C:\temp" | Out-Null
}

# Export the current security policy to a temporary file
secedit /export /cfg $policyPath

# Read the exported policy content
$content = Get-Content $policyPath

# Check if warpbuild is already in the SeServiceLogonRight policy
if ($content -match "SeServiceLogonRight = .*?($UserName)") {
    Write-Host "User $UserName already has 'Log on as a Service' permission."
} else {
    # Add the username to the "Log on as a service" policy (SeServiceLogonRight)
    $updatedContent = $content -replace "(SeServiceLogonRight = .*)", "`$1,$UserName"

    # Save the modified policy file
    $updatedContent | Set-Content $policyPath

    # Apply the modified policy non-interactively
    echo y | secedit /configure /db secedit.sdb /cfg $policyPath /overwrite

    # Force policy update
    gpupdate /force

    Write-Host "User $UserName has been granted the 'Log on as a Service' permission."
}