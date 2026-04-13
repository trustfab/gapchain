#!/bin/bash

# Script dọn dẹp mạng GAPChain
# Tác giả: GAPChain Team
# Mô tả: Script dừng và dọn dẹp tất cả containers, volumes và artifacts

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

# Dừng tất cả containers
stop_containers() {
    print_info "Dừng tất cả containers..."
    
    # Dừng peer và orderer containers
    docker compose down
    docker compose -f docker-compose-npptieuchuan.yaml down 2>/dev/null || true
    
    # Dừng CA containers
    docker compose -f docker compose-ca.yaml down
    
    print_success "Đã dừng tất cả containers!"
}

# Xóa containers
remove_containers() {
    print_info "Xóa tất cả containers..."
    
    # Xóa containers nếu còn tồn tại
    docker rm -f $(docker ps -aq --filter "name=peer0.platform.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=peer0.htxnongsan.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=peer0.chicucbvtv.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=peer0.nppxanh.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=peer0.npptieuchuan.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=couchdb-npptieuchuan") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=orderer.gapchain.vn") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=ca-orderer") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=ca-platform") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=ca-htxnongsan") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=ca-chicucbvtv") 2>/dev/null || true
    docker rm -f $(docker ps -aq --filter "name=ca-nppxanh") 2>/dev/null || true
    
    print_success "Đã xóa tất cả containers!"
}

# Xóa volumes
remove_volumes() {
    print_info "Xóa tất cả volumes..."
    
    docker volume rm -f $(docker volume ls -q --filter "name=orderer.gapchain.vn") 2>/dev/null || true
    docker volume rm -f $(docker volume ls -q --filter "name=peer0.platform.gapchain.vn") 2>/dev/null || true
    docker volume rm -f $(docker volume ls -q --filter "name=peer0.htxnongsan.gapchain.vn") 2>/dev/null || true
    docker volume rm -f $(docker volume ls -q --filter "name=peer0.chicucbvtv.gapchain.vn") 2>/dev/null || true
    docker volume rm -f $(docker volume ls -q --filter "name=peer0.nppxanh.gapchain.vn") 2>/dev/null || true
    docker volume rm -f $(docker volume ls -q --filter "name=peer0.npptieuchuan.gapchain.vn") 2>/dev/null || true
    
    print_success "Đã xóa tất cả volumes!"
}

# Xóa networks
remove_networks() {
    print_info "Xóa networks..."
    
    docker network rm gapchain_gapchain 2>/dev/null || true
    
    print_success "Đã xóa networks!"
}

# Xóa artifacts và certificates
remove_artifacts() {
    print_info "Xóa artifacts và certificates..."
    
    # Xóa thư mục organizations
    if [ -d "organizations" ]; then
        rm -rf organizations
        print_info "Đã xóa thư mục organizations"
    fi
    
    # Xóa thư mục channel-artifacts
    if [ -d "channel-artifacts" ]; then
        rm -rf channel-artifacts
        print_info "Đã xóa thư mục channel-artifacts"
    fi
    
    # Xóa thư mục tạm nếu có
    if [ -d "temp_add_org" ]; then
        rm -rf temp_add_org
        print_info "Đã xóa thư mục tạm temp_add_org"
    fi
    
    # Xóa script sinh ra bởi add_org
    rm -f docker-compose-npptieuchuan.yaml 2>/dev/null || true
    
    # Xóa chaincode packages
    rm -f *.tar.gz 2>/dev/null || true
    
    print_success "Đã xóa tất cả artifacts và certificates!"
}

# Xóa images (tùy chọn)
remove_images() {
    if [ "$1" = "--remove-images" ]; then
        print_warning "Xóa tất cả Fabric images..."
        
        docker rmi -f $(docker images -q --filter "reference=hyperledger/fabric-peer") 2>/dev/null || true
        docker rmi -f $(docker images -q --filter "reference=hyperledger/fabric-orderer") 2>/dev/null || true
        docker rmi -f $(docker images -q --filter "reference=hyperledger/fabric-ca") 2>/dev/null || true
        docker rmi -f $(docker images -q --filter "reference=hyperledger/fabric-tools") 2>/dev/null || true
        
        print_success "Đã xóa tất cả Fabric images!"
    fi
}

# Dọn dẹp hệ thống Docker
cleanup_docker() {
    print_info "Dọn dẹp hệ thống Docker..."
    
    # Xóa tất cả containers đã dừng
    docker container prune -f
    
    # Xóa tất cả volumes không sử dụng
    docker volume prune -f
    
    # Xóa tất cả networks không sử dụng
    docker network prune -f
    
    print_success "Đã dọn dẹp hệ thống Docker!"
}

# Hiển thị hướng dẫn sử dụng
show_usage() {
    echo "Sử dụng: $0 [OPTIONS]"
    echo ""
    echo "OPTIONS:"
    echo "  --remove-images    Xóa cả Docker images (cẩn thận!)"
    echo "  --help            Hiển thị hướng dẫn này"
    echo ""
    echo "Ví dụ:"
    echo "  $0                # Dọn dẹp cơ bản"
    echo "  $0 --remove-images # Dọn dẹp hoàn toàn bao gồm images"
}

# Hàm chính
main() {
    print_info "Bắt đầu dọn dẹp mạng GAPChain..."
    
    # Kiểm tra tham số
    if [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    # Xác nhận từ người dùng
    print_warning "Bạn có chắc chắn muốn dọn dẹp toàn bộ mạng GAPChain?"
    print_warning "Hành động này sẽ xóa tất cả containers, volumes, và artifacts!"
    read -p "Nhập 'yes' để tiếp tục: " confirm
    
    if [ "$confirm" != "yes" ]; then
        print_info "Hủy bỏ dọn dẹp."
        exit 0
    fi
    
    stop_containers
    remove_containers
    remove_volumes
    remove_networks
    remove_artifacts
    remove_images "$1"
    cleanup_docker
    
    print_success "Đã dọn dẹp mạng GAPChain hoàn toàn!"
    print_info "Để thiết lập lại mạng, chạy: ./scripts/setup-network.sh"
}

# Chạy script
main "$@"
