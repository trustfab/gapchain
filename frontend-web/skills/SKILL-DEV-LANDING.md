Bạn là một Senior Frontend Developer chuyên xây dựng 
Developer-facing Landing Page cho các dự án Open Source.

Hãy tạo một landing page hoàn chỉnh bằng Vue 3 + Tailwind CSS 
cho dự án GAPChain — một MVP truy xuất nguồn gốc nông sản sạch 
đang chạy thực tế ở local, xây dựng trên Hyperledger Fabric 3.1.1 
+ Golang Backend + Flutter Mobile.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 🎯 MỤC TIÊU DUY NHẤT CỦA TRANG

Không bán sản phẩm. Không kêu gọi nông dân hay doanh nghiệp.

Mục tiêu: Truyền cảm hứng và kêu gọi các nhà phát triển 
(Developer, DevOps, Blockchain Engineer) mạnh dạn đóng góp 
vào một dự án Blockchain thực chiến có giá trị xã hội thực tế 
tại Việt Nam — hiện đang cần community để scale từ local 
lên production.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 👥 ĐỐI TƯỢNG MỤC TIÊU

- Blockchain Engineer muốn thực chiến với Hyperledger Fabric
- Backend Developer (Golang) muốn làm dự án có impact
- DevOps Engineer muốn contribute CI/CD, cloud deployment
- Flutter Developer muốn build app có người dùng thật
- Tech Lead muốn mentor dự án xã hội ý nghĩa

Tâm lý của họ:
- Chán side project vô nghĩa, muốn portfolio có dự án thực tế
- Bị thu hút bởi tech stack enterprise-grade và bài toán thú vị
- Muốn được credit công khai khi contribute
- Thích peer learning trong cộng đồng open-source

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 🏗️ THÔNG TIN KỸ THUẬT THỰC TẾ CỦA DỰ ÁN

(Dùng chính xác các thông tin này trong trang — không bịa)

**Tech Stack đã build:**
- Hyperledger Fabric 3.1.1 (4 org: Platform, HTXNongSan, 
  ChiCucBVTV, NPPXanh — 1 orderer, 4 peer, 4 CouchDB)
- 2 Channel: nhatky-htx-channel / giaodich-channel
- 3 Chaincode viết bằng Go:
  · lohang_cc  — Lô hàng & chứng nhận VietGAP/GlobalGAP
  · nhatky_cc  — Nhật ký canh tác gắn lô hàng, audit trail
  · giao_dich_cc — Giao dịch B2B, hoa hồng NPP, Smart Contract
- Backend: Golang + Gin + Fabric Gateway SDK v1.x (REST API)
- Mobile: Flutter (offline-first, GPS, ảnh SHA256 hash)
- Web: Vue 3 + QR Consumer Portal
- State DB: CouchDB (rich query) + GetHistoryForKey audit trail

**Bài toán kỹ thuật thú vị đã giải:**
1. Token hóa sản lượng: 1 tấn nông sản = số token tương đương,
   không thể xin thêm QR để trà trộn hàng giả
2. Bất biến tuyệt đối: dữ liệu sau kiểm duyệt bị lock trên 
   Hyperledger Fabric, GetHistoryForKey cho audit trail đầy đủ
3. Anti-clone QR địa lý: phát hiện 1 mã bị photocopy 
   và quét ở 2 địa điểm khác nhau trong 60 giây → cảnh báo đỏ
4. Burn token on-chain: người mua xác nhận → token bị đốt (Burn),
   tem không thể tái kích hoạt lần nào nữa
5. MSP identity enforcement: quyền thực thi từng hàm chaincode 
   được kiểm soát chặt theo org (HTX/ChiCuc/NPP/Platform)

**Hiện trạng thực tế:**
- MVP đang chạy ở môi trường local
- 9 Docker containers đang chạy ổn định
- Chaincode đã commit trên cả 2 channel
- Go client demo cho HTX và NPP đã hoạt động

**Cần community giúp:**
- CI/CD pipeline (GitHub Actions + Docker Compose)
- Deploy lên cloud (AWS / GCP / VPS bare metal)
- Monitoring & Alerting (Prometheus + Grafana)
- Multi-org network mở rộng thực tế
- Unit test & Integration test Chaincode (coverage > 70%)
- Performance benchmark Fabric network
- Flutter app hoàn thiện cho nông dân
- Fabric Explorer tích hợp

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 📐 CẤU TRÚC CÁC SECTION (theo thứ tự)

### SECTION 1 — NAVBAR
- Logo: chữ "G" gradient xanh lá + "GAPChain"
- Menu: Kiến Trúc / Tech Stack / Roadmap / Đóng Góp
- Badge nhỏ: "Local MVP · Open Source"
- CTA: "Star on GitHub →" (link #)

---

### SECTION 2 — HERO
Tone: Dev nói chuyện với dev — không phải quảng cáo

Headline:
"Blockchain Thực Chiến.
Bài Toán Thật. Cộng Đồng Thật."

Subheadline:
"GAPChain MVP đang chạy local trên Hyperledger Fabric 3.1.1
với 4 org, 2 channel, 3 chaincode — và cần những developer
dũng cảm đưa nó lên production."

2 CTA button:
- Primary: "Khám Phá Kiến Trúc" (scroll xuống)
- Secondary: "⭐ Star on GitHub" (link #)

3 badge nhỏ dưới CTA:
[MIT License] [Hyperledger Fabric 3.1.1] [Local → Production]

Thêm fake terminal block hiển thị lệnh thực tế:

  $ ./scripts/setup-network.sh
  ✅ 9 containers running (orderer + 4 peer + 4 couchdb)

  $ peer chaincode invoke -C nhatky-htx-channel -n lohang_cc \
    -c '{"function":"TaoLotHang","Args":["LH-HTX001-2025-001",
    "HTX001","Gạo ST25","lua","1000","kg","Đông Xuân 2025","Cần Thơ"]}'

  ✅ Chaincode invoke successful. TxID: a3f9c2...

---

### SECTION 3 — BÀI TOÁN KỸ THUẬT
Tiêu đề: "Tại Sao Bài Toán Này Xứng Đáng Với Thời Gian Của Bạn?"

5 card — mỗi card là 1 bài toán kỹ thuật thật đã giải:

Card 1 — Token hóa sản lượng:
Thu hoạch đúng 1 tấn → hệ thống chỉ phát hành số token 
tương đương 1 tấn. Không thể xin thêm QR để trà trộn hàng giả.
→ Bài toán: Chaincode lifecycle + quantity enforcement

Card 2 — Audit trail bất biến:
GetHistoryForKey trả về toàn bộ lịch sử thay đổi của 
từng lô hàng. Không ai — kể cả admin — có thể sửa 
ngược dữ liệu sau khi đã commit lên Fabric.
→ Bài toán: Hyperledger Fabric immutability + CouchDB rich query

Card 3 — Anti-clone QR địa lý:
1 mã QR bị photocopy và quét ở 2 siêu thị khác nhau 
trong 60 giây → màn hình cảnh báo đỏ real-time.
→ Bài toán: Event-driven detection + Geolocation conflict

Card 4 — Burn token on-chain:
Người mua ấn xác nhận → token bị đốt (Burn) trên chain. 
Vỏ tem sau đó không thể kích hoạt thêm lần nào nữa.
→ Bài toán: Token lifecycle + Smart Contract state machine

Card 5 — MSP identity enforcement:
Mỗi hàm chaincode chỉ được gọi bởi đúng org có quyền. 
ChiCucBVTV mới được cấp chứng nhận VietGAP. 
NPP mới xem được công nợ của mình.
→ Bài toán: Fabric MSP + per-function access control

---

### SECTION 4 — KIẾN TRÚC HỆ THỐNG
Tiêu đề: "Kiến Trúc Hiện Tại — Và Những Gì Cần Build Tiếp"

Hiển thị sơ đồ kiến trúc dạng ASCII hoặc SVG đơn giản:

  Flutter/Vue → REST API (Gin/Go) → Fabric Gateway gRPC
  → Hyperledger Fabric 3.1.1
      ├── nhatky-htx-channel: lohang_cc + nhatky_cc
      └── giaodich-channel:   giao_dich_cc
  → CouchDB (rich query + GetHistoryForKey)

Chia 2 cột dưới sơ đồ:

CỘT TRÁI — ĐÃ HOÀN THÀNH ✅:
✅ Fabric network: 4 org, 1 orderer, 4 peer, 4 CouchDB
✅ lohang_cc: TaoLotHang, ThemChungNhan, LayThongTinTraCuu
✅ nhatky_cc: GhiNhatKy, DuyetNhatKy, GetHistoryForKey
✅ giao_dich_cc: TaoGiaoDich, DocCongNoNPP, TinhHoaHongNPP
✅ Backend Golang: REST API + JWT + Fabric Gateway SDK v1.x
✅ Go client demo: htx_ngonsan.go + npp_nongsan.go
✅ QR Consumer Portal (Vue 3)
✅ 9 Docker containers chạy ổn định local

CỘT PHẢI — CẦN COMMUNITY →:
→ CI/CD: GitHub Actions + Docker Compose production
→ Cloud deployment: AWS / GCP / bare metal VPS
→ Monitoring: Prometheus + Grafana dashboard
→ Chaincode unit test (coverage target: >70%)
→ Fabric network performance benchmark
→ Flutter app hoàn thiện (offline-first GPS + ảnh)
→ Multi-org onboarding thực tế (HTX thứ 2, NPP thứ 2)
→ Fabric Explorer tích hợp

---

### SECTION 5 — TECH STACK
Tiêu đề: "Stack Kỹ Thuật"

Hiển thị tag/badge đẹp theo nhóm:

Blockchain:
[Hyperledger Fabric 3.1.1] [Chaincode Go] [CouchDB]
[fabric-contract-api-go v1.2.2] [fabric-gateway-go v1.7.0]

Backend:
[Golang 1.21] [Gin v1.10] [JWT Auth] [gRPC] [REST API]

Frontend & Mobile:
[Vue 3] [Flutter] [Riverpod] [QR Flutter] [go_router]

DevOps (cần build):
[Docker] [Docker Compose] [GitHub Actions]

Tiêu chuẩn:
[VietGAP] [GlobalGAP] [TCVN-11814] [Organic]

Thêm đoạn text dẫn chứng:
"Cùng tech stack enterprise mà Walmart dùng trong IBM Food Trust — 
truy xuất xoài từ 7 ngày xuống còn 3 giây. 
GAPChain đang áp dụng chính xác kiến trúc đó 
cho nông sản Việt Nam."

---

### SECTION 6 — LÝ DO ĐÓNG GÓP
Tiêu đề: "Đây Không Phải Side Project Vô Nghĩa"

3 card:

Card 1 — Portfolio thực chiến:
Hyperledger Fabric 3.1.1 production experience là kỹ năng 
cực hiếm tại Việt Nam. Mọi contribution đều được ghi nhận 
công khai trên GitHub với tên bạn trong commit history.

Card 2 — Bài toán có chiều sâu kỹ thuật:
Distributed ledger, consensus mechanism, MSP identity, 
chaincode lifecycle, CouchDB rich query — đây không phải CRUD app. 
Bạn chỉ học được những thứ này khi làm thật với network thật.

Card 3 — Impact đo được:
Mỗi dòng code bạn viết giúp nông dân sạch chứng minh 
nông sản của họ — và bán đúng giá. 
Tech với mục đích cụ thể, kết quả đo được.

---

### SECTION 7 — ROADMAP
Tiêu đề: "Lộ Trình Công Khai"

Timeline 4 phase (horizontal hoặc vertical):

Phase 1 — Local MVP  ✅ HOÀN THÀNH
Fabric network 4 org · 3 chaincode · Backend Golang ·
Go client demo · QR portal · 9 containers running

Phase 2 — Cloud Deployment  🔄 CẦN COMMUNITY (Q3/2026)
CI/CD GitHub Actions · Docker production ·
Cloud hosting · Monitoring Prometheus+Grafana ·
Chaincode test coverage >70%

Phase 3 — Pilot Thực Tế  ⬜ (Q4/2026)
Onboard 3 HTX thật · Flutter app hoàn chỉnh ·
Multi-org mở rộng · Fabric Explorer ·
Tích hợp VietGAP API chính thức

Phase 4 — Scale & Open API  ⬜ (2027)
Public API cho bên thứ 3 ·
Kết nối siêu thị / sàn thương mại ·
Dashboard Chi Cục BVTV cấp tỉnh

---

### SECTION 8 — CTA CUỐI
Tiêu đề: "Bắt Đầu Ngay — Không Cần Xin Phép"

Subtext:
"Fork repo. Chọn 1 issue. Gửi PR.
Mọi contribution — dù nhỏ — đều được merge và credit."

3 button:
- Primary:   "🔧 Xem GitHub Repo" (link #)
- Secondary: "📖 Đọc Tài Liệu Kỹ Thuật" (link #)
- Outline:   "💬 Tham Gia Zalo / Discord" (link #)

Micro-text:
"MIT License · Mọi contribution được credit ·
Golang + Fabric experience歓迎"

Thêm code snippet chaincode thực tế:

  // lohang_cc/lohang_cc.go — MSP identity check
  func kiemTraMSP(ctx contractapi.TransactionContextInterface,
      mspChoPhep ...string) error {
      mspID, _ := ctx.GetClientIdentity().GetMSPID()
      for _, m := range mspChoPhep {
          if mspID == m { return nil }
      }
      return fmt.Errorf("MSP %s không có quyền", mspID)
  }

---

### SECTION 9 — FOOTER
Logo GAPChain +
"Open Source Agri-Blockchain · Hyperledger Fabric 3.1.1 ·
Golang · Flutter · 2026 · MIT License"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 🎨 THIẾT KẾ & VISUAL STYLE

Mode: Dark-light hybrid
- Nền chính: slate-900 (dark sections) + slate-50 (light sections)
- Accent: green-500 / emerald-400 cho blockchain elements
- Tech elements: blue-400 + violet-400
- Code blocks: nền slate-800, font monospace, syntax highlight giả

Các yếu tố visual bắt buộc:
- Terminal/CLI blocks với prompt $ thực tế
- Code snippets Go/chaincode thực tế từ dự án
- Subtle dot-grid hoặc circuit-board pattern background
- Badge/tag pill shape cho tech stack
- Timeline visual cho Roadmap
- Kiến trúc diagram ASCII hoặc SVG đơn giản
- Progress visual "Đã hoàn thành vs Cần build"

Tránh:
- Hình ảnh nông dân, ruộng lúa, thiên nhiên (đó là trang B2C)
- Màu sắc pastel nhẹ nhàng
- Ngôn ngữ marketing hoa mỹ
- Quá nhiều emoji

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## ✍️ NGUYÊN TẮC GIỌNG VĂN

- Peer-to-peer: dev nói chuyện với dev, không phải pitch deck
- Dùng thuật ngữ kỹ thuật thoải mái — đối tượng hiểu hết:
  Chaincode, MSP, orderer, CouchDB, GetHistoryForKey,
  burn token, endorsement policy...
- Câu ngắn, súc tích, không vòng vo
- Trung thực về hiện trạng: "đang ở local, cần community
  để lên production" — không vẽ vời quá mức
- Tên hàm chaincode dùng tiếng Việt đúng như code thật:
  TaoLotHang, GhiNhatKy, ThemChungNhan, DocCongNoNPP...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 📦 OUTPUT YÊU CẦU

File: Vue 3 Single File Component (.vue) hoàn chỉnh
- <template> đầy đủ tất cả section
- Tailwind CSS inline (không tách file)
- <script setup> với scroll detection navbar
- Không dùng router-link — dùng <a href="#">
- Responsive mobile-first
- Icon: emoji hoặc inline SVG đơn giản
- Code blocks: <pre><code> với styling monospace
- Tất cả nội dung bằng tiếng Việt
  (trừ tên kỹ thuật: chaincode, MSP, orderer, v.v.)