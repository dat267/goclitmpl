param (
    [Parameter(Mandatory=$true)]
    [string]$NewModule
)

$OldModule = "github.com/dat267/goclitmpl"
$OldApp = "goclitmpl"
$NewApp = Split-Path -Leaf $NewModule

$OldEnv = "GOCLITMPL"
$NewEnv = $NewApp.ToUpper()

Write-Host "Renaming module from '$OldModule' to '$NewModule'..." -ForegroundColor Cyan
Write-Host "Renaming application name from '$OldApp' to '$NewApp'..." -ForegroundColor Cyan
Write-Host "Renaming environment prefix from '$OldEnv' to '$NewEnv'..." -ForegroundColor Cyan

# Gather all target files in the project
$files = Get-ChildItem -Recurse -File -Include *.go, go.mod, go.sum, Makefile, *.md, *.yml, *.yaml, *.ps1, *.sh

foreach ($file in $files) {
    # Skip standard hidden directories (e.g. .git)
    if ($file.FullName -like "*\.git\*") {
        continue
    }

    $content = Get-Content $file.FullName -Raw
    $modified = $false

    # 1. Replace module namespace
    if ($content -match [regex]::Escape($OldModule)) {
        $content = $content -replace [regex]::Escape($OldModule), $NewModule
        $modified = $true
    }

    # 2. Replace AppName/binary name
    if ($content -match [regex]::Escape($OldApp)) {
        $content = $content -replace [regex]::Escape($OldApp), $NewApp
        $modified = $true
    }

    # 3. Replace EnvPrefix
    if ($content -match [regex]::Escape($OldEnv)) {
        $content = $content -replace [regex]::Escape($OldEnv), $NewEnv
        $modified = $true
    }

    if ($modified) {
        Set-Content $file.FullName $content
    }
}

# Rename the main entrypoint cmd directory
$cmdPath = Join-Path "cmd" $OldApp
if (Test-Path $cmdPath) {
    $newCmdPath = Join-Path "cmd" $NewApp
    Write-Host "Renaming $cmdPath directory to $newCmdPath..." -ForegroundColor Cyan
    Rename-Item -Path $cmdPath -NewName $NewApp
}

Write-Host "Successfully renamed module to '$NewModule'!" -ForegroundColor Green
Write-Host "Running go mod tidy..." -ForegroundColor Cyan
go mod tidy
