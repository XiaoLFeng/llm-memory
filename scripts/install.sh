#!/bin/bash
# llm-memory 安装脚本
# 自动检测系统和架构，下载、校验并安装 llm-memory
#
# 使用方法：
#   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.sh | bash
#   或指定版本：
#   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.sh | bash -s v0.0.2

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
        print_error "   请检查网络连接或手动指定版本: bash install.sh v0.0.2"
        exit 1
    fi

    echo "$version"
}

# 主函数
main() {
    print_info "🚀 llm-memory 安装脚本"
    print_info ""

    # 检查依赖
    check_dependencies

    # 检测系统
    OS=$(detect_os)
    ARCH=$(detect_arch)
    print_success "✅ 检测到系统: $OS-$ARCH"

    # 获取版本
    VERSION="${1:-latest}"
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(get_latest_version)
    else
        # 去掉可能的 v 前缀
        VERSION="${VERSION#v}"
    fi
    print_success "✅ 目标版本: v$VERSION"
    print_info ""

    # 设置下载 URL
    BINARY_NAME="llm-memory-${OS}-${ARCH}"
    DOWNLOAD_URL="https://github.com/XiaoLFeng/llm-memory/releases/download/v${VERSION}/${BINARY_NAME}"
    CHECKSUM_URL="https://github.com/XiaoLFeng/llm-memory/releases/download/v${VERSION}/checksums.txt"

    # 创建临时目录
    TMP_DIR=$(mktemp -d)
    trap "rm -rf '$TMP_DIR'" EXIT

    # 下载二进制
    print_info "📥 正在下载 llm-memory v${VERSION} for ${OS}-${ARCH}..."
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

    # 安装二进制
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"

    print_info "📦 正在安装到 ${INSTALL_DIR}..."
    chmod +x "$TMP_DIR/$BINARY_NAME"
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/llm-memory"
    print_success "✅ 安装成功！"

    print_info ""

    # 检查 PATH 配置
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        print_warning "⚠️  注意：${INSTALL_DIR} 不在 PATH 中"
        print_info ""
        print_info "请将以下内容添加到你的 shell 配置文件："
        print_info ""

        # 检测 shell 类型并给出建议
        if [ -n "$ZSH_VERSION" ]; then
            SHELL_RC="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            if [ -f "$HOME/.bashrc" ]; then
                SHELL_RC="$HOME/.bashrc"
            else
                SHELL_RC="$HOME/.bash_profile"
            fi
        else
            SHELL_RC="$HOME/.profile"
        fi

        print_info "    ${CYAN}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> $SHELL_RC${NC}"
        print_info "    ${CYAN}source $SHELL_RC${NC}"
        print_info ""
        print_info "或者直接运行（临时生效）："
        print_info "    ${CYAN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
        print_info ""
    else
        print_success "🎉 安装完成！你现在可以运行: ${CYAN}llm-memory --version${NC}"
        print_info ""
        print_info "使用帮助："
        print_info "  llm-memory --help       # 查看帮助"
        print_info "  llm-memory tui          # 启动 TUI 界面"
        print_info "  llm-memory mcp          # 启动 MCP 服务"
    fi
}

# 执行主函数
main "$@"
