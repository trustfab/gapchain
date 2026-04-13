# GAPChain - Phát triển Go Client

Hỗ trợ phát triển ứng dụng Go client kết nối với mạng GAPChain.

## Cấu trúc fabric-go-client

Thư mục hiện tại chỉ còn lại 2 file demo chính sử dụng Fabric Gateway SDK:
```
fabric-go-client/
├── go.mod
├── go.sum
├── htx_ngonsan.go   - Client demo cho HTX Nông Sản
└── npp_nongsan.go   - Client demo cho Nhà Phân Phối
```

## Cách kết nối (Fabric Gateway SDK)

Không sử dụng file cấu hình `.yaml` trung gian, thay vào đó khai báo trực tiếp tới các file crypto (cert/key) của user và gọi qua gRPC tới Peer:

```go
// 1. Tạo kết nối gRPC Client TLS
clientConnection, err := newGrpcConnection()

// 2. Nạp thông tin identity (X509) và Private key (Sign) 
id := newIdentity()
sign := newSign()

// 3. Khởi tạo Gateway Client
gateway, err := client.Connect(
	id,
	client.WithSign(sign),
	client.WithHash(hash.SHA256),
	client.WithClientConnection(clientConnection),
)

// 4. Lấy đối tượng Network và Contract
network := gateway.GetNetwork("nhatky-htx-channel")
contract := network.GetContract("nhatky_cc")

// Giao dịch thay đổi trạng thái (SubmitTransaction)
_, err = contract.SubmitTransaction("TaoNhatKy", args...)

// Giao dịch truy vấn dữ liệu (EvaluateTransaction)
result, err := contract.EvaluateTransaction("DocTatCaNhatKy")
```

## Nhiệm vụ

Người dùng muốn phát triển / chạy Go client. Hãy:

1. Đọc nội dung file `htx_ngonsan.go` hoặc `npp_nongsan.go` để làm mẫu code viết đúng.
2. Kiểm tra đường dẫn biến số `cryptoPath` để đảm bảo certs trên máy đúng chuẩn và còn cấp phép.
3. Chú ý port kết nối gRPC (ví dụ HTXNongSan là `localhost:8051`, NPPXanh là `localhost:10051`).
4. Nhắc kiểm tra các Docker container của mạng GAPChain phải đang chạy.

## Chạy client

Di chuyển vào thư mục client và dùng `go run`:

```bash
cd fabric-go-client

# Chạy demo HTX Nông Sản
go run htx_ngonsan.go

# Chạy demo Nhà Phân Phối
go run npp_nongsan.go
```

$ARGUMENTS
