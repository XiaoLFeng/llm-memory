# llm-memory 安装脚本 (Windows PowerShell)
# 自动检测架构、下载、校验并安装 llm-memory
#
# 使用方法：
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1 | iex
#   或指定版本：
#   & ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1))) -Version v0.0.2

param(
    [string]$Version = "latest"
)

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

# 检测架构
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" {
            Write-Warning "[!] Windows ARM64 架构可能不支持完整 SQLite 功能（CGO 限制）"
            return "arm64"
        }
        default {
            Write-Error "[x] 不支持的架构: $arch"
            Write-Error "    支持的架构: AMD64, ARM64"
            exit 1
        }
    }
}

# 下载文件（带重试）
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
                Write-Warning "[!] 下载失败，尝试 $attempt/$MaxAttempts，等3 秒后重试..."
                Start-Sleep -Seconds 3
            }
        }
    }

    Write-Error "[x] 下载失败，已重试 $MaxAttempts 次"
    Write-Error "    URL: $Url"
    return $false
}

# 获取最新版本
function Get-LatestVersion {
    Write-Info "[*] 正在获取最新版本..."

    try {
        $releaseUrl = "https://api.github.com/repos/XiaoLFeng/llm-memory/releases/latest"
        $release = Invoke-RestMethod -Uri $releaseUrl -UseBasicParsing
        $version = $release.tag_name -replace '^v', ''

        if ([string]::IsNullOrEmpty($version)) {
            throw "无法解析版本号"
        }

        return $version
    } catch {
        Write-Error "[x] 无法获取最新版本"
        Write-Error "    $_"
        Write-Error "    请检查网络连接或手动指定版本: install.ps1 -Version v0.0.2"
        exit 1
    }
}

# 主函数
function Main {
    Write-Info "[*] llm-memory 安装脚本"
    Write-Info ""

    # 检测架构
    $Arch = Get-Architecture
    Write-Success "[+] 检测到系统: windows-$Arch"

    # 获取版本
    if ($Version -eq "latest") {
        $Version = Get-LatestVersion
    } else {
        # 去掉可能的 v 前缀
        $Version = $Version -replace '^v', ''
    }
    Write-Success "[+] 目标版本: v$Version"
    Write-Info ""

    # 构建下载 URL
    $BinaryName = "llm-memory-windows-$Arch.exe"
    $DownloadUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$Version/$BinaryName"
    $ChecksumUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$Version/checksums.txt"

    # 创建临时目录
    $TmpDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

    try {
        # 下载二进制
        Write-Info "[*] 正在下载 llm-memory v$Version for windows-$Arch..."
        $BinaryPath = Join-Path $TmpDir $BinaryName

        if (-not (Download-WithRetry -Url $DownloadUrl -OutputPath $BinaryPath)) {
            Write-Error "    提示：请检查版本号是否正确，或访问 GitHub Release 页面手动下载"
            Write-Error "    https://github.com/XiaoLFeng/llm-memory/releases"
            exit 1
        }
        Write-Success "[+] 下载完成"

        # 下载并验证校验和
        Write-Info "[*] 验证文件完整性..."
        $ChecksumPath = Join-Path $TmpDir "checksums.txt"

        if (Download-WithRetry -Url $ChecksumUrl -OutputPath $ChecksumPath) {
            try {
                # 读取期望校验和
                $ChecksumContent = Get-Content $ChecksumPath
                $ExpectedLine = $ChecksumContent | Where-Object { $_ -match $BinaryName }

                if ($ExpectedLine) {
                    $ExpectedChecksum = ($ExpectedLine -split '\s+')[0].ToLower()

                    # 计算实际校验和
                    $ActualChecksum = (Get-FileHash -Path $BinaryPath -Algorithm SHA256).Hash.ToLower()

                    if ($ExpectedChecksum -ne $ActualChecksum) {
                        Write-Error "[x] 文件校验失败，文件可能已损坏或被篡改"
                        Write-Error "    期望: $ExpectedChecksum"
                        Write-Error "    实际: $ActualChecksum"
                        exit 1
                    }

                    Write-Success "[+] 文件校验通过"
                } else {
                    Write-Warning "[!] 未找到对应的校验和，跳过校验"
                }
            } catch {
                Write-Warning "[!] 校验和验证失败: $_"
                Write-Warning "[!] 跳过校验，继续安装"
            }
        } else {
            Write-Warning "[!] 无法下载校验和文件，跳过校验"
        }

        # 安装二进制
        $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        Write-Info "[*] 正在安装到 $InstallDir..."
        $DestPath = Join-Path $InstallDir "llm-memory.exe"
        Copy-Item -Path $BinaryPath -Destination $DestPath -Force
        Write-Success "[+] 安装成功！"

        Write-Info ""

        # 检查 PATH 环境
        $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($UserPath -notlike "*$InstallDir*") {
            Write-Warning "[!] 注意：$InstallDir 不在 PATH 中"
            Write-Info ""
            Write-Info "你可以通过以下方式添加到 PATH："
            Write-Info ""
            Write-Info "方式 1：手动添加（推荐）" -ForegroundColor Cyan
            Write-Info "  1. 右键 '此电脑' -> '属性' -> '高级系统设置'"
            Write-Info "  2. 点击 '环境变量'"
            Write-Info "  3. 在 '用户变量' 中找到 'Path'"
            Write-Info "  4. 点击 '编辑'，添加: $InstallDir"
            Write-Info ""
            Write-Info "方式 2：命令行添加（需要重启 PowerShell 生效）" -ForegroundColor Cyan
            Write-Host "  " -NoNewline
            Write-Host "[Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$InstallDir', 'User')" -ForegroundColor Gray
            Write-Info ""
            Write-Info "方式 3：临时添加（仅当前会话有效）" -ForegroundColor Cyan
            Write-Host "  " -NoNewline
            Write-Host "`$env:Path += ';$InstallDir'" -ForegroundColor Gray
            Write-Info ""
        } else {
            Write-Success "[+] 安装完成！现在可以运行: " -NoNewline
            Write-Host "llm-memory --version" -ForegroundColor Cyan
            Write-Info ""
            Write-Info "使用帮助："
            Write-Info "  llm-memory --help       # 查看帮助"
            Write-Info "  llm-memory tui          # 启动 TUI 界面"
            Write-Info "  llm-memory mcp          # 启动 MCP 服务"
        }

    } finally {
        # 清理临时文件
        if (Test-Path $TmpDir) {
            Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# 执行主函数
Main
