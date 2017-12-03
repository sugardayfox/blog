#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

go build -o blog_app
exitIfFailed

./blog_app -netlify-build
exitIfFailed

$origDir = Get-Location
Set-Location -Path netlify_static
Write-Host "About to deploy"
netlifyctl deploy
Set-Location -Path $origDir
