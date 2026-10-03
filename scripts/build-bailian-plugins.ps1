[CmdletBinding()]
param([string[]]$PackageIDs = @('dashscope-image','dashscope-video'))
$ErrorActionPreference='Stop'
$repositoryRoot=Split-Path -Parent $PSScriptRoot
Push-Location $repositoryRoot
try {
    foreach ($id in $PackageIDs) { if ($id -notmatch '^[a-z0-9]+(?:-[a-z0-9]+)*$') { throw 'Invalid package ID' } }
    $v4PackageIDs=@('dashscope-image','dashscope-video')
    $legacyPackageIDs=@($PackageIDs | Where-Object { $_ -notin $v4PackageIDs })
    if ($legacyPackageIDs.Count -gt 0) {
        & node 'plugin-packages/generate-catalog.mjs' @legacyPackageIDs
        if ($LASTEXITCODE -ne 0) { throw 'Plugin generation failed' }
    }
    Add-Type -AssemblyName System.IO.Compression
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    foreach ($id in $PackageIDs) {
        $source=Join-Path $repositoryRoot "plugin-packages/$id"
        if ($id -in $v4PackageIDs) {
            $manifest=Get-Content -LiteralPath (Join-Path $source 'manifest.json') -Raw -Encoding UTF8 | ConvertFrom-Json
            if ($manifest.apiVersion -ne 'yingce.plugin/v4') { throw "Expected v4 manifest: $id" }
        }
        $temporary=Join-Path $repositoryRoot "plugin-packages/$id-$([Guid]::NewGuid().ToString('N')).zip"
        $archive=[System.IO.Compression.ZipFile]::Open($temporary,[System.IO.Compression.ZipArchiveMode]::Create)
        try {
            foreach($relative in @('manifest.json','README.md','docs/interface.md')) {
                [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile($archive,(Join-Path $source $relative),$relative) | Out-Null
            }
        } finally { $archive.Dispose() }
        Move-Item -LiteralPath $temporary -Destination (Join-Path $repositoryRoot "plugin-packages/$id.yingce-plugin") -Force
    }
} finally { Pop-Location }
