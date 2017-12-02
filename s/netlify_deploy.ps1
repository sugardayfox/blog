$origDir = Get-Location

go run .\s\netlify_build.go
Set-Location -Path .\netlify_static
Write-Host "About to deploy"
# netlifyctl deploy

Set-Location -Path $origDir
