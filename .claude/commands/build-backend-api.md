# GAPChain - Build Backend API

Hỗ trợ xây dựng Backend API Golang (Gin) kết nối Hyperledger Fabric Gateway SDK cho hệ thống GAPChain.
Dự án được cấu trúc theo chuẩn **Clean Architecture** và sử dụng **`go.uber.org/fx`** để quản lý Dependency Injection.

## Kiến trúc Multi-Org Gateway

**Nguyên tắc cốt lõi**: Mỗi tổ chức (Org) sử dụng bộ kết nối Fabric Gateway riêng với identity (cert/key) của chính tổ chức đó. Khi user đăng nhập, JWT chứa `msp_id` → middleware resolve → handler/usecase/repo sử dụng đúng Gateway của org tương ứng.

```text
Flutter App / Vue 3 Web
        ↓ REST API (JWT chứa msp_id)
  [ JWT Middleware ]                    ← Extract msp_id → set vào gin.Context
        ↓
  [ HTTP Handler ]                     ← Lấy msp_id từ context, gọi GatewayRegistry
        ↓
  [ GatewayRegistry ]                  ← Map msp_id → OrgGateway tương ứng
        ↓                                 (khởi tạo lúc startup, cache sẵn)
  [ OrgGateway ]                       ← Chứa Contracts cho org cụ thể
        ↓ gRPC (dùng identity + peer của org đó)
  Hyperledger Fabric 3.1.1
```

### Flow chi tiết:

1. **Startup**: `GatewayRegistry` đọc config tất cả orgs, khởi tạo Gateway cho từng org → cache trong map
2. **Request**: JWT middleware extract `msp_id` → set vào `c.Set("mspID", ...)`
3. **Handler**: Lấy `mspID` từ context → gọi `registry.GetOrgGateway(mspID)` → lấy Contract tương ứng
4. **Submit/Evaluate**: Transaction được ký bằng identity của org đó → Fabric endorsement đúng policy

### Ý nghĩa:
- **HTXNongSan** user → dùng cert/key của HTXNongSan → kết nối qua peer `peer0.htxnongsan.gapchain.vn:8051`
- **ChiCucBVTV** user → dùng cert/key của ChiCucBVTV → kết nối qua peer `peer0.chicucbvtv.gapchain.vn:9051`  
- **NPPXanh** user → dùng cert/key của NPPXanh → kết nối qua peer `peer0.nppxanh.gapchain.vn:10051`
- **Platform** user → dùng cert/key của Platform → kết nối qua peer `peer0.platform.gapchain.vn:7051`

## Cấu trúc thư mục

```
gapchain/backend/
├── cmd/server/
│   └── main.go                 # Entry point: fx.New(fx.Provide(...), fx.Invoke(...))
├── internal/
│   ├── config/
│   │   └── config.go           # Cấu hình đa org từ .env
│   ├── infrastructure/
│   │   └── fabric/
│   │       └── gateway.go      # GatewayRegistry: quản lý multi-org connections
│   ├── repository/
│   │   └── fabric/             # Repositories nhận Contract từ handler (không giữ state)
│   │       ├── lohang_repo.go
│   │       ├── nhatky_repo.go
│   │       └── giaodich_repo.go
│   ├── usecase/                # Usecases nhận mspID, resolve Contract qua Registry
│   │   ├── lohang_uc.go
│   │   ├── nhatky_uc.go
│   │   └── giaodich_uc.go
│   ├── handler/
│   │   └── http/
│   │       ├── router.go
│   │       ├── auth_handler.go
│   │       ├── lohang_handler.go
│   │       ├── nhatky_handler.go
│   │       └── giaodich_handler.go
│   ├── middleware/
│   │   └── auth.go             # JWT middleware: extract msp_id vào context
│   └── model/
│       └── dto.go
├── go.mod
└── .env
```

## Dependencies (go.mod)

```go
module github.com/your-org/gapchain/backend

go 1.21

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/hyperledger/fabric-gateway v1.7.0
    github.com/joho/godotenv v1.5.1
    go.uber.org/fx v1.22.2
    google.golang.org/grpc v1.65.0
)
```

## Config Multi-Org (config.go)

```go
package config

type OrgConfig struct {
    MspID       string
    CertPath    string // Path tới cert PEM của org
    KeyPath     string // Path tới private key directory của org
    PeerEndpoint string // vd: localhost:7051
    TLSCertPath string // TLS cert của peer
    GatewayPeer string // vd: peer0.platform.gapchain.vn
    // Org tham gia channel nào
    Channels    []string // vd: ["nhatky-htx-channel"] hoặc ["nhatky-htx-channel", "giaodich-channel"]
}

type Config struct {
    Port    string
    GinMode string
    Orgs    map[string]*OrgConfig // key = MspID
}
```

Mỗi org được cấu hình qua biến môi trường theo pattern:

```bash
# .env
PORT=8080
GIN_MODE=debug
JWT_SECRET=your-secret-key-here

# Danh sách org (comma-separated)
FABRIC_ORGS=PlatformOrgMSP,HTXNongSanOrgMSP,ChiCucBVTVOrgMSP,NPPXanhOrgMSP

# Platform
FABRIC_PlatformOrgMSP_CERT_PATH=./fabric/platform/cert.pem
FABRIC_PlatformOrgMSP_KEY_PATH=./fabric/platform/keystore
FABRIC_PlatformOrgMSP_PEER_ENDPOINT=localhost:7051
FABRIC_PlatformOrgMSP_TLS_CERT=./fabric/platform/tls-cert.pem
FABRIC_PlatformOrgMSP_GATEWAY_PEER=peer0.platform.gapchain.vn
FABRIC_PlatformOrgMSP_CHANNELS=nhatky-htx-channel,giaodich-channel

# HTXNongSan
FABRIC_HTXNongSanOrgMSP_CERT_PATH=./fabric/htxnongsan/cert.pem
FABRIC_HTXNongSanOrgMSP_KEY_PATH=./fabric/htxnongsan/keystore
FABRIC_HTXNongSanOrgMSP_PEER_ENDPOINT=localhost:8051
FABRIC_HTXNongSanOrgMSP_TLS_CERT=./fabric/htxnongsan/tls-cert.pem
FABRIC_HTXNongSanOrgMSP_GATEWAY_PEER=peer0.htxnongsan.gapchain.vn
FABRIC_HTXNongSanOrgMSP_CHANNELS=nhatky-htx-channel,giaodich-channel

# ChiCucBVTV
FABRIC_ChiCucBVTVOrgMSP_CERT_PATH=./fabric/chicucbvtv/cert.pem
FABRIC_ChiCucBVTVOrgMSP_KEY_PATH=./fabric/chicucbvtv/keystore
FABRIC_ChiCucBVTVOrgMSP_PEER_ENDPOINT=localhost:9051
FABRIC_ChiCucBVTVOrgMSP_TLS_CERT=./fabric/chicucbvtv/tls-cert.pem
FABRIC_ChiCucBVTVOrgMSP_GATEWAY_PEER=peer0.chicucbvtv.gapchain.vn
FABRIC_ChiCucBVTVOrgMSP_CHANNELS=nhatky-htx-channel

# NPPXanh
FABRIC_NPPXanhOrgMSP_CERT_PATH=./fabric/nppxanh/cert.pem
FABRIC_NPPXanhOrgMSP_KEY_PATH=./fabric/nppxanh/keystore
FABRIC_NPPXanhOrgMSP_PEER_ENDPOINT=localhost:10051
FABRIC_NPPXanhOrgMSP_TLS_CERT=./fabric/nppxanh/tls-cert.pem
FABRIC_NPPXanhOrgMSP_GATEWAY_PEER=peer0.nppxanh.gapchain.vn
FABRIC_NPPXanhOrgMSP_CHANNELS=giaodich-channel
```

## GatewayRegistry (gateway.go) — Core thay đổi

```go
package fabric

import (
    "crypto/x509"
    "fmt"
    "os"
    "sync"

    "github.com/hyperledger/fabric-gateway/pkg/client"
    "github.com/hyperledger/fabric-gateway/pkg/identity"
    "github.com/your-org/gapchain/backend/internal/config"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

// OrgGateway chứa tất cả contracts mà một org có quyền truy cập
type OrgGateway struct {
    Gateway          *client.Gateway
    LotHangContract  *client.Contract // nil nếu org không tham gia nhatky-htx-channel
    NhatKyContract   *client.Contract // nil nếu org không tham gia nhatky-htx-channel
    GiaoDichContract *client.Contract // nil nếu org không tham gia giaodich-channel
}

// GatewayRegistry quản lý Gateway cho tất cả orgs, khởi tạo lúc startup
type GatewayRegistry struct {
    mu       sync.RWMutex
    gateways map[string]*OrgGateway // key = MspID
    // fallbackMspID dùng cho public endpoint (consumer tra cứu)
    FallbackMspID string
}

func NewGatewayRegistry(cfg *config.Config) (*GatewayRegistry, error) {
    registry := &GatewayRegistry{
        gateways:      make(map[string]*OrgGateway),
        FallbackMspID: "PlatformOrgMSP", // Public queries dùng Platform
    }

    for mspID, orgCfg := range cfg.Orgs {
        orgGw, err := connectOrg(mspID, orgCfg)
        if err != nil {
            return nil, fmt.Errorf("khong the ket noi org %s: %w", mspID, err)
        }
        registry.gateways[mspID] = orgGw
    }

    return registry, nil
}

// GetOrgGateway trả về OrgGateway theo mspID, error nếu không tồn tại
func (r *GatewayRegistry) GetOrgGateway(mspID string) (*OrgGateway, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    gw, ok := r.gateways[mspID]
    if !ok {
        return nil, fmt.Errorf("khong tim thay gateway cho org: %s", mspID)
    }
    return gw, nil
}

// GetFallbackGateway trả về gateway mặc định (Platform) cho public endpoints
func (r *GatewayRegistry) GetFallbackGateway() (*OrgGateway, error) {
    return r.GetOrgGateway(r.FallbackMspID)
}

// CloseAll đóng tất cả gateway connections
func (r *GatewayRegistry) CloseAll() {
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, gw := range r.gateways {
        if gw.Gateway != nil {
            gw.Gateway.Close()
        }
    }
}

func connectOrg(mspID string, orgCfg *config.OrgConfig) (*OrgGateway, error) {
    id, err := newIdentity(mspID, orgCfg.CertPath)
    if err != nil {
        return nil, fmt.Errorf("identity %s: %w", mspID, err)
    }
    sign, err := newSigner(orgCfg.KeyPath)
    if err != nil {
        return nil, fmt.Errorf("signer %s: %w", mspID, err)
    }

    tlsCert, err := os.ReadFile(orgCfg.TLSCertPath)
    if err != nil {
        return nil, fmt.Errorf("tls cert %s: %w", mspID, err)
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(tlsCert)
    creds := credentials.NewClientTLSFromCert(certPool, orgCfg.GatewayPeer)

    conn, err := grpc.NewClient(orgCfg.PeerEndpoint, grpc.WithTransportCredentials(creds))
    if err != nil {
        return nil, fmt.Errorf("grpc %s: %w", mspID, err)
    }

    gw, err := client.Connect(id, client.WithSign(sign), client.WithClientConnection(conn))
    if err != nil {
        return nil, err
    }

    orgGw := &OrgGateway{Gateway: gw}

    // Kết nối channels mà org tham gia
    channelSet := make(map[string]bool)
    for _, ch := range orgCfg.Channels {
        channelSet[ch] = true
    }

    if channelSet["nhatky-htx-channel"] {
        network := gw.GetNetwork("nhatky-htx-channel")
        orgGw.LotHangContract = network.GetContract("lohang_cc")
        orgGw.NhatKyContract = network.GetContract("nhatky_cc")
    }

    if channelSet["giaodich-channel"] {
        network := gw.GetNetwork("giaodich-channel")
        orgGw.GiaoDichContract = network.GetContract("giao_dich_cc")
    }

    return orgGw, nil
}

// newIdentity, newSigner giữ nguyên logic cũ
```

## Pattern Handler: Lấy mspID từ context

Thay đổi quan trọng: Handler không inject trực tiếp repo nữa, mà inject `GatewayRegistry` để resolve contract theo user.

```go
package http

type LohangHandler struct {
    registry *fabric.GatewayRegistry
}

func NewLohangHandler(registry *fabric.GatewayRegistry) *LohangHandler {
    return &LohangHandler{registry: registry}
}

func (h *LohangHandler) TaoLotHang(c *gin.Context) {
    // 1. Lấy mspID từ JWT context (đã set bởi middleware)
    mspID, _ := c.Get("mspID")

    // 2. Resolve gateway cho org này
    orgGw, err := h.registry.GetOrgGateway(mspID.(string))
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Org khong co quyen truy cap"})
        return
    }

    // 3. Kiểm tra org có contract cần thiết
    if orgGw.LotHangContract == nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Org khong tham gia channel nhatky-htx"})
        return
    }

    // 4. Gọi transaction với identity của org
    var req model.TaoLotHangReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err = orgGw.LotHangContract.SubmitTransaction("TaoLotHang",
        req.MaLo, req.MaHTX, req.TenSanPham, req.LoaiSanPham,
        fmt.Sprintf("%g", req.SoLuong), req.DonViTinh, req.VuMua, req.DiaDiem,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"ma_lo": req.MaLo, "message": "Tao lo hang thanh cong"})
}

// Public endpoint: dùng FallbackGateway (Platform)
func (h *LohangHandler) LayThongTinTraCuu(c *gin.Context) {
    orgGw, err := h.registry.GetFallbackGateway()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "He thong khong san sang"})
        return
    }
    maLo := c.Param("ma_lo")
    result, err := orgGw.LotHangContract.EvaluateTransaction("LayThongTinTraCuu", maLo)
    // ...
}
```

## Pattern với Usecase Layer (tuỳ chọn)

Nếu muốn giữ Usecase layer, có 2 cách:

### Cách 1: Handler resolve Contract, truyền vào Usecase method

```go
// Usecase nhận Contract thay vì Repo
type LohangUsecase interface {
    TaoLotHang(contract *client.Contract, req *model.TaoLotHangReq) (string, error)
}

// Handler
func (h *LohangHandler) TaoLotHang(c *gin.Context) {
    orgGw, _ := h.registry.GetOrgGateway(c.GetString("mspID"))
    result, err := h.uc.TaoLotHang(orgGw.LotHangContract, &req)
}
```

### Cách 2: Usecase nhận Registry, tự resolve (khuyến nghị cho MVP)

```go
type lohangUsecase struct {
    registry *fabric.GatewayRegistry
}

func (uc *lohangUsecase) TaoLotHang(mspID string, req *model.TaoLotHangReq) (string, error) {
    orgGw, err := uc.registry.GetOrgGateway(mspID)
    if err != nil {
        return "", err
    }
    if orgGw.LotHangContract == nil {
        return "", fmt.Errorf("org %s khong co quyen tao lo hang", mspID)
    }
    _, err = orgGw.LotHangContract.SubmitTransaction("TaoLotHang", ...)
    return req.MaLo, err
}
```

## Pattern DI với Fx (main.go) — Cập nhật

```go
func main() {
    fx.New(
        fx.Provide(
            config.LoadConfig,
            infra.NewGatewayRegistry, // THAY ĐỔI: Registry thay vì single Gateway

            // Handlers inject trực tiếp Registry
            httpserver.NewLohangHandler,
            httpserver.NewNhatkyHandler,
            httpserver.NewGiaodichHandler,
            httpserver.NewAuthHandler,
        ),
        fx.Invoke(httpserver.SetupRouter),
    ).Run()
}
```

**Lưu ý**: Nếu chọn cách đơn giản (handler gọi trực tiếp contract), có thể bỏ layer usecase và repository cho MVP. Nếu giữ usecase, usecase nhận `*GatewayRegistry` thay vì repo cũ.

## Router cập nhật (router.go)

```go
func SetupRouter(
    lc fx.Lifecycle,
    cfg *config.Config,
    registry *fabric.GatewayRegistry, // THAY ĐỔI
    lohangH *LohangHandler,
    nhatkyH *NhatkyHandler,
    giaodichH *GiaodichHandler,
    authH *AuthHandler,
) *gin.Engine {
    // ... routes giữ nguyên ...
    
    lc.Append(fx.Hook{
        OnStop: func(context.Context) error {
            registry.CloseAll() // Đóng tất cả org gateways
            return nil
        },
    })
    return r
}
```

## Bảng Channel Access theo Org

| Org MSP | nhatky-htx-channel | giaodich-channel | Contracts |
|---------|-------------------|-----------------|-----------|
| PlatformOrgMSP | ✅ | ✅ | lohang_cc, nhatky_cc, giao_dich_cc |
| HTXNongSanOrgMSP | ✅ | ✅ | lohang_cc, nhatky_cc, giao_dich_cc |
| ChiCucBVTVOrgMSP | ✅ | ❌ | lohang_cc, nhatky_cc |
| NPPXanhOrgMSP | ❌ | ✅ | giao_dich_cc |

## Bảng đầy đủ API endpoint MVP

**Lô hàng (`lohang_cc` — nhatky-htx-channel):**

| Method | Endpoint | UseCase Fn | Submit/Eval | Org cho phép |
|--------|----------|-------------|-------------|-------------|
| POST | `/api/v1/lohang` | TaoLotHang | Submit | HTXNongSan, Platform |
| PUT | `/api/v1/lohang/:ma_lo/trangthai` | CapNhatTrangThaiLo | Submit | HTXNongSan, Platform |
| POST | `/api/v1/lohang/:ma_lo/chungnhan` | ThemChungNhan | Submit | ChiCucBVTV, Platform |
| GET | `/api/v1/lohang/:ma_lo` | DocLotHang | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/lohang/:ma_lo/lichsu` | LichSuLotHang | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/lohang/htx/:ma_htx` | DocLotHangTheoHTX | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/consumer/:ma_lo` | LayThongTinTraCuu | Evaluate | **Public — dùng Platform fallback** |

**Nhật ký (`nhatky_cc` — nhatky-htx-channel):**

| Method | Endpoint | UseCase Fn | Submit/Eval | Org cho phép |
|--------|----------|-------------|-------------|-------------|
| POST | `/api/v1/nhatky` | GhiNhatKy | Submit | HTXNongSan, Platform |
| PUT | `/api/v1/nhatky/:id/duyet` | DuyetNhatKy | Submit | ChiCucBVTV, Platform |
| GET | `/api/v1/nhatky/lo/:ma_lo` | DocNhatKyTheoLo | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/nhatky/htx/:ma_htx` | DocNhatKyTheoHTX | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/nhatky/:id/lichsu` | LichSuNhatKy | Evaluate | Tất cả org trên nhatky-htx-channel |
| GET | `/api/v1/nhatky/thongke` | ThongKeNhatKy | Evaluate | Tất cả org trên nhatky-htx-channel |

**Giao dịch (`giao_dich_cc` — giaodich-channel):**

| Method | Endpoint | UseCase Fn | Submit/Eval | Org cho phép |
|--------|----------|-------------|-------------|-------------|
| POST | `/api/v1/giaodich` | TaoGiaoDich | Submit | HTXNongSan, Platform |
| PUT | `/api/v1/giaodich/:id/duyet` | DuyetGiaoDich | Submit | Platform |
| PUT | `/api/v1/giaodich/:id/trangthai` | CapNhatTrangThai | Submit | HTXNongSan, NPPXanh, Platform |
| GET | `/api/v1/giaodich/:id` | DocGiaoDich | Evaluate | Tất cả org trên giaodich-channel |
| GET | `/api/v1/giaodich/:id/lichsu` | LichSuGiaoDich | Evaluate | Tất cả org trên giaodich-channel |
| GET | `/api/v1/giaodich/npp/:ma_npp/congno` | DocCongNoNPP | Evaluate | NPPXanh, Platform |
| GET | `/api/v1/giaodich/npp/:ma_npp/hoahong` | TinhHoaHongNPP | Evaluate | Platform |

## Lưu ý khi Code

1. **GatewayRegistry là singleton**: Khởi tạo 1 lần lúc startup với tất cả org configs, dùng chung cho mọi request
2. **Mỗi request resolve đúng org**: Từ `mspID` trong JWT → `registry.GetOrgGateway(mspID)` → dùng Contract có sẵn
3. **Public endpoint**: Dùng `registry.GetFallbackGateway()` (Platform) vì không có JWT
4. **Kiểm tra nil Contract**: Trước khi gọi, phải check contract != nil (org có thể không tham gia channel đó)
5. **Phân quyền 2 lớp**: 
   - **Lớp 1 (Backend)**: Check mspID có contract tương ứng không (channel membership)
   - **Lớp 2 (Chaincode)**: `GetClientIdentity().GetMSPID()` check business rule (vd: chỉ BVTV mới duyệt)
6. **Graceful shutdown**: `registry.CloseAll()` đóng tất cả gRPC connections

$ARGUMENTS
