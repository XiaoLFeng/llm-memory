# llm-memory Uninstall Script (Windows PowerShell)
# Automatically find and remove installed binary files and configuration
#
# Usage:
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/uninstall.ps1 | iex
#   Or download and execute:
#   .\uninstall.ps1

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

# Ask user for confirmation
function Confirm-Action {
    param(
        [string]$Prompt,
        [bool]$DefaultYes = $false
    )

    if ($DefaultYes) {
        $choices = "&Yes", "&No"
        $defaultChoice = 0
    } else {
        $choices = "&Yes", "&No"
        $defaultChoice = 1
    }

    $decision = $Host.UI.PromptForChoice("", $Prompt, $choices, $defaultChoice)
    return $decision -eq 0
}

# Get file size in human-readable format
function Get-HumanReadableSize {
    param([string]$Path)

    try {
        $size = (Get-ChildItem -Path $Path -Recurse -ErrorAction SilentlyContinue |
                 Measure-Object -Property Length -Sum).Sum

        if ($size -gt 1GB) {
            return "{0:N2} GB" -f ($size / 1GB)
        } elseif ($size -gt 1MB) {
            return "{0:N2} MB" -f ($size / 1MB)
        } elseif ($size -gt 1KB) {
            return "{0:N2} KB" -f ($size / 1KB)
        } else {
            return "$size B"
        }
    } catch {
        return "Unknown"
    }
}

# Main function
function Main {
    Write-Info "[!] llm-memory Uninstall Script"
    Write-Info ""

    # Define installation locations
    $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    $BinaryPath = Join-Path $InstallDir "llm-memory.exe"
    $ConfigDir = Join-Path $env:USERPROFILE ".llm-memory"

    # Check if installed
    if (-not (Test-Path $BinaryPath)) {
        Write-Warn "[!] llm-memory installation not found"
        Write-Info "    Expected location: $BinaryPath"

        # Check if exists in other location
        $foundPath = Get-Command llm-memory -ErrorAction SilentlyContinue
        if ($foundPath) {
            Write-Warn "    Found in PATH: $($foundPath.Source)"
            Write-Info ""

            if (Confirm-Action "Do you want to remove llm-memory from this location?" $false) {
                $BinaryPath = $foundPath.Source
            } else {
                Write-Info "Uninstall cancelled"
                return
            }
        } else {
            Write-Info ""
            Write-Info "llm-memory may already be uninstalled or installed in a different location"
            return
        }
    }

    Write-Info "[*] Found installation:"
    Write-Info "    Binary file: $BinaryPath"

    # Get version
    if (Test-Path $BinaryPath) {
        try {
            $version = & $BinaryPath --version 2>$null
            Write-Info "    Current version: $version"
        } catch {
            Write-Info "    Current version: Unknown"
        }
    }

    # Check config directory
    if (Test-Path $ConfigDir) {
        Write-Info "    Config directory: $ConfigDir"
        $configSize = Get-HumanReadableSize -Path $ConfigDir
        Write-Info "    Config size: $configSize"
    }

    Write-Info ""

    # Ask for confirmation
    if (-not (Confirm-Action "Are you sure you want to uninstall llm-memory?" $false)) {
        Write-Info "Uninstall cancelled"
        return
    }

    Write-Info ""

    # Remove binary file
    Write-Info "[!] Removing binary file..."
    try {
        Remove-Item -Path $BinaryPath -Force -ErrorAction Stop
        Write-Success "[+] Removed: $BinaryPath"
    } catch {
        Write-Err "[x] Failed to remove: $BinaryPath"
        Write-Err "    Error: $_"
        Write-Err "    You may need to manually delete this file"
    }

    # Ask whether to remove config
    if (Test-Path $ConfigDir) {
        Write-Info ""
        Write-Warn "[!] Note: Config directory contains user data (memories, plans, todos)"

        if (Confirm-Action "Do you want to also remove the config directory and all data?" $false) {
            Write-Info "[!] Removing config directory..."
            try {
                Remove-Item -Path $ConfigDir -Recurse -Force -ErrorAction Stop
                Write-Success "[+] Removed: $ConfigDir"
            } catch {
                Write-Err "[x] Failed to remove: $ConfigDir"
                Write-Err "    Error: $_"
                Write-Err "    You may need to manually delete this directory"
            }
        } else {
            Write-Info "Keeping config directory: $ConfigDir"
            Write-Info "To remove it later, run:"
            Write-Host "  Remove-Item -Path '$ConfigDir' -Recurse -Force" -ForegroundColor Gray
        }
    }

    Write-Info ""
    Write-Success "[+] llm-memory uninstall complete!"

    # Check if still in PATH
    $stillInPath = Get-Command llm-memory -ErrorAction SilentlyContinue
    if ($stillInPath) {
        Write-Info ""
        Write-Warn "[!] Note: llm-memory still found in PATH"
        Write-Warn "    Location: $($stillInPath.Source)"
        Write-Warn "    This may be another installation, please handle manually"
    }

    # Ask whether to remove from PATH
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -like "*$InstallDir*") {
        Write-Info ""
        if (Confirm-Action "Do you want to remove $InstallDir from PATH environment variable?" $false) {
            try {
                $NewPath = $UserPath -replace [regex]::Escape(";$InstallDir"), ""
                $NewPath = $NewPath -replace [regex]::Escape("$InstallDir;"), ""
                $NewPath = $NewPath -replace [regex]::Escape($InstallDir), ""

                [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
                Write-Success "[+] Removed from PATH"
                Write-Info "    Requires new PowerShell session to take effect"
            } catch {
                Write-Err "[x] Failed to remove from PATH: $_"
                Write-Info "    You may need to manually remove it from 'Environment Variables'"
            }
        }
    }

    Write-Info ""
    Write-Info "Thank you for using llm-memory! (^_^)/"
    Write-Info ""
    Write-Info "If you have any questions or suggestions, please visit:"
    Write-Info "  https://github.com/XiaoLFeng/llm-memory/issues"
}

# Execute main function
Main
