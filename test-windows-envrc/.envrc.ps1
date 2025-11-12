# windows.envrc - PowerShell-specific envrc for direnv
# Safe execution test
$env:HELLO = "from-windows-envrc"
$env:PWSH_SAFE = "temp-wrapper-ok"
# Avoid stray output; only set env vars.
