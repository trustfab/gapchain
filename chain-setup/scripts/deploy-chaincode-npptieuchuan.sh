#!/bin/bash

# Script triển khai chaincode giao_dich_cc cho tổ chức mới (NPPTieuChuanOrg)
# Tác giả: GAPChain Team (Agentic Assistant)
# Mô tả: Chu trình Package -> Install -> Approve -> Commit để cấp quyền
#        cho NPPTieuChuan tương tác với chaincode giao_dich_cc.
# Chú ý: Vì có thêm thành viên, cần nâng Sequence lên +1 (sequence=9) và 
#        Commit lại trên toàn network.

set -e

cd "$(dirname "$0")/.."
CHAINSETUP_DIR=$(pwd)

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Cấu hình Chaincode (Giao Dịch CC)
CC_NAME="giao_dich_cc"
CC_CHANNEL="giaodich-channel"
CC_VERSION="1.1"
CC_SEQUENCE="2"  # Nâng từ 1 lên 2 để update commit definition

PACKAGE_ID=""
ORDERER_CA="/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem"

print_info "1. Đóng gói (Package) chaincode giao_dich_cc..."
docker volume create chaincode_pkg_vol || true
docker run --rm -v $(pwd)/chaincode/giao_dich_cc:/chaincode -w /chaincode golang:1.21-alpine sh -c "go mod tidy && go mod download"
docker run --rm -v $(pwd)/chaincode/giao_dich_cc:/src -v chaincode_pkg_vol:/dest -w /src -e GOCACHE=/tmp/.cache hyperledger/fabric-tools:latest peer lifecycle chaincode package /dest/giao_dich_cc.tar.gz --path /src --lang golang --label giao_dich_cc_${CC_VERSION}
docker container create --name dummy_pkg_new -v chaincode_pkg_vol:/dest alpine
docker cp dummy_pkg_new:/dest/giao_dich_cc.tar.gz ./giao_dich_cc.tar.gz
docker rm dummy_pkg_new

print_info "2. Copy tới các Peer container..."
# Chúng ta copy vào NPPTieuChuan (tổ chức mới)
docker cp giao_dich_cc.tar.gz peer0.npptieuchuan.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
# Copy cả vào các org đã có vì cần Approve lại chuỗi sequence mới
docker cp giao_dich_cc.tar.gz peer0.platform.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp giao_dich_cc.tar.gz peer0.htxnongsan.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp giao_dich_cc.tar.gz peer0.nppxanh.gapchain.vn:/opt/gopath/src/github.com/hyperledger/fabric/peer/

print_info "3. Install (Cài đặt) chaincode lên NPPTieuChuanOrg..."
docker exec \
  -e CORE_PEER_LOCALMSPID="NPPTieuChuanOrgMSP" \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/users/Admin@npptieuchuan.gapchain.vn/msp \
  peer0.npptieuchuan.gapchain.vn \
  peer lifecycle chaincode install giao_dich_cc.tar.gz || true

# Các org cũ (có thể đã có, cứ chạy lại ignor lỗi)
docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp peer0.platform.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true
docker exec -e CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp peer0.htxnongsan.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true
docker exec -e CORE_PEER_LOCALMSPID="NPPXanhOrgMSP" -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/users/Admin@nppxanh.gapchain.vn/msp peer0.nppxanh.gapchain.vn peer lifecycle chaincode install giao_dich_cc.tar.gz || true


print_info "4. Truy xuất Package ID..."
PACKAGE_ID=$(docker exec -e CORE_PEER_LOCALMSPID="NPPTieuChuanOrgMSP" -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/users/Admin@npptieuchuan.gapchain.vn/msp peer0.npptieuchuan.gapchain.vn peer lifecycle chaincode queryinstalled --output json | jq -r '.installed_chaincodes[] | select(.label=="'${CC_NAME}'_'${CC_VERSION}'") | .package_id')

if [ -z "$PACKAGE_ID" ]; then
    print_error "Không lấy được PACKAGE_ID. Cài đặt có thể bị lỗi."
fi
print_success "PACKAGE_ID: $PACKAGE_ID"

print_info "5. Approve (Phê duyệt) chaincode sequence mới ($CC_SEQUENCE) từ cả 4 Orgs..."

approve_org() {
    local org=$1
    local domain=$2
    docker exec \
      -e CORE_PEER_LOCALMSPID="${org}MSP" \
      -e CORE_PEER_TLS_ENABLED=true \
      -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/${domain}/peers/peer0.${domain}/tls/ca.crt \
      -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/${domain}/users/Admin@${domain}/msp \
      peer0.${domain} \
      peer lifecycle chaincode approveformyorg \
      -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn \
      --channelID ${CC_CHANNEL} --name ${CC_NAME} --version ${CC_VERSION} \
      --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} --tls --cafile $ORDERER_CA
}

approve_org "PlatformOrg" "platform.gapchain.vn"
approve_org "HTXNongSanOrg" "htxnongsan.gapchain.vn"
approve_org "NPPXanhOrg" "nppxanh.gapchain.vn"
approve_org "NPPTieuChuanOrg" "npptieuchuan.gapchain.vn"

print_success "Cả 4 tổ chức đã Approve."

print_info "6. Commit cập nhật Chaincode giaodich-channel..."
docker exec \
  -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp \
  peer0.platform.gapchain.vn \
  peer lifecycle chaincode commit \
  -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn \
  --channelID ${CC_CHANNEL} --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --tls --cafile $ORDERER_CA \
  --peerAddresses peer0.platform.gapchain.vn:7051 \
  --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  --peerAddresses peer0.htxnongsan.gapchain.vn:7051 \
  --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt \
  --peerAddresses peer0.nppxanh.gapchain.vn:7051 \
  --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/nppxanh.gapchain.vn/peers/peer0.nppxanh.gapchain.vn/tls/ca.crt \
  --peerAddresses peer0.npptieuchuan.gapchain.vn:11051 \
  --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/tls/ca.crt

print_success "Chaincode giao_dich_cc (seq $CC_SEQUENCE) đã được Commit!"

print_info "7. Query kiểm tra..."
docker exec \
  -e CORE_PEER_LOCALMSPID="NPPTieuChuanOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/npptieuchuan.gapchain.vn/users/Admin@npptieuchuan.gapchain.vn/msp \
  peer0.npptieuchuan.gapchain.vn \
  peer lifecycle chaincode querycommitted --channelID ${CC_CHANNEL} --name ${CC_NAME}

print_success "Hoàn tất! NPPTieuChuanOrg giờ đây đã có thể Invoke và Query giao_dich_cc."
