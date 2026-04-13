# GAPChain - Claude Code Context

## Tổng quan dự án

**GAPChain** là mạng Hyperledger Fabric v3.1.1 quản lý chuỗi cung ứng nông sản cho các HTX Việt Nam. Hệ thống phục vụ 2 luồng nghiệp vụ chính: Nhật ký & Lô hàng (truy xuất nguồn gốc) và Giao dịch mua bán giữa HTX và Nhà phân phối.

## Cấu trúc

```
gapchain/
├── bin/                    - Fabric binaries v3.1.1 (không edit)
├── chain-setup/            - Cấu hình mạng và chaincode
│   ├── chaincode/
│   │   ├── lohang_cc/     - (MỚI) Chaincode lô hàng & chứng nhận
│   │   ├── nhatky_cc/     - Chaincode nhật ký hoạt động HTX gắn với lô hàng
│   │   └── giao_dich_cc/  - Chaincode giao dịch mua bán, tính hoa hồng, công nợ
│   ├── scripts/            - Scripts setup/deploy/cleanup
│   ├── organizations/      - MSP certs (generated, không commit)
│   ├── channel-artifacts/  - Genesis blocks (generated)
│   ├── configtx.yaml       - Topology mạng
│   ├── crypto-config.yaml  - Cấu hình certificates
│   └── docker-compose.yaml - Container orchestration
├── fabric-go-client/       - Go SDK client application demo
│   ├── htx_ngonsan.go      - Demo HTX (lohang + nhatky)
│   └── npp_nongsan.go      - Demo NPP (giaodich)
├── backend/                - Backend API (Golang + Gin + Fabric Gateway v1.x)
└── flutter_app/            - App mobile cho HTX (Flutter)
```

## Tổ chức & MSP

| Org | MSP ID | Domain | Peer Port | CouchDB Port |
|-----|--------|--------|-----------|--------------|
| Platform | PlatformOrgMSP | platform.gapchain.vn | 7051 | 5984 |
| HTXNongSan | HTXNongSanOrgMSP | htxnongsan.gapchain.vn | 8051 | 6984 |
| ChiCucBVTV | ChiCucBVTVOrgMSP | chicucbvtv.gapchain.vn | 9051 | 7984 |
| NPPXanh | NPPXanhOrgMSP | nppxanh.gapchain.vn | 10051 | 8984 |
| Orderer | OrdererOrgMSP | gapchain.vn | 7050 | — |

## Channel & Chaincode

| Channel | Chaincode | Thành viên | Chức năng |
|---------|-----------|-----------|-----------|
| nhatky-htx-channel | lohang_cc | Platform, HTXNongSan, ChiCucBVTV | Quản lý Lô hàng, chứng nhận, dùng làm QR truy xuất |
| nhatky-htx-channel | nhatky_cc | Platform, HTXNongSan, ChiCucBVTV | Nhật ký canh tác có tham chiếu đến mã lô hàng |
| giaodich-channel | giao_dich_cc | Platform, HTXNongSan, NPPXanh | Quản lý Giao dịch, tính hoa hồng và công nợ NPP |

## Commands thường dùng

```bash
# Setup từ đầu
cd chain-setup && ./scripts/setup-network.sh

# Deploy chaincode (lohang_cc, nhatky_cc, giao_dich_cc)
cd chain-setup && ./scripts/deploy-chaincode.sh

# Cleanup
cd chain-setup && ./scripts/cleanup-network.sh

# Run Go client (Demo CLI)
cd fabric-go-client
go run htx_ngonsan.go   # Demo HTX: tạo lô + ghi nhật ký
go run npp_nongsan.go   # Demo NPP: xem giao dịch, công nợ

# Run Backend API
cd backend && go run cmd/server/main.go

# Run Flutter App
cd flutter_app && flutter run

# Kiểm tra containers
docker ps --format "table {{.Names}}\t{{.Status}}"
```

## Quy tắc kỹ thuật chaincode

- Fabric binaries ở `bin/` — thêm vào PATH: `export PATH=$PATH:$(pwd)/../bin`
- Truy cập API Gateway bằng `fabric-gateway-go` (`v1.x`), KHÔNG dùng `fabric-sdk-go` cũ.
- Chaincode viết bằng `fabric-contract-api-go v1.2.2`, Go 1.19+, key lưu dạng plain string (`maLo`, `maNhatKy`, `maGiaoDich`).
- Tên hàm chaincode tiếng Việt (`TaoLotHang`, `DocGiaoDich`), struct field và JSON tags kiểu `snake_case`.
- **MSP per-transition**: Chaincode check MSP theo từng trạng thái đích (`kiemTraMSP()`). VD: `dinh_chi` lô hàng chỉ `PlatformOrgMSP` + `ChiCucBVTVOrgMSP`; phục hồi từ `dinh_chi` chỉ `PlatformOrgMSP`.
- **Lifecycle enforce**: `lohang_cc` và `giao_dich_cc` dùng `chuyenTrangThaiHopLe` map — không cho "nhảy cóc" trạng thái.
- **Timestamp deterministic**: Dùng `GetTxTimestamp()` (qua helper `txNow()`), KHÔNG dùng `time.Now()`.
- **Danh mục động**: `loai_hoat_dong` trong `nhatky_cc` KHÔNG hardcode enum — chaincode chỉ check `!= ""`, backend validate theo danh mục.
- **Audit Trail**: `GetHistoryForKey()` cho `LichSuLotHang`, `LichSuNhatKy`, `LichSuGiaoDich`.
- **CouchDB sort**: Mọi query có sort PHẢI thêm `"field":{"$gt":null}` để CouchDB index đúng.

## Trạng thái chính (tham chiếu BA skill để có chi tiết đầy đủ)

| Entity | Trạng thái | Terminal |
|--------|-----------|----------|
| lohang_cc | `dang_trong` → `da_thu_hoach` → `cho_chung_nhan` → `san_sang_ban` → `het_hang` | `het_hang` |
| lohang_cc | Bất kỳ (trừ het_hang) → `dinh_chi` → phục hồi | `dinh_chi` (tạm) |
| nhatky_cc | `cho_duyet` → `da_duyet` hoặc `tu_choi` | `da_duyet`, `tu_choi` |
| giao_dich_cc | `cho_duyet` → `da_duyet` → `dang_giao` → `da_giao` → `cho_thanh_toan` → `da_thanh_toan` | `da_thanh_toan`, `huy_bo` |

## Skills & Commands

**Skills (auto-trigger theo context):**
- `gapchain-mvp` — Kiến trúc kỹ thuật MVP: stack, data model, API endpoints, directory structure
- `gapchain-ba` — **Single Source of Truth nghiệp vụ**: vòng đời trạng thái (6 trạng thái lô hàng, 3 nhật ký, 7 giao dịch), ma trận quyền hạn MSP, cross-entity rules R1-R7, frontend display rules. **BẮT BUỘC tham chiếu khi viết/review code**

**Commands (gọi thủ công):**
- `/build-backend-api` - Xây dựng Backend API (Golang + Gin + Fabric Gateway, multi-org)
- `/build-flutter-htx` - Phát triển Flutter app cho HTX (offline-first, ghi nhật ký, QR)
- `/build-frontend-multi-tenant` - Xây dựng Frontend multi-tenant (Vue 3 + Vite, dashboard Platform/HTX/BVTV/NPP)
- `/deploy-chaincode` - Deploy 3 chaincode (lohang → nhatky → giaodich)
- `/chaincode-dev` - Phát triển chaincode
- `/go-client-dev` - Phát triển Go client demo
- `/install-fabric` - Cài đặt môi trường Fabric v3.x
- `/setup-network` - Khởi tạo mạng
- `/debug-network` - Debug lỗi mạng
