# Hướng dẫn thiết lập mạng GAPChain

## Tổng quan dự án

**GAPChain** là mạng Hyperledger Fabric v3.1.1 quản lý chuỗi cung ứng nông sản cho các HTX Việt Nam.
Mạng bao gồm 4 tổ chức chính:

- **PlatformOrg**: Tổ chức quản lý platform và ứng dụng (`platform.gapchain.vn`, Port: `7051`)
- **HTXNongSanOrg**: Tổ chức HTX Nông Sản (`htxnongsan.gapchain.vn`, Port: `8051`)
- **ChiCucBVTVOrg**: Tổ chức Chi Cục Bảo Vệ Thực Vật (`chicucbvtv.gapchain.vn`, Port: `9051`)
- **NPPXanhOrg**: Tổ chức Nhà Phân Phối Xanh (`nppxanh.gapchain.vn`, Port: `10051`)

## Kiến trúc mạng

### Kênh (Channels) và Chaincode

| Channel | Chaincode | Ngôn ngữ | Chức năng | Thành viên |
|---------|-----------|----------|-----------|-----------|
| `nhatky-htx-channel` | `nhatky_cc` | Go | Nhật ký hoạt động HTX | Platform, HTXNongSan, ChiCucBVTV |
| `giaodich-channel` | `giao_dich_cc` | Go | Quản lý giao dịch | Platform, HTXNongSan, NPPXanh |

### Cấu trúc thư mục

```text
chain-setup/
├── configtx.yaml              # Cấu hình mạng (Topology)
├── crypto-config.yaml         # Cấu hình certificates
├── docker-compose-ca.yaml     # CA servers
├── docker-compose.yaml        # Peers và Orderer
├── chaincode/                 # Source code chaincode Go
│   ├── nhatky_cc/             
│   └── giao_dich_cc/          
├── scripts/                   # Scripts quản lý vòng đời mạng
│   ├── setup-network.sh       # Khởi động toàn bộ mạng
│   ├── cleanup-network.sh     # Dọn dẹp/reset
│   ├── deploy-chaincode.sh    # Đóng gói, cài đặt và commit
│   ├── verify-network.sh      # Kiểm tra trạng thái mạng
│   └── verify-anchor-peers.sh # Kiểm tra cấu hình anchor peers
├── organizations/             # MSP certs (do cryptogen sinh ra)
├── channel-artifacts/         # Genesis blocks trống (sinh tự động)
└── README.md                  # Tài liệu này
```

## Yêu cầu hệ thống

1. **Docker và Docker Compose**
   - Đảm bảo Docker Daemon đang chạy.
   - Các hình ảnh mạng yêu cầu: `hyperledger/fabric-peer:3.1.1`, `hyperledger/fabric-orderer:3.1.1`, `hyperledger/fabric-ca:1.5.15`.
2. **Fabric Binaries (v3.1.1)**
   - Các file thực thi CLI mạng phải nằm trong thư mục `bin/` ngang hàng thư mục `chain-setup/` ở cấp độ thư mục dự án. Để tương tác qua CLI, phải thêm `bin/` vào PATH: `export PATH=$PATH:$(pwd)/../bin`.
   
   *Ví dụ lệnh tải bộ binaries nếu bị thiếu:*
   ```bash
   curl -sSL https://bit.ly/2ysbOFE | bash -s -- 3.1.1 1.5.15
   ```

## Hướng dẫn thiết lập mạng

### Bước 1: Khởi động mạng từ đầu

```bash
# Di chuyển vào thư mục quản lý mạng
cd /path/to/gapchain/chain-setup

# Đảm bảo bin/ nằm trong PATH
export PATH=$PATH:$(pwd)/../bin

# Chạy script tự động thiết lập
./scripts/setup-network.sh
```

Script này sẽ:
- Tạo certificates (bởi cryptogen) vào thư mục `organizations/`
- Tạo genesis block và channel artifacts
- Khởi động mạng blockchain bằng các container Docker
- Tạo các kênh `nhatky-htx-channel` và `giaodich-channel`
- Join peers vào kênh và thiết lập anchor peers

### Bước 2: Triển khai Chaincode

Fabric v3.1.1 áp dụng quy trình deploy 4 bước (Package, Install, Approve, Commit). Tuy nhiên, mọi việc đã được tự động thông qua script:

```bash
# Deploy cả hai chaincode (nhatky_cc & giao_dich_cc)
./scripts/deploy-chaincode.sh
```

**Kiểm tra chaincode đã được commit thành công:**
```bash
export FABRIC_CFG_PATH=$(pwd)/config/platform
peer lifecycle chaincode querycommitted \
  --channelID nhatky-htx-channel \
  --name nhatky_cc \
  --output json
```

### Bước 3: Dọn dẹp & Reset mạng

```bash
# Xóa cấu hình, container và certs
./scripts/cleanup-network.sh

# Hoặc dọn dẹp hoàn toàn (xóa cả images của chaincode cũ)
./scripts/cleanup-network.sh --remove-images
```

## Debug và Kiểm tra mạng

```bash
# Kiểm tra tổng quát số lượng containers
docker ps -a --filter "name=gapchain" --filter "name=orderer" --filter "name=peer"

# Xem logs của Orderer hoặc Peer (vd Peer Platform)
docker logs peer0.platform.gapchain.vn --tail 50
docker logs orderer.gapchain.vn --tail 50
```

## Phát triển Go Client (`fabric-go-client`)

Dự án GAPChain bao gồm một thư mục `fabric-go-client/` hỗ trợ phát triển các ứng dụng Client Go tùy chỉnh tương tác với các Node.
- Hiện tại project cung cấp 2 file demo chính sử dụng **Fabric Gateway SDK**: `htx_ngonsan.go` và `npp_nongsan.go`.
- Code kết nối trực tiếp đến cấu hình bảo mật thông qua đường dẫn chỉ định sẵn (thường trỏ tới thư mục cấu hình `chain-setup/organizations/`).
- Khuyến nghị sử dụng các file này làm mẫu gốc để xây dựng và mở rộng logic ứng dụng mới.

```bash
# Chạy demo Client của HTX Nông Sản
cd ../fabric-go-client
go run htx_ngonsan.go

# Chạy demo Client của Nhà Phân Phối
go run npp_nongsan.go
```

## Gọi giao dịch thử nghiệm sau khi Deploy

Mẫu xuất đường dẫn SSL/TLS thông dụng:
```bash
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
export HTX_TLS_CERT=${PWD}/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt
```

### Chaincode Nhật Ký (`nhatky_cc`)

Tạo mới một nhật ký:
```bash
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 \
  --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C nhatky-htx-channel \
  -n nhatky_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:8051 \
  --tlsRootCertFiles $HTX_TLS_CERT \
  -c '{"function":"TaoNhatKy","Args":["NK003","HTX001","san_xuat","Gieo trồng lúa","Cánh đồng C","kt_vien_001"]}'
```

Gọi đọc (Query) - DocTatCaNhatKy:
```bash
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 \
  --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C nhatky-htx-channel \
  -n nhatky_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:8051 \
  --tlsRootCertFiles $HTX_TLS_CERT \
  -c '{"function":"DocTatCaNhatKy","Args":[]}'
```

### Danh mục chức năng trong Chaincodes

**Nhật ký HTX (`nhatky_cc`):**
1. **TaoNhatKy**: Tạo nhật ký mới
2. **DuyetNhatKy**: Phê duyệt nhật ký
3. **DocNhatKy**: Đọc thông tin nhật ký
4. **DocNhatKyTheoHTX**: Đọc nhật ký theo HTX
5. **DocNhatKyTheoTrangThai**: Đọc nhật ký theo trạng thái
6. **ThongKeNhatKy**: Thống kê nhật ký

**Giao Dịch (`giao_dich_cc`):**
1. **TaoGiaoDich**: Tạo giao dịch mới
2. **DuyetGiaoDich**: Phê duyệt giao dịch
3. **CapNhatTrangThai**: Cập nhật trạng thái giao dịch
4. **DocGiaoDich**: Đọc thông tin giao dịch
5. **DocGiaoDichTheoHTX**: Đọc giao dịch theo HTX
6. **DocGiaoDichTheoNPP**: Đọc giao dịch theo NPP
7. **ThongKeGiaoDich**: Thống kê giao dịch
8. **TinhHoaHong**: Tính hoa hồng cho NPP
9. **DocGiaoDichPhaiThu**: Đọc các giao dịch ở trạng thái "cho_thanh_toan" đối với NPP
