#!/bin/bash
# llm-memory 卸载脚本
# 自动清理安装的二进制文件和相关配置
#
# 使用方法：
#   curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/uninstall.sh | bash
#   或者下载后执行：
#   bash uninstall.sh

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

# 询问用户确认
confirm() {
    local prompt="$1"
    local default="${2:-n}"

    if [ "$default" = "y" ]; then
        prompt="$prompt [Y/n] "
        default_value="y"
    else
        prompt="$prompt [y/N] "
        default_value="n"
    fi

    read -p "$(echo -e ${CYAN}${prompt}${NC})" -r response
    response=${response:-$default_value}

    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# 主函数
main() {
    print_info "🗑️  llm-memory 卸载脚本"
    print_info ""

    # 定义安装位置
    INSTALL_DIR="$HOME/.local/bin"
    BINARY_PATH="$INSTALL_DIR/llm-memory"
    CONFIG_DIR="$HOME/.llm-memory"

    # 检查是否已安装
    if [ ! -f "$BINARY_PATH" ]; then
        print_warning "⚠️  未找到已安装的 llm-memory"
        print_info "   预期位置: $BINARY_PATH"

        # 检查是否在其他位置
        if command -v llm-memory &> /dev/null; then
            local found_path=$(which llm-memory)
            print_warning "   但在 PATH 中找到: $found_path"
            print_info ""

            if confirm "是否删除该位置的 llm-memory？" "n"; then
                BINARY_PATH="$found_path"
            else
                print_info "取消卸载"
                exit 0
            fi
        else
            print_info ""
            print_info "llm-memory 可能已经卸载，或安装在其他位置"
            exit 0
        fi
    fi

    print_info "📍 找到安装位置："
    print_info "   二进制文件: $BINARY_PATH"

    # 检查版本
    if [ -x "$BINARY_PATH" ]; then
        print_info "   当前版本: $($BINARY_PATH --version 2>/dev/null || echo '未知')"
    fi

    # 检查配置目录
    if [ -d "$CONFIG_DIR" ]; then
        print_info "   配置目录: $CONFIG_DIR"

        # 计算配置目录大小
        local config_size=$(du -sh "$CONFIG_DIR" 2>/dev/null | cut -f1)
        print_info "   配置大小: $config_size"
    fi

    print_info ""

    # 询问是否继续
    if ! confirm "确定要卸载 llm-memory 吗？" "n"; then
        print_info "取消卸载"
        exit 0
    fi

    print_info ""

    # 删除二进制文件
    print_info "🗑️  正在删除二进制文件..."
    if rm -f "$BINARY_PATH"; then
        print_success "✅ 已删除: $BINARY_PATH"
    else
        print_error "❌ 删除失败: $BINARY_PATH"
        print_error "   你可能需要手动删除该文件"
    fi

    # 询问是否删除配置
    if [ -d "$CONFIG_DIR" ]; then
        print_info ""
        print_warning "⚠️  注意：配置目录包含你的所有数据（记忆、计划、待办）"

        if confirm "是否同时删除配置目录和所有数据？" "n"; then
            print_info "🗑️  正在删除配置目录..."
            if rm -rf "$CONFIG_DIR"; then
                print_success "✅ 已删除: $CONFIG_DIR"
            else
                print_error "❌ 删除失败: $CONFIG_DIR"
                print_error "   你可能需要手动删除该目录"
            fi
        else
            print_info "保留配置目录: $CONFIG_DIR"
            print_info "如果将来需要删除，可以运行："
            print_info "  ${CYAN}rm -rf $CONFIG_DIR${NC}"
        fi
    fi

    print_info ""
    print_success "🎉 llm-memory 卸载完成！"

    # 检查是否还在 PATH 中
    if command -v llm-memory &> /dev/null; then
        print_info ""
        print_warning "⚠️  注意：llm-memory 仍在 PATH 中"
        print_warning "   位置: $(which llm-memory)"
        print_warning "   这可能是另一个安装位置，请手动检查"
    fi

    print_info ""
    print_info "感谢使用 llm-memory！(´∀｀)💖"
    print_info ""
    print_info "如果你遇到了问题或有建议，欢迎反馈："
    print_info "  https://github.com/XiaoLFeng/llm-memory/issues"
}

# 执行主函数
main "$@"
