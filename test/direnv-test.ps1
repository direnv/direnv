#!/usr/bin/env pwsh
$OriginalLocation = $PWD
$TestDir = Split-Path $MyInvocation.MyCommand.Path -Parent
Set-Location $TestDir

$env:XDG_CONFIG_HOME = "$TestDir/config"
$env:XDG_DATA_HOME = "$TestDir/data"

# Save original PATH value
$OriginalPATH = $env:PATH
$env:PATH = "$(Split-Path $TestDir -Parent):$($env:PATH)"

# Reset the direnv loading if any
$env:DIRENV_CONFIG = $TestDir
$envVarsToReset = @(
  "DIRENV_BASH",
  "DIRENV_DIR",
  "DIRENV_FILE",
  "DIRENV_WATCHES",
  "DIRENV_DIFF"
)
foreach ($var in $envVarsToReset) {
  if (Get-Item "env:/$var" -ErrorAction SilentlyContinue) {
    Remove-Item "env:/$var"
  }
}

if (!(Get-Item "${Env:XDG_CONFIG_HOME}/direnv/direnvrc" -ErrorVariable SilentlyContinue)) {
  New-Item -Path "${Env:XDG_CONFIG_HOME}/direnv/direnvrc" -ItemType File
}

function Stop-TestSuite {
  [CmdletBinding()]
  param ([int]$ExitCode = 0)
  Stop-DirenvTest
  $env:PATH = $OriginalPATH
  Set-Location $OriginalLocation
  exit $ExitCode
}

function Invoke-Test {
  param (
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]$Name,
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [scriptblock]$Test
  )

  try {
    Start-DirenvTest $Name
    $Test.Invoke()
    Stop-DirenvTest
  }
  catch {
    Write-Error "${Name}: $_"
    Stop-TestSuite 1
  }
}

function Invoke-DirenvEval {
  $export = direnv export pwsh
  if ($export) {
    Invoke-Expression $export
  }
}

function Test-Equal {
  [CmdletBinding()]
  param (
    [string]$Expect,
    [string]$Actual
  )

  if ($Expect -ne $Actual) {
    throw "FAILED: '$Expected' == '$Actual'"
  }
}

function Test-NotEqual {
  [CmdletBinding()]
  param (
    [string]$Expect,
    [string]$Actual
  )

  if ($Expect -eq $Actual) {
    throw "FAILED: '$Expected' != '$Actual'"
  }
}

function Test-Empty {
  [CmdletBinding()]
  param ([string]$Actual)

  if (-not [string]::IsNullOrEmpty($Actual)) {
    throw "FAILED: '$Actual' not empty"
  }
}

function Test-NonEmpty {
  [CmdletBinding()]
  param ([string]$Actual)

  if ([string]::IsNullOrEmpty($Actual)) {
    throw "FAILED: '$Actual' empty"
  }
}


function Start-DirenvTest {
  [CmdletBinding()]
  param ([string]$Scenario)

  Set-Location "$TestDir/scenarios/$Scenario"
  direnv allow
  if ($env:DIRENV_DEBUG -eq "1") {
    Write-Host
  }
  Write-Host "## Testing $Scenario ##" -ForegroundColor Green
  if ($env:DIRENV_DEBUG -eq "1") {
    Write-Host
  }
}

function Stop-DirenvTest {
  Remove-Item "${Env:XDG_CONFIG_HOME}/direnv/direnv.toml" -ErrorAction SilentlyContinue
  cd\ # Built-in function to move to root.
  Invoke-DirenvEval
}

#region Test
direnv allow
Invoke-DirenvEval

Invoke-Test "base" -Test {
  Write-Host "Setting up"
  Invoke-DirenvEval
  Test-Equal $env:HELLO "World"

  $env:WATCHES = $env:DIRENV_WATCHES

  Write-Host "Reloading (should be no-op)"
  Invoke-DirenvEval
  Test-Equal $env:WATCHES $env:DIRENV_WATCHES

  Write-Host "Updating envrc and reloading (should reload)"
  touch .envrc
  Invoke-DirenvEval
  Test-NotEqual $env:WATCHES $env:DIRENV_WATCHES

  Write-Host "Leaving dir (should clear env set by dir's envrc)"
  Set-Location ..
  Invoke-DirenvEval
  Write-Host $env:HELLO
  Test-Empty $env:HELLO

  Remove-Item "env:/WATCHES"
}

Invoke-Test "inherit" -Test {
  Copy-Item "../base/.envrc" "../inherited/.envrc"
  Invoke-DirenvEval
  Test-Equal $env:HELLO "World"

  Start-Sleep 1
  Write-Output "export HELLO=goodbye" | Out-File -FilePath "../inherited/.envrc"
  Invoke-DirenvEval
  Test-Equal $env:HELLO "goodbye"
}

#region ruby scenario
#TODO:
#endregion

Invoke-Test "space dir" -Test {
  Invoke-DirenvEval
  Test-Equal $env:SPACE_DIR "true"
}

Invoke-Test "child-env" -Test {
  Invoke-DirenvEval
  Test-Equal $env:PARENT_PRE "1"
  Test-Equal $env:CHILD "1"
  Test-Equal $env:PARENT_POST "1"
  Test-Empty $env:REMOVE_ME
}

Invoke-Test "special-vars" -Test {
  $env:DIRENV_BASH = "$(command -v bash)"
  $env:DIRENV_CONFIG = "foobar"
  Invoke-DirenvEval
  Test-NonEmpty $env:DIRENV_BASH
  Test-Equal $env:DIRENV_CONFIG "foobar"
  Remove-Item env:/DIRENV_BASH
  Remove-Item env:/DIRENV_CONFIG
}

Invoke-Test "dump" -Test {
  Invoke-DirenvEval
  Test-Equal $env:LS_COLORS "*.ogg=38;5;45:*.wav=38;5;45"
  Test-Equal $env:THREE_BACKSLASHES '\\\'
  Test-Equal $env:LESSOPEN "||/usr/bin/lesspipe.sh %s"
}

#region empty-var scenario

# PowerShell assumes an empty environment variable is unset.

#endregion

Invoke-Test "empty-var-unset" -Test {
  $env:FOO = ""
  Invoke-DirenvEval

  if (-not $env:FOO) {
    $env:FOO = "unset"
  }
  Test-Equal $env:FOO "unset"
  Remove-Item env:/FOO
}

Invoke-Test "in-envrc" -Test {
  Invoke-DirenvEval
  ./test-in-envrc
  Test-Equal $LASTEXITCODE "1"
}

Invoke-Test "missing-file-source-env" -Test {
  Invoke-DirenvEval
}

Invoke-Test "symlink-changed" -Test {
  ln -fs ./state-A ./symlink
  Invoke-DirenvEval
  Test-Equal $env:STATE "A"
  Start-Sleep 1

  ln -fs ./state-B ./symlink
  Invoke-DirenvEval
  Test-Equal $env:STATE "B"
}

Invoke-Test "symlink-dir" -Test {
  # we can allow and deny the target
  direnv allow foo
  direnv deny foo
  # we can allow and deny the symlink
  direnv allow bar
  direnv deny bar
}

Invoke-Test "utf-8" -Test {
  Invoke-DirenvEval
  Test-Equal $env:UTFSTUFF "♀♂"
}

Invoke-Test "failure" -Test {
  # Test that DIRENV_DIFF and DIRENV_WATCHES are set even after a failure.
  #
  # This is needed so that direnv doesn't go into a loop when the loading
  # fails.
  Test-Empty $env:DIRENV_DIFF
  Test-Empty $env:DIRENV_WATCHES

  Invoke-DirenvEval

  Test-NonEmpty $env:DIRENV_DIFF
  Test-NonEmpty $env:DIRENV_WATCHES
}

Invoke-Test "watch-dir" -Test {
  Write-Host "no watches by default"
  Test-Equal $env:DIRENV_WATCHES $env:WATCHES

  Invoke-DirenvEval

  if (-not (direnv show_dump $env:DIRENV_WATCHES | Select-String "testfile")) {
    throw "FAILED: testfile not added to DIRENV_WATCHES"
  }

  Write-Host "After eval, watches have changed"
  Test-NotEqual $env:DIRENV_WATCHES $env:WATCHES
  Remove-Item "env:/WATCHES"
}

Invoke-Test "load-envrc-before-env" -Test {
  Invoke-DirenvEval
  Test-Equal $env:HELLO "bar"
}

Invoke-Test "load-env" -Test {
  Write-Output @"
[global]
load_dotenv = true
"@ | Out-File "${env:XDG_CONFIG_HOME}/direnv/direnv.toml"
  direnv allow
  Invoke-DirenvEval
  Test-Equal $env:HELLO "world"
}

Invoke-Test "skip-env" -Test {
  Invoke-DirenvEval
  Test-Empty $env:SKIPPED
}

if (Get-Command python -ErrorAction SilentlyContinue) {
  Invoke-Test "python-layout" -Test {
    Remove-Item .direnv -Force

    Invoke-DirenvEval
    Test-NonEmpty $env:VIRTUAL_ENV

    # TODO: Layout is currently bash-only. Must solve for this.
    if (($env:PATH -split ":") -notcontains "${env:VIRTUAL_ENV}/bin") {
      throw "FAILED: VIRTUAL_ENV/bin not added to PATH"
    }
  }
}

#endregion

Stop-TestSuite
