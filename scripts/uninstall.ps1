# llm-memory 卸载脚本 (Windows PowerShell)
# 自动清理已安装的二进制文件和配置数据
#
# 使用方法：
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/uninstall.ps1 | iex
#   或下载后执行：
#   .\uninstall.ps1

# 设置错误处理
$ErrorActionPreference = "Stop"

# 打印函数
function Write-Info {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Red
}

# 询问用户确认
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

# 获取文件大小（人类可读格式）
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
        return "未知"
    }
}

# 主函数
function Main {
    Write-Info "[!] llm-memory 卸载脚本"
    Write-Info ""

    # 定义安装位置
    $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    $BinaryPath = Join-Path $InstallDir "llm-memory.exe"
    $ConfigDir = Join-Path $env:USERPROFILE ".llm-memory"

    # 检查是否已安装
    if (-not (Test-Path $BinaryPath)) {
        Write-Warning "[!] 未找到已安装的 llm-memory"
        Write-Info "    预期位置: $BinaryPath"

        # 检查是否在其他位置
        $foundPath = Get-Command llm-memory -ErrorAction SilentlyContinue
        if ($foundPath) {
            Write-Warning "    但在 PATH 中找到: $($foundPath.Source)"
            Write-Info ""

            if (Confirm-Action "是否删除该位置的 llm-memory？" $false) {
                $BinaryPath = $foundPath.Source
            } else {
                Write-Info "取消卸载"
                return
            }
        } else {
            Write-Info ""
            Write-Info "llm-memory 可能已经卸载，或安装在其他位置"
            return
        }
    }

    Write-Info "[*] 找到安装位置："
    Write-Info "    二进制文件: $BinaryPath"

    # 检测版本
    if (Test-Path $BinaryPath) {
        try {
            $version = & $BinaryPath --version 2>$null
            Write-Info "    当前版本: $version"
        } catch {
            Write-Info "    当前版本: 未知"
        }
    }

    # 检查配置目录
    if (Test-Path $ConfigDir) {
        Write-Info "    配置目录: $ConfigDir"
        $configSize = Get-HumanReadableSize -Path $ConfigDir
        Write-Info "    配置大小: $configSize"
    }

    Write-Info ""

    # 询问是否继续
    if (-not (Confirm-Action "确定要卸载 llm-memory 吗？" $false)) {
        Write-Info "取消卸载"
        return
    }

    Write-Info ""

    # 删除二进制文件
    Write-Info "[!] 正在删除二进制文件..."
    try {
        Remove-Item -Path $BinaryPath -Force -ErrorAction Stop
        Write-Success "[+] 已删除: $BinaryPath"
    } catch {
        Write-Error "[x] 删除失败: $BinaryPath"
        Write-Error "    错误信息: $_"
        Write-Error "    你可能需要手动删除该文件"
    }

    # 询问是否删除配置
    if (Test-Path $ConfigDir) {
        Write-Info ""
        Write-Warning "[!] 注意：配置目录包含用户数据（记忆、计划、待办）"

        if (Confirm-Action "是否同时删除配置目录及所有数据？" $false) {
            Write-Info "[!] 正在删除配置目录..."
            try {
                Remove-Item -Path $ConfigDir -Recurse -Force -ErrorAction Stop
                Write-Success "[+] 已删除: $ConfigDir"
            } catch {
                Write-Error "[x] 删除失败: $ConfigDir"
                Write-Error "    错误信息: $_"
                Write-Error "    你可能需要手动删除该目录"
            }
        } else {
            Write-Info "保留配置目录: $ConfigDir"
            Write-Info "如果以后需要删除，请运行："
            Write-Host "  Remove-Item -Path '$ConfigDir' -Recurse -Force" -ForegroundColor Gray
        }
    }

    Write-Info ""
    Write-Success "[+] llm-memory 卸载完成！"

    # 检查是否还在 PATH 中
    $stillInPath = Get-Command llm-memory -ErrorAction SilentlyContinue
    if ($stillInPath) {
        Write-Info ""
        Write-Warning "[!] 注意：llm-memory 仍在 PATH 中"
        Write-Warning "    位置: $($stillInPath.Source)"
        Write-Warning "    这可能是另一个安装位置，请手动处理"
    }

    # 询问是否从 PATH 中移除安装目录
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -like "*$InstallDir*") {
        Write-Info ""
        if (Confirm-Action "是否从 PATH 环境变量中移除 $InstallDir ？" $false) {
            try {
                $NewPath = $UserPath -replace [regex]::Escape(";$InstallDir"), ""
                $NewPath = $NewPath -replace [regex]::Escape("$InstallDir;"), ""
                $NewPath = $NewPath -replace [regex]::Escape($InstallDir), ""

                [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
                Write-Success "[+] 已从 PATH 中移除"
                Write-Info "    需要重启 PowerShell 或终端才能生效"
            } catch {
                Write-Error "[x] 移除失败: $_"
                Write-Info "    你可以手动在 '环境变量' 中移除"
            }
        }
    }

    Write-Info ""
    Write-Info "感谢使用 llm-memory！(^_^)/"
    Write-Info ""
    Write-Info "如果你有任何问题或建议，欢迎反馈："
    Write-Info "  https://github.com/XiaoLFeng/llm-memory/issues"
}

# 执行主函数
Main
