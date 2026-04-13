#!/bin/bash

# Script thiết lập mạng GAPChain
# Tác giả: GAPChain Team
# Mô tả: Script tự động thiết lập mạng blockchain cho các tổ chức

# Di chuyển đến thư mục gốc của chain-setup
cd "$(dirname "$0")/.."
CHAINSETUP_DIR=$(pwd)

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

# Kiểm tra các công cụ cần thiết
check_prerequisites() {
    print_info "Kiểm tra các công cụ cần thiết..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker chưa được cài đặt!"
        exit 1
    fi
    
    if ! command -v docker compose &> /dev/null; then
        print_error "Docker Compose chưa được cài đặt!"
        exit 1
    fi
    
    # Kiểm tra các công cụ Fabric trên host
    for tool in cryptogen configtxgen osnadmin; do
        if ! command -v $tool &> /dev/null; then
            print_error "$tool chưa được cài đặt hoặc không có trong PATH! Vui lòng cài đặt Hyperledger Fabric binaries v3.1.1."
            exit 1
        fi
    done
    
    print_warning "Đảm bảo rằng các công cụ Fabric (cryptogen, configtxgen, osnadmin) trong PATH của bạn là phiên bản 3.1.1."
    print_success "Tất cả công cụ cần thiết đã sẵn sàng!"
}

# Tạo thư mục cần thiết
create_directories() {
    print_info "Tạo các thư mục cần thiết..."
    
    # Xóa các thư mục cũ để đảm bảo làm lại từ đầu
    rm -rf organizations channel-artifacts
    
    mkdir -p organizations
    mkdir -p channel-artifacts
    
    print_success "Đã tạo các thư mục cần thiết!"
}

# Tạo certificates cho các tổ chức
generate_certificates() {
    print_info "Tạo certificates cho các tổ chức bằng cryptogen trên host..."
    
    if [ -f "crypto-config.yaml" ]; then
        cryptogen generate --config=./crypto-config.yaml --output="./organizations"
        print_success "Đã tạo certificates thành công!"
    else
        print_error "Không tìm thấy file crypto-config.yaml!"
        exit 1
    fi
}

# Tạo channel artifacts
generate_artifacts() {
    print_info "Tạo channel artifacts bằng configtxgen trên host..."
    
    if [ -f "configtx.yaml" ]; then
        export FABRIC_CFG_PATH=${CHAINSETUP_DIR}

        # Tạo genesis block cho Orderer
        configtxgen -profile GapChainGenesis \
            -outputBlock ${CHAINSETUP_DIR}/channel-artifacts/gapchain-genesis.block \
            -channelID gapchain-genesis

        # Tạo genesis block cho kênh nhật ký HTX
        configtxgen -profile NhatKyHTXChannel \
            -outputBlock ${CHAINSETUP_DIR}/channel-artifacts/nhatky-htx-channel.block \
            -channelID nhatky-htx-channel
        # Tạo genesis block cho kênh giao dịch
        configtxgen -profile GiaoDichChannel \
            -outputBlock ${CHAINSETUP_DIR}/channel-artifacts/giaodich-channel.block \
            -channelID giaodich-channel

        
        # # Tạo channel transaction cho kênh nhật ký HTX
        # configtxgen -profile NhatKyHTXChannel -outputCreateChannelTx "${CHAINSETUP_DIR}/channel-artifacts/nhatky-htx-channel.tx" -channelID nhatky-htx-channel
        
        # # Tạo channel transaction cho kênh giao dịch
        # configtxgen -profile GiaoDichChannel -outputCreateChannelTx "${CHAINSETUP_DIR}/channel-artifacts/giaodich-channel.tx" -channelID giaodich-channel
        
        # Kiểm tra xem file artifact có được tạo và không rỗng
        if [ -s "${CHAINSETUP_DIR}/channel-artifacts/gapchain-genesis.block" ] && [ -s "${CHAINSETUP_DIR}/channel-artifacts/nhatky-htx-channel.block" ] && [ -s "${CHAINSETUP_DIR}/channel-artifacts/giaodich-channel.block" ]; then
            print_success "Đã tạo tất cả channel artifacts thành công!"
        else
            print_error "Tạo channel artifacts thất bại! Một hoặc nhiều tệp .tx bị rỗng. Vui lòng kiểm tra phiên bản configtxgen của bạn."
            exit 1
        fi
    else
        print_error "Không tìm thấy file configtx.yaml!"
        exit 1
    fi
}

# Khởi động mạng
start_network() {
    print_info "Khởi động mạng blockchain..."
    docker compose up -d
    print_info "Đợi các node khởi động trong 15 giây..."
    sleep 10
    print_success "Mạng blockchain đã được khởi động!"
}
# Tạo và join kênh (Fabric v3.x)
create_and_join_channels() {
    print_info "Tạo và join các kênh (quy trình Fabric v3.x)..."

    # --- Cấu hình biến môi trường cho osnadmin (Orderer Admin) ---
    export FABRIC_CFG_PATH=$PWD
    export FABRIC_LOGGING_SPEC=INFO

    export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/gapchain.vn/users/Admin@gapchain.vn/tls/client.crt
    export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/gapchain.vn/users/Admin@gapchain.vn/tls/client.key
    export ORDERER_CA=${PWD}/organizations/ordererOrganizations/gapchain.vn/tlsca/tlsca.gapchain.vn-cert.pem
    
    print_info "Join Orderer vào kênh 'nhatky-htx-channel'..."
    osnadmin channel join \
        --channelID nhatky-htx-channel \
        --config-block ${PWD}/channel-artifacts/nhatky-htx-channel.block \
        -o orderer.gapchain.vn:9441 \
        --ca-file $ORDERER_CA \
        --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT \
        --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

    print_info "Join Orderer vào kênh 'giaodich-channel'..."
    osnadmin channel join \
        --channelID giaodich-channel \
        --config-block ${PWD}/channel-artifacts/giaodich-channel.block \
        -o orderer.gapchain.vn:9441 \
        --ca-file $ORDERER_CA \
        --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT \
        --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

    print_info "Check channel list..."
    osnadmin channel list -o orderer.gapchain.vn:9441 \
        --ca-file $ORDERER_CA \
        --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT \
        --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

    print_info "Đợi Orderer xử lý... (5 giây)"
    sleep 5

    print_info "Join các peer vào kênh..."

    # --- Join peer bằng CLI (Admin của từng org) ---
    export CORE_PEER_TLS_ENABLED=true

    # PlatformOrg
    export FABRIC_CFG_PATH=${PWD}/config/platform.gapchain.vn    
    export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt
    export CORE_PEER_ADDRESS=peer0.platform.gapchain.vn:7051
    peer channel join -b ${PWD}/channel-artifacts/nhatky-htx-channel.block
    peer channel join -b ${PWD}/channel-artifacts/giaodich-channel.block

    # HTXNongSanOrg
    export FABRIC_CFG_PATH=${PWD}/config/htxnongsan.gapchain.vn
    export CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP"
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt
    export CORE_PEER_ADDRESS=peer0.htxnongsan.gapchain.vn:8051
    peer channel join -b ${PWD}/channel-artifacts/nhatky-htx-channel.block
    peer channel join -b ${PWD}/channel-artifacts/giaodich-channel.block

    # ChiCucBVTVOrg
    export FABRIC_CFG_PATH=${PWD}/config/chicucbvtv.gapchain.vn
    export CORE_PEER_LOCALMSPID="ChiCucBVTVOrgMSP"
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/chicucbvtv.gapchain.vn/users/Admin@chicucbvtv.gapchain.vn/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt
    export CORE_PEER_ADDRESS=peer0.chicucbvtv.gapchain.vn:9051
    peer channel join -b ${PWD}/channel-artifacts/nhatky-htx-channel.block

    # NPPXanhOrg
    export FABRIC_CFG_PATH=${PWD}/config/nppxanh.gapchain.vn
    export CORE_PEER_LOCALMSPID="NPPXanhOrgMSP"
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/nppxanh.gapchain.vn/users/Admin@nppxanh.gapchain.vn/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/nppxanh.gapchain.vn/peers/peer0.nppxanh.gapchain.vn/tls/ca.crt
    export CORE_PEER_ADDRESS=peer0.nppxanh.gapchain.vn:10051
    peer channel join -b ${PWD}/channel-artifacts/giaodich-channel.block

    print_success "Đã tạo và join các kênh thành công!"
}


# Hàm chính
main() {
    print_info "Bắt đầu thiết lập mạng GAPChain..."
    export PATH=$PATH:${CHAINSETUP_DIR}/tools/bin
    
    check_prerequisites
    create_directories
    generate_certificates
    generate_artifacts
    start_network
    create_and_join_channels
    
    
    print_success "Mạng GAPChain đã được thiết lập thành công!"
    print_info "Các kênh đã tạo:"
    print_info "  - nhatky-htx-channel (Platform, HTX Nông Sản, Chi Cục BVTV)"
    print_info "  - giaodich-channel (Platform, HTX Nông Sản, NPP Xanh)"
    print_info "Các tổ chức đã tham gia:"
    print_info "  - PlatformOrg (peer0.platform.gapchain.vn:7051)"
    print_info "  - HTXNongSanOrg (peer0.htxnongsan.gapchain.vn:8051)"
    print_info "  - ChiCucBVTVOrg (peer0.chicucbvtv.gapchain.vn:9051)"
    print_info "  - NPPXanhOrg (peer0.nppxanh.gapchain.vn:10051)"
}

# Chạy script
main "$@"