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
    Invoke-Command -Command $Test
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

function Assert-Equal {
  [CmdletBinding()]
  param (
    [string]$Expect,
    [string]$Actual
  )

  if ($Expect -ne $Actual) {
    throw "FAILED: '$Expected' == '$Actual'"
  }
}

function Assert-NotEqual {
  [CmdletBinding()]
  param (
    [string]$Expect,
    [string]$Actual
  )

  if ($Expect -eq $Actual) {
    throw "FAILED: '$Expected' != '$Actual'"
  }
}

function Assert-Empty {
  [CmdletBinding()]
  param ([string]$Actual)

  if (-not [string]::IsNullOrEmpty($Actual)) {
    throw "FAILED: '$Actual' not empty"
  }
}

function Assert-NotEmpty {
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
  Assert-Equal $env:HELLO "World"

  $env:WATCHES = $env:DIRENV_WATCHES

  Write-Host "Reloading (should be no-op)"
  Invoke-DirenvEval
  Assert-Equal $env:WATCHES $env:DIRENV_WATCHES

  Write-Host "Updating envrc and reloading (should reload)"
  touch .envrc
  Invoke-DirenvEval
  Assert-NotEqual $env:WATCHES $env:DIRENV_WATCHES

  Write-Host "Leaving dir (should clear env set by dir's envrc)"
  Set-Location ..
  Invoke-DirenvEval
  Write-Host $env:HELLO
  Assert-Empty $env:HELLO

  Remove-Item "env:/WATCHES"
}

Invoke-Test "inherit" -Test {
  Copy-Item "../base/.envrc" "../inherited/.envrc"
  Invoke-DirenvEval
  Assert-Equal $env:HELLO "World"

  Start-Sleep 1
  Write-Output "export HELLO=goodbye" | Out-File -FilePath "../inherited/.envrc"
  Invoke-DirenvEval
  Assert-Equal $env:HELLO "goodbye"
}

if (Get-Command ruby -ErrorAction SilentlyContinue) {
  Invoke-Test "ruby-layout" -Test {
    Invoke-DirenvEval
    Assert-NotEmpty $env:GEM_HOME
  }
}

Invoke-Test "space dir" -Test {
  Invoke-DirenvEval
  Assert-Equal $env:SPACE_DIR "true"
}

Invoke-Test "child-env" -Test {
  Invoke-DirenvEval
  Assert-Equal $env:PARENT_PRE "1"
  Assert-Equal $env:CHILD "1"
  Assert-Equal $env:PARENT_POST "1"
  Assert-Empty $env:REMOVE_ME
}

Invoke-Test "special-vars" -Test {
  $env:DIRENV_BASH = (Get-Command bash).Source
  $env:DIRENV_CONFIG = "foobar"
  Invoke-DirenvEval
  Assert-NotEmpty $env:DIRENV_BASH
  Assert-Equal $env:DIRENV_CONFIG "foobar"
  Remove-Item env:/DIRENV_BASH
  Remove-Item env:/DIRENV_CONFIG
}

Invoke-Test "dump" -Test {
  Invoke-DirenvEval
  Assert-Equal $env:LS_COLORS "*.ogg=38;5;45:*.wav=38;5;45"
  Assert-Equal $env:THREE_BACKSLASHES '\\\'
  Assert-Equal $env:LESSOPEN "||/usr/bin/lesspipe.sh %s"
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
  Assert-Equal $env:FOO "unset"
  Remove-Item env:/FOO
}

Invoke-Test "in-envrc" -Test {
  Invoke-DirenvEval
  ./test-in-envrc
  Assert-Equal $LASTEXITCODE "1"
}

Invoke-Test "missing-file-source-env" -Test {
  Invoke-DirenvEval
}

Invoke-Test "symlink-changed" -Test {
  ln -fs ./state-A ./symlink
  Invoke-DirenvEval
  Assert-Equal $env:STATE "A"
  Start-Sleep 1

  ln -fs ./state-B ./symlink
  Invoke-DirenvEval
  Assert-Equal $env:STATE "B"
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
  Assert-Equal $env:UTFSTUFF "♀♂"
}

Invoke-Test "failure" -Test {
  # Test that DIRENV_DIFF and DIRENV_WATCHES are set even after a failure.
  #
  # This is needed so that direnv doesn't go into a loop when the loading
  # fails.
  Assert-Empty $env:DIRENV_DIFF
  Assert-Empty $env:DIRENV_WATCHES

  Invoke-DirenvEval

  Assert-NotEmpty $env:DIRENV_DIFF
  Assert-NotEmpty $env:DIRENV_WATCHES
}

Invoke-Test "watch-dir" -Test {
  Write-Host "no watches by default"
  Assert-Equal $env:DIRENV_WATCHES $env:WATCHES

  Invoke-DirenvEval

  if (-not (direnv show_dump $env:DIRENV_WATCHES | Select-String "testfile")) {
    throw "FAILED: testfile not added to DIRENV_WATCHES"
  }

  Write-Host "After eval, watches have changed"
  Assert-NotEqual $env:DIRENV_WATCHES $env:WATCHES
}

Invoke-Test "load-envrc-before-env" -Test {
  Invoke-DirenvEval
  Assert-Equal $env:HELLO "bar"
}

Invoke-Test "load-env" -Test {
  Write-Output @"
[global]
load_dotenv = true
"@ | Out-File "${env:XDG_CONFIG_HOME}/direnv/direnv.toml"
  direnv allow
  Invoke-DirenvEval
  Assert-Equal $env:HELLO "world"
}

Invoke-Test "skip-env" -Test {
  Invoke-DirenvEval
  Assert-Empty $env:SKIPPED
}

if (Get-Command python -ErrorAction SilentlyContinue) {
  Invoke-Test "python-layout" -Test {
    if (Get-Item .direnv -ErrorAction SilentlyContinue) {
      Remove-Item .direnv -Force -Recurse -ErrorAction SilentlyContinue
    }

    Invoke-DirenvEval
    Assert-NotEmpty $env:VIRTUAL_ENV

    if (($env:PATH -split ":") -notcontains "${env:VIRTUAL_ENV}/bin") {
      throw "FAILED: VIRTUAL_ENV/bin not added to PATH"
    }

    if (-not (Get-Item ./.direnv/CACHEDIR.TAG -ErrorAction SilentlyContinue)) {
      throw "the layout dir should contain that file to filter that folder out of backups"
    }
  }

  Invoke-Test "python-custom-virtual-env" -Test {
    Invoke-DirenvEval
    if (-not (Get-Item $env:VIRTUAL_ENV)) {
      throw "${env:VIRTUAL_ENV} does not exist"
    }

    if (($env:PATH -split ":") -notcontains "$PWD/foo/bin") {
      throw "FAILED: VIRTUAL_ENV/bin not added to PATH"
    }
  }
}

Invoke-Test "aliases" -Test {
  direnv deny
  # check that allow/deny aliases work
  Write-Host "direnv permit"
  direnv permit
  Invoke-DirenvEval
  Assert-NotEmpty $env:HELLO

  Write-Host "direnv block"
  direnv block
  Invoke-DirenvEval
  Assert-Empty $env:HELLO

  Write-Host "direnv grant"
  direnv grant
  Invoke-DirenvEval
  Assert-NotEmpty $env:HELLO

  Write-Host "direnv revoke"
  direnv revoke
  Invoke-DirenvEval
  Assert-Empty $env:HELLO
}

Invoke-Test '$test' -Test {
  Invoke-DirenvEval
  Assert-Equal $env:FOO "bar"
}

Invoke-Test "special-characters/backspace/return" -Test {
  Invoke-DirenvEval
  Assert-Equal $env:HI "there"
}

#endregion

Stop-TestSuite
