# GAPChain Business Rules — Single Source of Truth (Version 1.2 - Full Edition)

> **Mục đích**: Tài liệu này là nguồn sự thật duy nhất cho mọi trạng thái, quyền hạn và luồng xử lý của hệ thống GAPChain. 
> **BẮT BUỘC** tham chiếu khi: Viết/sửa chaincode, xây dựng API, thiết kế UI/UX và xử lý tranh chấp dữ liệu.

---

## 1. Actors & Phân quyền Chi tiết (Identity & Access Control)

| Actor | MSP ID | Vai trò | Cơ chế xác thực |
| :--- | :--- | :--- | :--- |
| **Platform** | `PlatformOrgMSP` | Quản trị hệ thống, duyệt GD, tính hoa hồng, xử lý sự cố/phục hồi lô. | MSP Identity |
| **HTX** | `HTXNongSanOrgMSP` | Tạo lô hàng, ghi nhật ký, tạo giao dịch bán, quản lý nông dân. | MSP Identity |
| **Nông dân** | *(Thuộc HTX)* | Thực hiện canh tác, xác nhận công việc thực tế tại vườn. | OTP SMS / Chữ ký ảnh |
| **Chi Cục BVTV** | `ChiCucBVTVOrgMSP` | Duyệt nhật ký, cấp chứng nhận (VietGAP...), **Đình chỉ lô hàng**. | MSP Identity |
| **NPP** | `NPPXanhOrgMSP` | Mua hàng, xác nhận nhận hàng, thanh toán, theo dõi công nợ. | MSP Identity |
| **Consumer** | *(Public)* | Quét QR truy xuất nguồn gốc — chỉ đọc dữ liệu đã duyệt. | No Auth (Public API) |

---

## 2. Vòng đời Lô Hàng & Logic Tách Lô (Batch Management)

### 2.1 Sơ đồ trạng thái Lô Hàng (`lohang_cc`)
1.  **dang_trong**: Trạng thái khởi tạo khi xuống giống.
2.  **da_thu_hoach**: Đã thu hoạch thực tế, đang lưu kho tại HTX.
3.  **cho_chung_nhan**: Đang chờ Chi cục BVTV kiểm định/cấp chứng chỉ.
4.  **san_sang_ban**: Đã có chứng nhận, có thể tạo Giao dịch hoặc Tách lô con.
5.  **het_hang**: Trạng thái kết thúc khi `so_luong_con_lai == 0`.
6.  **dinh_chi**: Trạng thái khóa tạm thời khi phát hiện vi phạm (chỉ Platform/BVTV có quyền).

### 2.2 Quy tắc Tách Lô & Inventory (Batch Splitting)
* **Lô mẹ (Parent Batch)**: Lưu tổng sản lượng ban đầu.
* **Lô con / Giao dịch**: Khi bán hàng hoặc tách nhỏ, hệ thống phải kiểm tra: `Tổng_Lượng_Tách <= so_luong_con_lai` của lô mẹ.
* **Truy vết ngược**: Lô con phải lưu `ma_lo_me` để người dùng có thể truy xuất ngược về toàn bộ nhật ký canh tác của lô gốc.

---

## 3. Nhật ký Sản xuất & Farmer Bridge (Lớp chặn dữ liệu rác)

Để đảm bảo dữ liệu "sạch" trước khi lên Blockchain, áp dụng quy tắc **Xác nhận 2 lớp**:

1.  **Lớp 1 (Ghi nhận)**: Kỹ thuật viên (KTV) nhập nhật ký thay nông dân $\rightarrow$ Trạng thái: `cho_nong_dan_xac_nhan`.
2.  **Lớp 2 (Xác thực thực địa)**: 
    * Nông dân dùng App xác nhận **HOẶC** KTV chụp ảnh phiếu làm việc có chữ ký tay nông dân $\rightarrow$ Tính Hash ảnh $\rightarrow$ Trạng thái: `cho_duyet`.
3.  **Lớp 3 (Thẩm định)**: Chi cục BVTV kiểm tra nội dung $\rightarrow$ Trạng thái: `da_duyet`.

> **Quy tắc hiển thị QR**: Chỉ những nhật ký đạt trạng thái `da_duyet` mới được hiển thị cho Consumer.

---

## 4. Vòng đời Giao dịch Thương mại (`giao_dich_cc`)

| Trạng thái | Bên thực hiện | Mô tả |
| :--- | :--- | :--- |
| **cho_duyet** | HTX | Tạo đơn hàng bán cho NPP. Hệ thống trừ tạm kho ở `lohang_cc`. |
| **da_duyet** | Platform | Xác nhận đơn hàng hợp lệ về giá và chính sách. |
| **dang_giao** | Platform | Xác nhận hàng đã rời kho HTX. |
| **da_giao** | NPP | Xác nhận đã nhận đủ hàng và đúng chất lượng. |
| **cho_thanh_toan**| Platform | Chốt công nợ. NPP thấy khoản phải trả trong danh sách công nợ. |
| **da_thanh_toan** | NPP | Xác nhận chuyển khoản thành công. Kết thúc vòng đời GD. |
| **huy_bo** | Platform | Chỉ Platform có quyền. Hệ thống cộng hoàn lại số lượng vào kho lô hàng. |

---

## 5. Cơ chế Thu hồi & Phản ứng nhanh (Recall Logic)

Khi phát hiện lô hàng có dư lượng thuốc trừ sâu hoặc vi phạm VietGAP:
1.  **Action**: BVTV hoặc Platform gọi hàm `CapNhatTrangThaiLo` $\rightarrow$ `dinh_chi`.
2.  **Impact 1 (Giao dịch)**: Tất cả GD liên quan chưa hoàn thành (`da_giao`, `dang_giao`) sẽ bị treo và hiển thị cảnh báo **"HÀNG CẦN THU HỒI"**.
3.  **Impact 2 (Truy xuất)**: QR code của lô hàng đó khi quét sẽ hiện thông báo màu đỏ: **"SẢN PHẨM KHÔNG AN TOÀN - ĐANG THU HỒI"**.
4.  **Recovery**: Chỉ Platform mới có quyền chuyển từ `dinh_chi` về lại trạng thái cũ sau khi khắc phục.

---

## 6. Quy tắc Kỹ thuật Cross-Channel (Data Consistency)

Hệ thống vận hành trên 2 channel (`nhatky` và `giaodich`), Backend API phải đảm bảo:
* **Atomic Transaction (Saga)**: Khi tạo GD ở Channel Giao dịch, phải gọi đồng thời hàm trừ kho ở Channel Nhật ký.
* **Retry Mechanism**: Nếu cập nhật kho thất bại do lỗi mạng, hệ thống phải tự động thực hiện lại (Retry) cho đến khi dữ liệu khớp nhau.
* **Data Mapping**: `ma_lo` là khóa chính (Primary Key) để Join dữ liệu giữa các Channel khi trả kết quả cho API Truy xuất.

---

## 7. Checklist Review cho MVP

- [ ] **Chaincode**: Đã kiểm tra `GetTxTimestamp()` để đảm bảo tính deterministic?
- [ ] **Quyền hạn**: KTV có thể sửa nhật ký ở trạng thái `cho_duyet` không? (Đáp án: Không, phải ghi bản mới).
- [ ] **Validation**: `ty_le_hoa_hong` có nằm trong khoảng 0-100%?
- [ ] **Frontend**: Badge màu trạng thái có đúng quy định (Xanh: An toàn, Đỏ: Đình chỉ, Vàng: Chờ duyệt)?
- [ ] **Consumer**: Trang QR đã ẩn các thông tin nhạy cảm (Giá bán, Chiết khấu)?

---