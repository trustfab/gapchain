# GAPChain - Phát triển Chaincode

Hỗ trợ viết, sửa, và test chaincode Go cho GAPChain.

## Cấu trúc Chaincode

```
chain-setup/chaincode/
├── nhatky_cc/          - Nhật ký hoạt động HTX
│   ├── nhatky.go       - Logic chính
│   └── go.mod          - Go 1.19, fabric-contract-api-go v1.2.2
└── giao_dich_cc/       - Quản lý giao dịch
    ├── giao_dich.go    - Logic chính
    └── go.mod
```

## Dependencies

```
github.com/hyperledger/fabric-contract-api-go v1.2.2
github.com/hyperledger/fabric-chaincode-go
github.com/hyperledger/fabric-protos-go
```

## Pattern chuẩn của project

### Struct data

```go
type NhatKy struct {
    ID          string `json:"id"`
    MaHTX       string `json:"ma_htx"`
    LoaiHoatDong string `json:"loai_hoat_dong"`
    TrangThai   string `json:"trang_thai"` // cho_duyet, da_duyet, tu_choi
    NguoiTao    string `json:"nguoi_tao"`
    ThoiGianTao string `json:"thoi_gian_tao"`
}
```

### Hàm invoke mẫu

```go
func (s *SmartContract) TaoNhatKy(ctx contractapi.TransactionContextInterface, id string, maHTX string, ...) error {
    exists, err := s.nhatKyExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("Nhật ký %s đã tồn tại", id)
    }

    nhatKy := NhatKy{
        ID:    id,
        MaHTX: maHTX,
        // ...
        TrangThai: "cho_duyet",
    }

    nhatKyJSON, err := json.Marshal(nhatKy)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, nhatKyJSON)
}
```

### Query với selector (CouchDB)

```go
func (s *SmartContract) DocNhatKyTheoHTX(ctx contractapi.TransactionContextInterface, maHTX string) ([]*NhatKy, error) {
    queryString := fmt.Sprintf(`{"selector":{"ma_htx":"%s"}}`, maHTX)
    return getQueryResultForQueryString(ctx, queryString)
}
```

## Nhiệm vụ

Người dùng muốn thêm/sửa chaincode. Hãy:

1. Đọc file chaincode liên quan trước khi đề xuất thay đổi
2. Giữ đúng conventions của project (tên hàm tiếng Việt, JSON tags snake_case)
3. Sau khi sửa code, nhắc người dùng cần redeploy chaincode với version mới
4. Kiểm tra: tất cả hàm public phải có error handling đúng
5. Với query phức tạp dùng CouchDB selector - đảm bảo index được tạo nếu cần

## Redeploy sau khi sửa

Khi chaincode thay đổi, cần bump version và redeploy:
```bash
# Trong deploy-chaincode.sh, tăng CC_VERSION
# Sau đó chạy lại deploy
cd chain-setup && ./scripts/deploy-chaincode.sh
```

$ARGUMENTS
