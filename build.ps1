<#
.SYNOPSIS
    Build and validation script for Windows PowerShell.
.DESCRIPTION
    Mimics Makefile targets to provide a native build and test utility on Windows.
.PARAMETER Task
    The task to run. Options: all, build, test, lint, vulncheck, fmt, clean. Defaults to 'all'.
.EXAMPLE
    .\build.ps1 -Task build
#>
param (
    [ValidateSet("all", "build", "test", "lint", "vulncheck", "fmt", "clean")]
    [string]$Task = "all"
)

$BinaryName = "goclitmpl"
$MainPath = ".\cmd\goclitmpl"
$BuildDir = "bin"

# Determine Version, Commit, and Date
$Version = "dev"
$Commit = "none"
$Date = [DateTime]::UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ")

if (Get-Command git -ErrorAction SilentlyContinue) {
    $gitVersion = git describe --tags --always --dirty 2>$null
    if ($gitVersion) { $Version = $gitVersion }
    
    $gitCommit = git rev-parse --short HEAD 2>$null
    if ($gitCommit) { $Commit = $gitCommit }
}

$LDFlags = "-s -w " + `
            "-X github.com/dat267/goclitmpl/internal/cli.Version=$Version " + `
            "-X github.com/dat267/goclitmpl/internal/cli.Commit=$Commit " + `
            "-X github.com/dat267/goclitmpl/internal/cli.Date=$Date"

function Run-Fmt {
    Write-Host "==> Formatting code..." -ForegroundColor Cyan
    go fmt ./...
}

function Run-Lint {
    Write-Host "==> Running golangci-lint..." -ForegroundColor Cyan
    if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
        golangci-lint run
    } else {
        Write-Warning "golangci-lint is not installed. Skipping. Install via: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    }
}

function Run-Vulncheck {
    Write-Host "==> Running govulncheck..." -ForegroundColor Cyan
    if (Get-Command govulncheck -ErrorAction SilentlyContinue) {
        govulncheck ./...
    } else {
        Write-Warning "govulncheck is not installed. Skipping. Install via: go install golang.org/x/vuln/cmd/govulncheck@latest"
    }
}

function Run-Test {
    Write-Host "==> Running unit tests..." -ForegroundColor Cyan
    go test -v -race ./...
}

function Run-Build {
    Write-Host "==> Building binary to $BuildDir\$BinaryName.exe..." -ForegroundColor Cyan
    if (-not (Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }
    go build -ldflags $LDFlags -o "$BuildDir\$BinaryName.exe" $MainPath
}

function Run-Clean {
    Write-Host "==> Cleaning build artifacts..." -ForegroundColor Cyan
    if (Test-Path $BuildDir) {
        Remove-Item -Recurse -Force $BuildDir
    }
    if (Test-Path "coverage.out") { Remove-Item "coverage.out" }
    if (Test-Path "coverage.html") { Remove-Item "coverage.html" }
}

# Run the selected task
switch ($Task) {
    "clean" {
        Run-Clean
    }
    "fmt" {
        Run-Fmt
    }
    "lint" {
        Run-Lint
    }
    "vulncheck" {
        Run-Vulncheck
    }
    "test" {
        Run-Test
    }
    "build" {
        Run-Build
    }
    "all" {
        Run-Fmt
        Run-Lint
        Run-Vulncheck
        Run-Test
        Run-Build
    }
}
