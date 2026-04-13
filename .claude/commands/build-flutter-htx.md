---
name: build-flutter-htx
description: Hướng dẫn xây dựng Flutter Mobile App cho Hợp tác xã (HTX) thuộc hệ sinh thái GAPChain. Sử dụng skill này để tạo/thiết kế UI, code logic bằng Riverpod, hoặc API ghi nhật ký offline-first trên Mobile.
---

# GAPChain - Build Flutter Mobile App (HTX)

Hỗ trợ xây dựng ứng dụng di động Flutter chuyên dụng cho người dùng thuộc vai trò **Hợp tác xã (HTX)**. Ứng dụng tập trung vào việc quản lý Lô hàng, Ghi nhật ký canh tác (có ảnh, định vị GPS, hash ảnh), chia sẻ mã QR truy xuất và theo dõi Giao dịch.

## 1. Cấu trúc thư mục (Features-First)

Khi khởi tạo hoặc thêm mới màn hình, ứng dụng thiết lập theo cấu trúc Features-First kết hợp Clean Architecture:

```text
gapchain/flutter_app/
├── android/
├── ios/
├── assets/
│   └── loai_hoat_dong.json            # Data seeding tĩnh (offline)
├── lib/
│   ├── core/                          # Setup chung toàn app
│   │   ├── api/                       # Dio client & Interceptors (JWT)
│   │   ├── db/                        # Cấu hình sqflite (Local Database)
│   │   ├── theme/                     # AppTheme, colors, text styles
│   │   └── utils/                     # Helpers (DateFormatter, GPS, Hash)
│   ├── features/
│   │   ├── auth/                      # Login & Cache Token
│   │   ├── lohang/                    # Lô hàng
│   │   │   ├── data/                  # API client, Models DTO
│   │   │   └── presentation/          # Screens List/Create, Riverpod Controllers
│   │   ├── nhatky/                    # Nhật ký canh tác (Form nhập offline, timeline)
│   │   ├── giaodich/                  # Tạo và theo dõi tiến độ giao dịch bán
│   │   └── qr/                        # Generate QR (chia sẻ) & Scan QR
│   ├── shared/widgets/                # Custom widgets dùng nhiều nơi
│   │   ├── trang_thai_badge.dart
│   │   └── photo_capture.dart
│   └── main.dart                      # ProviderScope, GoRouter config
└── pubspec.yaml
```

## 2. Dependencies Bắt Buộc

Trong `pubspec.yaml` cần đảm bảo tối thiểu các packages:
- **State & DI**: `flutter_riverpod` (v2+)
- **Navigation**: `go_router`
- **Network**: `dio`
- **Storage/Offline**: `sqflite`, `flutter_secure_storage`
- **Hardware**: `geolocator`, `image_picker`, `connectivity_plus`
- **Utils**: `crypto` (hash ảnh SHA-256), `intl` (date format)
- **QR Code**: `qr_flutter` (tạo mã), `qr_code_scanner` (quét mã)

## 3. Core Logic & Offline-First UX

Khác với app thông thường, nông nghiệp thường diễn ra ngoài đồng hoặc vùng sóng 3G/4G chập chờn. Ứng dụng cần hỗ trợ "Offline-First".

### 3.1 Flow Nhập Nhật Ký
1. Nông dân chọn Lô hàng.
2. Khi mở Form Nhật Ký, ứng dụng gọi `geolocator` lấy GPS ngầm tự động và khóa input.
3. Chụp ảnh bắt buộc -> tính hash SHA256 ngay lập tức.
4. Bấm Lưu -> Lưu bản Record xuống bảng `nhatky_pending` trong `sqflite` với trạng thái chưa đồng bộ. Cập nhật thẳng UI.

### 3.2 Tự Động Đồng Bộ (Background/Manual Sync)
Xây dựng một `SyncService` check trạng thái mạng với `connectivity_plus`. Nếu online:
```dart
class SyncService {
  final LocalDb localDb;
  final ApiClient apiClient;

  Future<void> syncPendingNhatKy() async {
    final pending = await localDb.getPendingNhatKy();
    for (final nk in pending) {
      try {
        // Post nhật ký lên API Backend
        await apiClient.ghiNhatKy(nk); 
        // Bắn khỏi danh sách pending local
        await localDb.markSynced(nk.id); 
      } catch (e) { 
        // Bỏ qua, retry khi restart lại app hoặc có mạng ổn định
      }
    }
  }
}
```

## 4. Giao diện (UI/UX Guidelines)
- **Tối giản thao tác**: Form điền nên xài Dropdown list có sẵn các hoạt động (bón phân hữu cơ, bắt sâu,...). Tránh bắt người dùng gõ bàn phím vì thao tác gõ tại vườn/ruộng rất khó.
- **Font & Size**: Font chữ phải to, đậm (>= 16sp). Buttons dễ chạm với Action khu vực lớn.
- **Màu sắc**: Màu chủ đạo màu Xanh lá (tượng trưng nông nghiệp/phát triển bền vững, VD: Green, Forest Green). 
- **QR Truy xuất**: Tự động sinh giao diện chia sẻ mã QR theo chuẩn bằng `qr_flutter` truyền chuỗi `ma_lo`.

## 5. Vai Trò Claude (Nhiệm vụ của Agent)
Khi user gọi lệnh thiết lập tính năng cho Mobile, vui lòng:
1. Xác định User muốn code phần `Auth`, `Lô hàng`, `Nhật ký`, hay cấu hình `Utils offline`.
2. Generate mã nguồn `Dart` bám chặt state bằng `StateNotifierProvider` hoặc `Notifier` trên thư viện `riverpod` (tránh boilerplate).
3. Đảm bảo mọi giao tiếp Http/Dio đều truyền JWT qua interceptors.
4. Trẻ trung, năng động, mang đến UI "wow", dynamic effects khi có thể. Đưa đầy đủ snippet.
