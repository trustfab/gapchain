---
name: gapchain-mvp
description: >
  Hướng dẫn xây dựng MVP hệ thống truy xuất nguồn gốc nông sản GAPChain sử dụng
  Hyperledger Fabric v3.1.1, Backend Golang và Mobile Flutter. Sử dụng skill này bất cứ
  khi nào người dùng muốn: xây dựng hoặc lên kế hoạch hệ thống GAPChain, viết chaincode
  Go cho lohang_cc / nhatky_cc / giao_dich_cc, xây dựng REST API Golang kết nối Fabric
  Gateway SDK, phát triển Flutter app cho HTX ghi nhật ký canh tác, tạo QR Code truy xuất
  sản phẩm, hoặc tương tác với mạng gapchain-network. Trigger khi người dùng đề cập
  "gapchain", "lohang_cc", "nhatky_cc", "giao_dich_cc", "htxnongsan", "lô hàng", "chaincode
  nhật ký", "truy xuất nguồn gốc", "giao dịch nông sản".
---

# Skill: Xây dựng MVP GAPChain — Truy Xuất Nguồn Gốc Nông Sản
## Stack: Hyperledger Fabric 3.1.1 · Golang · Flutter

## Tổng quan hệ thống

**GAPChain** quản lý chuỗi cung ứng nông sản cho các HTX Việt Nam với 2 luồng nghiệp vụ chính:
1. **Nhật ký & Lô hàng** — HTX tạo lô hàng, ghi nhật ký canh tác, Chi Cục BVTV kiểm duyệt, tạo QR truy xuất
2. **Giao dịch** — HTX bán lô hàng cho NPP, Platform quản lý và tính hoa hồng

Kiến trúc gồm 4 lớp:

```
┌─────────────────────────────────────────────────────┐
│  FRONTEND                                           │
│  Flutter (iOS/Android)  │  Vue 3 (Web/QR Portal)    │
└──────────────┬──────────────────────────────────────┘
               │ REST API
┌──────────────▼──────────────────────────────────────┐
│  BACKEND API — Golang (Gin framework)               │
│  Fabric Gateway SDK · JWT Auth                      │
└──────────────┬──────────────────────────────────────┘
               │ Fabric Gateway gRPC
┌──────────────▼──────────────────────────────────────┐
│  HYPERLEDGER FABRIC 3.1.1 — gapchain-network        │
│  Platform │ HTXNongSan │ ChiCucBVTV │ NPPXanh        │
│  nhatky-htx-channel: lohang_cc + nhatky_cc          │
│  giaodich-channel:   giao_dich_cc                   │
└──────────────┬──────────────────────────────────────┘
               │ state db
┌──────────────▼──────────────────────────────────────┐
│  CouchDB (mỗi org 1 instance, rich query)           │
└─────────────────────────────────────────────────────┘
```

---

## Kiến trúc mạng GAPChain

### 4 Tổ chức

| Org | MSP ID | Domain | Peer Port | CouchDB Port |
|-----|--------|--------|-----------|--------------|
| Platform | PlatformOrgMSP | platform.gapchain.vn | 7051 | 5984 |
| HTXNongSan | HTXNongSanOrgMSP | htxnongsan.gapchain.vn | 8051 | 6984 |
| ChiCucBVTV | ChiCucBVTVOrgMSP | chicucbvtv.gapchain.vn | 9051 | 7984 |
| NPPXanh | NPPXanhOrgMSP | nppxanh.gapchain.vn | 10051 | 8984 |
| Orderer | OrdererOrgMSP | gapchain.vn | 7050 | — |

### 2 Channel và 3 Chaincode (MVP)

| Channel | Chaincode | Thành viên | Chức năng |
|---------|-----------|-----------|-----------|
| `nhatky-htx-channel` | `lohang_cc` | Platform, HTXNongSan, ChiCucBVTV | **MỚI** — Lô hàng, chứng nhận, QR truy xuất |
| `nhatky-htx-channel` | `nhatky_cc` | Platform, HTXNongSan, ChiCucBVTV | Nhật ký canh tác gắn với lô hàng |
| `giaodich-channel` | `giao_dich_cc` | Platform, HTXNongSan, NPPXanh | Giao dịch mua bán, hoa hồng NPP |

> **Tại sao thêm `lohang_cc`?** Chaincode demo hiện tại thiếu khái niệm "lô hàng" — không thể liên kết nhật ký canh tác với giao dịch bán hàng, không thể làm QR truy xuất end-to-end. `lohang_cc` là cầu nối giữa nhật ký và giao dịch.

### Cấu trúc thư mục cập nhật

```
gapchain/
├── bin/                           # Fabric binaries v3.1.1
├── chain-setup/
│   ├── chaincode/
│   │   ├── lohang_cc/             # MỚI — Lô hàng & chứng nhận
│   │   ├── nhatky_cc/             # Cập nhật — link lô hàng, MSP check
│   │   └── giao_dich_cc/          # Cập nhật — link lô hàng, DocCongNoNPP
│   ├── scripts/
│   │   ├── setup-network.sh
│   │   ├── deploy-chaincode.sh    # Cập nhật: deploy 3 chaincode
│   │   └── cleanup-network.sh
│   └── ...
└── fabric-go-client/
    ├── htx_ngonsan.go             # Demo HTX (lohang + nhatky)
    └── npp_nongsan.go             # Demo NPP (giaodich)
```

---

## Luồng nghiệp vụ End-to-End

```
HTX tạo lô hàng (lohang_cc)
    ↓
HTX ghi nhật ký canh tác theo lô (nhatky_cc)
    ↓
ChiCucBVTV duyệt nhật ký + cấp chứng nhận (nhatky_cc + lohang_cc)
    ↓
Platform cấp QR Code → Consumer quét → xem toàn bộ hành trình
    ↓
HTX tạo giao dịch bán lô hàng cho NPP (giao_dich_cc, tham chiếu ma_lo)
    ↓
Platform duyệt → NPP xác nhận nhận hàng → thanh toán → tính hoa hồng
```

---

## Bước 1: Hyperledger Fabric 3.1.1 Network

### 1.1 Yêu cầu môi trường

| Thành phần | Phiên bản | Ghi chú |
|------------|-----------|---------|
| Docker Engine | v24+ | `docker --version` |
| Go (chaincode) | 1.19+ | Module path `fabric-chaincode-go/v2` |
| fabric-chaincode-go | v2.x | Bắt buộc cho tất cả chaincode MVP |
| fabric-gateway-go | v1.x | SDK client, thay thế fabric-sdk-go cũ |
| Fabric binaries | 3.1.1 | Trong `bin/`, `export PATH=$PATH:$(pwd)/../bin` |

### 1.2 Khởi động & deploy

```bash
cd gapchain/chain-setup
export PATH=$PATH:$(pwd)/../bin

./scripts/setup-network.sh        # Tạo certs + channels + join peers
./scripts/deploy-chaincode.sh     # Deploy lohang_cc + nhatky_cc + giao_dich_cc

# Verify
docker ps --format "table {{.Names}}\t{{.Status}}"
# 9 containers: 1 orderer + 4 peer + 4 couchdb
```

---

## Bước 2: Chaincode `lohang_cc` — Lô hàng & Chứng nhận

Channel: `nhatky-htx-channel`

### 2.1 Data model

```go
// chain-setup/chaincode/lohang_cc/lohang_cc.go
// Quy ước GAPChain: tên hàm tiếng Việt, JSON tags snake_case

type LotHang struct {
    MaLo          string     `json:"ma_lo"`           // VD: LH-HTX001-2025-001
    MaHTX         string     `json:"ma_htx"`
    TenSanPham    string     `json:"ten_san_pham"`    // VD: Gạo ST25, Cà phê Lâm Đồng
    LoaiSanPham   string     `json:"loai_san_pham"`   // lua | ca_phe | rau | qua...
    SoLuong       float64    `json:"so_luong"`
    DonViTinh     string     `json:"don_vi_tinh"`     // kg | tan | thung
    VuMua         string     `json:"vu_mua"`          // VD: Đông Xuân 2025
    DiaDiem       string     `json:"dia_diem"`        // Tên vùng trồng
    ChungNhan     []ChungNhan `json:"chung_nhan"`     // Chứng nhận đính kèm
    TrangThai     string     `json:"trang_thai"`      // dang_trong | da_thu_hoach | san_sang_ban | da_ban
    NgayTao       string     `json:"ngay_tao"`
    NgayCapNhat   string     `json:"ngay_cap_nhat"`
}

type ChungNhan struct {
    LoaiChungNhan string `json:"loai_chung_nhan"` // VietGAP | GlobalGAP | Organic | TCVN
    MaChungNhan   string `json:"ma_chung_nhan"`
    CoQuanCap     string `json:"co_quan_cap"`     // VD: ChiCucBVTV Cần Thơ
    NgayCap       string `json:"ngay_cap"`
    NgayHetHan    string `json:"ngay_het_han"`
    GhiChu        string `json:"ghi_chu"`
}
```

### 2.2 Các hàm của `lohang_cc`

| Hàm | Mô tả | MSP được phép gọi |
|-----|-------|------------------|
| `TaoLotHang` | Tạo lô hàng mới | HTXNongSanOrgMSP, PlatformOrgMSP |
| `CapNhatTrangThaiLo` | Cập nhật trạng thái lô | HTXNongSanOrgMSP, PlatformOrgMSP |
| `ThemChungNhan` | Thêm chứng nhận VietGAP/GlobalGAP | ChiCucBVTVOrgMSP, PlatformOrgMSP |
| `DocLotHang` | Đọc thông tin lô hàng | Tất cả |
| `DocLotHangTheoHTX` | Query lô hàng theo HTX | Tất cả |
| `DocLotHangTheoTrangThai` | Query theo trạng thái | Tất cả |
| `LichSuLotHang` | Toàn bộ lịch sử thay đổi (GetHistoryForKey) | Tất cả |
| `LayThongTinTraCuu` | Tổng hợp: lô hàng + nhật ký + chứng nhận (dùng cho QR) | Public |

### 2.3 Pattern MSP identity check (áp dụng cho tất cả chaincode MVP)

```go
// Kiểm tra quyền dựa trên MSP — không cần ActorCC riêng
func kiemTraMSP(ctx contractapi.TransactionContextInterface, mspChoPhep ...string) error {
    mspID, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("không thể lấy MSP ID: %v", err)
    }
    for _, m := range mspChoPhep {
        if mspID == m {
            return nil
        }
    }
    return fmt.Errorf("MSP %s không có quyền thực hiện thao tác này", mspID)
}

// Sử dụng trong hàm ThemChungNhan:
func (c *LotHangContract) ThemChungNhan(ctx contractapi.TransactionContextInterface,
    maLo string, loaiChungNhan, maChungNhan, coQuanCap, ngayCap, ngayHetHan string) error {

    // Chỉ ChiCucBVTV và Platform mới được cấp chứng nhận
    if err := kiemTraMSP(ctx, "ChiCucBVTVOrgMSP", "PlatformOrgMSP"); err != nil {
        return err
    }
    // ... xử lý tiếp
}
```

### 2.4 Composite key cho query hiệu quả

```go
// Dùng composite key thay vì plain string key — tránh xung đột giữa các chaincode
// và hỗ trợ GetStateByPartialCompositeKey

// Khi tạo lô hàng:
key, err := ctx.GetStub().CreateCompositeKey("LOHANG", []string{maHTX, maLo})
ctx.GetStub().PutState(key, lotHangJSON)

// Khi query theo HTX:
iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("LOHANG", []string{maHTX})
```

---

## Bước 3: Chaincode `nhatky_cc` — Nhật Ký Canh Tác (cập nhật)

Channel: `nhatky-htx-channel`

### 3.1 Data model (cập nhật — thêm liên kết lô hàng)

```go
// Thêm trường MaLo để liên kết với lohang_cc
type NhatKyHTX struct {
    MaNhatKy     string `json:"ma_nhat_ky"`
    MaLo         string `json:"ma_lo"`          // QUAN TRỌNG: liên kết lô hàng
    MaHTX        string `json:"ma_htx"`
    LoaiHoatDong string `json:"loai_hoat_dong"` // san_xuat | bao_ve_thuc_vat | thu_hoach | van_chuyen
    ChiTiet      string `json:"chi_tiet"`
    ViTri        string `json:"vi_tri"`
    NguoiThucHien string `json:"nguoi_thuc_hien"` // ID kỹ thuật viên
    MinhChungHash string `json:"minh_chung_hash"` // SHA256 của ảnh/file bằng chứng
    TrangThai    string `json:"trang_thai"`       // cho_duyet | da_duyet | tu_choi
    LyDoTuChoi   string `json:"ly_do_tu_choi"`   // Điền khi tu_choi
    NguoiDuyet   string `json:"nguoi_duyet"`
    NgayGhi      string `json:"ngay_ghi"`
    NgayTao      string `json:"ngay_tao"`
    NgayCapNhat  string `json:"ngay_cap_nhat"`
}
```

### 3.2 Các hàm của `nhatky_cc` (cập nhật)

| Hàm | Mô tả | MSP được phép gọi |
|-----|-------|------------------|
| `GhiNhatKy` | Tạo nhật ký mới gắn với lô hàng | HTXNongSanOrgMSP, PlatformOrgMSP |
| `DuyetNhatKy` | Phê duyệt / từ chối | ChiCucBVTVOrgMSP, PlatformOrgMSP |
| `DocNhatKy` | Đọc 1 nhật ký | Tất cả |
| `DocNhatKyTheoLo` | Query nhật ký theo lô hàng | Tất cả |
| `DocNhatKyTheoHTX` | Query nhật ký theo HTX | Tất cả |
| `DocNhatKyTheoTrangThai` | Query theo trạng thái | Tất cả |
| `LichSuNhatKy` | Lịch sử thay đổi của 1 bản ghi (GetHistoryForKey) | Tất cả |
| `ThongKeNhatKy` | Thống kê tổng hợp | Tất cả |

### 3.3 GetHistoryForKey — audit trail thực sự

```go
// LichSuNhatKy trả về toàn bộ lịch sử thay đổi của 1 nhật ký
func (c *NhatKyContract) LichSuNhatKy(ctx contractapi.TransactionContextInterface, maNhatKy string) ([]LichSuRecord, error) {
    iterator, err := ctx.GetStub().GetHistoryForKey(maNhatKy)
    if err != nil {
        return nil, err
    }
    defer iterator.Close()

    var lichSu []LichSuRecord
    for iterator.HasNext() {
        record, err := iterator.Next()
        if err != nil {
            return nil, err
        }
        lichSu = append(lichSu, LichSuRecord{
            TxID:      record.TxId,
            Timestamp: record.Timestamp.String(),
            IsDeleted: record.IsDelete,
            Data:      string(record.Value),
        })
    }
    return lichSu, nil
}
```

---

## Bước 4: Chaincode `giao_dich_cc` — Giao Dịch (cập nhật)

Channel: `giaodich-channel`

### 4.1 Data model (cập nhật — thêm liên kết lô hàng)

```go
type GiaoDich struct {
    MaGiaoDich   string  `json:"ma_giao_dich"`
    MaLo         string  `json:"ma_lo"`          // Liên kết lô hàng từ nhatky-htx-channel
    MaHTX        string  `json:"ma_htx"`
    MaNPP        string  `json:"ma_npp"`
    SanPham      string  `json:"san_pham"`
    SoLuong      float64 `json:"so_luong"`
    DonViTinh    string  `json:"don_vi_tinh"`
    DonGia       float64 `json:"don_gia"`
    TongTien     float64 `json:"tong_tien"`
    TyLeHoaHong  float64 `json:"ty_le_hoa_hong"` // % — thiết lập khi tạo GD
    TienHoaHong  float64 `json:"tien_hoa_hong"`  // Tính sẵn = TongTien * TyLeHoaHong / 100
    TrangThai    string  `json:"trang_thai"`     // cho_duyet | da_duyet | dang_giao | da_giao | cho_thanh_toan | da_thanh_toan | huy_bo
    GhiChu       string  `json:"ghi_chu"`
    NgayTao      string  `json:"ngay_tao"`
    NgayCapNhat  string  `json:"ngay_cap_nhat"`
}
```

### 4.2 Vòng đời trạng thái giao dịch

```
cho_duyet → da_duyet → dang_giao → da_giao → cho_thanh_toan → da_thanh_toan
                                                               ↑
                                                          Kết thúc, tính hoa hồng
Bất kỳ trạng thái nào → huy_bo (chỉ Platform)
```

### 4.3 Các hàm của `giao_dich_cc` (cập nhật)

| Hàm | Mô tả | MSP được phép gọi |
|-----|-------|------------------|
| `TaoGiaoDich` | Tạo giao dịch mới kèm tỷ lệ hoa hồng | HTXNongSanOrgMSP, PlatformOrgMSP |
| `DuyetGiaoDich` | Phê duyệt → `da_duyet` | PlatformOrgMSP |
| `CapNhatTrangThai` | Chuyển trạng thái theo vòng đời | Tùy trạng thái |
| `DocGiaoDich` | Đọc 1 giao dịch | Tất cả |
| `DocGiaoDichTheoHTX` | Query theo HTX | Tất cả |
| `DocGiaoDichTheoNPP` | Query theo NPP | Tất cả |
| `DocCongNoNPP` | **MỚI** — Giao dịch `cho_thanh_toan` của NPP | NPPXanhOrgMSP, PlatformOrgMSP |
| `TinhHoaHongNPP` | Tổng hoa hồng theo NPP và kỳ | PlatformOrgMSP |
| `ThongKeGiaoDich` | Thống kê tổng hợp | Tất cả |

### 4.4 DocCongNoNPP — hàm còn thiếu trong demo

```go
// DocCongNoNPP trả về các giao dịch ở trạng thái "cho_thanh_toan" của 1 NPP
func (c *GiaoDichContract) DocCongNoNPP(ctx contractapi.TransactionContextInterface, maNPP string) ([]*GiaoDich, error) {
    if err := kiemTraMSP(ctx, "NPPXanhOrgMSP", "PlatformOrgMSP"); err != nil {
        return nil, err
    }
    queryString := fmt.Sprintf(`{"selector":{"ma_npp":"%s","trang_thai":"cho_thanh_toan"}}`, maNPP)
    return c.queryGiaoDich(ctx, queryString)
}
```

---

## Bước 5: go.mod chuẩn cho tất cả chaincode MVP

```go
// Áp dụng cho lohang_cc, nhatky_cc, giao_dich_cc
module lohang_cc  // đổi tên theo từng chaincode

go 1.19

require (
    github.com/hyperledger/fabric-contract-api-go/v2 v2.x.x
    github.com/hyperledger/fabric-chaincode-go/v2 v2.x.x
    // KHÔNG dùng v1 — đã deprecated
)
```

---

## Bước 6: Backend API (Golang + Gin)

### 6.0 go.mod cho Backend

```go
module github.com/your-org/gapchain/backend

go 1.21

require (
    github.com/gin-gonic/gin                          v1.10.0
    github.com/hyperledger/fabric-gateway/pkg/client  v1.7.0
    google.golang.org/grpc                            v1.65.0
    // KHÔNG dùng fabric-sdk-go — đã deprecated
)
```

### 6.1 Cấu trúc thư mục (Clean Architecture + fx)

```
backend/
├── cmd/server/main.go          # Entry point — setup uber-go/fx
├── internal/
│   ├── config/
│   │   └── config.go           # Load .env
│   ├── infrastructure/
│   │   └── fabric/
│   │       └── gateway.go      # Kết nối nhatky-htx-channel & giaodich-channel
│   ├── repository/
│   │   └── fabric/             # Giao tiếp Fabric Gateway (Submit/Evaluate)
│   │       ├── lohang_repo.go
│   │       ├── nhatky_repo.go
│   │       └── giaodich_repo.go
│   ├── usecase/                # Logic xử lý/chuyển tiếp
│   │   ├── lohang_uc.go
│   │   ├── nhatky_uc.go
│   │   └── giaodich_uc.go
│   ├── handler/
│   │   └── http/               # Controllers cho Gin
│   │       ├── router.go       # Setup engine và mapping route
│   │       ├── auth.go         # Login endpoint mock
│   │       ├── lohang_handler.go
│   │       ├── nhatky_handler.go
│   │       └── giaodich_handler.go
│   ├── middleware/
│   │   └── auth.go             # JWT middleware
│   └── model/
│       └── dto.go              # Request structs
├── fabric/config/              # Trỏ vào chain-setup/organizations/
├── go.mod
└── .env
```

### 6.2 Kết nối Channel với Infrastructure Layer

```go
// internal/infrastructure/fabric/gateway.go
type GatewayService struct {
    Gw1 *client.Gateway // nhatky-htx-channel
    Gw2 *client.Gateway // giaodich-channel
    LotHangContract  *client.Contract
    NhatKyContract   *client.Contract
    GiaoDichContract *client.Contract
}

func NewFabricGateway(cfg *config.Config) (*GatewayService, error) {
    // Config và trả về instance gateway
}
```

### 6.3 Các API endpoint MVP

**Lô hàng (`lohang_cc`):**

| Method | Endpoint | Chaincode fn | Mô tả |
|--------|----------|--------------|-------|
| POST | `/api/v1/lohang` | TaoLotHang | Tạo lô hàng mới |
| PUT | `/api/v1/lohang/:ma_lo/trangthai` | CapNhatTrangThaiLo | Cập nhật trạng thái |
| POST | `/api/v1/lohang/:ma_lo/chungnhan` | ThemChungNhan | Thêm chứng nhận |
| GET | `/api/v1/lohang/:ma_lo` | DocLotHang | Đọc lô hàng |
| GET | `/api/v1/lohang/:ma_lo/lichsu` | LichSuLotHang | Lịch sử thay đổi |
| GET | `/api/v1/lohang/htx/:ma_htx` | DocLotHangTheoHTX | Lô hàng của HTX |
| GET | `/api/v1/consumer/:ma_lo` | LayThongTinTraCuu | **Public** — QR scan |

**Nhật ký (`nhatky_cc`):**

| Method | Endpoint | Chaincode fn | Mô tả |
|--------|----------|--------------|-------|
| POST | `/api/v1/nhatky` | GhiNhatKy | Ghi nhật ký |
| PUT | `/api/v1/nhatky/:id/duyet` | DuyetNhatKy | Duyệt / từ chối |
| GET | `/api/v1/nhatky/lo/:ma_lo` | DocNhatKyTheoLo | Nhật ký theo lô |
| GET | `/api/v1/nhatky/htx/:ma_htx` | DocNhatKyTheoHTX | Nhật ký theo HTX |
| GET | `/api/v1/nhatky/thongke` | ThongKeNhatKy | Thống kê |

**Giao dịch (`giao_dich_cc`):**

| Method | Endpoint | Chaincode fn | Mô tả |
|--------|----------|--------------|-------|
| POST | `/api/v1/giaodich` | TaoGiaoDich | Tạo giao dịch |
| PUT | `/api/v1/giaodich/:id/duyet` | DuyetGiaoDich | Duyệt |
| PUT | `/api/v1/giaodich/:id/trangthai` | CapNhatTrangThai | Cập nhật trạng thái |
| GET | `/api/v1/giaodich/:id` | DocGiaoDich | Đọc giao dịch |
| GET | `/api/v1/giaodich/npp/:ma_npp/congno` | DocCongNoNPP | Công nợ NPP |
| GET | `/api/v1/giaodich/npp/:ma_npp/hoahong` | TinhHoaHongNPP | Hoa hồng NPP |

---

## Bước 7: Flutter Mobile App (cho HTX)

### 7.1 Cấu trúc thư mục

```
flutter_app/
├── lib/
│   ├── features/
│   │   ├── auth/login_screen.dart
│   │   ├── lohang/
│   │   │   ├── lohang_list_screen.dart
│   │   │   ├── lohang_create_screen.dart  # Tạo lô hàng mới
│   │   │   └── lohang_detail_screen.dart  # Xem lịch sử + nhật ký
│   │   ├── nhatky/
│   │   │   ├── nhatky_form_screen.dart    # Ghi nhật ký theo lô
│   │   │   └── nhatky_list_screen.dart
│   │   └── qr/
│   │       ├── qr_generate_screen.dart    # Tạo QR từ ma_lo
│   │       └── qr_scan_screen.dart        # Quét QR xem truy xuất
│   └── shared/widgets/
│       ├── hoat_dong_picker.dart          # Dropdown loại hoạt động
│       ├── photo_capture.dart             # Chụp ảnh, tính hash SHA256
│       └── trang_thai_badge.dart
├── pubspec.yaml
└── assets/loai_hoat_dong.json
```

### 7.2 Dependencies

```yaml
dependencies:
  dio: ^5.4.0
  sqflite: ^2.3.0               # Offline cache
  flutter_secure_storage: ^9.0.0
  image_picker: ^1.0.7
  qr_flutter: ^4.1.0            # Tạo QR Code từ ma_lo
  qr_code_scanner: ^1.0.1       # Scan QR truy xuất
  geolocator: ^11.0.0
  connectivity_plus: ^6.0.1
  crypto: ^3.0.3                # SHA256 hash ảnh
  riverpod: ^2.5.1
  go_router: ^13.2.0
  intl: ^0.19.0
```

### 7.3 Luồng UX chính cho HTX

```
1. Đăng nhập → Chọn lô hàng (hoặc tạo mới)
2. Trong lô hàng → Ghi nhật ký (chọn loại, chụp ảnh, GPS tự động)
3. Xem lịch sử lô hàng: timeline nhật ký + chứng nhận
4. Tạo QR Code từ lô hàng → Chia sẻ cho người mua
```

### 7.4 Offline-first

```dart
// Khi offline: lưu nhật ký vào SQLite với status = "PENDING", tính hash ảnh local
// Khi online: sync tự động, upload ảnh, submit lên Fabric
class SyncService {
  Future<void> syncPendingNhatKy() async {
    final pending = await localDb.getPendingNhatKy();
    for (final nk in pending) {
      try {
        await apiClient.ghiNhatKy(nk);
        await localDb.markSynced(nk.id);
      } catch (_) { /* Retry lần sau */ }
    }
  }
}
```

**UX principles:**
- Tối đa 3 tap để ghi 1 nhật ký trong lô đang chọn
- Dropdown loại hoạt động, GPS tự động, ảnh bắt buộc
- Offline-first, sync khi có mạng
- Tiếng Việt, font ≥ 16sp

---

## Bước 8: Web Frontend (Vue 3 + Vite)

### 8.1 Dashboard quản lý (Platform)

- **Lô hàng**: Danh sách + filter + phê duyệt chứng nhận
- **Nhật ký**: Danh sách nhật ký chờ duyệt của ChiCucBVTV
- **Giao dịch**: Quản lý vòng đời + duyệt
- **Công nợ NPP**: Bảng cho_thanh_toan theo NPP
- **Hoa hồng**: Tính theo kỳ
- **Xuất báo cáo**: CSV/PDF

### 8.2 QR Consumer Page

```
[Tên sản phẩm + Ảnh]
[Badge chứng nhận: VietGAP ✅ | GlobalGAP ✅]

🏘️ HTX: [Tên HTX] — [Địa phương]
🌾 Lô hàng: LH-HTX001-2025-001
📅 Thu hoạch: 15/03/2025

[Timeline hành trình]
  ✅ 01/01 — Gieo hạt ST25         (HTXNongSan)
  ✅ 15/02 — Bón phân hữu cơ       (HTXNongSan)
  ✅ 10/03 — Kiểm định VietGAP     (ChiCucBVTV — Đã cấp chứng nhận)
  ✅ 20/03 — Thu hoạch              (HTXNongSan)
  ✅ 25/03 — Giao NPP Xanh         (Platform)

[Nút] Xem TX trên Fabric Explorer
[Nút] Tải chứng nhận PDF
```

**Quan trọng**: Load < 3s, không đăng nhập, mobile-first.

---

## Bước 9: Quy trình Deploy MVP

```bash
# ── BƯỚC 1: Mạng ──
cd gapchain/chain-setup
export PATH=$PATH:$(pwd)/../bin
./scripts/setup-network.sh

# ── BƯỚC 2: Deploy 3 chaincode ──
./scripts/deploy-chaincode.sh
# Verify:
peer lifecycle chaincode querycommitted --channelID nhatky-htx-channel --name lohang_cc
peer lifecycle chaincode querycommitted --channelID nhatky-htx-channel --name nhatky_cc
peer lifecycle chaincode querycommitted --channelID giaodich-channel --name giao_dich_cc

# ── BƯỚC 3: Test CLI ──
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/gapchain.vn/orderers/orderer.gapchain.vn/msp/tlscacerts/tlsca.gapchain.vn-cert.pem
export HTX_TLS=${PWD}/organizations/peerOrganizations/htxnongsan.gapchain.vn/peers/peer0.htxnongsan.gapchain.vn/tls/ca.crt

# Tạo lô hàng
peer chaincode invoke \
  -o orderer.gapchain.vn:7050 --ordererTLSHostnameOverride orderer.gapchain.vn \
  --tls --cafile $ORDERER_CA \
  -C nhatky-htx-channel -n lohang_cc \
  --peerAddresses peer0.htxnongsan.gapchain.vn:8051 --tlsRootCertFiles $HTX_TLS \
  -c '{"function":"TaoLotHang","Args":["LH-HTX001-2025-001","HTX001","Gạo ST25","lua","1000","kg","Đông Xuân 2025","Cần Thơ"]}'

# Ghi nhật ký theo lô
peer chaincode invoke ... -n nhatky_cc \
  -c '{"function":"GhiNhatKy","Args":["NK001","LH-HTX001-2025-001","HTX001","san_xuat","Gieo hạt ST25","Cánh đồng A","kt001",""]}'

# ── BƯỚC 4: Chạy Go client demo ──
cd ../fabric-go-client
go run htx_ngonsan.go   # Demo HTX: tạo lô + ghi nhật ký
go run npp_nongsan.go   # Demo NPP: xem giao dịch, công nợ

# ── BƯỚC 5: Backend API ──
cd ../backend
go mod download
go run cmd/server/main.go

# ── BƯỚC 6: Flutter ──
cd ../flutter_app
flutter pub get && flutter run
```

---

## Checklist MVP

### Fabric Network
- [ ] 9 containers running (orderer + 4 peer + 4 couchdb)
- [ ] `nhatky-htx-channel`: `lohang_cc` và `nhatky_cc` đã commit
- [ ] `giaodich-channel`: `giao_dich_cc` đã commit
- [ ] Anchor peers thiết lập đúng 2 channel

### Chaincode
- [ ] `go.mod` dùng `fabric-chaincode-go/v2` và `fabric-contract-api-go/v2`
- [ ] MSP identity check trong TaoLotHang, ThemChungNhan, DuyetNhatKy, DuyetGiaoDich, DocCongNoNPP
- [ ] Composite key cho tất cả PutState (tránh xung đột key)
- [ ] `GetHistoryForKey` hoạt động trong LichSuLotHang và LichSuNhatKy
- [ ] `lohang_cc`: TaoLotHang, ThemChungNhan, LayThongTinTraCuu (QR) hoạt động
- [ ] `nhatky_cc`: GhiNhatKy, DuyetNhatKy, DocNhatKyTheoLo hoạt động
- [ ] `giao_dich_cc`: TaoGiaoDich, DocCongNoNPP, TinhHoaHongNPP hoạt động
- [ ] Unit test coverage > 70%

### Backend Golang
- [ ] `go.mod` dùng `fabric-gateway-go v1.x`
- [ ] Kết nối 2 channel (Channel1Service + Channel2Service)
- [ ] Tất cả endpoint lô hàng, nhật ký, giao dịch hoạt động
- [ ] Consumer endpoint (QR) public — không cần JWT
- [ ] JWT middleware cho tất cả endpoint còn lại

### Flutter & Web
- [ ] Tạo lô hàng → ghi nhật ký → xem timeline
- [ ] Tạo QR Code từ ma_lo, scan và xem truy xuất
- [ ] Offline-first nhật ký, sync khi có mạng
- [ ] Dashboard Platform: duyệt nhật ký + giao dịch + công nợ NPP

---

## Lỗi thường gặp & cách xử lý

| Vấn đề | Nguyên nhân | Giải pháp |
|--------|-------------|-----------|
| MSP không có quyền | Gọi hàm với org sai | Kiểm tra bảng MSP cho phép của từng hàm |
| Composite key không tìm thấy | Dùng GetState thay vì GetStateByPartialCompositeKey | Đổi sang `GetStateByPartialCompositeKey("LOHANG", []string{maHTX})` |
| Channel not found | Sai tên channel | `nhatky-htx-channel` hoặc `giaodich-channel` |
| MSP ID mismatch | Gọi nhầm MSP | Xem bảng: Platform=`PlatformOrgMSP`, HTX=`HTXNongSanOrgMSP`, BVTV=`ChiCucBVTVOrgMSP`, NPP=`NPPXanhOrgMSP` |
| GetHistoryForKey rỗng | CouchDB chưa bật history | Thêm `CORE_LEDGER_STATE_HISTORY_ENABLED=true` vào peer env |
| lohang_cc và nhatky_cc key trùng | Dùng plain key thay composite | Dùng `CreateCompositeKey` với prefix riêng |
| DocCongNoNPP trả về rỗng | Giao dịch chưa ở `cho_thanh_toan` | Kiểm tra vòng đời trạng thái |
| `cannot find package fabric-protos-go` | Module v1 | Đổi sang `fabric-protos-go-apiv2` |
| Peer crash sau deploy | Thiếu endorsement đủ org | Kiểm tra endorsement policy trong deploy-chaincode.sh |
| CouchDB connection refused | Container chưa ready | Đợi 5-10s sau `docker-compose up` |

---

## Tham khảo

- [chain-setup/README.md](gapchain/chain-setup/README.md) — Hướng dẫn vận hành mạng
- [CLAUDE.md](gapchain/CLAUDE.md) — Context dự án cho Claude Code
- [chain-setup/chaincode/nhatky_cc/](gapchain/chain-setup/chaincode/nhatky_cc/) — Chaincode nhật ký (cần cập nhật theo MVP)
- [chain-setup/chaincode/giao_dich_cc/](gapchain/chain-setup/chaincode/giao_dich_cc/) — Chaincode giao dịch (cần cập nhật theo MVP)
- [fabric-go-client/htx_ngonsan.go](gapchain/fabric-go-client/htx_ngonsan.go) — Mẫu Go client cho HTX
- [fabric-go-client/npp_nongsan.go](gapchain/fabric-go-client/npp_nongsan.go) — Mẫu Go client cho NPP
