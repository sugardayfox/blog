#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if (-Not $?) { exit } }

$sha1 = (git rev-parse HEAD) | Out-String
$sha1 = $sha1.Replace([System.Environment]::NewLine,"")

$exe = ".\blog_app.exe"
$plat = $PSVersionTable["Platform"]
if ($plat = "Unix") {
    $exe = "./blog_app"
}
go build -o $exe -ldflags "-X main.sha1ver=$sha1"
exitIfFailed

Remove-Item $exe
