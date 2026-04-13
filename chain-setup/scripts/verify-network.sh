#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${DIR}/.."

export FABRIC_CFG_PATH=${PWD}/config/platform.gapchain.vn
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="PlatformOrgMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "[INFO] Verifying nhatky_cc on nhatky-htx-channel..."
docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp \
  peer0.platform.gapchain.vn peer chaincode invoke -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem -C nhatky-htx-channel -n nhatky_cc -c '{"function":"InitLedger","Args":[]}'

echo "[INFO] Waiting 3s for transaction commit..."
sleep 3

echo "[INFO] Querying nhatky_cc..."
docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp \
  peer0.platform.gapchain.vn peer chaincode query -C nhatky-htx-channel -n nhatky_cc -c '{"function":"DocTatCaNhatKy","Args":[]}'

echo -e "\n[INFO] Verifying giao_dich_cc on giaodich-channel..."
docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp \
  peer0.platform.gapchain.vn peer chaincode invoke -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem -C giaodich-channel -n giao_dich_cc -c '{"function":"InitLedger","Args":[]}'

echo "[INFO] Waiting 3s for transaction commit..."
sleep 3

echo "[INFO] Querying giao_dich_cc..."
docker exec -e CORE_PEER_LOCALMSPID="PlatformOrgMSP" \
  -e CORE_PEER_TLS_ENABLED=true \
  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt \
  -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp \
  peer0.platform.gapchain.vn peer chaincode query -C giaodich-channel -n giao_dich_cc -c '{"function":"DocTatCaGiaoDich","Args":[]}'

echo -e "\n[SUCCESS] Network verification complete!"
