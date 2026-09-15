$ErrorActionPreference = 'Stop'
$Repo = "ju4n97/esquema"
$InstallDir = "$env:LOCALAPPDATA\Programs\esquema"

$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name
$TagNoV = $Tag.TrimStart('v')
$FileName = "esquema_${TagNoV}_windows_${Arch}.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Tag/$FileName"

Write-Host "Downloading esquema $Tag for windows/$Arch..."
$ZipPath = "$env:TEMP\$FileName"
Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Expand-Archive -Path $ZipPath -DestinationPath $InstallDir -Force
Remove-Item -Path $ZipPath

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
}

Write-Host "esquema was installed successfully to $InstallDir\esquema.exe"