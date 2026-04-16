<template>
  <div class="min-h-screen font-sans text-slate-300 bg-slate-900 selection:bg-emerald-500/30 overflow-x-hidden">
    
    <!-- SECTION 1 — NAVBAR -->
    <nav class="fixed w-full z-50 transition-all duration-300 border-b" :class="scrolled ? 'bg-slate-900/90 backdrop-blur-md border-slate-800/80 py-3 shadow-lg' : 'bg-transparent border-transparent py-5'">
      <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded bg-gradient-to-br from-green-500 to-emerald-700 flex items-center justify-center font-bold text-white shadow-lg shadow-green-500/20">G</div>
          <span class="font-bold text-xl text-white tracking-tight">GAPChain</span>
          <span class="ml-3 hidden sm:inline-flex items-center px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 tracking-widest uppercase">Local MVP · Open Source</span>
        </div>
        <div class="hidden md:flex items-center gap-8 text-sm font-medium text-slate-400">
          <a href="#kien-truc" class="hover:text-emerald-400 transition-colors">Kiến Trúc</a>
          <a href="#tech-stack" class="hover:text-emerald-400 transition-colors">Tech Stack</a>
          <a href="#roadmap" class="hover:text-emerald-400 transition-colors">Roadmap</a>
          <a href="#dong-gop" class="hover:text-emerald-400 transition-colors">Đóng Góp</a>
        </div>
        <div class="flex items-center gap-3">
          <router-link to="/business" class="px-5 py-2.5 text-emerald-400 text-sm font-bold rounded-lg border border-emerald-500/30 hover:bg-emerald-500/10 hover:border-emerald-500/60 transition-colors focus:outline-none focus:ring-2 focus:ring-emerald-500/20">
            Dành cho Doanh Nghiệp &rarr;
          </router-link>
          <a href="https://github.com/trustfab/gapchain" class="hidden lg:flex px-5 py-2.5 bg-slate-800 text-slate-200 text-sm font-bold rounded-lg border border-slate-700 hover:bg-slate-700 hover:border-slate-500 transition-colors items-center gap-2 shadow-sm">
            ⭐ Star on GitHub <span class="text-slate-400">&rarr;</span>
          </a>
        </div>
      </div>
    </nav>

    <!-- SECTION 2 — HERO -->
    <section class="pt-32 md:pt-40 pb-20 relative overflow-hidden">
      <!-- Decor -->
      <div class="absolute top-0 right-0 w-[40rem] h-[40rem] bg-emerald-900/20 rounded-full blur-3xl pointer-events-none"></div>
      <div class="absolute bottom-0 left-0 w-[30rem] h-[30rem] bg-blue-900/10 rounded-full blur-3xl pointer-events-none border border-red-500"></div>
      
      <div class="max-w-7xl mx-auto px-6 flex flex-col lg:flex-row gap-12 items-center relative z-10">
        <!-- Text -->
        <div class="w-full lg:w-1/2">
          <h1 class="text-5xl md:text-6xl font-extrabold text-white mb-6 leading-[1.15] tracking-tight">
            Blockchain Thực Chiến.<br/>
            Bài Toán Thật.<br/>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-green-500">Cộng Đồng Thật.</span>
          </h1>
          <p class="text-lg text-slate-400 mb-8 leading-relaxed max-w-xl">
            GAPChain MVP đang chạy local trên <span class="text-slate-200 font-mono text-sm bg-slate-800 px-1 py-0.5 rounded border border-slate-700">Hyperledger Fabric v3</span> với 4 org, 2 channel, 3 chaincode — và cần những developer dũng cảm đưa nó lên production.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 mb-8 flex-wrap">            
            <a href="#kien-truc" class="px-6 py-3.5 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg shadow-lg shadow-emerald-900/50 transition-colors text-center flex items-center justify-center">
              Khám Phá Kiến Trúc ↓
            </a>
            <router-link to="/qr/consumer/LH-HTX001-1775545225847" class="px-6 py-3.5 bg-white text-slate-900 hover:bg-slate-200 font-bold rounded-lg shadow-lg shadow-emerald-900/20 transition-colors flex items-center justify-center gap-2">
              📱 Quét Thẻ QR Thử
            </router-link>
          </div>
          <div class="flex flex-wrap gap-4 text-xs font-mono font-medium text-slate-500">
            <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-slate-600"></span> MIT License</span>
            <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-blue-500"></span> Fabric v3</span>
            <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span> Local → Production</span>
          </div>
        </div>

        <!-- Terminal -->
        <div class="w-full lg:w-1/2 relative perspective-1000">
          <div class="absolute inset-0 bg-gradient-to-b from-blue-500/10 to-violet-500/10 rounded-xl blur-2xl transform -rotate-1"></div>
          <div class="bg-[#0D1117] rounded-xl border border-slate-700 shadow-2xl relative z-10 font-mono text-sm overflow-hidden transform hover:-translate-y-1 transition-transform duration-500">
            <!-- Header -->
            <div class="h-10 bg-slate-800/80 border-b border-slate-700 flex items-center px-4 gap-2">
              <div class="w-3 h-3 rounded-full bg-rose-500/80"></div>
              <div class="w-3 h-3 rounded-full bg-amber-500/80"></div>
              <div class="w-3 h-3 rounded-full bg-emerald-500/80"></div>
              <span class="ml-2 text-slate-500 text-xs font-medium w-full text-center pr-8 pointer-events-none">gapchain-network ~ bash</span>
            </div>
            <!-- Body -->
            <div class="p-6 text-slate-300 leading-relaxed overflow-x-auto whitespace-pre">
<span class="text-emerald-400 font-bold">$</span> ./scripts/setup-network.sh
<span class="text-slate-500">Starting GAPChain network (4 Orgs, 1 Orderer)...</span>
<span class="text-green-400">✅ 9 containers running (orderer + 4 peer + 4 couchdb)</span>

<span class="text-emerald-400 font-bold">$</span> peer chaincode invoke -C nhatky-htx-channel -n lohang_cc \
  -c '{"function":"TaoLotHang","Args":["LH-HTX001-2025-001",
  "HTX001","Gạo ST25","lua","1000","kg","Đông Xuân 2025","Cần Thơ"]}'

<span class="text-slate-500">2026-04-13 14:05:21.104 ICT [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001</span>
<span class="text-green-400">✅ Chaincode invoke successful. TxID: a3f9c21b92...</span>
<span class="animate-pulse">_</span></div>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 3 — BÀI TOÁN KỸ THUẬT -->
    <section class="py-24 bg-slate-50 text-slate-800 border-t border-slate-200">
      <div class="max-w-7xl mx-auto px-6">
        <div class="mb-16">
          <h2 class="text-3xl md:text-4xl font-extrabold text-slate-900 tracking-tight mb-4">Tại Sao Bài Toán Này Xứng Đáng Với Thời Gian Của Bạn?</h2>
          <p class="text-slate-500 text-lg">Những kỹ thuật Distributed Ledger Enterprise-grade đang được áp dụng thay cho các CRUD app thông thường.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Card 1 -->
          <div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow flex flex-col">
            <div class="w-12 h-12 bg-blue-100 text-blue-600 rounded-xl flex items-center justify-center text-xl mb-6 font-mono font-bold">1</div>
            <h3 class="text-xl font-bold mb-3 text-slate-900">Token hóa sản lượng</h3>
            <p class="text-slate-600 text-sm mb-6 leading-relaxed flex-1">
              Thu hoạch đúng 1 tấn hệ thống chỉ phát hành số token tương đương 1 tấn. Không thể xin thêm mã QR để trà trộn hàng giả.
            </p>
            <div class="pt-4 border-t border-slate-100 mt-auto">
              <p class="text-xs font-mono text-emerald-600 font-semibold"><span class="text-slate-300">→</span> Chaincode lifecycle + quantity enforcement</p>
            </div>
          </div>

          <!-- Card 2 -->
          <div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow flex flex-col">
            <div class="w-12 h-12 bg-violet-100 text-violet-600 rounded-xl flex items-center justify-center text-xl mb-6 font-mono font-bold">2</div>
            <h3 class="text-xl font-bold mb-3 text-slate-900">Audit trail bất biến</h3>
            <p class="text-slate-600 text-sm mb-6 leading-relaxed flex-1">
              <code class="text-xs font-bold text-rose-500 bg-slate-100 px-1.5 py-0.5 rounded">GetHistoryForKey</code> trả về toàn bộ lịch sử thay đổi của lô hàng. Không quản trị viên nào có thể sửa ngược dữ liệu sau khi đã commit lên Fabric.
            </p>
            <div class="pt-4 border-t border-slate-100 mt-auto">
              <p class="text-xs font-mono text-emerald-600 font-semibold"><span class="text-slate-300">→</span> Fabric immutability + CouchDB rich query</p>
            </div>
          </div>

          <!-- Card 3 -->
          <div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow flex flex-col">
            <div class="w-12 h-12 bg-rose-100 text-rose-600 rounded-xl flex items-center justify-center text-xl mb-6 font-mono font-bold">3</div>
            <h3 class="text-xl font-bold mb-3 text-slate-900">Anti-clone QR địa lý</h3>
            <p class="text-slate-600 text-sm mb-6 leading-relaxed flex-1">
              1 mã QR bị photocopy và quét ở 2 siêu thị khác nhau trong 60 giây xuất hiện cảnh báo đỏ real-time.
            </p>
            <div class="pt-4 border-t border-slate-100 mt-auto">
              <p class="text-xs font-mono text-emerald-600 font-semibold"><span class="text-slate-300">→</span> Event-driven detection + Geolocation</p>
            </div>
          </div>

          <!-- Card 4 -->
          <div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow flex flex-col">
            <div class="w-12 h-12 bg-amber-100 text-amber-600 rounded-xl flex items-center justify-center text-xl mb-6 font-mono font-bold">4</div>
            <h3 class="text-xl font-bold mb-3 text-slate-900">Burn token on-chain</h3>
            <p class="text-slate-600 text-sm mb-6 leading-relaxed flex-1">
              Người mua ấn xác nhận trên App → token bị đốt (Burn) trên chain. Tem vỏ sau đó không thể kích hoạt lại lần nào nữa.
            </p>
            <div class="pt-4 border-t border-slate-100 mt-auto">
              <p class="text-xs font-mono text-emerald-600 font-semibold"><span class="text-slate-300">→</span> Token lifecycle + Smart Contract State</p>
            </div>
          </div>

          <!-- Card 5 -->
          <div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow flex flex-col">
            <div class="w-12 h-12 bg-emerald-100 text-emerald-600 rounded-xl flex items-center justify-center text-xl mb-6 font-mono font-bold">5</div>
            <h3 class="text-xl font-bold mb-3 text-slate-900">MSP identity enforcement</h3>
            <p class="text-slate-600 text-sm mb-6 leading-relaxed flex-1">
              Mỗi hàm chaincode được kiểm soát chặt theo org. Chi Cục BVTV cấp chứng nhận, NPP xem công nợ, phân tầng bảo mật tuyệt đối.
            </p>
            <div class="pt-4 border-t border-slate-100 mt-auto">
              <p class="text-xs font-mono text-emerald-600 font-semibold"><span class="text-slate-300">→</span> Fabric MSP + per-function Identity Check</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 4 — KIẾN TRÚC HỆ THỐNG -->
    <section id="kien-truc" class="py-24 relative bg-[#0B0F19] border-t border-slate-800">
      <div class="absolute top-0 right-0 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl pointer-events-none"></div>

      <div class="max-w-7xl mx-auto px-6 relative z-10">
        <h2 class="text-3xl md:text-4xl font-extrabold text-white tracking-tight mb-16">Kiến Trúc Hiện Tại —<br/>Và Những Gì Cần Build Tiếp</h2>
        
        <!-- Architecture Diagram -->
        <div class="mb-16 p-8 bg-[#0D1117] rounded-xl border border-slate-800 shadow-xl overflow-x-auto">
<pre class="font-mono text-sm text-slate-300 leading-loose">
[Flutter / Vue 3 Client App] 
            │
            ▼
[REST API (Golang / Gin)]  ────── JWT Auth Validation
            │
            ▼
[Fabric Gateway SDK v1.x]  ────── Query ────► [CouchDB] (State DB)
            │                                     • Rich JSON Query
            ▼                                     • GetHistoryForKey
[Hyperledger Fabric v3 Network]
            ├── <span class="text-emerald-400 font-bold">nhatky-htx-channel</span>
            │     ├── <span class="text-blue-400 font-bold">lohang_cc</span>  (Lô hàng & chứng nhận VietGAP)
            │     └── <span class="text-blue-400 font-bold">nhatky_cc</span>  (Nhật ký canh tác, audit trail)
            │
            └── <span class="text-violet-400 font-bold">giaodich-channel</span>
                  └── <span class="text-blue-400 font-bold">giao_dich_cc</span> (Giao dịch B2B, hoa hồng NPP)
</pre>
        </div>

        <!-- 2 Columns Checklist -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Left: Done -->
          <div class="bg-slate-800/40 border border-slate-700/50 rounded-2xl p-8 backdrop-blur-sm">
            <h3 class="text-xl font-bold text-white mb-6 flex items-center gap-3">
              <span class="w-8 h-8 rounded bg-emerald-500/20 text-emerald-400 flex items-center justify-center font-bold">✓</span> 
              Đã Hoàn Thành
            </h3>
            <ul class="space-y-4 text-slate-400 text-sm font-medium">
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span>Fabric network: 4 org, 1 orderer, 4 peer, 4 CouchDB</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span><code class="text-slate-300">lohang_cc</code>: TaoLotHang, ThemChungNhan, LayThongTinTraCuu</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span><code class="text-slate-300">nhatky_cc</code>: GhiNhatKy, DuyetNhatKy, GetHistoryForKey</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span><code class="text-slate-300">giao_dich_cc</code>: TaoGiaoDich, DocCongNoNPP, TinhHoaHongNPP</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span>Backend: Golang + REST API + JWT + Fabric Gateway SDK v1.x</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span>Go client demo: <code class="text-slate-300">htx_nongsan.go</code>, <code class="text-slate-300">npp_nongsan.go</code></span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span>QR Consumer Portal (Vue 3) hoạt động</span></li>
              <li class="flex gap-3 items-start"><span class="text-green-500 mt-0.5">✅</span> <span>9 Docker containers đang chạy ổn định (local MVP)</span></li>
            </ul>
          </div>

          <!-- Right: Todo -->
          <div class="bg-blue-900/10 border border-blue-500/20 rounded-2xl p-8 backdrop-blur-sm relative overflow-hidden shadow-[0_0_30px_rgba(59,130,246,0.05)]">
            <div class="absolute top-0 right-0 px-4 py-1 bg-blue-500 text-white text-xs font-bold rounded-bl-lg tracking-widest">CẦN COMMUNITY</div>
            <h3 class="text-xl font-bold text-white mb-6 flex items-center gap-3">
              <span class="w-8 h-8 rounded bg-blue-500/20 text-blue-400 flex items-center justify-center font-bold">→</span> 
              Cơ Hội Đóng Góp
            </h3>
            <ul class="space-y-4 text-slate-300 text-sm font-medium">
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">CI/CD:</strong> Xây dựng GitHub Actions + Docker Compose production</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Cloud Deployment:</strong> Phân bổ kiến trúc lên AWS / GCP / VPS bare metal</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Monitoring:</strong> Cài đặt Prometheus + Grafana dashboard</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Chaincode Testing:</strong> Unit/Integration test bằng ginkgo/go test (mục tiêu >70%)</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Performance:</strong> Benchmark network Fabric & tối ưu hóa</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Flutter App:</strong> Hoàn thiện module offline-first, GPS tracking & SHA256 ảnh</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Multi-Org Scale:</strong> Setup config onboard HTX/Nhà phân phối thật lên network</span></li>
              <li class="flex gap-3 items-start"><span class="text-blue-400 font-bold mt-0.5">→</span> <span><strong class="text-white">Fabric Explorer:</strong> Tích hợp UI theo dõi giao dịch chuỗi</span></li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 5 — TECH STACK -->
    <section id="tech-stack" class="py-24 bg-slate-50 border-y border-slate-200">
      <div class="max-w-7xl mx-auto px-6">
        <h2 class="text-3xl md:text-3xl font-extrabold text-slate-900 tracking-tight mb-8">Stack Kỹ Thuật</h2>
        
        <blockquote class="border-l-4 border-emerald-500 bg-emerald-50/50 p-6 rounded-r-lg max-w-4xl mb-12 shadow-sm">
          <p class="text-lg text-slate-700 font-medium italic leading-relaxed">
            "Cùng tech stack enterprise mà Walmart dùng trong IBM Food Trust — truy xuất xoài từ 7 ngày xuống còn 3 giây. GAPChain đang áp dụng chính xác kiến trúc đó cho bảo vệ hệ sinh thái nông sản Việt Nam."
          </p>
        </blockquote>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-10">
          <div>
            <h4 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Blockchain Core</h4>
            <div class="flex flex-wrap gap-2">
              <span class="px-3 py-1.5 bg-white text-slate-800 rounded-lg text-sm font-bold border border-slate-200 shadow-sm">Hyperledger Fabric v3</span>
              <span class="px-3 py-1.5 bg-blue-50 text-blue-800 rounded-lg text-sm font-bold border border-blue-200 shadow-sm">Chaincode Go</span>
              <span class="px-3 py-1.5 bg-red-50 text-red-800 rounded-lg text-sm font-bold border border-red-200 shadow-sm">CouchDB</span>
              <span class="px-3 py-1.5 bg-white font-mono text-slate-600 rounded-lg text-xs font-medium border border-slate-200">fabric-contract-api-go v1.2.2</span>
              <span class="px-3 py-1.5 bg-white font-mono text-slate-600 rounded-lg text-xs font-medium border border-slate-200">fabric-gateway-go v1.7.0</span>
            </div>
          </div>
          <div>
            <h4 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Backend Ecosystem</h4>
            <div class="flex flex-wrap gap-2">
              <span class="px-3 py-1.5 bg-cyan-50 text-cyan-800 rounded-lg text-sm font-bold border border-cyan-200 shadow-sm">Golang 1.21</span>
              <span class="px-3 py-1.5 bg-white text-slate-800 rounded-lg text-sm font-bold border border-slate-200 shadow-sm">Gin v1.10</span>
              <span class="px-3 py-1.5 bg-indigo-50 text-indigo-800 rounded-lg text-sm font-bold border border-indigo-200 shadow-sm">JWT Auth</span>
              <span class="px-3 py-1.5 bg-purple-50 text-purple-800 rounded-lg text-sm font-bold border border-purple-200 shadow-sm">gRPC</span>
              <span class="px-3 py-1.5 bg-white text-slate-800 rounded-lg text-sm font-bold border border-slate-200 shadow-sm">REST API</span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-10 mt-10 p-6 bg-slate-100 rounded-2xl border border-slate-200/60">
          <div>
            <h4 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Frontend & Mobile</h4>
            <div class="flex flex-wrap gap-2">
              <span class="px-3 py-1.5 bg-emerald-100 text-emerald-900 rounded text-sm font-bold shadow-sm">Vue 3</span>
              <span class="px-3 py-1.5 bg-sky-100 text-sky-900 rounded text-sm font-bold shadow-sm">Flutter</span>
              <span class="px-3 py-1.5 bg-white text-slate-700 rounded text-sm font-bold border border-slate-200">Riverpod</span>
              <span class="px-3 py-1.5 bg-white text-slate-700 rounded text-sm font-bold border border-slate-200">go_router</span>
            </div>
          </div>
          <div>
            <h4 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4 flex items-center gap-2">
              DevOps 
              <span class="bg-amber-100 text-amber-800 px-2 py-0.5 rounded text-[10px] uppercase font-bold tracking-widest">Cần Build</span>
            </h4>
            <div class="flex flex-wrap gap-2">
              <span class="px-3 py-1.5 bg-blue-50 text-blue-900 rounded text-sm font-bold border border-blue-200 border-dashed">Docker Compose</span>
              <span class="px-3 py-1.5 bg-slate-200 text-slate-900 rounded text-sm font-bold border border-slate-300 border-dashed">GitHub Actions</span>
              <span class="px-3 py-1.5 bg-orange-50 text-orange-900 rounded text-sm font-bold border border-orange-200 border-dashed">AWS / GCP</span>
            </div>
          </div>
        </div>
        
        <div class="mt-10 flex gap-3 text-sm font-bold text-slate-600 items-center justify-center bg-white py-3 rounded-lg border border-slate-200 shadow-sm inline-flex px-8 mx-auto w-full md:w-auto">
          <span class="text-xs uppercase tracking-widest mr-2 text-slate-400">Tiêu chuẩn nghiệp vụ:</span>
          VietGAP <span class="text-emerald-400">&bull;</span> GlobalGAP <span class="text-emerald-400">&bull;</span> TCVN-11814 <span class="text-emerald-400">&bull;</span> Organic
        </div>
      </div>
    </section>

    <!-- SECTION 6 — LÝ DO ĐÓNG GÓP -->
    <section id="dong-gop" class="py-32 relative overflow-hidden bg-slate-900 border-b border-slate-800">
      <div class="max-w-7xl mx-auto px-6 relative z-10">
        <h2 class="text-3xl md:text-5xl font-extrabold text-white tracking-tight mb-20 text-center">Đây Không Phải Side Project Vô Nghĩa</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Card 1 -->
          <div class="p-8 bg-slate-800/40 border border-slate-700/60 rounded-3xl hover:bg-slate-800/60 transition-colors">
            <div class="w-14 h-14 bg-slate-800 rounded-xl flex items-center justify-center text-2xl mb-6 shadow-inner border border-slate-700">🏆</div>
            <h3 class="text-xl font-bold text-white mb-4">Portfolio Thực Chiến</h3>
            <p class="text-slate-400 text-sm leading-relaxed">
              Hyperledger Fabric v3 production experience là kỹ năng cực hiếm tại môi trường công nghệ Việt Nam. Mọi contribution đều được ghi nhận công khai trên GitHub với tên tuổi của bạn trong lịch sử commit.
            </p>
          </div>
          <!-- Card 2 -->
          <div class="p-8 bg-slate-800/40 border border-slate-700/60 rounded-3xl hover:bg-slate-800/60 transition-colors transform md:-translate-y-6">
            <div class="w-14 h-14 bg-slate-800 rounded-xl flex items-center justify-center text-2xl mb-6 shadow-inner border border-slate-700">⚔️</div>
            <h3 class="text-xl font-bold text-white mb-4">Bài Toán Có Chiều Sâu KT</h3>
            <p class="text-slate-400 text-sm leading-relaxed">
              Distributed ledger, consensus mechanism, MSP identity, lifecycle chaincode, CouchDB query — đây rõ ràng không phải CRUD app. Bạn chỉ học được những thứ này thông qua cọ xát với mạng lưới phân tán.
            </p>
          </div>
          <!-- Card 3 -->
          <div class="p-8 bg-slate-800/40 border border-slate-700/60 rounded-3xl hover:bg-slate-800/60 transition-colors">
            <div class="w-14 h-14 bg-slate-800 rounded-xl flex items-center justify-center text-2xl mb-6 shadow-inner border border-slate-700">🌍</div>
            <h3 class="text-xl font-bold text-white mb-4">Impact Đo Lường Được</h3>
            <p class="text-slate-400 text-sm leading-relaxed">
              Mỗi dòng code bạn viết giúp những người nông dân sạch dễ dàng chứng minh nguồn gốc nông sản của họ — qua đó bán được giá xứng đáng. Cống hiến kỹ thuật vào những kết quả đong đếm được.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 7 — ROADMAP -->
    <section id="roadmap" class="py-32 bg-slate-50">
      <div class="max-w-7xl mx-auto px-6">
        <h2 class="text-3xl md:text-5xl font-extrabold text-slate-900 tracking-tight mb-20 text-center">Lộ Trình Công Khai</h2>
        
        <div class="space-y-6 relative before:absolute before:inset-y-0 before:left-[19px] md:before:left-1/2 md:before:-ml-[1px] before:w-[2px] before:bg-slate-200">
          
          <!-- Phase 1 -->
          <div class="relative flex flex-col md:flex-row gap-6 md:gap-12 md:even:flex-row-reverse w-full">
            <div class="absolute w-10 h-10 bg-emerald-500 rounded-full left-0 md:left-1/2 md:-ml-5 top-0 flex items-center justify-center text-white font-bold border-4 border-slate-50 shadow-sm z-10">✓</div>
            <div class="w-full md:w-1/2 pl-16 md:pl-0 md:pr-16 text-left md:text-right">
              <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
                <span class="inline-block px-2 py-1 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded mb-3">HOÀN THÀNH</span>
                <h3 class="text-xl font-bold text-slate-900 mb-3">Phase 1 — Local MVP</h3>
                <p class="text-sm text-slate-600 font-medium leading-relaxed">
                  Thiết lập Fabric network 4 org <span class="text-slate-300">|</span> Viết 3 chaincode lõi bằng Go <span class="text-slate-300">|</span> API Backend Golang <span class="text-slate-300">|</span> Client demo CLI <span class="text-slate-300">|</span> Portal web quét QR (Vue) <span class="text-slate-300">|</span> Setup thành công chạy ổn định 9 Docker containers ở local.
                </p>
              </div>
            </div>
            <div class="hidden md:block w-1/2"></div>
          </div>

          <!-- Phase 2 -->
          <div class="relative flex flex-col md:flex-row gap-6 md:gap-12 md:even:flex-row-reverse w-full -mt-2">
            <div class="absolute w-10 h-10 bg-blue-500 rounded-full left-0 md:left-1/2 md:-ml-5 top-0 flex items-center justify-center text-white font-bold border-4 border-slate-50 shadow-sm z-10 animate-pulse">+</div>
            <div class="w-full md:w-1/2 pl-16 md:pl-0 md:pl-16 text-left">
              <div class="bg-white p-6 rounded-2xl border-2 border-blue-400 shadow-lg relative overflow-hidden group hover:border-blue-500 transition-colors">
                <div class="absolute top-0 right-0 w-32 h-32 bg-blue-500/10 rounded-full blur-2xl -mr-10 -mt-10 pointer-events-none group-hover:bg-blue-500/20 transition-colors"></div>
                <span class="inline-block px-2 py-1 bg-blue-500 text-white text-[10px] tracking-widest font-bold rounded mb-3">CẦN COMMUNITY (Q3/2026)</span>
                <h3 class="text-xl font-bold text-slate-900 mb-3">Phase 2 — Cloud Deployment</h3>
                <p class="text-sm text-slate-600 font-medium leading-relaxed">
                  Thiết lập CI/CD GitHub Actions <span class="text-slate-300">|</span> Đóng gói Docker Compose production <span class="text-slate-300">|</span> Hosting linh hoạt <span class="text-slate-300">|</span> Giám sát bằng Prometheus + Grafana <span class="text-slate-300">|</span> Đẩy độ phủ test Chaincode >70%.
                </p>
              </div>
            </div>
            <div class="hidden md:block w-1/2"></div>
          </div>

          <!-- Phase 3 -->
          <div class="relative flex flex-col md:flex-row gap-6 md:gap-12 md:even:flex-row-reverse w-full -mt-2">
            <div class="absolute w-10 h-10 bg-slate-300 rounded-full left-0 md:left-1/2 md:-ml-5 top-0 flex items-center justify-center text-slate-500 font-bold border-4 border-slate-50 shadow-sm z-10">3</div>
            <div class="w-full md:w-1/2 pl-16 md:pl-0 md:pr-16 text-left md:text-right">
              <div class="bg-slate-50 border border-slate-200 border-dashed p-6 rounded-2xl opacity-80 hover:opacity-100 transition-opacity">
                <span class="inline-block px-2 py-1 bg-slate-200 text-slate-600 text-[10px] font-bold rounded mb-3">Q4 / 2026</span>
                <h3 class="text-xl font-bold text-slate-700 mb-3">Phase 3 — Pilot Thực Tế</h3>
                <p class="text-sm text-slate-500 font-medium leading-relaxed">
                  Onboard 3 dự án HTX vào hệ thống <span class="text-slate-300">|</span> Hoàn tất phiên bản Flutter App cho bà con <span class="text-slate-300">|</span> Network Topology mở rộng thêm Org <span class="text-slate-300">|</span> Cài đặt Fabric Explorer view <span class="text-slate-300">|</span> Sync API VietGAP chính phủ.
                </p>
              </div>
            </div>
            <div class="hidden md:block w-1/2"></div>
          </div>

          <!-- Phase 4 -->
          <div class="relative flex flex-col md:flex-row gap-6 md:gap-12 md:even:flex-row-reverse w-full -mt-2">
            <div class="absolute w-10 h-10 bg-slate-200 rounded-full left-0 md:left-1/2 md:-ml-5 top-0 flex items-center justify-center text-slate-400 font-bold border-4 border-slate-50 shadow-sm z-10">4</div>
            <div class="w-full md:w-1/2 pl-16 md:pl-0 md:pl-16 text-left">
              <div class="bg-slate-50 border border-slate-200 border-dashed p-6 rounded-2xl opacity-60 hover:opacity-100 transition-opacity">
                <span class="inline-block px-2 py-1 bg-slate-200 text-slate-500 text-[10px] font-bold rounded mb-3">NĂM 2027</span>
                <h3 class="text-xl font-bold text-slate-600 mb-3">Phase 4 — Scale & Open API</h3>
                <p class="text-sm text-slate-500 font-medium leading-relaxed">
                  Thiết kế module Open/Public API cho 3rd-party Apps <span class="text-slate-300">|</span> Automation Match-making siêu thị bán lẻ <span class="text-slate-300">|</span> Cấp quyền module Dashboard tổng cho cấp Tỉnh.
                </p>
              </div>
            </div>
            <div class="hidden md:block w-1/2"></div>
          </div>
          
        </div>
      </div>
    </section>

    <!-- SECTION 8 — CTA CUỐI -->
    <section class="py-32 relative bg-[#090C15] border-t border-slate-800">
      <div class="max-w-6xl mx-auto px-6 relative z-10 flex flex-col md:flex-row gap-12 items-center">
        <!-- Text -->
        <div class="w-full md:w-1/2">
          <h2 class="text-4xl md:text-5xl font-extrabold text-white mb-8 leading-tight">Bắt Đầu Ngay —<br/> Không Cần Xin Phép.</h2>
          <p class="text-lg text-slate-400 mb-10 font-mono leading-relaxed">
            <span class="text-emerald-400 font-bold">1.</span> Fork repo.<br/>
            <span class="text-emerald-400 font-bold">2.</span> Chọn 1 block issue.<br/>
            <span class="text-emerald-400 font-bold">3.</span> Submit pull request trực tiếp.<br/><br/>
            <span class="text-slate-300 font-sans">Mọi contribution — cấu hình docker nhỏ nhất cho đến module Flutter lớn — đều được review cẩn trọng và cấp credit vĩnh viễn trên README.</span>
          </p>
          <div class="flex flex-col sm:flex-row gap-4 mb-6">
            <a href="https://github.com/trustfab/gapchain" class="px-6 py-4 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg text-center transition-colors shadow-lg shadow-emerald-600/30">
              🔧 Khám Phá Repo GitHub
            </a>
            <a href="#" class="px-6 py-4 bg-slate-800 hover:bg-slate-700 text-white font-bold rounded-lg border border-slate-700 text-center transition-colors">
              📖 Đọc Document Cốt Lõi
            </a>
          </div>
          <p class="text-xs font-mono font-medium text-slate-500 flex flex-wrap gap-x-4 gap-y-2 mt-8">
            <span class="flex items-center gap-1.5"><span class="text-emerald-400">✓</span> MIT License v3</span> 
            <span class="flex items-center gap-1.5"><span class="text-emerald-400">✓</span> Open Community Peer-to-Peer</span>
            <span class="flex items-center gap-1.5"><span class="text-emerald-400">✓</span> Golang + Fabric Welcome</span>
          </p>
        </div>

        <!-- Code block snippet -->
        <div class="w-full md:w-1/2 relative">
          <div class="absolute inset-0 bg-blue-500/10 rounded-2xl blur-xl transform rotate-3 z-0"></div>
          <div class="bg-[#121622] rounded-2xl border border-slate-700/80 shadow-2xl p-6 font-mono text-sm text-slate-300 relative z-10 overflow-x-auto">
            <div class="absolute top-0 right-0 px-3 py-1 bg-slate-800/80 backdrop-blur rounded-bl-lg text-slate-400 text-xs font-bold border-b border-l border-slate-700">lohang_cc/lohang_cc.go</div>
<pre class="pt-4 leading-loose">
<span class="text-slate-500">// Thực thi giới hạn identity MSP ở tầng chaincode</span>
<span class="text-rose-400 font-bold">func</span> <span class="text-blue-300 font-bold">kiemTraMSP</span>(ctx contractapi.TransactionContextInterface,
    mspChoPhep ...<span class="text-emerald-400">string</span>) <span class="text-emerald-400">error</span> {

    mspID, _ := ctx.GetClientIdentity().GetMSPID()
    
    <span class="text-rose-400 font-bold">for</span> _, m := <span class="text-rose-400 font-bold">range</span> mspChoPhep {
        <span class="text-rose-400 font-bold">if</span> mspID == m { 
            <span class="text-rose-400 font-bold">return</span> <span class="text-amber-300">nil</span> 
        }
    }
    <span class="text-rose-400 font-bold">return</span> fmt.Errorf(<span class="text-amber-200">"Từ chối: MSP %s không có quyền!"</span>, mspID)
}
</pre>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 9 — FOOTER -->
    <footer class="bg-[#05080f] py-16 border-t border-slate-800 text-center">
      <div class="max-w-7xl mx-auto px-6">
        <div class="flex items-center justify-center gap-3 mb-6 opacity-60">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-slate-600 to-slate-800 flex items-center justify-center font-bold text-white text-sm">G</div>
          <span class="font-bold text-xl text-slate-400 tracking-tight">GAPChain Project</span>
        </div>
        <p class="text-sm font-medium text-slate-500 mb-4 inline-flex items-center flex-wrap justify-center gap-2">
          <span class="px-2 py-1 bg-slate-800 rounded">Open Source Agri-Blockchain</span>
          <span class="px-2 py-1 bg-slate-800 rounded">Hyperledger Fabric v3</span>
          <span class="px-2 py-1 bg-slate-800 rounded">Golang</span>
          <span class="px-2 py-1 bg-slate-800 rounded">Flutter</span>
          <span class="px-2 py-1 bg-slate-800 rounded">Vue</span>
        </p>
        <p class="text-xs text-slate-600 font-mono mt-4">&copy; 2026. Code under MIT License. All right preserved.</p>
        <div class="mt-8 flex justify-center space-x-6">
          <a href="https://github.com/trustfab/gapchain" class="text-slate-500 hover:text-white transition-colors">GitHub</a>
          <a href="#" class="text-slate-500 hover:text-white transition-colors">Discord</a>
          <a href="https://github.com/trustfab/gapchain/blob/main/README.md" class="text-slate-500 hover:text-white transition-colors">Documents</a>
          <a href="/gapchain/login" class="text-slate-500 hover:text-white transition-colors">Demo ORG</a>          
          <a href="/gapchain/qr/consumer/LH-HTX001-1775545225847" class="text-slate-500 hover:text-white transition-colors">Demo QR</a>          
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const scrolled = ref(false)

const handleScroll = () => {
  scrolled.value = window.scrollY > 20
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
html {
  scroll-behavior: smooth;
}
</style>
