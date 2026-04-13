---
name: gapchain-landing-page
description: >
  Hướng dẫn xây dựng Landing Page cho dự án GapChain trên nền tảng TrustFab (Hyperledger Fabric v3.1.x).
  Tập trung vào trải nghiệm người dùng nông nghiệp (HTX, NPP) với phong cách Eco-Modern, 
  ngôn ngữ đơn giản hóa công nghệ và tối ưu hóa chuyển đổi.
---

# Skill: Xây dựng Landing Page GAPChain — Gắn kết Niềm tin Nông sản
## Design Style: Eco-Modern · Clean & Trust · High Conversion

## 1. Tư duy Thiết kế (Design Thinking)
* **Tone & Mood:** Sử dụng màu **Xanh lá (Emerald/Forest Green)** đại diện cho nông nghiệp sạch, kết hợp với **Xanh dương đậm (Navy Blue)** của công nghệ Fabric và màu **Trắng/Xám nhạt** để tạo sự thông thoáng.
* **Eco-Modern Style:** Sử dụng hình ảnh thực tế về nông sản Việt (Gạo ST25, Cà phê...) kết hợp với các hiệu ứng đường kẻ, node mạng (network nodes) tinh tế để thể hiện sự kết nối Blockchain.
* **Bình dân hóa công nghệ:** Tuyệt đối không dùng thuật ngữ quá hàn lâm. 
    * Thay vì "Decentralized Ledger" $\rightarrow$ Hãy nói "Sổ cái điện tử không thể tẩy xóa".
    * Thay vì "Smart Contract" $\rightarrow$ Hãy nói "Hợp đồng điện tử tự động".

## 2. Cấu trúc nội dung (Information Architecture)

| Section | Mục tiêu (BA Perspective) | Nội dung chi tiết |
| :--- | :--- | :--- |
| **Hero Section** | Tạo ấn tượng đầu tiên | Headline: "Minh bạch hóa hành trình Nông sản Việt". Sub: Giới thiệu TrustFab (Fabric v3.1.x). CTA: "Xem Demo MVP". |
| **The Pain Point** | Chạm vào nỗi đau | Nêu bật vấn đề: Hàng giả, thiếu tin tưởng, gánh nặng hồ sơ giấy VietGAP. |
| **Solution (TrustFab)** | Khẳng định sức mạnh | Giới thiệu TrustFab - Hạ tầng Blockchain doanh nghiệp giúp dữ liệu bất biến, bảo mật và tin cậy. |
| **Benefits by Role** | Cá nhân hóa đối tượng | Cột HTX: Số hóa nhật ký, tăng giá bán. Cột NPP: Kiểm soát chất lượng, minh bạch nguồn cung. |
| **Core Features** | Tính năng MVP | 1. Nhật ký số VietGAP; 2. QR Code truy xuất hành trình; 3. Quản lý giao dịch & Công nợ. |
| **Social Proof** | Xây dựng niềm tin | Logo các đối tác tiềm năng, các chứng chỉ tiêu chuẩn hệ thống hỗ trợ (VietGAP, GlobalGAP). |
| **Call to Action** | Chốt hạ chuyển đổi | Form đăng ký "Thí điểm không rủi ro" (Pilot Program). |

## 3. Tech Stack đề xuất (MVP Speed)
* **Frontend:** Vue 3 + Vite & Tailwind CSS.
* **Animations:** Framer Motion (hiệu ứng mượt mà khi cuộn trang).
* **Icons:** Lucide Icons (phong cách line-art hiện đại).
* **Assets:** Ảnh thực tế từ HTX kết hợp minh họa luồng dữ liệu Blockchain trực quan.

## 4. Checklist triển khai

### Nội dung & Nghiệp vụ (BA)
- [ ] Headline tập trung vào lợi ích (Tăng lợi nhuận, Giảm rủi ro) thay vì tính năng kỹ thuật.
- [ ] Có phần giới thiệu riêng cho nền tảng **TrustFab** và Hyperledger Fabric v3.1.x.
- [ ] Các thuật ngữ kỹ thuật đã được "nông thôn hóa" để dễ hiểu nhất.

### Thiết kế & Trải nghiệm (UX/UI)
- [ ] **Mobile-first:** Giao diện tối ưu cho điện thoại (thiết bị chính của nông dân).
- [ ] **Interactive QR:** Có mẫu QR Code để khách truy cập quét thử xem trang truy xuất nguồn gốc thực tế.
- [ ] Tốc độ load trang < 2s để đảm bảo trải nghiệm ở vùng sóng yếu.

### Kỹ thuật (Dev)
- [ ] Tối ưu SEO cho từ khóa: "truy xuất nguồn gốc blockchain", "vietgap digital diary", "gapchain".
- [ ] Gắn Tracking (Google Analytics/Hotjar) để phân tích hành vi người dùng (HTX hay NPP quan tâm phần nào hơn).

## 5. Mẫu Copywriting cho Hero Section
> **Headline:** Nâng tầm giá trị Nông sản sạch với Công nghệ Blockchain hàng đầu.
>
> **Sub-headline:** GAPChain chạy trên nền tảng **TrustFab** (Hyperledger Fabric v3.1.x) giúp Hợp tác xã số hóa quy trình VietGAP và giúp Nhà phân phối kiểm soát nguồn gốc thực phẩm chỉ bằng một cú quét mã.
>
> **CTA Button:** [Đăng ký Thí điểm Miễn phí]