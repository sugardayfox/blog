#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

$sha1 = (git rev-parse HEAD) | Out-String
exitIfFailed
$sha1 = $sha1.Replace([System.Environment]::NewLine,"")

$exe = ".\blog_app.exe"
$plat = $PSVersionTable["Platform"]
if ($plat = "Unix") {
    $exe = "./blog_app"
}
go build -o $exe -ldflags "-X main.sha1ver=$sha1"
exitIfFailed

Start-Process -FilePath $exe -ArgumentList "-addr=localhost:5020" -Wait
Remove-Item $exe
