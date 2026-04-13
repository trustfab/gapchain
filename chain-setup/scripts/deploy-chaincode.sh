#!/bin/bash

# Script triển khai chaincode cho mạng GAPChain
# Tác giả: GAPChain Team
# Mô tả: Script tự động triển khai chaincode nhatky_cc và giao_dich_cc

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

# Biến cấu hình
CHAINCODE_VERSION="1.2"
CHAINCODE_SEQUENCE="1"
PACKAGE_ID=""

# Triển khai chaincode lô hàng (lohang_cc) — phải deploy TRƯỚC nhatky_cc
deploy_lohang_chaincode() {
    print_info "Triển khai chaincode lohang_cc cho kênh nhatky-htx-channel..."

    print_info "Building and packaging lohang_cc..."
    docker volume create chaincode_pkg_vol || true
    docker run --rm -v $(pwd)/chaincode/lohang_cc:/chaincode -w /chaincode golang:1.21-alpine sh -c "go mod tidy && go mod download"
    docker run --rm -v $(pwd)/chaincode/lohang_cc:/src -v chaincode_pkg_vol:/dest -w /src -e GOCACHE=/tmp/.cache hyperledger/fabric-tools:latest peer lifecycle chaincode package /dest/lohang_cc.tar.gz --path /src --lang golang --label lohang_cc_${CHAINCODE_VERSION}
    docker container create --name dummy_pkg_l -v chaincode_pkg_vol:/dest alpine
    docker cp dummy_pkg_l:/dest/lohang_cc.tar.gz ./lohang_cc.tar.gz
    docker rm dummy_pkg_l

    docker cp lohang_cc.tar.gz peer0.platform.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp lohang_cc.tar.gz peer0.htxnongsan.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp lohang_cc.tar.gz peer0.chicucbvtv.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/

    # Install chaincode trên tất cả peers trong kênh nhatky-htx-channel
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode install lohang_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode install lohang_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="ChiCucBVTVOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/users/Admin@chicucbvtv.gapchain.vn/msp peer0.chicucbvtv.gapchain.vn peer lifecycle chaincode install lohang_cc.tar.gz || true

    # Lấy package ID
    PACKAGE_ID=$(docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode queryinstalled --output json | jq -r '.installed_chaincodes[] | select(.label=="lohang_cc_'${CHAINCODE_VERSION}'") | .package_id')

    if [ -z "$PACKAGE_ID" ]; then
        print_error "Không thể lấy package ID cho lohang_cc"
        exit 1
    fi

    print_info "Package ID cho lohang_cc: $PACKAGE_ID"

    # Approve chaincode cho từng organization
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name lohang_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem

    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name lohang_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem

    docker exec -e CORE_PEER_LOCALMSPID="ChiCucBVTVOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/users/Admin@chicucbvtv.gapchain.vn/msp peer0.chicucbvtv.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name lohang_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem

    # Commit chaincode
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode commit -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name lohang_cc --version ${CHAINCODE_VERSION} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem --peerAddresses peer0.platform.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt --peerAddresses peer0.chicucbvtv.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt

    print_success "Đã triển khai chaincode lohang_cc thành công!"
}

# Triển khai chaincode nhật ký HTX
deploy_nhatky_chaincode() {
    print_info "Triển khai chaincode nhatky_cc cho kênh nhatky-htx-channel..."
    
    # Build and package chaincode using a standalone golang container
    print_info "Building and packaging nhatky_cc..."
    docker volume create chaincode_pkg_vol || true
    docker run --rm -v $(pwd)/chaincode/nhatky_cc:/chaincode -w /chaincode golang:1.21-alpine sh -c "go mod tidy && go mod download"
    docker run --rm -v $(pwd)/chaincode/nhatky_cc:/src -v chaincode_pkg_vol:/dest -w /src -e GOCACHE=/tmp/.cache hyperledger/fabric-tools:latest peer lifecycle chaincode package /dest/nhatky_cc.tar.gz --path /src --lang golang --label nhatky_cc_${CHAINCODE_VERSION}
    docker container create --name dummy_pkg -v chaincode_pkg_vol:/dest alpine
    docker cp dummy_pkg:/dest/nhatky_cc.tar.gz ./nhatky_cc.tar.gz
    docker rm dummy_pkg
    
    docker cp nhatky_cc.tar.gz peer0.platform.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp nhatky_cc.tar.gz peer0.htxnongsan.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp nhatky_cc.tar.gz peer0.chicucbvtv.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    
    # Install chaincode trên tất cả peers trong kênh
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode install nhatky_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode install nhatky_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="ChiCucBVTVOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/users/Admin@chicucbvtv.gapchain.vn/msp peer0.chicucbvtv.gapchain.vn peer lifecycle chaincode install nhatky_cc.tar.gz || true
    
    # Lấy package ID
    PACKAGE_ID=$(docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode queryinstalled --output json | jq -r '.installed_chaincodes[] | select(.label=="nhatky_cc_'${CHAINCODE_VERSION}'") | .package_id')
    
    if [ -z "$PACKAGE_ID" ]; then
        print_error "Không thể lấy package ID cho nhatky_cc"
        exit 1
    fi
    
    print_info "Package ID cho nhatky_cc: $PACKAGE_ID"
    
    # Approve chaincode cho từng organization
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name nhatky_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name nhatky_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    docker exec -e CORE_PEER_LOCALMSPID="ChiCucBVTVOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/users/Admin@chicucbvtv.gapchain.vn/msp peer0.chicucbvtv.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name nhatky_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    # Commit chaincode
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode commit -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID nhatky-htx-channel --name nhatky_cc --version ${CHAINCODE_VERSION} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem --peerAddresses peer0.platform.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt --peerAddresses peer0.chicucbvtv.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/chicucbvtv.gapchain.vn/peers/peer0.chicucbvtv.gapchain.vn/tls/ca.crt
    
    print_success "Đã triển khai chaincode nhatky_cc thành công!"
}

# Triển khai chaincode giao dịch
deploy_giaodich_chaincode() {
    print_info "Triển khai chaincode giao_dich_cc cho kênh giaodich-channel..."
    
    # Build and package chaincode using a standalone golang container
    print_info "Building and packaging giao_dich_cc..."
    docker volume create chaincode_pkg_vol || true
    docker run --rm -v $(pwd)/chaincode/giao_dich_cc:/chaincode -w /chaincode golang:1.21-alpine sh -c "go mod tidy && go mod download"
    docker run --rm -v $(pwd)/chaincode/giao_dich_cc:/src -v chaincode_pkg_vol:/dest -w /src -e GOCACHE=/tmp/.cache hyperledger/fabric-tools:latest peer lifecycle chaincode package /dest/giao_dich_cc.tar.gz --path /src --lang golang --label giao_dich_cc_${CHAINCODE_VERSION}
    docker container create --name dummy_pkg_g -v chaincode_pkg_vol:/dest alpine
    docker cp dummy_pkg_g:/dest/giao_dich_cc.tar.gz ./giao_dich_cc.tar.gz
    docker rm dummy_pkg_g
    
    docker cp giao_dich_cc.tar.gz peer0.platform.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp giao_dich_cc.tar.gz peer0.htxnongsan.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    docker cp giao_dich_cc.tar.gz peer0.nppxanh.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
    
    # Install chaincode trên tất cả peers trong kênh
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true
    docker exec -e CORE_PEER_LOCALMSPID="NPPXanhOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/users/Admin@nppxanh.gapchain.vn/msp peer0.nppxanh.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true
    
    # Lấy package ID
    PACKAGE_ID=$(docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode queryinstalled --output json | jq -r '.installed_chaincodes[] | select(.label=="giao_dich_cc_'${CHAINCODE_VERSION}'") | .package_id')
    
    if [ -z "$PACKAGE_ID" ]; then
        print_error "Không thể lấy package ID cho giao_dich_cc"
        exit 1
    fi
    
    print_info "Package ID cho giao_dich_cc: $PACKAGE_ID"
    
    # Approve chaincode cho từng organization
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID giaodich-channel --name giao_dich_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID giaodich-channel --name giao_dich_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    docker exec -e CORE_PEER_LOCALMSPID="NPPXanhOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/peers/peer0.nppxanh.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/users/Admin@nppxanh.gapchain.vn/msp peer0.nppxanh.gapchain.vn peer lifecycle chaincode approveformyorg -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID giaodich-channel --name giao_dich_cc --version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
    
    # Commit chaincode
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode commit -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --channelID giaodich-channel --name giao_dich_cc --version ${CHAINCODE_VERSION} --sequence ${CHAINCODE_SEQUENCE} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem --peerAddresses peer0.platform.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt --peerAddresses peer0.htxnongsan.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt --peerAddresses peer0.nppxanh.gapchain.vn:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/peers/peer0.nppxanh.gapchain.vn/tls/ca.crt
    
    print_success "Đã triển khai chaincode giao_dich_cc thành công!"
}

# Kiểm tra trạng thái chaincode
check_chaincode_status() {
    print_info "Kiểm tra trạng thái chaincode..."
    
    # Kiểm tra lohang_cc
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode querycommitted --channelID nhatky-htx-channel --name lohang_cc

    # Kiểm tra nhatky_cc
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode querycommitted --channelID nhatky-htx-channel --name nhatky_cc

    # Kiểm tra giao_dich_cc
    docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode querycommitted --channelID giaodich-channel --name giao_dich_cc
    
    print_success "Kiểm tra trạng thái chaincode hoàn tất!"
}

# Hàm chính
main() {
    print_info "Bắt đầu triển khai chaincode cho mạng GAPChain..."
    
    # Kiểm tra xem mạng đã chạy chưa
    if ! docker ps | grep -q "peer0.platform.gapchain.vn"; then
        print_error "Mạng blockchain chưa được khởi động! Vui lòng chạy setup-network.sh trước."
        exit 1
    fi
    
    deploy_lohang_chaincode
    deploy_nhatky_chaincode
    deploy_giaodich_chaincode
    check_chaincode_status

    print_success "Đã triển khai tất cả chaincode thành công!"
    print_info "Chaincode đã triển khai:"
    print_info "  - lohang_cc  trên kênh nhatky-htx-channel (Platform, HTXNongSan, ChiCucBVTV)"
    print_info "  - nhatky_cc  trên kênh nhatky-htx-channel (Platform, HTXNongSan, ChiCucBVTV)"
    print_info "  - giao_dich_cc trên kênh giaodich-channel  (Platform, HTXNongSan, NPPXanh)"
}

# Chạy script
main "$@"
