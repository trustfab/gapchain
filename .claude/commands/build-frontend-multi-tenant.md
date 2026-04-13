---
name: build-frontend-multi-tenant
description: >
  Hướng dẫn xây dựng Web Frontend sử dụng Vue.js 3 kết nối với Backend API GAPChain.
  Đặc tả chuyên sâu về cơ chế nhận diện Multi-tenant (Đa Khách Hàng) để phân quyền 
  và hiển thị dữ liệu Lô Hàng, Nhật Ký, Giao dịch theo từng Role (Platform, HTX, NPP, BVTV).
  Kết hợp bộ tiêu chuẩn UI/UX "WOW Factor" cao cấp và tuân thủ chặt chẽ Business Rules (SKILL.md).
---

# Skill: Xây dựng GAPChain Web Frontend (Vue 3)

Tài liệu này là bộ nguyên tắc chỉ đạo dành cho AI để xây dựng ứng dụng Web dành cho mạng GAPChain bằng Vue 3. Ứng dụng này cung cấp 2 mảng không gian mạng: **Dashboard quản trị cho các tác nhân** và một **Portal Truy Xuất cho Người dùng cuối** (Public).

Tất cả các logic giao diện phải tuân thủ nghiêm ngặt **GAPChain Business Rules (`gapchain-ba/SKILL.md`)**.

## 1. Yêu Cầu Căn Bản (Tech Stack)

* **Framework Core:** Vue 3 (Composition API, `<script setup>`) + Vite
* **State Management:** Pinia
* **Routing:** Vue Router 4
* **Networking:** Axios (Kèm cấu hình Interceptors)
* **Aesthetics CSS:** Tailwind CSS (Đòi hỏi khả năng dựng layout sắc nét)
* **Bản Đồ Nâng Cao:** Sử dụng `leaflet` kết hợp `vue-leaflet`

## 2. Tiêu Chuẩn Thẩm Mỹ Giao Diện (WOW Effect cho Demo MVP)

Vì sản phẩm nhắm tới việc chạy Pitching & Demo:
1. **Thiết kế Glassmorphism**: Trang không dùng cấu trúc khối hộp trắng cứng nhắc. Lạm dụng kỹ thuật nền mờ (Backdrop Blur), Gradient chuyển sắc giữa Xanh lá cây - Vàng chanh nông nghiệp.
2. **Micro-animations Chuyển Trang**: Bọc `<router-view>` bằng `<Transition name="fade">`. Bo góc mượt `rounded-2xl`, thêm shadow sáng glow khi Hover vào các nút CTA (Call-to-action).

## 3. Kiến Trúc Multi-tenant ở State Store (Pinia)

API Backend trả về Token JWT đính kèm Payload định tuyến đa khách hàng:
```json
{
  "username": "htx001",
  "role": "htx",                 // platform | htx | bvtv | npp
  "msp_id": "HTXNongSanOrgMSP",  
  "tenant_id": "HTX001"          // Mã mapping cụ thể cho dữ liệu hệ thống
}
```
**Kiến trúc Pinia Auth Store (`stores/auth.js`)**:
- Trích xuất toàn bộ `token`, `role`, `tenantId`, `mspId` sau khi login thành công để đưa vào `localStorage`.
- Khai báo các Getter: `isLoggedIn`, `isPlatform`, `isHTX`, `isNPP`, `isBVTV`.

## 4. Trải Nghiệm Đăng Nhập "One-Click Demo" & Vue Router Guard

Hệ thống Router cần tuân thủ cấu trúc Guard (`router.beforeEach`) chặt chẽ:
* Nếu gọi REST bị lỗi mã 401 Unauthorized -> Call `authStore.logout()` và `router.push('/login')`.
* Route phân tán: `/htx/...` (Cho HTX), `/npp/...` (Cho NPP), `/platform/...` (Cho Admin hệ thống). Base URL khi fetch Axios dựa vào `tenantId`.

⚡ **Màn Hình Login (WOW Factor)**: 
Cấm thiết kế Form nhập liệu gõ text trống trơn nhàm chán. Màn hình login sẽ hiển thị **4 Thẻ Giao Diện (Role Cards)**:
1. 🧑‍🌾 Thẻ **HTX**: Cảnh nông dân (On click -> Tự điền `htx001` và Login tự động).
2. 🚚 Thẻ **NPP**: Cảnh xe tải kho bãi (On click -> Tự điền `npp001`).
3. 🛡️ Thẻ **BVTV**: Cảnh lab kiểm định (On click -> Tự điền `bvtv`).
4. 🌐 Thẻ **Platform**: Cảnh bản đồ mạng lưới (On click -> Tự điền `platform`).

## 5. Tổ Chức Cấu Hình Interceptor Mạng (Axios)
Tạo `src/plugins/axios.js`:
- Tự động gắn kèm Header `Authorization: Bearer <token>`.
- Chú ý nối param linh động `tenant_id` lấy từ Pinia Store với các API trích xuất đặc cách (Ví dụ: `GET /api/v1/giaodich/htx/${authStore.tenantId}`).

## 6. Các Phân Hệ Tính Năng Cốt Lõi (Modules) & Business Rules

UI phải thiết kế dựa trên các quy tắc trong `SKILL.md`:

1. **Quản lý Nông Trại (`Role: HTX`)**: 
   - Danh sách **Lô Hàng** hiển thị theo Badge trạng thái (`dang_trong` 🟢, `da_thu_hoach` 🟡, `cho_chung_nhan` 🟠, `san_sang_ban` 🔵, `het_hang` ⚫, `dinh_chi` 🔴).
   - Hỗ trợ luồng **Tách Lô** khi lô đang ở `san_sang_ban` (Lô con kế thừa dữ liệu lô mẹ).
   - Module ghi nhật ký phải tích hợp giao diện **Farmer Bridge**: yêu cầu tải/chụp ảnh kèm chữ ký (Lớp xác thực 2) đối với các hành động nhạy cảm (bón phân, phun thuốc).
   - Hỗ trợ Offline Mode: Lưu cache LocalStorage/IndexedDB khi mất mạng và đồng bộ (Sync) lại khi có kết nối.

2. **Kênh Chấp Thuận (`Role: BVTV`)**: 
   - Duyệt nhật ký: Check minh chứng (hash ảnh). Chuyển trạng thái từ `cho_duyet` sang `da_duyet` (🟢) hoặc `tu_choi` (🔴 - kèm lý do).
   - Duyệt Lô hàng và cấp chứng nhận (VietGAP).
   - Có nút **[Đình chỉ]** khẩn cấp cho Lô hàng vi phạm.

3. **Đơn Hàng & Công Nợ (`Role: NPP`)**: 
   - Cảnh báo tức thì (Banner đỏ) nếu có giao dịch trỏ tới lô hàng bị chuyển sang `dinh_chi` (Recall Logic).
   - Xác nhận chuyển trạng thái giao dịch (`da_giao`, `da_thanh_toan`).

4. **Platform & Ecosystem Overview (`Role: Platform`)**: 
   - Live Activity Feed. Thống kê số lượng.
   - Thẩm quyền cao nhất: Phục hồi trạng thái từ `dinh_chi`, hủy giao dịch (`huy_bo`).

## 7. Đỉnh Cao Trải Nghiệm (Public QR Portal)

Route cho ngưởi dùng cuối không đăng nhập: `/qr/consumer/:maLo`
Màn Render này quyết định 90% sự thuyết phục dự án:
* **Huy Hiệu Xác Thực Blockchain**: Icon Shield lớn trên cùng nhấp nháy dòng *"Được bảo vệ bởi Hyperledger Fabric"*. Rê chuột hiển thị TxID hash tường minh.
* **Cảnh Báo Thu Hồi Khẩn Cấp (Recall Logic)**: Nếu trạng thái Lô (hoặc Lô mẹ) là `dinh_chi`, CHE toàn màn hình bằng Banner Đỏ: **"SẢN PHẨM KHÔNG AN TOÀN - ĐANG THU HỒI KIỂM TRA"** (Ẩn chứng nhận).
* **Bản Đồ Canh Tác (Leaflet Interactive Map)**: App gọi API từ `vi_tri`, cắm `Pin Marker` tọa độ GPS.
* **Timeline Kể Chuyện Kéo Dài**: Chỉ hiển thị nhật ký `da_duyet`. Kéo timeline các hoạt động. Đặc biệt hỗ trợ truy ngược (Back-tracing) lịch sử nếu là Lô Con tách ra từ Lô Cha.

## 8. Lệnh Khởi tạo Mẫu (Run CLI)

AI phải tuân thủ việc thi công các tính năng với các Command mặc định sau, triển khai tại `gapchain/`:
```bash
npm create vue@latest frontend-web -- --default
cd frontend-web
npm install axios vue-router pinia tailwindcss postcss autoprefixer
npm install leaflet vue-leaflet
npx tailwindcss init -p
```
Triển khai kỹ lưỡng logic `stores/auth.js` trước, sau đó xây Interceptors của Axios, sau cùng mới đâm sâu vào các Trang giao diện theo đúng tinh thần "Multi-tenant + Wow-factor" đã mô tả ở trên.
