# llm-memory Update Script (Windows PowerShell)
# Automatically detect current version, download, verify and update llm-memory
#
# Usage:
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.ps1 | iex
#   To specify version:
#   & ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.ps1))) -Version v0.0.3

param(
    [string]$Version = "latest"
)

# Set error handling
$ErrorActionPreference = "Stop"

# Print functions
function Write-Info {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Err {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Red
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Err "[x] Unsupported architecture: $arch"
            Write-Err "    Supported architectures: AMD64, ARM64"
            exit 1
        }
    }
}

# Download file with retry
function Download-WithRetry {
    param(
        [string]$Url,
        [string]$OutputPath,
        [int]$MaxAttempts = 3
    )

    for ($attempt = 1; $attempt -le $MaxAttempts; $attempt++) {
        try {
            Invoke-WebRequest -Uri $Url -OutFile $OutputPath -UseBasicParsing
            return $true
        } catch {
            if ($attempt -lt $MaxAttempts) {
                Write-Warn "[!] Download failed, attempt $attempt/$MaxAttempts, retrying in 3 seconds..."
                Start-Sleep -Seconds 3
            }
        }
    }

    Write-Err "[x] Download failed after $MaxAttempts attempts"
    Write-Err "    URL: $Url"
    return $false
}

# Get latest version
function Get-LatestVersion {
    Write-Info "[*] Fetching latest version..."

    try {
        $releaseUrl = "https://api.github.com/repos/XiaoLFeng/llm-memory/releases/latest"
        $release = Invoke-RestMethod -Uri $releaseUrl -UseBasicParsing
        $version = $release.tag_name -replace '^v', ''

        if ([string]::IsNullOrEmpty($version)) {
            throw "Unable to parse version number"
        }

        return $version
    } catch {
        Write-Err "[x] Unable to fetch latest version"
        Write-Err "    $_"
        Write-Err "    Please check your network connection or specify version manually: update.ps1 -Version v0.0.3"
        exit 1
    }
}

# Get current installed version
function Get-CurrentVersion {
    param([string]$BinaryPath)

    if (Test-Path $BinaryPath) {
        try {
            $output = & $BinaryPath --version 2>&1 | Select-Object -First 1
            if ($output -match '(\d+\.\d+\.\d+)') {
                return $Matches[1]
            }
        } catch {
            # Ignore error
        }
    }
    return $null
}

# Compare version numbers
# Returns: 0=equal, 1=v1>v2, 2=v1<v2
function Compare-Versions {
    param(
        [string]$v1,
        [string]$v2
    )

    if ($v1 -eq $v2) {
        return 0
    }

    $v1Parts = $v1.Split('.')
    $v2Parts = $v2.Split('.')

    $maxLen = [Math]::Max($v1Parts.Length, $v2Parts.Length)

    for ($i = 0; $i -lt $maxLen; $i++) {
        $n1 = if ($i -lt $v1Parts.Length) { [int]$v1Parts[$i] } else { 0 }
        $n2 = if ($i -lt $v2Parts.Length) { [int]$v2Parts[$i] } else { 0 }

        if ($n1 -gt $n2) {
            return 1
        } elseif ($n1 -lt $n2) {
            return 2
        }
    }

    return 0
}

# Main function
function Main {
    Write-Info "[*] llm-memory Update Script"
    Write-Info ""

    # Detect architecture
    $Arch = Get-Architecture
    Write-Success "[+] Detected system: windows-$Arch"

    # Define installation path
    $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    $BinaryPath = Join-Path $InstallDir "llm-memory.exe"

    # Check if installed
    if (-not (Test-Path $BinaryPath)) {
        # Try to find in PATH
        $foundPath = Get-Command "llm-memory" -ErrorAction SilentlyContinue
        if ($foundPath) {
            $BinaryPath = $foundPath.Source
            $InstallDir = Split-Path $BinaryPath -Parent
        } else {
            Write-Warn "[!] llm-memory installation not found"
            Write-Info "    Please use install.ps1 to install first"
            Write-Info ""
            Write-Info "    Install command:"
            Write-Info "    iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1 | iex"
            exit 1
        }
    }

    # Get current version
    $CurrentVersion = Get-CurrentVersion -BinaryPath $BinaryPath
    if ($CurrentVersion) {
        Write-Success "[+] Current version: v$CurrentVersion"
    } else {
        Write-Warn "[!] Unable to get current version"
        $CurrentVersion = "0.0.0"
    }

    # Get target version
    $TargetVersion = $Version
    if ($TargetVersion -eq "latest") {
        $TargetVersion = Get-LatestVersion
    } else {
        # Remove possible v prefix
        $TargetVersion = $TargetVersion -replace '^v', ''
    }
    Write-Success "[+] Latest version: v$TargetVersion"
    Write-Info ""

    # Compare versions
    $VersionCmp = Compare-Versions -v1 $CurrentVersion -v2 $TargetVersion

    if ($VersionCmp -eq 0) {
        Write-Success "[*] Already at latest version v$CurrentVersion, no update needed"
        exit 0
    } elseif ($VersionCmp -eq 1) {
        Write-Warn "[!] Current version v$CurrentVersion is newer than target version v$TargetVersion"
        $response = Read-Host "Do you want to downgrade? [y/N]"
        if ($response -notmatch '^[yY]') {
            Write-Info "Update cancelled"
            exit 0
        }
        Write-Info "Proceeding with downgrade..."
    }

    # Build download URL
    $BinaryName = "llm-memory-windows-$Arch.exe"
    $DownloadUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$TargetVersion/$BinaryName"
    $ChecksumUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$TargetVersion/checksums.txt"

    # Create temp directory
    $TmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

    try {
        # Download binary
        Write-Info "[*] Downloading llm-memory v$TargetVersion for windows-$Arch..."
        $TmpBinaryPath = Join-Path $TmpDir $BinaryName

        if (-not (Download-WithRetry -Url $DownloadUrl -OutputPath $TmpBinaryPath)) {
            Write-Err "    Hint: Please verify the version is correct, or download manually from GitHub Releases"
            Write-Err "    https://github.com/XiaoLFeng/llm-memory/releases"
            exit 1
        }
        Write-Success "[+] Download complete"

        # Download and verify checksum
        Write-Info "[*] Verifying file integrity..."
        $ChecksumPath = Join-Path $TmpDir "checksums.txt"

        if (Download-WithRetry -Url $ChecksumUrl -OutputPath $ChecksumPath) {
            try {
                # Get checksum
                $ChecksumContent = Get-Content $ChecksumPath
                $ExpectedLine = $ChecksumContent | Where-Object { $_ -match $BinaryName }

                if ($ExpectedLine) {
                    $ExpectedChecksum = ($ExpectedLine -split '\s+')[0].ToLower()

                    # Calculate actual checksum
                    $ActualChecksum = (Get-FileHash -Path $TmpBinaryPath -Algorithm SHA256).Hash.ToLower()

                    if ($ExpectedChecksum -ne $ActualChecksum) {
                        Write-Err "[x] Checksum verification failed: file may be corrupted or tampered"
                        Write-Err "    Expected: $ExpectedChecksum"
                        Write-Err "    Actual:   $ActualChecksum"
                        exit 1
                    }

                    Write-Success "[+] Checksum verification passed"
                } else {
                    Write-Warn "[!] Corresponding checksum not found, skipping verification"
                }
            } catch {
                Write-Warn "[!] Checksum verification failed: $_"
                Write-Warn "[!] Skipping verification, proceeding with update"
            }
        } else {
            Write-Warn "[!] Unable to download checksum file, skipping verification"
        }

        # Backup old version
        $BackupPath = $null
        if (Test-Path $BinaryPath) {
            $BackupPath = "$BinaryPath.backup"
            Write-Info "[*] Backing up old version to $BackupPath..."
            Copy-Item -Path $BinaryPath -Destination $BackupPath -Force
        }

        # Install new version
        Write-Info "[*] Updating to $InstallDir..."
        Copy-Item -Path $TmpBinaryPath -Destination $BinaryPath -Force

        # Verify installation
        $NewVersion = Get-CurrentVersion -BinaryPath $BinaryPath
        if ($NewVersion -eq $TargetVersion) {
            Write-Success "[+] Update successful!"

            # Remove backup
            if ($BackupPath -and (Test-Path $BackupPath)) {
                Remove-Item -Path $BackupPath -Force -ErrorAction SilentlyContinue
            }
        } else {
            Write-Err "[x] Version verification failed after update"
            Write-Err "    Expected: v$TargetVersion"
            Write-Err "    Actual:   v$NewVersion"

            # Restore backup
            if ($BackupPath -and (Test-Path $BackupPath)) {
                Write-Info "[*] Restoring old version..."
                Move-Item -Path $BackupPath -Destination $BinaryPath -Force
            }
            exit 1
        }

        Write-Info ""
        Write-Success "[*] Update complete: v$CurrentVersion -> v$TargetVersion"
        Write-Info ""
        Write-Info "Usage:"
        Write-Info "  llm-memory --help       # View help"
        Write-Info "  llm-memory tui          # Start TUI interface"
        Write-Info "  llm-memory mcp          # Start MCP server"

    } finally {
        # Cleanup temp files
        if (Test-Path $TmpDir) {
            Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Execute main function
Main
