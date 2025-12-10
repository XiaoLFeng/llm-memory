# llm-memory 更新脚本 (Windows PowerShell)
# 自动检测当前版本，下载、校验并更新 llm-memory
#
# 使用方法：
#   iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.ps1 | iex
#   或指定版本：
#   & ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.ps1))) -Version v0.0.3

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

function Write-Warn {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Err {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Red
}

# 检测架构
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Err "[x] 不支持的架构: $arch"
            Write-Err "    支持的架构: AMD64, ARM64"
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
                Write-Warn "[!] 下载失败（尝试 $attempt/$MaxAttempts），3 秒后重试..."
                Start-Sleep -Seconds 3
            }
        }
    }

    Write-Err "[x] 下载失败，已重试 $MaxAttempts 次"
    Write-Err "    URL: $Url"
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
        Write-Err "[x] 无法获取最新版本"
        Write-Err "    $_"
        Write-Err "    请检查网络连接或手动指定版本: update.ps1 -Version v0.0.3"
        exit 1
    }
}

# 获取当前安装版本
function Get-CurrentVersion {
    param([string]$BinaryPath)

    if (Test-Path $BinaryPath) {
        try {
            $output = & $BinaryPath --version 2>&1 | Select-Object -First 1
            if ($output -match '(\d+\.\d+\.\d+)') {
                return $Matches[1]
            }
        } catch {
            # 忽略错误
        }
    }
    return $null
}

# 比较版本号
# 返回: 0=相等, 1=v1>v2, 2=v1<v2
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

# 主函数
function Main {
    Write-Info "[*] llm-memory 更新脚本"
    Write-Info ""

    # 检测架构
    $Arch = Get-Architecture
    Write-Success "[+] 检测到系统: windows-$Arch"

    # 定义安装路径
    $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    $BinaryPath = Join-Path $InstallDir "llm-memory.exe"

    # 检查是否已安装
    if (-not (Test-Path $BinaryPath)) {
        # 尝试在 PATH 中查找
        $foundPath = Get-Command "llm-memory" -ErrorAction SilentlyContinue
        if ($foundPath) {
            $BinaryPath = $foundPath.Source
            $InstallDir = Split-Path $BinaryPath -Parent
        } else {
            Write-Warn "[!] 未找到已安装的 llm-memory"
            Write-Info "    请先使用 install.ps1 进行安装"
            Write-Info ""
            Write-Info "    安装命令："
            Write-Info "    iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1 | iex"
            exit 1
        }
    }

    # 获取当前版本
    $CurrentVersion = Get-CurrentVersion -BinaryPath $BinaryPath
    if ($CurrentVersion) {
        Write-Success "[+] 当前版本: v$CurrentVersion"
    } else {
        Write-Warn "[!] 无法获取当前版本"
        $CurrentVersion = "0.0.0"
    }

    # 获取目标版本
    $TargetVersion = $Version
    if ($TargetVersion -eq "latest") {
        $TargetVersion = Get-LatestVersion
    } else {
        # 去掉可能的 v 前缀
        $TargetVersion = $TargetVersion -replace '^v', ''
    }
    Write-Success "[+] 最新版本: v$TargetVersion"
    Write-Info ""

    # 比较版本
    $VersionCmp = Compare-Versions -v1 $CurrentVersion -v2 $TargetVersion

    if ($VersionCmp -eq 0) {
        Write-Success "[*] 已经是最新版本 v$CurrentVersion，无需更新"
        exit 0
    } elseif ($VersionCmp -eq 1) {
        Write-Warn "[!] 当前版本 v$CurrentVersion 比目标版本 v$TargetVersion 更新"
        $response = Read-Host "是否要降级？[y/N]"
        if ($response -notmatch '^[yY]') {
            Write-Info "取消更新"
            exit 0
        }
        Write-Info "继续降级..."
    }

    # 设置下载 URL
    $BinaryName = "llm-memory-windows-$Arch.exe"
    $DownloadUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$TargetVersion/$BinaryName"
    $ChecksumUrl = "https://github.com/XiaoLFeng/llm-memory/releases/download/v$TargetVersion/checksums.txt"

    # 创建临时目录
    $TmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

    try {
        # 下载二进制
        Write-Info "[*] 正在下载 llm-memory v$TargetVersion for windows-$Arch..."
        $TmpBinaryPath = Join-Path $TmpDir $BinaryName

        if (-not (Download-WithRetry -Url $DownloadUrl -OutputPath $TmpBinaryPath)) {
            Write-Err "    提示：请检查版本号是否正确，或访问 GitHub Release 页面手动下载"
            Write-Err "    https://github.com/XiaoLFeng/llm-memory/releases"
            exit 1
        }
        Write-Success "[+] 下载完成"

        # 下载并验证校验和
        Write-Info "[*] 验证文件完整性..."
        $ChecksumPath = Join-Path $TmpDir "checksums.txt"

        if (Download-WithRetry -Url $ChecksumUrl -OutputPath $ChecksumPath) {
            try {
                # 读取校验和
                $ChecksumContent = Get-Content $ChecksumPath
                $ExpectedLine = $ChecksumContent | Where-Object { $_ -match $BinaryName }

                if ($ExpectedLine) {
                    $ExpectedChecksum = ($ExpectedLine -split '\s+')[0].ToLower()

                    # 计算实际校验和
                    $ActualChecksum = (Get-FileHash -Path $TmpBinaryPath -Algorithm SHA256).Hash.ToLower()

                    if ($ExpectedChecksum -ne $ActualChecksum) {
                        Write-Err "[x] 文件校验失败！文件可能已损坏或被篡改"
                        Write-Err "    期望: $ExpectedChecksum"
                        Write-Err "    实际: $ActualChecksum"
                        exit 1
                    }

                    Write-Success "[+] 文件校验通过"
                } else {
                    Write-Warn "[!] 未找到对应的校验和，跳过校验"
                }
            } catch {
                Write-Warn "[!] 校验和验证失败: $_"
                Write-Warn "[!] 跳过校验，继续更新"
            }
        } else {
            Write-Warn "[!] 无法下载校验和文件，跳过校验"
        }

        # 备份旧版本
        $BackupPath = $null
        if (Test-Path $BinaryPath) {
            $BackupPath = "$BinaryPath.backup"
            Write-Info "[*] 备份旧版本到 $BackupPath..."
            Copy-Item -Path $BinaryPath -Destination $BackupPath -Force
        }

        # 安装新版本
        Write-Info "[*] 正在更新到 $InstallDir..."
        Copy-Item -Path $TmpBinaryPath -Destination $BinaryPath -Force

        # 验证安装
        $NewVersion = Get-CurrentVersion -BinaryPath $BinaryPath
        if ($NewVersion -eq $TargetVersion) {
            Write-Success "[+] 更新成功！"

            # 删除备份
            if ($BackupPath -and (Test-Path $BackupPath)) {
                Remove-Item -Path $BackupPath -Force -ErrorAction SilentlyContinue
            }
        } else {
            Write-Err "[x] 更新后版本验证失败"
            Write-Err "    期望: v$TargetVersion"
            Write-Err "    实际: v$NewVersion"

            # 恢复备份
            if ($BackupPath -and (Test-Path $BackupPath)) {
                Write-Info "[*] 正在恢复旧版本..."
                Move-Item -Path $BackupPath -Destination $BinaryPath -Force
            }
            exit 1
        }

        Write-Info ""
        Write-Success "[*] 更新完成！v$CurrentVersion -> v$TargetVersion"
        Write-Info ""
        Write-Info "使用帮助："
        Write-Info "  llm-memory --help       # 查看帮助"
        Write-Info "  llm-memory tui          # 启动 TUI 界面"
        Write-Info "  llm-memory mcp          # 启动 MCP 服务"

    } finally {
        # 清理临时文件
        if (Test-Path $TmpDir) {
            Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# 执行主函数
Main
