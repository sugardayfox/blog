#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

go build -o blog_app
exitIfFailed

./blog_app -netlify-build
exitIfFailed

netlifyctl deploy
