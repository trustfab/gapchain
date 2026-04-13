#!/bin/bash

# Script cập nhật Docker images lên phiên bản mới nhất
# Tác giả: GAPChain Team
# Mô tả: Script tải xuống và cập nhật các Docker images cho Fabric 3.1.1

set -e

# Màu sắc cho output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Hàm in thông báo
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Phiên bản mới
FABRIC_VERSION="3.1.1"
FABRIC_CA_VERSION="1.5.15"

# Tải xuống Fabric images
download_fabric_images() {
    print_info "Tải xuống Hyperledger Fabric images phiên bản $FABRIC_VERSION..."
    
    # Tải xuống orderer image
    print_info "Tải xuống fabric-orderer:$FABRIC_VERSION..."
    docker pull hyperledger/fabric-orderer:$FABRIC_VERSION
    
    # Tải xuống peer image
    print_info "Tải xuống fabric-peer:$FABRIC_VERSION..."
    docker pull hyperledger/fabric-peer:$FABRIC_VERSION
    
    # Tải xuống tools image
    print_info "Tải xuống fabric-tools:$FABRIC_VERSION..."
    docker pull hyperledger/fabric-tools:$FABRIC_VERSION
    
    print_success "Đã tải xuống tất cả Fabric images!"
}

# Tải xuống Fabric CA images
download_fabric_ca_images() {
    print_info "Tải xuống Hyperledger Fabric CA images phiên bản $FABRIC_CA_VERSION..."
    
    # Tải xuống fabric-ca image
    print_info "Tải xuống fabric-ca:$FABRIC_CA_VERSION..."
    docker pull hyperledger/fabric-ca:$FABRIC_CA_VERSION
    
    print_success "Đã tải xuống Fabric CA images!"
}

# Kiểm tra images đã tải xuống
check_images() {
    print_info "Kiểm tra các images đã tải xuống..."
    
    echo "=== Fabric Images ==="
    docker images | grep "hyperledger/fabric" | grep -E "($FABRIC_VERSION|$FABRIC_CA_VERSION)"
    
    print_success "Kiểm tra hoàn tất!"
}

# Dọn dẹp images cũ (tùy chọn)
cleanup_old_images() {
    if [ "$1" = "--cleanup" ]; then
        print_warning "Dọn dẹp các images cũ..."
        
        # Xóa images cũ của Fabric
        docker rmi $(docker images "hyperledger/fabric-peer:2.5.4" -q) 2>/dev/null || true
        docker rmi $(docker images "hyperledger/fabric-orderer:2.5.4" -q) 2>/dev/null || true
        docker rmi $(docker images "hyperledger/fabric-tools:2.5.4" -q) 2>/dev/null || true
        
        # Xóa images cũ của Fabric CA
        docker rmi $(docker images "hyperledger/fabric-ca:1.5.4" -q) 2>/dev/null || true
        
        print_success "Đã dọn dẹp các images cũ!"
    fi
}

# Hiển thị hướng dẫn sử dụng
show_usage() {
    echo "Sử dụng: $0 [OPTIONS]"
    echo ""
    echo "OPTIONS:"
    echo "  --cleanup    Dọn dẹp các images cũ sau khi tải xuống"
    echo "  --help       Hiển thị hướng dẫn này"
    echo ""
    echo "Ví dụ:"
    echo "  $0                # Chỉ tải xuống images mới"
    echo "  $0 --cleanup      # Tải xuống và dọn dẹp images cũ"
}

# Hàm chính
main() {
    print_info "Bắt đầu cập nhật Docker images lên Fabric $FABRIC_VERSION..."
    
    # Kiểm tra tham số
    if [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    # Kiểm tra Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker chưa được cài đặt!"
        exit 1
    fi
    
    download_fabric_images
    download_fabric_ca_images
    check_images
    cleanup_old_images "$1"
    
    print_success "Cập nhật Docker images hoàn tất!"
    print_info "Các images mới:"
    print_info "  - hyperledger/fabric-orderer:$FABRIC_VERSION"
    print_info "  - hyperledger/fabric-peer:$FABRIC_VERSION"
    print_info "  - hyperledger/fabric-tools:$FABRIC_VERSION"
    print_info "  - hyperledger/fabric-ca:$FABRIC_CA_VERSION"
    print_info ""
    print_info "Bây giờ bạn có thể chạy:"
    print_info "  ./scripts/setup-network.sh"
}

# Chạy script
main "$@"
