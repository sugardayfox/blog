#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

git pull
exitIfFailed

$exe = ".\blog_app.exe"
$plat = $PSVersionTable["Platform"]
if ($plat = "Unix") {
    $exe = "./blog_app"
}
go build -o $exe
exitIfFailed

Start-Process -Wait -FilePath $exe -ArgumentList "-update-notes"
exitIfFailed

git commit articles/notes.txt -m "update notes"
exitIfFailed

git push
exitIfFailed

Remove-Item $exe
