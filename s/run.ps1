#!/usr/bin/env pwsh
Set-StrictMode -Version Latest

$sha1 = (git rev-parse HEAD) | Out-String
$sha1 = $sha1.Replace([System.Environment]::NewLine,"")

$exe = ".\blog_app.exe"
$plat = $PSVersionTable["Platform"]
if ($plat = "Unix") {
    $exe = "./blog_app"
}
go build -o $exe -ldflags "-X main.sha1ver=$sha1"

Write-Host "exe: $exe"

Start-Process -FilePath $exe -ArgumentList "-addr=localhost:5020" -Wait
Remove-Item $exe
