<#
Environment variables for the project - DO NOT ADD SENSITIVE INFO HERE
PowerShell version of prior bash/direnv logic.

Usage:
  # Dot-source in PowerShell
  . ./windows.envrc

Behavior:
  - If ARTIFACTORY_USERNAME and ARTIFACTORY_PASSWORD are set, they are copied
    into UV_INDEX_ARTIFACTORY_USERNAME / UV_INDEX_ARTIFACTORY_PASSWORD.
  - Otherwise if either UV_INDEX_ variable missing, a warning is shown.
#>

$artUser = $env:ARTIFACTORY_USERNAME
$artPass = $env:ARTIFACTORY_PASSWORD
$uvUser  = $env:UV_INDEX_ARTIFACTORY_USERNAME
$uvPass  = $env:UV_INDEX_ARTIFACTORY_PASSWORD

if ([string]::IsNullOrWhiteSpace($artUser) -eq $false -and [string]::IsNullOrWhiteSpace($artPass) -eq $false) {
    $env:UV_INDEX_ARTIFACTORY_USERNAME = $artUser
    $env:UV_INDEX_ARTIFACTORY_PASSWORD = $artPass
}
elseif ([string]::IsNullOrWhiteSpace($uvUser) -or [string]::IsNullOrWhiteSpace($uvPass)) {
    Write-Error "No Artifactory credentials found." 
    Write-Error "Set either ARTIFACTORY_USERNAME/PASSWORD or UV_INDEX_ARTIFACTORY_USERNAME/PASSWORD" 
    
}

#Write-Host "UV_INDEX_ARTIFACTORY_USERNAME: $($env:UV_INDEX_ARTIFACTORY_USERNAME)"
#Write-Host "UV_INDEX_ARTIFACTORY_PASSWORD: $(if ($env:UV_INDEX_ARTIFACTORY_PASSWORD) { '********' } else { '(empty)' })"
