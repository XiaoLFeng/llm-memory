#!/bin/bash
# llm-memory 更新脚本
# 自动检测当前版本，下载、校验并更新 llm-memory
#
# 使用方法：
#   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.sh | bash
#   或指定版本：
#   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/update.sh | bash -s v0.0.3

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${CYAN}$1${NC}"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}" >&2
}

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        FreeBSD*)   echo "freebsd" ;;
        *)
            print_error "❌ 不支持的操作系统: $(uname -s)"
            print_error "   支持的系统: Linux, macOS (Darwin), FreeBSD"
            exit 1
            ;;
    esac
}

# 检测架构
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)
            print_error "❌ 不支持的架构: $(uname -m)"
            print_error "   支持的架构: x86_64 (amd64), aarch64 (arm64)"
            exit 1
            ;;
    esac
}

# 检查依赖命令
check_dependencies() {
    local missing_deps=()

    for cmd in curl sha256sum; do
        if ! command -v "$cmd" &> /dev/null; then
            missing_deps+=("$cmd")
        fi
    done

    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "❌ 缺少必要的命令: ${missing_deps[*]}"
        print_error "   请先安装这些工具"
        exit 1
    fi
}

# 下载文件（带重试）
download_with_retry() {
    local url="$1"
    local output="$2"
    local max_attempts=3
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -fsSL "$url" -o "$output"; then
            return 0
        fi

        print_warning "⚠️  下载失败（尝试 $attempt/$max_attempts），3 秒后重试..."
        sleep 3
        attempt=$((attempt + 1))
    done

    print_error "❌ 下载失败，已重试 $max_attempts 次"
    print_error "   URL: $url"
    return 1
}

# 获取最新版本
get_latest_version() {
    local release_url="https://api.github.com/repos/XiaoLFeng/llm-memory/releases/latest"
    local version

    print_info "🔍 正在获取最新版本..."

    version=$(curl -fsSL "$release_url" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        print_error "❌ 无法获取最新版本"
        print_error "   请检查网络连接或手动指定版本: bash update.sh v0.0.3"
        exit 1
    fi

    echo "$version"
}

# 获取当前安装版本
get_current_version() {
    local install_path="$1"

    if [ -x "$install_path" ]; then
        local version=$("$install_path" --version 2>/dev/null | head -n1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "")
        echo "$version"
    else
        echo ""
    fi
}

# 比较版本号
# 返回: 0=相等, 1=v1>v2, 2=v1<v2
compare_versions() {
    local v1="$1"
    local v2="$2"

    if [ "$v1" = "$v2" ]; then
        echo 0
        return
    fi

    local IFS=.
    local i v1_parts=($v1) v2_parts=($v2)

    for ((i=0; i<${#v1_parts[@]} || i<${#v2_parts[@]}; i++)); do
        local n1=${v1_parts[i]:-0}
        local n2=${v2_parts[i]:-0}

        if [ $n1 -gt $n2 ]; then
            echo 1
            return
        elif [ $n1 -lt $n2 ]; then
            echo 2
            return
        fi
    done

    echo 0
}

# 主函数
main() {
    print_info "🔄 llm-memory 更新脚本"
    print_info ""

    # 检查依赖
    check_dependencies

    # 检测系统
    OS=$(detect_os)
    ARCH=$(detect_arch)
    print_success "✅ 检测到系统: $OS-$ARCH"

    # 定义安装路径
    INSTALL_DIR="$HOME/.local/bin"
    BINARY_PATH="$INSTALL_DIR/llm-memory"

    # 检查是否已安装
    if [ ! -f "$BINARY_PATH" ]; then
        # 尝试在 PATH 中查找
        if command -v llm-memory &> /dev/null; then
            BINARY_PATH=$(which llm-memory)
            INSTALL_DIR=$(dirname "$BINARY_PATH")
        else
            print_warning "⚠️  未找到已安装的 llm-memory"
            print_info "   请先使用 install.sh 进行安装"
            print_info ""
            print_info "   安装命令："
            print_info "   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.sh | bash"
            exit 1
        fi
    fi

    # 获取当前版本
    CURRENT_VERSION=$(get_current_version "$BINARY_PATH")
    if [ -n "$CURRENT_VERSION" ]; then
        print_success "✅ 当前版本: v$CURRENT_VERSION"
    else
        print_warning "⚠️  无法获取当前版本"
        CURRENT_VERSION="0.0.0"
    fi

    # 获取目标版本
    TARGET_VERSION="${1:-latest}"
    if [ "$TARGET_VERSION" = "latest" ]; then
        TARGET_VERSION=$(get_latest_version)
    else
        # 去掉可能的 v 前缀
        TARGET_VERSION="${TARGET_VERSION#v}"
    fi
    print_success "✅ 最新版本: v$TARGET_VERSION"
    print_info ""

    # 比较版本
    VERSION_CMP=$(compare_versions "$CURRENT_VERSION" "$TARGET_VERSION")

    if [ "$VERSION_CMP" = "0" ]; then
        print_success "🎉 已经是最新版本 v$CURRENT_VERSION，无需更新"
        exit 0
    elif [ "$VERSION_CMP" = "1" ]; then
        print_warning "⚠️  当前版本 v$CURRENT_VERSION 比目标版本 v$TARGET_VERSION 更新"
        read -p "$(echo -e ${CYAN}是否要降级？[y/N] ${NC})" -r response
        case "$response" in
            [yY][eE][sS]|[yY])
                print_info "继续降级..."
                ;;
            *)
                print_info "取消更新"
                exit 0
                ;;
        esac
    fi

    # 设置下载 URL
    BINARY_NAME="llm-memory-${OS}-${ARCH}"
    DOWNLOAD_URL="https://github.com/XiaoLFeng/llm-memory/releases/download/v${TARGET_VERSION}/${BINARY_NAME}"
    CHECKSUM_URL="https://github.com/XiaoLFeng/llm-memory/releases/download/v${TARGET_VERSION}/checksums.txt"

    # 创建临时目录
    TMP_DIR=$(mktemp -d)
    trap "rm -rf '$TMP_DIR'" EXIT

    # 下载二进制
    print_info "📥 正在下载 llm-memory v${TARGET_VERSION} for ${OS}-${ARCH}..."
    if ! download_with_retry "$DOWNLOAD_URL" "$TMP_DIR/$BINARY_NAME"; then
        print_error "   提示：请检查版本号是否正确，或访问 GitHub Release 页面手动下载"
        print_error "   https://github.com/XiaoLFeng/llm-memory/releases"
        exit 1
    fi
    print_success "✅ 下载完成"

    # 下载并验证校验和
    print_info "🔍 验证文件完整性..."
    if download_with_retry "$CHECKSUM_URL" "$TMP_DIR/checksums.txt"; then
        # 提取对应文件的校验和
        EXPECTED_CHECKSUM=$(grep "$BINARY_NAME" "$TMP_DIR/checksums.txt" | awk '{print $1}')

        if [ -z "$EXPECTED_CHECKSUM" ]; then
            print_warning "⚠️  未找到对应的校验和，跳过校验"
        else
            # 计算实际校验和
            ACTUAL_CHECKSUM=$(sha256sum "$TMP_DIR/$BINARY_NAME" | awk '{print $1}')

            if [ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]; then
                print_error "❌ 文件校验失败！文件可能已损坏或被篡改"
                print_error "   期望: $EXPECTED_CHECKSUM"
                print_error "   实际: $ACTUAL_CHECKSUM"
                exit 1
            fi

            print_success "✅ 文件校验通过"
        fi
    else
        print_warning "⚠️  无法下载校验和文件，跳过校验"
    fi

    # 备份旧版本
    if [ -f "$BINARY_PATH" ]; then
        BACKUP_PATH="${BINARY_PATH}.backup"
        print_info "📦 备份旧版本到 ${BACKUP_PATH}..."
        cp "$BINARY_PATH" "$BACKUP_PATH"
    fi

    # 安装新版本
    print_info "📦 正在更新到 ${INSTALL_DIR}..."
    chmod +x "$TMP_DIR/$BINARY_NAME"
    mv "$TMP_DIR/$BINARY_NAME" "$BINARY_PATH"

    # 验证安装
    NEW_VERSION=$(get_current_version "$BINARY_PATH")
    if [ "$NEW_VERSION" = "$TARGET_VERSION" ]; then
        print_success "✅ 更新成功！"

        # 删除备份
        if [ -f "$BACKUP_PATH" ]; then
            rm -f "$BACKUP_PATH"
        fi
    else
        print_error "❌ 更新后版本验证失败"
        print_error "   期望: v$TARGET_VERSION"
        print_error "   实际: v$NEW_VERSION"

        # 恢复备份
        if [ -f "$BACKUP_PATH" ]; then
            print_info "🔄 正在恢复旧版本..."
            mv "$BACKUP_PATH" "$BINARY_PATH"
        fi
        exit 1
    fi

    print_info ""
    print_success "🎉 更新完成！v${CURRENT_VERSION} → v${TARGET_VERSION}"
    print_info ""
    print_info "使用帮助："
    print_info "  llm-memory --help       # 查看帮助"
    print_info "  llm-memory tui          # 启动 TUI 界面"
    print_info "  llm-memory mcp          # 启动 MCP 服务"
}

# 执行主函数
main "$@"
