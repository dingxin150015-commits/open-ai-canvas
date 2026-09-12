[CmdletBinding()]
param()
$ErrorActionPreference='Stop'
$repositoryRoot=Split-Path -Parent $PSScriptRoot
Push-Location $repositoryRoot
try {
    & node 'plugin-packages/generate-catalog.mjs' 'dashscope-image' 'dashscope-video'
    if ($LASTEXITCODE -ne 0) { throw 'Plugin generation failed' }
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    foreach ($id in @('dashscope-image','dashscope-video')) {
        $source=Join-Path $repositoryRoot "plugin-packages/$id"
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
