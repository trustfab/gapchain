# Hướng Dẫn Kỹ Năng Tận Dụng AI Agent (Claude, Antigravity) Trong Dự Án

Tài liệu này tổng hợp các kỹ năng và phương pháp tốt nhất dành cho Developer để sử dụng các AI Agent (như Claude Code, Antigravity, Cursor, v.v.) trong quá trình xây dựng, bảo trì và phát triển tính năng cho dự án GAPChain dựa trên cơ chế khai báo **SKILL** (kỹ năng/nghiệp vụ đặc thù dự án).

---

## 1. Khái niệm SKILL trong tương tác với AI Agent

Trong bối cảnh làm việc với AI Agent, **SKILL** (được định nghĩa trong các thư mục như `.claude/skills/` hoặc các file document trung tâm như `CLAUDE.md`, `SKILL.md`) là một tập hợp các tài liệu Markdown chứa:
- **Định nghĩa tác vụ (Business Rules - BA):** Các luồng trạng thái, vòng đời của đối tượng, ma trận phân quyền.
- **Quy tắc kỹ thuật (Technical Guidelines):** Kiến trúc, tech stack, thư viện bắt buộc (ví dụ: dùng `fabric-gateway-go` thay cho bản cũ).
- **Bối cảnh dự án (Context):** Sơ đồ quan hệ, các trường dữ liệu quan trọng, cách thức tổ chức cấu trúc file.

**Tại sao cần thiết lập SKILL?**
Các mô hình AI rất mạnh trong việc sinh code, nhưng lại dễ bị "ảo giác" (hallucination) — tức là sáng tạo ra các trạng thái (states) không tồn tại, sử dụng sai thư viện, hoặc viết logic bỏ qua các quy tắc nghiệp vụ ngầm. 
Tập tin SKILL đóng vai trò là **SSOT (Single Source of Truth - Nguồn sự thật duy nhất)** để "neo" tư duy của AI vào sát nhất với yêu cầu nghiệp vụ thực tế.

---

## 2. Cách thiết lập và tổ chức SKILL hiệu quả

### 2.1. File System Context (`CLAUDE.md`, `README.md`)
- **Vị trí:** Đặt tại gốc của dự án (`gapchain/CLAUDE.md`).
- **Nội dung:** Chứa các bức tranh tổng view toàn hệ thống: Tech Stack là gì, danh sách Commands build/deploy, Architecture cơ bản.
- **Quy trình sử dụng:** AI Agent sẽ tự động nạp các file nhận diện ngữ cảnh dự án (như `CLAUDE.md`) ngay khi load khởi tạo phiên làm việc.

### 2.2. Folder SKILL chuyên biệt (`.claude/skills/<tên-kỹ-năng>`)
- **Cấu hình YAML Headers (Metadata Trigger):** Việc khai báo dòng `name:` và `description:` kết hợp keywords thông minh (như "nghiệp vụ", "trạng thái", "business rule") tại đầu file `SKILL.md` cho phép AI (đặc biệt là Claude Code/Cursor) tự phân tích dự đoán và inject document này vào ngữ cảnh mỗi khi trao đổi đi qua các từ khoá đó.
- Lợi ích của việc này là khi user chat tự do, AI sẽ có bộ nhớ tự động lôi file SKILL ra đọc mã không cần nhắc thủ công ở mỗi prompt.

---

## 3. Kỹ năng giao tiếp Command/Prompt cho Developer

Để tận dụng tối đa `SKILL`, Developer nên thay đổi thói quen viết Prompt từ "Yêu cầu hành động" sang "Yêu cầu đọc context trước khi hành động".

### Kỹ năng 1: Prompts Ép buộc (Enforce) Đối chiếu Tài liệu
Thay vì để AI tự đoán logic theo common sense, hãy "ép" nó dùng SKILL chuẩn:
> ❌ **Sai:** "Viết hàm chuyển trạng thái lô hàng sang đình chỉ."
> ✅ **Đúng:** "Dựa vào `SKILL.md`, trích xuất các điều kiện để chuyển Lô Hàng sang `dinh_chi` (Ai có thể thực hiện? Lúc nào có thể thực hiện?), sau đó viết tiếp logic `CapNhatTrangThaiLo` bảo đảm quy tắc kiểm tra chứng thực."

### Kỹ năng 2: Sử dụng Checklists để Review Code
Trong `SKILL.md` (phần 8) có đính kèm Checklist Review Code. Có thể yêu cầu AI tự động rà soát file thông qua checklist này.
> **Prompt:** "Review tệp `nhatky_cc.go` và đánh giá chéo theo mục *8. Checklist Review Code* trong file BA SKILL. Đảm bảo việc dùng `time.Now()` đã được gỡ bỏ và dùng đúng `GetTxTimestamp()`."

### Kỹ năng 3: Ràng buộc Kiến trúc kỹ thuật
Để tránh AI code tuỳ hứng:
> **Prompt:** "Tạo dự án mới backend API cho đối tượng HTX. Tham khảo cấu trúc trong mục *Quy tắc kỹ thuật chaincode* của `@gapchain/CLAUDE.md`, sử dụng đúng phiên bản `fabric-contract-api-go v1.2.2`. Tường minh các JSON struct theo chuẩn `snake_case`."

---

## 4. Ứng dụng thực chiến SKILL trong GAPChain

Dưới đây là một số phương thức Workflow ứng dụng cho hệ gen AI:

- **Back-Tracing / State Machine Validation:**
  Yêu cầu AI dựa trên sơ đồ thiết kế luồng lô hàng trong file referance để vá lại logic chaincode: *"Kiểm tra luồng `cho_duyet` sang `da_duyet` của giao dịch. Hệ thống có cho phép vượt cấp từ `dang_trong` không? Dựa vào sơ đồ Lô Hàng STRIKE, code bảo vệ cho tôi cơ chế này."*

- **Frontend & UX Rules Sync:**
  Khi gen UI Frontend: *"Tạo một component Badge trạng thái trong Vue 3. Hãy tự động tham chiếu bảng màu (Green, Yellow, Orange...) cho từng trạng thái Lô Hàng và Nhật ký từ chương 7 phần Frontend Display Rules trong SKILL."*

- **Tạo Saga / Fallback:**
  Yêu cầu AI tự đọc hiểu cấu trúc Microservice/Blockchain để xử lý: *"Đọc phần 5.1 (Saga Pattern & Redundancy) trong SKILL. Hãy viết middleware Golang đảm nhận logic rollback tồn kho lô hàng khi giao dịch `cho_duyet` bị từ chối."*

---

## 5. Đúc kết kinh nghiệm cho Đội ngũ Developer

1. **Docs as Code (Sửa Document trước khi sửa Code):** Khi xảy ra sự kiện thay đổi Business Rule từ phía BA, hành động đầu tiên của dev là CẬP NHẬT vào `SKILL.md` và `referance.md`. Sau đó chia sẻ file cập nhật này ra cho AI tự bảo trì mã nguồn để sửa tất cả backend/frontend đồng bộ luồng đó. Bản thân Document lúc này điều hướng hành vi Code.
2. **Sử dụng `@` Explicit Context:** Trong Antigravity hoặc IDE Agent, dùng cú pháp định danh tệp `@gapchain/.claude/skills/gapchain-ba/SKILL.md` một cách cụ thể mỗi lúc xử lý ticket phức tạp, đảm bảo session luôn chứa context.
3. **Giữ SKILL sắc bén:** Tránh thêm các thông tin rườm rà không phải quy luật/luật lệ vào file SKILL. Giữ Enum rõ ràng, bảng biểu tường minh vì AI giải mã thông tin có cấu trúc bảng biểu tốt hơn văn xuôi (Ví dụ: Các bảng ma trận phân quyền, matrix trạng thái là định dạng tốt nhất cho AI).
