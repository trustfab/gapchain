#!/bin/bash

# Script thêm tổ chức NPPTieuChuanOrg vào mạng GAPChain (kênh giaodich-channel)
# Tác giả: GAPChain Team (Agentic Assistant)
# Yêu cầu: Chạy ở thư mục gốc của chain-setup. Đảm bảo mạng đang chạy.

cd "$(dirname "$0")/.."
CHAINSETUP_DIR=$(pwd)
export PATH=$PATH:${CHAINSETUP_DIR}/../bin
export FABRIC_CFG_PATH=${CHAINSETUP_DIR}
export ORDERER_CA=${CHAINSETUP_DIR}/organizations/ordererOrganizations/gapchain.vn/tlsca/tlsca.gapchain.vn-cert.pem

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

TEMP_DIR="temp_add_org"
mkdir -p ${TEMP_DIR}

print_info "1. Tạo cấu hình crypto-config cho NPPTieuChuanOrg..."
cat <<EOF > ${TEMP_DIR}/crypto-config-npptieuchuan.yaml
PeerOrgs:
  - Name: NPPTieuChuanOrg
    Domain: npptieuchuan.gapchain.vn
    EnableNodeOUs: true
    Template:
      Count: 1
      SANS:
        - "localhost"
        - "127.0.0.1"
    Users:
      Count: 2
EOF

cryptogen generate --config=${TEMP_DIR}/crypto-config-npptieuchuan.yaml --output="./organizations"
print_success "Đã tạo certificates cho NPPTieuChuanOrg"

print_info "2. Chuẩn bị configtx cục bộ để printOrg..."
cat <<EOF > ${TEMP_DIR}/configtx.yaml
Organizations:
  - &NPPTieuChuanOrg
    Name: NPPTieuChuanOrgMSP
    ID: NPPTieuChuanOrgMSP
    MSPDir: ../organizations/peerOrganizations/npptieuchuan.gapchain.vn/msp
    Policies:
      Readers:
        Type: Signature
        Rule: "OR('NPPTieuChuanOrgMSP.admin', 'NPPTieuChuanOrgMSP.peer', 'NPPTieuChuanOrgMSP.client')"
      Writers:
        Type: Signature
        Rule: "OR('NPPTieuChuanOrgMSP.admin', 'NPPTieuChuanOrgMSP.client')"
      Admins:
        Type: Signature
        Rule: "OR('NPPTieuChuanOrgMSP.admin')"
      Endorsement:
        Type: Signature
        Rule: "OR('NPPTieuChuanOrgMSP.peer')"
    AnchorPeers:
      - Host: peer0.npptieuchuan.gapchain.vn
        Port: 11051
EOF

configtxgen -printOrg NPPTieuChuanOrgMSP -configPath ./${TEMP_DIR} > ${TEMP_DIR}/npptieuchuan.json
print_success "Đã trích xuất cấu hình JSON của NPPTieuChuanOrg"

print_info "3. Lấy block cấu hình mới nhất từ giaodich-channel..."
# Tạm thiết lập context là Admin của PlatformOrg (Đã nằm trong channel)
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_MSPCONFIGPATH=${CHAINSETUP_DIR}/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${CHAINSETUP_DIR}/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt
export CORE_PEER_ADDRESS=peer0.platform.gapchain.vn:7051

peer channel fetch config ${TEMP_DIR}/config_block.pb -o localhost:7050 -c giaodich-channel --tls --cafile $ORDERER_CA

print_info "4. Chuyển đổi và tính toán bản cập nhật (Delta)..."
configtxlator proto_decode --input ${TEMP_DIR}/config_block.pb --type common.Block --output ${TEMP_DIR}/config_block.json
jq .data.data[0].payload.data.config ${TEMP_DIR}/config_block.json > ${TEMP_DIR}/config.json

# Gắn NPPTieuChuanOrg vào nhóm Application
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"NPPTieuChuanOrgMSP":.[1]}}}}}' ${TEMP_DIR}/config.json ${TEMP_DIR}/npptieuchuan.json > ${TEMP_DIR}/modified_config.json

configtxlator proto_encode --input ${TEMP_DIR}/config.json --type common.Config --output ${TEMP_DIR}/original_config.pb
configtxlator proto_encode --input ${TEMP_DIR}/modified_config.json --type common.Config --output ${TEMP_DIR}/modified_config.pb
configtxlator compute_update --channel_id giaodich-channel --original ${TEMP_DIR}/original_config.pb --updated ${TEMP_DIR}/modified_config.pb --output ${TEMP_DIR}/config_update.pb

configtxlator proto_decode --input ${TEMP_DIR}/config_update.pb --type common.ConfigUpdate --output ${TEMP_DIR}/config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"giaodich-channel", "type":2}},"data":{"config_update":'$(cat ${TEMP_DIR}/config_update.json)'}}}' | jq . > ${TEMP_DIR}/config_update_in_envelope.json
configtxlator proto_encode --input ${TEMP_DIR}/config_update_in_envelope.json --type common.Envelope --output ${TEMP_DIR}/npptieuchuan_update_in_envelope.pb
print_success "Đã tạo file cấu hình cập nhật: npptieuchuan_update_in_envelope.pb"

print_info "5. Ký và Submit bản cập nhật lên Orderer..."
# Hiện tại có 3 Org trong giaodich-channel (Platform, HTX, NPPXanh). Majority = 2 Orgs.
# Đang đứng ở context PlatformOrg admin, ta ký trước:
peer channel signconfigtx -f ${TEMP_DIR}/npptieuchuan_update_in_envelope.pb

# Đổi sang context HTXNongSanOrg Admin để gửi update (gửi cũng mang ý nghĩa là đã ký)
export CORE_PEER_LOCALMSPID="HTXNongSanOrgMSP"
export CORE_PEER_MSPCONFIGPATH=${CHAINSETUP_DIR}/organizations/peerOrganizations/htxnongsan.gapchain.vn/users/Admin@htxnongsan.gapchain.vn/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${CHAINSETUP_DIR}/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt
export CORE_PEER_ADDRESS=peer0.htxnongsan.gapchain.vn:8051

peer channel update -f ${TEMP_DIR}/npptieuchuan_update_in_envelope.pb -c giaodich-channel -o localhost:7050 --tls --cafile $ORDERER_CA
print_success "Submit cập nhật thành công! giaodich-channel đã có NPPTieuChuanOrg."

print_info "6. Tạo docker-compose chạy Peer mới..."
cat <<EOF > docker-compose-npptieuchuan.yaml
networks:
  gapchain-network:
    external: true
    name: gapchain-network

volumes:
  peer0.npptieuchuan.gapchain.vn:

services:
  peer0.npptieuchuan.gapchain.vn:
    container_name: peer0.npptieuchuan.gapchain.vn
    image: hyperledger/fabric-peer:3.1.1
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=gapchain-network
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_PROFILE_ENABLED=false
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_PEER_ID=peer0.npptieuchuan.gapchain.vn
      - CORE_PEER_ADDRESS=peer0.npptieuchuan.gapchain.vn:11051
      - CORE_PEER_LISTENADDRESS=0.0.0.0:11051
      - CORE_PEER_CHAINCODEADDRESS=peer0.npptieuchuan.gapchain.vn:11052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:11052
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.npptieuchuan.gapchain.vn:11051
      - CORE_PEER_LOCALMSPID=NPPTieuChuanOrgMSP
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb-npptieuchuan:5984
      - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
      - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw
    volumes:
      - /var/run/docker.sock:/host/var/run/docker.sock
      - ./organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/msp:/etc/hyperledger/fabric/msp
      - ./organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/tls:/etc/hyperledger/fabric/tls
      - ./organizations:/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations
      - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
      - peer0.npptieuchuan.gapchain.vn:/var/hyperledger/production
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    ports:
      - 11051:11051
    networks:
      - gapchain-network
    depends_on:
      - couchdb-npptieuchuan

  couchdb-npptieuchuan:
    container_name: couchdb-npptieuchuan
    image: couchdb:3.3
    environment:
      - COUCHDB_USER=admin
      - COUCHDB_PASSWORD=adminpw
    ports:
      - "9984:5984"
    networks:
      - gapchain-network
EOF

docker compose -f docker-compose-npptieuchuan.yaml up -d
print_info "Đợi peer khởi động..."
sleep 5

print_info "7. Lấy genesis block mới nhất & Join peer vào kênh..."

# Lấy genesis block bằng quyền của PlatformOrg (đảm bảo quyền read Block)
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_MSPCONFIGPATH=${CHAINSETUP_DIR}/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${CHAINSETUP_DIR}/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt
export CORE_PEER_ADDRESS=peer0.platform.gapchain.vn:7051

peer channel fetch 0 ${TEMP_DIR}/giaodich-channel.block -o localhost:7050 -c giaodich-channel --tls --cafile $ORDERER_CA

# Switch environment back to NPPTieuChuanOrg để Join channel
export CORE_PEER_LOCALMSPID="NPPTieuChuanOrgMSP"
export CORE_PEER_MSPCONFIGPATH=${CHAINSETUP_DIR}/organizations/peerOrganizations/npptieuchuan.gapchain.vn/users/Admin@npptieuchuan.gapchain.vn/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${CHAINSETUP_DIR}/organizations/peerOrganizations/npptieuchuan.gapchain.vn/peers/peer0.npptieuchuan.gapchain.vn/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:11051
export CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer0.npptieuchuan.gapchain.vn

peer channel join -b ${TEMP_DIR}/giaodich-channel.block

print_success "NPPTieuChuanOrg đã join thành công kênh giaodich-channel!"
print_info "Lưu ý: Để NPPTieuChuanOrg có thể thực thi chaincode, bạn cần đóng gói (package), install, approve và commit lại phân quyền chaincode giao_dich_cc với Endorsement Policy mới (nếu policy có chỉ định cụ thể Org này)."

# Dọn dẹp folder tạp
rm -rf ${TEMP_DIR}
