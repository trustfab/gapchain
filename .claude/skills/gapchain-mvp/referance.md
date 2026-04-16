# Tham Chiếu Chuyên Sâu: GAPChain MVP Architecture & Data Models

> **Mục đích**: Tài liệu này đóng vai trò là cơ sở lý thuyết và tiêu chuẩn thiết kế kiến trúc, cấu trúc dữ liệu cho MVP GAPChain. Được tích hợp và trích xuất từ "Báo cáo Chuyên sâu: Phát triển MVP Hệ thống Truy xuất Nguồn gốc Nông sản Sạch dựa trên Blockchain".
> **BẮT BUỘC** tham chiếu khi: Thiết kế Field (KDEs) trong Chaincode, chuẩn bị luồng dữ liệu cho API Backend và xây dựng logic đồng bộ Offline-to-Online trên Frontend Flutter.

---

## 1. Nguyên Tắc Thiết Kế Cốt Lõi của MVP

1. **Truy xuất lấy Lô Hàng làm trung tâm (Batch-based Traceability)**: 
   Sự ra đời của `lohang_cc` là nền tảng tối quan trọng. Mỗi lô hàng định danh cho hạt giống, thửa ruộng và vụ mùa. Hệ thống giải quyết truy xuất ngược và thu hồi (recall) cục bộ bằng cách cô lập dữ liệu giới hạn trong từng batch.
2. **Offline-First & Mobile-First cho Nông dân (Trọng tâm trải nghiệm)**:
   Do bối cảnh tại Việt Nam (hạ tầng mạng nông thôn chưa ổn định, kỹ năng số của nông dân hạn chế), ứng dụng di động phải có khả năng lưu log ngoại tuyến, tự động bắt GPS nền, sau đó đồng bộ (sync) khi có mạng ổn định nhằm giảm thiểu tối đa rào cản thao tác.
3. **Mô hình Dữ Liệu đáp ứng Tuân thủ & Xuất Khẩu**:
   Hệ thống phải có bộ khung cấu trúc đủ tốt để hỗ trợ xuất báo cáo kiểm định: VietGAP (nhật ký "4 đúng"), chuẩn Hữu cơ PGS, EUDR (chống phá rừng qua bằng chứng GPS canh tác) và FSMA 204 (truy xuất nguồn gốc quốc tế).
4. **Kiến Trúc Dữ Liệu Lai (Hybrid Architecture)**:
   - *On-chain*: Lưu trữ mã băm (SHA256 Hash), trạng thái giao dịch (`dang_trong`, `da_thu_hoach`), định danh lô, CTEs (Sự kiện theo dõi thiết yếu).
   - *Off-chain*: Ảnh chụp làm bằng chứng, tài liệu PDF vật lý, số liệu chi tiết mật độ cao của cảm biến IoT. Phối hợp với CouchDB để Query Rich nhanh chóng.

---

## 2. Hệ Thống Các Yếu Tố Dữ Liệu Chính (KDEs) Áp Dụng Chéo

Thiết kế Chaincode (`lohang_cc`, `nhatky_cc`, `giao_dich_cc`) cần phản ánh chính xác báo cáo thực địa theo thiết kế KDE:

### 2.1 Cấp Nông Trại / HTX (Khởi Trị Lô & Ghi Nhật Ký)
- **Định danh Không gian (Vị trí)**: ID Lô Đất, Tọa độ GPS gắn liền lô đó (Thiết yếu cho EUDR). Nhúng tại metadata của lô hoặc tọa độ chính xác lúc submit nhật ký.
- **Nguồn Gốc Cây Trồng/Vật Tư**: Phân bón, giống, thuốc trừ sâu. Các thuộc tính này nằm trong trường `chi_tiet` hoặc các struct chuyên biệt định nghĩa "Nông nghiệp sạch".
- **Nhật Ký Hoạt Động (Nhật ký Cấp bậc theo VietGAP)**: Loại can thiệp (tưới nước, bón phân, thu hoạch), Liều lượng, Cơ chế, Ngày giờ hệ thống. Báo cáo sâu bệnh, dịch hại.
- **Thông Số Thu Hoạch**: ID Lô xuất sau khi nhổ rễ, Ngày thu hoạch, Năng suất ước lượng kiểm đếm.

### 2.2 Cấp Chế Biến & Phân Phối (Giao Dịch & Lưu Kho)
- **Minh Bạch Hành Trình**: Ghi chép cho sơ chế, làm sạch, sấy, đóng gói.
- **Giám sát Môi Trường (IoT)**: Luồng dữ liệu cho thông số Nhiệt độ, Độ ẩm (đặc biệt quan trọng với kho lạnh nếu có).
- **Sự Kiện Chuyển Giao Quyền Sở Hữu**: Giao dịch Bán (`giao_dich_cc`) bắt buộc phải ghi ghim lại Ai (NPP nào), Khi nào, Trạng thái đơn và số Công nợ được tính tự động, hình thành 1 điểm CTE.

### 2.3 Cấp Siêu Thị / Người Tiêu Dùng (Consumer Endpoint - QR Trust Portal)
- **Tích Hợp Timeline Duy Nhất (Single Truth Story)**: Mã QR không chỉ trỏ đến 1 mã lô tĩnh, mà nó gọi hàm `LayThongTinTraCuu` để render "Câu chuyện Mã QR" linh động = `Thông tin Lô Hàng` + `Timeline Nhật Ký Đã Duyệt` + `Hồ Sơ Chứng Nhận HTX`.

---

## 3. Khả Năng Tương Tác & Tiêu Chuẩn Hóa
- Hệ thống nên có tư duy chừa khe hở hoặc map property theo **GS1 Standards** (e.g. EPCIS vocabulary) trong tương lai.
- Sử dụng danh pháp và mã chuẩn (State Dictionary, Event Types) để đảm bảo không gặp vướng mắc nếu nhà nước triển khai **Cổng thông tin Truy xuất Nguồn gốc Quốc gia** và yêu cầu kết nối với API quốc gia.
