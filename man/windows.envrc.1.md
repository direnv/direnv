% windows.envrc(1) direnv User Manuals

# NAME
windows.envrc - PowerShell-specific environment configuration for direnv (Windows)

# DESCRIPTION
On Windows systems, direnv will look for a file named `windows.envrc` in addition to the standard `.envrc` and optional `.env` files. If `windows.envrc` exists and the `pwsh` executable is available, direnv executes the file in a PowerShell subprocess to collect environment changes.

Unlike the bash-based `.envrc` execution model, the PowerShell loader computes a delta: only variables that were added or changed and variables that were removed are returned to direnv. This minimizes data transfer and reduces the chance of output interference.

# USAGE
Create a `windows.envrc` file in your project directory:

```powershell
# windows.envrc
$env:APP_ENV = "development"
$env:DEBUG = "1"
$env:PATH = "$PWD\bin;$env:PATH"
```

Authorize it:

```shell
direnv allow .
```

From then on entering the directory will apply the variable changes.

# GUIDELINES
* Only set environment variables using `$env:NAME = VALUE`.
* Normal informational output (`Write-Host`, `Write-Output`, warnings, verbose messages) is suppressed by the wrapper so it won't corrupt the JSON payload. Only use stdout deliberately if you are debugging and understand it will be discarded.
* To remove a variable, explicitly unset it using `Remove-Item Env:VAR -ErrorAction SilentlyContinue`; omission alone doesn't remove existing variables.
* Functions, aliases, and other non-exportable shell state are not propagated.
* If `pwsh` is not found in PATH or `enable_pwsh=false` in `direnv.toml`, the file is skipped.
* Precedence: `windows.envrc` > `.envrc` > `.env` (when enabled).

# SECURITY
`windows.envrc` follows the same authorization model as `.envrc`. Changes are ignored until the user runs `direnv allow` in the directory. Review its contents prior to allowing.

# CONFIGURATION
Set `enable_pwsh = false` under `[global]` in `direnv.toml` to disable loading of `windows.envrc` even if present.

# SEE ALSO
`direnv(1)`, `direnv.toml(1)`, `direnv-stdlib(1)`
