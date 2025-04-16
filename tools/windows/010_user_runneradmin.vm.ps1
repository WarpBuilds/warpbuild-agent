# Create runneradmin user and add the user to relevant groups
Write-Host 'Create runneradmin user'
$MACHINE_USER = "runneradmin"
$MACHINE_PASSWORD = "WARPPassw0rd!2024"
$Password = ConvertTo-SecureString $MACHINE_PASSWORD -AsPlainText -Force
New-LocalUser -Name $MACHINE_USER -Password $Password -FullName "Runner Admin" -Description "Runner admin user for CI/CD"
Add-LocalGroupMember -Group "Administrators" -Member $MACHINE_USER
Add-LocalGroupMember -Group "Users" -Member $MACHINE_USER

# Set runneradmin user to not require password change at next logon
Set-LocalUser -Name $MACHINE_USER -PasswordNeverExpires $true

# Enable auto-login for runneradmin user
Write-Host 'Enable auto-login for runneradmin user'
$RegistryPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon"
Set-ItemProperty -Path $RegistryPath -Name "AutoAdminLogon" -Value "1"
Set-ItemProperty -Path $RegistryPath -Name "DefaultUsername" -Value $MACHINE_USER
Set-ItemProperty -Path $RegistryPath -Name "DefaultPassword" -Value $MACHINE_PASSWORD

# Disable UAC
Write-Host 'Disable UAC'
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "EnableLUA" -Value 0

# Enable PowerShell remoting
Write-Host 'Enable PowerShell remoting'
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