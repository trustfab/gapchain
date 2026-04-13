# GAPChain Skill: Thêm Mới Node (Organization) Vào Mạng Đang Chạy

## 1. Mục Tiêu Nghiệp Vụ
Mở rộng mạng Hyperledger Fabric đang hoạt động (ví dụ: `giaodich-channel`) bằng cách bổ sung một Organization mới (ví dụ: `NPPTieuChuanOrg`) mà không làm gián đoạn Blockchain hay mất dữ liệu Ledger hiện hành.

## 2. Quy Trình Tổng Quan
Quy trình thêm một Node gồm hai phân hệ:
1. **Network Layer (Fabric):** Tạo Crypto material $\rightarrow$ Trích xuất configtx JSON $\rightarrow$ Tạo Config Update (Delta) $\rightarrow$ Ký bởi Majority Admin $\rightarrow$ Submit Orderer $\rightarrow$ Start Peer $\rightarrow$ Join Channel $\rightarrow$ Install & Approve Chaincode với Sequence mới.
2. **Application Layer:** Cập nhật Frontend (Vue Router, State Auth) và Backend (Go .env FABRIC_ORGS, auth_handler, gateway).

---

## 3. Các Lỗi Thường Gặp & Bài Học Xương Máu (Lessons Learned)

### Lỗi 1: `policy not satisfied` khi Ký Bản Cập Nhật
- **Hiện tượng:** Yêu cầu Admin của Tổ chức *Mới* ký vào config update để "xin gia nhập". Kết quả bị từ chối cập nhật.
- **Nguyên lý:** Việc thêm mới một Node vào Application Group được Fabric coi là hành vi **Sửa đổi cấp độ Channel/Application**. Do đó, nó BẮT BUỘC dùng `mod_policy` của Admins cũ (PlatformOrg + HTXNongSanOrg).
- **Giải pháp:** Tổ chức mới KHÔNG CẦN và KHÔNG THỂ ký vào block xin gia nhập. Chỉ mượn context của Admin Tổ chức cũ để ký (`peer channel signconfigtx`).

### Lỗi 2: `Unknown Application Org ConfigValue name: Endpoints` (Configtxlator)
- **Hiện tượng:** Biên dịch JSON ra Protobuf bị lỗi không tìm thấy ConfigValue `Endpoints`.
- **Nguyên lý:** Lệnh `configtxgen -printOrg` xuất ra tất cả JSON bao gồm `OrdererEndpoints`. Khi dùng `jq` merge trực tiếp vào `Application` Group thì bị sai schema. Application Group (Protobuf) chỉ nhận `AnchorPeers` và `MSP`.
- **Giải pháp:** Gỡ bỏ trường `OrdererEndpoints` ra khỏi định nghĩa cục bộ của Tổ chức mới nhưng dứt khoát **phải cấu hình `AnchorPeers`** (để đảm bảo Gossip Protocol giao tiếp chéo giữa các Peer).

### Lỗi 3: `x509: certificate is valid for peer0..., not localhost`
- **Hiện tượng:** Chạy `peer channel join` bằng môi trường local (`localhost:11051`) thì TLS Handshake bị rớt.
- **Giải pháp kép bền vững:** 
  1. Thêm mảng `SANS: - "localhost"` vào file `crypto-config.yaml` lúc build cert.
  2. Map cứng biến ghi đè vào Session bash: `export CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer0.tentoChuc.gapchain.vn`

### Lỗi 4: `Error: can't read the block: &{FORBIDDEN}` khi kéo Genesis Block
- **Hiện tượng:** Orderer từ chối không cho Admin của Node mới kéo Block 0 về (dù đã được xác nhận vào kênh).
- **Nguyên lý:** Có thể do độ trễ đồng bộ chứng chỉ TLS mới lên Orderer hoặc Orderer giới hạn quyền đọc block gốc.
- **Vá lỗi cực hay ho (Workaround):** Mượn context (chứng chỉ) của `PlatformOrg` (kẻ chắc chắn có quyền đọc) để gõ lệnh `peer channel fetch 0`, tải cục `.block` vật lý xuống. Sau đó, đổi biến môi trường về `Org Mới` và dùng chính file đó đập lệnh `peer channel join`.

### Lỗi 5: Lỗi Container thiếu Volume & `read/write on closed pipe` khi Install Chaincode
- **Hiện tượng:** Lệnh `docker cp` nén file chaincode vào thư mục `/opt/gopath...` báo lỗi không tồn tại đường dẫn.
- **Nguyên lý:** Do file `docker-compose` sinh ra quên mount Volume `organizations`, `channel-artifacts` và `var/hyperledger/production`.
- **Giải pháp:** File `docker-compose.yaml` cho Peer mới bắt buộc phải map đường dẫn Volumes hệ thống và chỉ định `working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer`.

### Lỗi 6: API Backend báo `khong tim thay gateway cho org`
- **Hiện tượng:** Gọi API trên Postman hoặc web thì báo lỗi gateway.
- **Nguyên lý:** Dù đã viết đủ File biến môi trường, Router và Gateway Go, nhưng lại **quên bổ sung MSP** của Node mới vào list mạng định danh.
- **Giải pháp:** Luôn nhớ chèn MSP ID vào chuỗi cấu hình `FABRIC_ORGS` bên trong file `backend/.env`.

---

## 4. Disaster Recovery (Hoàn Tác Rủi Ro Config)

*Hỏi: Làm sao khi thao tác nhầm hoặc Node mới bị lỗi mà không phải Wipe toàn cục Mạng (Cleanup)?*
- **Tình huống 1 (Mới sinh file pb):** File Config chưa submit hoặc bị Orderer chối lỗi JSON $\rightarrow$ Mạng chưa hề thay đổi. Chỉnh script, xóa file rác sinh lại.
- **Tình huống 2 (Submit Config thành công nhưng Docker Peer sập/lỗi cert):** Peer chết nhưng Config Blockchain đã lưu Org này $\rightarrow$ Sửa lại Docker/Cert rồi up peer lên lại, tự nó sẽ Join kênh và catch-up Ledger mà không cản trở Mạng gốc chạy.
- **Tình huống 3 (Thêm nhầm Org lỗi nặng vào Mạng):** $\rightarrow$ Kéo config block mới nhất về, dùng `jq` **XÓA** Node đó khỏi cây `Application.groups`, lấy Update Delta trích xuất ngược lại, xin chữ ký Admins và Submit $\rightarrow$ Orderer sinh Block tống cổ Node lỗi khỏi mạng an toàn.
