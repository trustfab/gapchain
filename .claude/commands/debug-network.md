# GAPChain - Debug Network

Hỗ trợ chẩn đoán và sửa lỗi mạng GAPChain.

## Nhiệm vụ

Người dùng gặp lỗi với mạng GAPChain. Hãy chạy các lệnh sau để chẩn đoán, phân tích kết quả, rồi đề xuất fix:

### Bước 1: Kiểm tra tổng quan

```bash
# Xem tất cả containers liên quan
docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "gapchain|orderer|peer|couchdb|ca"
```

### Bước 2: Kiểm tra logs nếu container lỗi

```bash
# Thay <container-name> bằng tên container lỗi
docker logs <container-name> --tail 100 2>&1

# Các container quan trọng:
# - orderer.gapchain.vn
# - peer0.platform.gapchain.vn
# - peer0.htxnongsan.gapchain.vn
# - peer0.chicucbvtv.gapchain.vn
# - peer0.nppxanh.gapchain.vn
# - ca.platform.gapchain.vn (và các CA khác)
# - couchdb.platform (và các CouchDB khác)
```

### Bước 3: Kiểm tra network connectivity

```bash
# Kiểm tra Docker network
docker network ls | grep gapchain
docker network inspect <network-name>
```

### Bước 4: Kiểm tra certificates

```bash
cd chain-setup
# Xem certificates đã được tạo chưa
ls organizations/peerOrganizations/ 2>/dev/null || echo "Chưa có certs - cần chạy cryptogen"
ls organizations/ordererOrganizations/ 2>/dev/null || echo "Chưa có orderer certs"

# Kiểm tra channel artifacts
ls channel-artifacts/ 2>/dev/null || echo "Chưa có channel artifacts"
```

### Bước 5: Kiểm tra channel

```bash
cd chain-setup
export PATH=$PATH:$(pwd)/../bin
export FABRIC_CFG_PATH=$(pwd)/config/platform
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="PlatformMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=organizations/peerOrganizations/platform.gapchain.vn/peers/peer0.platform.gapchain.vn/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=organizations/peerOrganizations/platform.gapchain.vn/users/Admin@platform.gapchain.vn/msp
export CORE_PEER_ADDRESS=localhost:7051

peer channel list
peer channel getinfo -c nhatky-htx-channel
```

## Lỗi phổ biến và cách fix

| Lỗi | Nguyên nhân | Fix |
|-----|-------------|-----|
| Container không start | Port conflict hoặc thiếu certs | Check `docker ps`, chạy cleanup rồi setup lại |
| `connection refused` | Peer/orderer chưa start | Kiểm tra docker ps, xem logs |
| `failed to connect to orderer` | TLS cert path sai | Kiểm tra path trong connection config |
| `chaincode not found` | Chưa deploy hoặc name sai | Chạy deploy-chaincode.sh |
| `endorsement policy failure` | Thiếu endorser org | Kiểm tra số org approve |
| CouchDB connection error | CouchDB container lỗi | `docker restart couchdb.platform` |

## Reset toàn bộ nếu cần

```bash
cd chain-setup
./scripts/cleanup-network.sh
# Chờ cleanup xong
./scripts/setup-network.sh
# Sau đó deploy chaincode
./scripts/deploy-chaincode.sh
```

$ARGUMENTS
