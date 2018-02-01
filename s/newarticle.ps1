#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

if ($args.Length -ne 1) {
    Write-Host "usage: ./s/newarticle.ps1 <title>"
    exit 1
}

$exe = ".\blog_app.exe"
$plat = $PSVersionTable["Platform"]
if ($plat = "Unix") {
    $exe = "./blog_app"
}
go build -o $exe
exitIfFailed

$title = $args[0]
$arglist = '-newarticle "' + $title + '"'
Write-Host "arglist : $arglist"
Start-Process -Wait -FilePath $exe -ArgumentList $arglist

Remove-Item -Path $exe
