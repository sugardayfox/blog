#!/usr/bin/env pwsh
go build -o blog_app
./blog_app -netlify-build

$origDir = Get-Location
Set-Location -Path netlify_static
Write-Host "About to deploy"
# netlifyctl deploy
Set-Location -Path $origDir
