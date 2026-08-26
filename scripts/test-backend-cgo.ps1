[CmdletBinding()]
param(
    [switch]$NoCache
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repositoryRoot = Split-Path -Parent $PSScriptRoot
$arguments = @(
    "build",
    "--file", (Join-Path $repositoryRoot "backend\Dockerfile"),
    "--target", "backend-test",
    "--tag", "open-ai-canvas-backend-test:local",
    "--progress", "plain"
)
if ($NoCache) {
    $arguments += "--no-cache"
}
$arguments += $repositoryRoot

Write-Host "Running isolated Linux CGO backend tests. No project database or Docker volume is mounted."
& docker @arguments
if ($LASTEXITCODE -ne 0) {
    throw "Backend CGO test image build failed with exit code $LASTEXITCODE"
}
