# llm-memory Install Script (Windows PowerShell)
# Automatically detect architecture, download, verify and install llm-memory
#
# Usage:
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1 | iex
#   To specify version:
#   & ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1))) -Version v0.0.2

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
        "ARM64" {
            Write-Warn "[!] Windows ARM64 architecture may have limited SQLite support (CGO limitation)"
            return "arm64"
        }
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
        Write-Err "    Please check your network connection or specify version manually: install.ps1 -Version v0.0.2"
        exit 1
    }
}

# Main function
function Main {
    Write-Info "[*] llm-memory Install Script"
    Write-Info ""

    # Detect architecture
    $Arch = Get-Architecture
    Write-Success "[+] Detected system: windows-$Arch"

    # Get version
    if ($Version -eq "latest") {
        $Version = Get-LatestVersion
    } else {
        # Remove possible v prefix
        $Version = $Version -replace '^v', ''
    }
    Write-Success "[+] Target version: v$Version"
    Write-Info ""

    # Build download URL
    $BinaryName = "llm-memory-windows-$Arch.exe"
    $DownloadUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$Version/$BinaryName"
    $ChecksumUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$Version/checksums.txt"

    # Create temp directory
    $TmpDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

    try {
        # Download binary
        Write-Info "[*] Downloading llm-memory v$Version for windows-$Arch..."
        $BinaryPath = Join-Path $TmpDir $BinaryName

        if (-not (Download-WithRetry -Url $DownloadUrl -OutputPath $BinaryPath)) {
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
                # Get expected checksum
                $ChecksumContent = Get-Content $ChecksumPath
                $ExpectedLine = $ChecksumContent | Where-Object { $_ -match $BinaryName }

                if ($ExpectedLine) {
                    $ExpectedChecksum = ($ExpectedLine -split '\s+')[0].ToLower()

                    # Calculate actual checksum
                    $ActualChecksum = (Get-FileHash -Path $BinaryPath -Algorithm SHA256).Hash.ToLower()

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
                Write-Warn "[!] Skipping verification, proceeding with installation"
            }
        } else {
            Write-Warn "[!] Unable to download checksum file, skipping verification"
        }

        # Install binary
        $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        Write-Info "[*] Installing to $InstallDir..."
        $DestPath = Join-Path $InstallDir "llm-memory.exe"
        Copy-Item -Path $BinaryPath -Destination $DestPath -Force
        Write-Success "[+] Installation successful!"

        Write-Info ""

        # Check PATH environment
        $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($UserPath -notlike "*$InstallDir*") {
            Write-Warn "[!] Note: $InstallDir is not in PATH"
            Write-Info ""
            Write-Info "You can add it to PATH using one of the following methods:"
            Write-Info ""
            Write-Info "Method 1: Manual addition (Recommended)" -ForegroundColor Cyan
            Write-Info "  1. Right-click 'This PC' -> 'Properties' -> 'Advanced system settings'"
            Write-Info "  2. Click 'Environment Variables'"
            Write-Info "  3. Find 'Path' in 'User variables'"
            Write-Info "  4. Click 'Edit' and add: $InstallDir"
            Write-Info ""
            Write-Info "Method 2: Command line (Requires new PowerShell session)" -ForegroundColor Cyan
            Write-Host "  " -NoNewline
            Write-Host "[Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$InstallDir', 'User')" -ForegroundColor Gray
            Write-Info ""
            Write-Info "Method 3: Temporary (Current session only)" -ForegroundColor Cyan
            Write-Host "  " -NoNewline
            Write-Host "`$env:Path += ';$InstallDir'" -ForegroundColor Gray
            Write-Info ""
        } else {
            Write-Success "[+] Installation complete! You can now run: " -NoNewline
            Write-Host "llm-memory --version" -ForegroundColor Cyan
            Write-Info ""
            Write-Info "Usage:"
            Write-Info "  llm-memory --help       # View help"
            Write-Info "  llm-memory tui          # Start TUI interface"
            Write-Info "  llm-memory mcp          # Start MCP server"
        }

    } finally {
        # Cleanup temp files
        if (Test-Path $TmpDir) {
            Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Execute main function
Main
