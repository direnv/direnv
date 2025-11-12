% envrc.ps1(1) direnv User Manuals

# NAME
envrc.ps1 - PowerShell-specific environment configuration for direnv (Windows)

# DESCRIPTION
On Windows systems, direnv will look for a file named `.envrc.ps1` (in addition to the standard `.envrc` and optional `.env`). If `.envrc.ps1` exists and the `pwsh` executable is available (and PowerShell support enabled), direnv executes the file in a PowerShell subprocess to collect environment changes.

The PowerShell loader computes a delta: only variables that were added or changed and variables that were removed are returned to direnv. This minimizes data transfer and prevents output interference.

# USAGE
Create a `.envrc.ps1` file in your project directory:

```powershell
# .envrc.ps1
$env:APP_ENV = "development"
$env:DEBUG = "1"
$env:PATH = "$PWD\bin;$env:PATH"
```

Authorize it:

```shell
direnv allow .
```

From then on, entering the directory will apply the variable changes.

# GUIDELINES
* Set environment variables using `$env:NAME = VALUE`.
* Normal informational output (Write-Host / Write-Output / Write-Warning / Write-Verbose) is suppressed; use `[Console]::Error.WriteLine()` for user-visible messages.
* To remove a variable, explicitly unset it with `Remove-Item Env:VAR -ErrorAction SilentlyContinue` or `$env:VAR = $null`; assigning an empty string leaves it defined.
* Functions, aliases, and other non-exportable shell state are not propagated.
* If `pwsh` is not found in PATH or `enable_pwsh=false` in `direnv.toml`, the file is ignored.
* Precedence (Windows): `.envrc.ps1` > `.envrc` > `.env`.

# SECURITY
`.envrc.ps1` follows the same authorization model as `.envrc`. Changes are ignored until the user runs `direnv allow` in the directory. Review its contents prior to allowing.

# CONFIGURATION
Set `enable_pwsh = false` under `[global]` in `direnv.toml` to disable loading of `.envrc.ps1` even if present.

# MIGRATION
The legacy filename `windows.envrc` has been removed. Rename any existing `windows.envrc` to `.envrc.ps1`.

# SEE ALSO
`direnv(1)`, `direnv.toml(1)`, `direnv-stdlib(1)`
