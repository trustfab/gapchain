---
name: gapchain-ba
description: >
  Quy tắc nghiệp vụ (Business Rules) cho hệ thống truy xuất nguồn gốc nông sản GAPChain MVP.
  Skill này là nguồn sự thật duy nhất (Single Source of Truth) cho tất cả trạng thái, quyền hạn,
  và luồng xử lý end-to-end. BẮT BUỘC tham chiếu khi: viết/sửa chaincode, xây dựng API endpoint,
  thiết kế UI/UX, review code, hoặc bất kỳ quyết định liên quan đến luồng nghiệp vụ.
  Trigger khi người dùng đề cập: "nghiệp vụ", "business rule", "trạng thái", "vòng đời",
  "quyền hạn", "ai được phép", "luồng xử lý", "BA", "kiểm tra logic", "validate".
---

# GAPChain Business Rules — Single Source of Truth

> **Mục đích**: Tất cả code (chaincode, backend, frontend) PHẢI tuân thủ tài liệu này.
> Khi có mâu thuẫn giữa code và tài liệu này → tài liệu này đúng, code cần sửa.
> *Bản cập nhật v1.2: Bổ sung logic Tách Lô, Cơ chế Thu hồi (Recall), và Xác thực 2 lớp (Farmer Bridge).*

---

## 1. Actors & Phân quyền

### 1.1 Bảng tổ chức & Định danh

| Actor | MSP ID | Vai trò | Cơ chế xác thực | Channel tham gia |
|-------|--------|---------|-----------------|-------------------|
| **Platform** | `PlatformOrgMSP` | Quản trị hệ thống, duyệt giao dịch, tính hoa hồng, hủy GD, xử lý sự cố/phục hồi lô. | MSP Identity | nhatky-htx-channel, giaodich-channel |
| **HTX Nông Sản** | `HTXNongSanOrgMSP` | Tạo lô hàng, tách lô, ghi nhật ký, tạo giao dịch bán, quản lý nông dân. | MSP Identity | nhatky-htx-channel, giaodich-channel |
| **Nông dân** | *(Thuộc HTX)* | Thực hiện canh tác, xác nhận công việc thực tế tại vườn. | OTP SMS / Chữ ký ảnh | *(Gián tiếp qua HTX)* |
| **Chi Cục BVTV** | `ChiCucBVTVOrgMSP` | Duyệt nhật ký, cấp chứng nhận (VietGAP...), **Đình chỉ lô hàng**. | MSP Identity | nhatky-htx-channel |
| **NPP Xanh** | `NPPXanhOrgMSP` | Mua hàng, xác nhận nhận hàng, thanh toán, theo dõi công nợ, yêu cầu thu hồi. | MSP Identity | giaodich-channel |
| **Consumer** | *(Public)* | Quét QR truy xuất nguồn gốc — chỉ đọc dữ liệu đã duyệt. | No Auth | *(Qua API public)* |

### 1.2 Ma trận quyền hạn chi tiết

| Thao tác | Platform | HTX | BVTV | NPP | Consumer |
|----------|:--------:|:---:|:----:|:---:|:--------:|
| **LÔ HÀNG** | | | | | |
| Tạo lô hàng / Tách lô con | ✅ | ✅ | ❌ | ❌ | ❌ |
| Chuyển trạng thái lô (bình thường) | ✅ | ✅ | ❌ | ❌ | ❌ |
| Đình chỉ lô hàng (`dinh_chi`) | ✅ | ❌ | ✅ | ❌ | ❌ |
| Phục hồi lô từ đình chỉ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Cập nhật số lượng lô | ✅ | ✅ | ❌ | ❌ | ❌ |
| Thêm chứng nhận | ✅ | ❌ | ✅ | ❌ | ❌ |
| Xem lô hàng | ✅ | ✅ | ✅ | ❌ | ✅ (QR) |
| Xem lịch sử lô hàng | ✅ | ✅ | ✅ | ❌ | ❌ |
| **NHẬT KÝ** | | | | | |
| Ghi nhận nhật ký (KTV - Lớp 1) | ✅ | ✅ | ❌ | ❌ | ❌ |
| Xác nhận thực địa (Nông dân - Lớp 2)| ❌ | *(Via App/KTV)* | ❌ | ❌ | ❌ |
| Duyệt / Từ chối (Thẩm định - Lớp 3) | ✅ | ❌ | ✅ | ❌ | ❌ |
| Xem nhật ký | ✅ | ✅ | ✅ | ❌ | ✅ (QR) |
| Xem lịch sử nhật ký | ✅ | ✅ | ✅ | ❌ | ❌ |
| **GIAO DỊCH** | | | | | |
| Tạo giao dịch | ✅ | ✅ | ❌ | ❌ | ❌ |
| Duyệt giao dịch | ✅ | ❌ | ❌ | ❌ | ❌ |
| Xác nhận giao hàng (dang_giao) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Xác nhận nhận hàng (da_giao) | ✅ | ❌ | ❌ | ✅ | ❌ |
| Chốt thanh toán (cho_thanh_toan) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Xác nhận đã thanh toán (da_thanh_toan) | ✅ | ❌ | ❌ | ✅ | ❌ |
| Hủy giao dịch | ✅ | ❌ | ❌ | ❌ | ❌ |
| Xem giao dịch | ✅ | ✅ | ❌ | ✅ | ❌ |
| Xem công nợ NPP | ✅ | ❌ | ❌ | ✅ | ❌ |
| Tính hoa hồng | ✅ | ❌ | ❌ | ❌ | ❌ |

---

## 2. Vòng đời Lô Hàng & Logic Tách Lô (`lohang_cc`)

### 2.1 Sơ đồ trạng thái (STRICT — enforce trong chaincode)

```
                    ┌─────────────┐
         HTX tạo    │ dang_trong  │  ← Trạng thái khởi tạo
                    └──────┬──────┘
                           │ HTX cập nhật
                    ┌──────▼──────┐
                    │da_thu_hoach │  ← Đã thu hoạch thực tế
                    └──────┬──────┘
                           │ HTX gửi kiểm định
                  ┌────────▼────────┐
                  │ cho_chung_nhan  │  ← BVTV đang kiểm tra
                  └────────┬────────┘
                           │ BVTV cấp chứng nhận xong
                    ┌──────▼──────┐
                    │san_sang_ban │  ← Có thể tạo Giao dịch hoặc Tách Lô (Partial batch)
                    └──────┬──────┘
                           │ Hết số lượng hoặc HTX đóng lô
                    ┌──────▼──────┐
                    │  het_hang   │  ← TRẠNG THÁI KẾT THÚC
                    └─────────────┘

     ┌────────────────────────────────────────────────────────┐
     │ Bất kỳ trạng thái nào (trừ het_hang)                   │
     │   → dinh_chi (chỉ PlatformOrgMSP hoặc ChiCucBVTVOrgMSP)│
     │                                                        │
     │ dinh_chi → trạng thái trước đó (CHỈ PlatformOrgMSP)    │
     └────────────────────────────────────────────────────────┘
```

### 2.2 Quy tắc Tách Lô & Inventory (Batch Splitting)
Để đảm bảo lưu thông hàng hóa mà không làm mất vết nguồn gốc, áp dụng quy tắc **Tách Lô**:
* **Lô mẹ (Parent Batch)**: Lưu tổng sản lượng ban đầu. Field `ma_lo_me` để trống.
* **Lô con / Giao dịch**: Khi tạo giao dịch bán phần hoặc tách lô nhỏ, hệ thống bắt buộc kiểm tra: `Tổng_Lượng_Tách <= so_luong_con_lai` của lô mẹ.
* **Truy vết ngược (Traceability Link)**: Lô con sẽ được sinh ra mang ID `ma_lo_me` trỏ về lô mẹ. QR Truy xuất của lô con hoặc các giao dịch con sẽ tự động được "kế thừa" và fetch toàn bộ **nhật ký canh tác** lẫn **chứng nhận** từ lô mẹ.

### 2.3 Quy tắc chuyển trạng thái (enforce trong chaincode)

| Từ | Đến | Ai thực hiện | Hàm chaincode | Điều kiện tiên quyết |
|----|-----|-------------|---------------|---------------------|
| _(mới tạo)_ | `dang_trong` | HTX, Platform | `TaoLotHang` | Tự động |
| `dang_trong` | `da_thu_hoach` | HTX, Platform | `CapNhatTrangThaiLo` | Nên có ≥1 nhật ký `thu_hoach` đã duyệt (backend check) |
| `da_thu_hoach` | `cho_chung_nhan` | HTX, Platform | `CapNhatTrangThaiLo` | HTX gửi hồ sơ cho BVTV kiểm tra |
| `cho_chung_nhan` | `san_sang_ban` | HTX, Platform | `CapNhatTrangThaiLo` | **NÊN** có ≥1 chứng nhận được cấp (backend check) |
| `san_sang_ban` | `het_hang` | HTX, Platform | `CapNhatTrangThaiLo` | Hết số lượng (`so_luong_con_lai == 0`) hoặc HTX đóng lô |
| _bất kỳ (trừ het_hang)_ | `dinh_chi` | **Platform, BVTV** | `CapNhatTrangThaiLo` | Phát hiện vi phạm chất lượng, dư lượng thuốc, mất chứng nhận |
| `dinh_chi` | _trạng thái phù hợp_ | **Platform** | `CapNhatTrangThaiLo` | Khắc phục xong. **Chỉ Platform** được phục hồi |

### 2.4 Chứng nhận (ChungNhan)

| Loại | Cấp bởi | Ý nghĩa |
|------|---------|---------|
| `VietGAP` | ChiCucBVTV | Tiêu chuẩn VietGAP — bắt buộc cho MVP |
| `GlobalGAP` | ChiCucBVTV | Tiêu chuẩn quốc tế — ưu tiên xuất khẩu |
| `Organic` | ChiCucBVTV | Hữu cơ — premium |
| `TCVN` | ChiCucBVTV | Tiêu chuẩn Việt Nam — cơ bản |
| `Khac` | ChiCucBVTV | Chứng nhận khác |

### 2.5 Validation rules Lô Hàng

| Field | Bắt buộc | Validation |
|-------|:--------:|-----------|
| `ma_lo` | ✅ | Unique, format: `LH-{maHTX}-{năm}-{seq}` |
| `ma_lo_me` | ❌ | Bỏ trống nếu là lô gốc; chứa ID lô mẹ nếu là lô tách/lô con để Traceability. |
| `ma_htx` | ✅ | Phải tồn tại trong hệ thống |
| `ten_san_pham` | ✅ | Không rỗng |
| `loai_san_pham` | ✅ | Enum: `lua`, `ca_phe`, `rau`, `qua`, `khac` |
| `so_luong` | ✅ | > 0 |
| `so_luong_con_lai` | ✅ | `0 <= so_luong_con_lai <= so_luong`. Khởi tạo bằng `so_luong`. |
| `don_vi_tinh` | ✅ | Enum: `kg`, `tan`, `thung`, `qua` |

---

## 3. Vòng đời Nhật Ký & Farmer Bridge (`nhatky_cc`)

### 3.1 Sơ đồ trạng thái (Farmer Bridge Integration)

```
                         ┌───────────────────────┐
             HTX ghi     │ cho_nong_dan_xac_nhan │ ← Ghi nhận qua App (Lớp 1)
                         └──────────┬────────────┘
                                    │ Nông dân xác nhận OTP
                                    │ hoặc KTV chụp chữ ký bản cứng (Hash lên chain)
                         ┌──────────▼────────────┐
                         │       cho_duyet       │ ← Xác thực thực địa (Lớp 2)
                         └──────────┬────────────┘
                 ┌──────────────────┴──────────────────┐
          BVTV duyệt                             BVTV từ chối
                 │                                     │
          ┌──────▼──────┐                       ┌──────▼──────┐
          │  da_duyet   │ ← Thẩm định (Lớp 3)   │  tu_choi    │
          └─────────────┘                       └──────┬──────┘
                                                       │
                                            HTX sửa → ghi lại bản mới
```

### 3.2 Farmer Bridge (Lớp chặn dữ liệu rác)
Để ngăn chặn tình trạng dữ liệu "ma", áp dụng **Xác nhận 2 lớp**:
1. **Lớp 1 (Ghi nhận)**: Kỹ thuật viên (KTV) nhập nhật ký thay nông dân → Trạng thái: `cho_nong_dan_xac_nhan`.
2. **Lớp 2 (Xác thực thực địa)**: 
   - *Nếu nông dân có App*: Nhập mã OTP hoặc bấm xác nhận.
   - *Nếu không có App*: KTV đính kèm ảnh chụp phiếu làm việc có chữ ký tay nông dân $\rightarrow$ Tính Hash `SHA256` ảnh và lưu vào chain $\rightarrow$ Trạng thái tự động chuyển sang `cho_duyet`.
3. **Lớp 3 (Thẩm định)**: BVTV thẩm định minh chứng (nội dung, ảnh chứng minh) $\rightarrow$ Quyết định chuyển sang `da_duyet` hoặc `tu_choi`.

> **BẮT BUỘC**: Phải đính kèm tệp minh chứng đối với các hoạt động nhạy cảm như (Phun thuốc, bón phân).
> **Quy tắc hiển thị QR**: Chỉ nhật ký đạt `da_duyet` mới hiển thị lên Portal Truy xuất nguồn gốc.

### 3.3 Quy tắc chuyển trạng thái

| Từ | Đến | Ai thực hiện | Điều kiện tiên quyết |
|----|-----|-------------|---------------------|
| _(mới tạo)_ | `cho_nong_dan_xac_nhan` | HTX | Tự động khi GhiNhatKy. Mặc định cần thêm bước ký. |
| `cho_nong_dan_xac_nhan` | `cho_duyet` | HTX (KTV) / Nông dân | Gửi kèm Hash hoặc mã OTP xác thực hợp lệ. |
| `cho_duyet` | `da_duyet` | ChiCucBVTV, Platform | Kiểm tra nội dung và Hash minh chứng hợp lệ. |
| `cho_duyet` | `tu_choi` | ChiCucBVTV, Platform | **BẮT BUỘC** điền `ly_do_tu_choi`. HTX phải ghi bản mới. |

### 3.4 Loại hoạt động — Danh mục động (Quản lý Backend/Frontend)

**Quy tắc**: Danh mục `loai_hoat_dong` do backend/frontend quản lý qua file mapping JSON/DB. Chaincode chỉ check `!= ""`.

| Mã (`value`) | Tên hiển thị | Nhóm | Giai đoạn lô hàng phù hợp |
|--------------|-------------|------|---------------------------|
| `gieo_hat` | Gieo hạt | Sản xuất | `dang_trong` |
| `bon_phan` | Bón phân | Sản xuất | `dang_trong` |
| `phun_thuoc` | Phun thuốc | Bảo vệ thực vật | `dang_trong` |
| `thu_hoach` | Thu hoạch | Thu hoạch | `dang_trong` → `da_thu_hoach` |
| `dong_goi` | Đóng gói | Sau thu hoạch | `da_thu_hoach` |
| `van_chuyen` | Vận chuyển | Logistics | `da_thu_hoach`, `san_sang_ban` |
| `khac` | Khác | Khác | Tất cả |

### 3.5 Validation rules Nhật Ký

| Field | Bắt buộc | Validation |
|-------|:--------:|-----------|
| `ma_nhat_ky` | ✅ | Unique, format: `NK-{seq}` |
| `ma_lo` | ✅ | Phải tồn tại trong `lohang_cc`. Lô hàng/lô con tương ứng. |
| `loai_hoat_dong` | ✅ | Không rỗng. Endpoint validate với file config. |
| `minh_chung_hash`| ✅ *(Có ĐK)* | Bắt buộc với hoạt động nhạy cảm. Chứa SHA256 Hex của ảnh chụp gửi từ app để xác thực thực địa (Farmer Bridge). |

---

## 4. Vòng đời Giao Dịch (`giao_dich_cc`)

### 4.1 Sơ đồ trạng thái (STRICT — enforce trong chaincode)

```
                    ┌───────────┐
       HTX tạo      │ cho_duyet │  ← Lô bắt buộc ở san_sang_ban, trừ kho tạm thời (Pending)
                    └─────┬─────┘
                          │ Platform duyệt
                    ┌─────▼─────┐
                    │ da_duyet  │
                    └─────┬─────┘
                          │ Platform xác nhận giao
                    ┌─────▼─────┐
                    │ dang_giao │  ← HTX đang giao hàng cho NPP
                    └─────┬─────┘
                          │ NPP / Platform xác nhận nhận hàng
                    ┌─────▼─────┐
                    │  da_giao  │  ← NPP đã nhận hàng
                    └─────┬─────┘
                          │ Platform chốt thanh toán
                 ┌────────▼────────┐
                 │ cho_thanh_toan  │  ← NPP cần thanh toán (hiện trong Công Nợ)
                 └────────┬────────┘
                          │ NPP / Platform xác nhận đã thanh toán
                 ┌────────▼────────┐
                 │ da_thanh_toan   │  ← TRẠNG THÁI KẾT THÚC. Hoa hồng hoàn tất tính.
                 └─────────────────┘

     ┌────────────────────────────────────────────────┐
     │ Bất kỳ trạng thái (ngoại trừ da_thanh_toan)   │
     │      → huy_bo (CHỈ PlatformOrgMSP)             │
     │  (Khi hủy, số lượng hoàn lại vào lô hàng)      │
     └────────────────────────────────────────────────┘
```

### 4.2 Cơ chế Thu hồi & Phản ứng nhanh (Recall Logic)
Hệ thống xử lý ngay lập tức khi phát hiện lô hàng kém chất lượng, nguy hại môi trường:
1. **Action (Đình Chỉ)**: BVTV hoặc Platform gọi hàm `CapNhatTrangThaiLo` đổi trạng thái Lô thành `dinh_chi`.
2. **Impact 1 (Giao dịch Treo)**: Backend tự động chặn/khóa mọi thao tác đối với các Giao dịch có trỏ tới Lô đó đang ở `dang_giao` hoặc `cho_duyet`. 
3. **Cảnh Báo (Tức thì)**: Cảnh báo NPP: **"HÀNG CẦN THU HỒI - LÔ HÀNG KHÔNG AN TOÀN"**.
4. **Impact 2 (Truy xuất đỏ)**: Consumer khi quét QR vào trực tiếp lô này (hoặc lô con từ nó) sẽ thấy màn hình Đỏ Toàn Phần cảnh báo: **"SẢN PHẨM KHÔNG AN TOÀN - ĐANG THU HỒI KIỂM TRA"** (Bị ẩn các chứng nhận đi).
5. **Recovery**: Phục hồi lô hàng từ `dinh_chi` chỉ được kích hoạt bởi chính Platform sau giám định pháp lý.

### 4.3 Validation rules Giao Dịch
| Field | Bắt buộc | Validation |
|-------|:--------:|-----------|
| `so_luong` | ✅ | > 0, ≤ `so_luong_con_lai` trong lô (Backend check & Trừ ngay khi tạo GD) |
| `don_gia` | ✅ | > 0, đơn vị VNĐ |
| `ty_le_hoa_hong` | ✅ | 0 ≤ x ≤ 100 (%) |

---

## 5. Liên kết giữa 3 Entity & Đồng bộ 

### 5.1 Xử lý Cross-Channel Failure & Atomic Traceability
Do hệ thống được phân ở 2 channels (`nhatky-htx-channel` & `giaodich-channel`), Backend phải đáp ứng:
* **Saga Pattern & Redundancy**: Khi tạo Giao dịch (Channel Giaodich) thành công, hệ thống phải thực hiện trừ kho ở `lohang_cc` (Channel Nhatky). Nếu xảy ra ngoại lệ rớt mạng/Timeout, Backend phải thiết lập tiến trình Retry ngầm tới khi khớp (`Eventual Consistency`). Tương tự cho Rollback khi `huy_bo`.
* **Atomic Query Truy vết**: Điểm nối (Join key) là `ma_lo` và `ma_lo_me`. Khi Consumer quét QR (chỉ cung cấp `ma_lo` của lô con), Backend API phải truy xuất đệ quy đọc ngược lên lấy toàn bộ nhật ký `da_duyet` và list giấy chứng nhận của lô gốc.

### 5.2 Quy tắc cross-entity
| # | Quy tắc | Layer enforce | Lý do |
|---|---------|:------------:|-------|
| R1| Nhật ký & Giao dịch phải gắn với 1 `ma_lo` tồn tại | Backend | Cross-check state 2 channel |
| R2| **Enforce inventory**: Tạo GD | Backend | Số lượng GD ≤ `so_luong_con_lai`. Trừ tồn kho đồng bộ. |
| R3| **Rollback inventory**: Hủy GD | Backend | Kích hoạt callback cộng lại hàng khi có GD bị rớt / hủy. |
| R4| **Auto het_hang** | Backend | Tự chuyển sang `het_hang` khi inventory về 0. |

---

## 6. Luồng End-to-End hoàn chỉnh

### 6.1 Happy Path — Từ gieo hạt đến thanh toán (Có Farmer Bridge)

```
Bước  Actor       Hành động                                 State sau
───────────────────────────────────────────────────────────────────────
 1    HTX         TaoLotHang("LH-001", "HTX001", ...)       LH: dang_trong
 2    HTX (KTV)   GhiNhatKy("NK-001", lo="LH-001",          NK: cho_nong_dan_xac_nhan
                    loai="gieo_hat", "Gieo hat ST25")
 3    Nông dân    Xác nhận SMS / App (Farmer Bridge)        NK: cho_duyet
 4    BVTV        DuyetNhatKy("NK-001", "da_duyet")          NK: da_duyet
 5    HTX         CapNhatTrangThaiLo("LH-001", ... )        LH: cho_chung_nhan
 6    BVTV        ThemChungNhan("LH-001", "VietGAP")        (Lô Hàng nhận Profile)
 7    HTX         CapNhatTrangThaiLo("LH-001", ... )        LH: san_sang_ban
 8    HTX         TaoGiaoDich("GD-001", lo="LH-001", ... )  GD: cho_duyet 
                                                            (LH: Bị trừ số lượng)
 9    Platform    DuyetGiaoDich("GD-001")                   GD: da_duyet
10    NPP         CapNhatTrangThai("GD-001", "da_giao")     GD: da_giao
11    NPP         CapNhatTrangThai("GD-001", "...")         GD: da_thanh_toan
12    [Backend]   Auto_het_hang (Khi tồn kho = 0)           LH: het_hang
13    Consumer    Quét QR → GET /api/v1/consumer/LH-001     (Hiện public profile)
```

---

## 7. Frontend Display Rules

### 7.1 Badge màu theo trạng thái

**Lô hàng:**
| Trạng thái | Màu badge | Icon | Hành động khả dụng |
|------------|-----------|------|--------------------|
| `dang_trong` | 🟢 Green | 🌱 | HTX: [Thu hoạch] / Platform: [Đình chỉ] |
| `da_thu_hoach` | 🟡 Yellow | 🌾 | HTX: [Gửi kiểm định] / Platform: [Đình chỉ] |
| `cho_chung_nhan`| 🟠 Orange | 🔍 | BVTV: [Cấp chứng nhận] / Platform,BVTV: [Đình chỉ] |
| `san_sang_ban` | 🔵 Blue | ✅ | HTX: [Tạo Giao Dịch] [Tách Lô] / Platform: [Đình chỉ] |
| `het_hang` | ⚫ Gray | 📦 | *(Không có hành động)* |
| `dinh_chi` | 🔴 Red | ⛔ | Platform: [Phục hồi → trạng thái phù hợp] |

**Nhật ký:**
| Trạng thái | Màu | | Trạng thái | Màu |
|------------|-----|-|------------|-----|
| `cho_nong_dan_xac_nhan`| 🟡 Vàng Nhạt | | `da_duyet` | 🟢 Xanh |
| `cho_duyet`| 🟠 Cam | | `tu_choi` | 🔴 Đỏ |

### 7.2 Màn hình QR Truy Xuất Nguồn Gốc (Consumer Page)
- 🔴 **TRƯỜNG HỢP KIỂM TRA/THU HỒI `dinh_chi`**: Hiện banner cảnh báo khẩn cấp màu đỏ thay cho VietGAP và ẨN tất cả dữ liệu xác minh.
- 🟢 **TRƯỜNG HỢP BÌNH THƯỜNG**: Load < 3 giây màn hình Mobile. Chỉ hiển thị thành phần `da_duyet`. List hoạt động sản xuất xếp Timeline Ascending. Nếu lô con, back-tracing hiển thị mọi thông tin kết nối từ lô cha. Tuyệt đối không hiển thị giá thành bán buôn/chiết khấu.

---

## 8. Checklist Review Code theo Business Rules

Dùng checklist này khi review chaincode, API, hoặc frontend:

- [ ] **Offline Mode (Frontend)**: Hỗ trợ SQLite Cache + Background Sync cho phép KTV ghi nhận nhật ký (trong rừng/núi) khi mất mạng và tự sync khi kết nối lại.
- [ ] **Farmer Bridge**: Bắt buộc có cơ chế yêu cầu hình ảnh đính kèm (Tính hash tải lên chain) với các công tác xử lý phân bón/xịt thuốc.
- [ ] **Data Integrity**: Các truy vấn trạng thái/tạo giao dịch phải đảm bảo `so_luong_con_lai >= 0`. Không bao giờ âm.
- [ ] **Back-Tracing**: Đảm bảo trường `ma_lo_me` được lưu trữ. Fetch QR API hoạt động bằng Join Query giữa Lô Cha (thông tin farm) và Lô Con.
- [ ] **Security**: Identity cho hàm `dinh_chi` bị hard-check vào `ChiCucBVTVOrgMSP` hoặc `PlatformOrgMSP`.
- [ ] **Cross-Check**: Giao dịch sinh ra trừ tạm kho (Logic Saga/Retry) thay vì đợi chốt lệnh GD xong mới trừ, tạo tình huống "overselling".
- [ ] **Timestamp Validity**: Sử dụng thuần `GetTxTimestamp()` phục vụ Deterministic block, cấm dùng `time.Now()`.
