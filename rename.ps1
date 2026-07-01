param (
    [Parameter(Mandatory=$true)]
    [string]$NewModule
)

$OldModule = "github.com/dat267/goclitmpl"

Write-Host "Renaming module from '$OldModule' to '$NewModule'..." -ForegroundColor Cyan

# Gather all target files in the project
$files = Get-ChildItem -Recurse -File -Include *.go, go.mod, go.sum, Makefile, *.md

foreach ($file in $files) {
    # Skip standard hidden directories (e.g. .git)
    if ($file.FullName -like "*\.git\*") {
        continue
    }

    $content = Get-Content $file.FullName -Raw
    if ($content -match [regex]::Escape($OldModule)) {
        $content = $content -replace [regex]::Escape($OldModule), $NewModule
        # Set-Content without -NoNewline handles raw string writing properly
        Set-Content $file.FullName $content
    }
}

Write-Host "Successfully renamed module to '$NewModule'!" -ForegroundColor Green
Write-Host "Running go mod tidy..." -ForegroundColor Cyan
go mod tidy
