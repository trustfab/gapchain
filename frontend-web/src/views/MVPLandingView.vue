<template>
  <div class="min-h-screen font-sans text-slate-700 bg-slate-50 selection:bg-emerald-500/30 overflow-x-hidden">
    
    <!-- SECTION 1 — NAVBAR -->
    <nav class="fixed w-full z-50 transition-all duration-300 border-b" :class="scrolled ? 'bg-white/95 backdrop-blur-md border-slate-200 py-3 shadow-md' : 'bg-transparent border-transparent py-5'">
      <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-emerald-500 to-green-700 flex items-center justify-center font-bold text-white shadow-lg shadow-emerald-500/30">G</div>
          <span class="font-extrabold text-xl text-slate-800 tracking-tight">GAPChain</span>
        </div>
        <div class="hidden md:flex items-center gap-8 text-sm font-bold text-slate-600">
          <a href="#giai-phap" class="hover:text-emerald-600 transition-colors">Giải Pháp</a>
          <a href="#loi-ich" class="hover:text-emerald-600 transition-colors">Lợi Ích</a>
          <a href="#tinh-nang" class="hover:text-emerald-600 transition-colors">Tính Năng</a>
        </div>
        <div class="flex items-center gap-3">
          <router-link to="/" class="hidden sm:flex px-4 py-2.5 text-slate-500 hover:text-emerald-700 text-sm font-bold transition-all items-center gap-2 border border-transparent hover:border-emerald-200 hover:bg-emerald-50 rounded-lg">
             <span class="font-mono text-xs bg-slate-200 text-slate-600 px-1.5 py-0.5 rounded">&lt;/&gt;</span>
             Dành cho Developer
          </router-link>
          <a href="#" class="px-5 py-2.5 bg-emerald-600 text-white text-sm font-bold rounded-lg hover:bg-emerald-500 hover:shadow-lg hover:shadow-emerald-500/30 transition-all flex items-center gap-2">
            Trải Nghiệm Pilot
          </a>
        </div>
      </div>
    </nav>

    <!-- SECTION 2 — HERO -->
    <section class="pt-32 md:pt-40 pb-24 relative overflow-hidden bg-white border-b border-slate-200">
      <!-- Decor -->
      <div class="absolute top-0 right-0 w-[40rem] h-[40rem] bg-emerald-500/10 rounded-full blur-3xl pointer-events-none -translate-y-1/2 translate-x-1/3"></div>
      <div class="absolute bottom-0 left-0 w-[30rem] h-[30rem] bg-blue-500/5 rounded-full blur-3xl pointer-events-none translate-y-1/2 -translate-x-1/2"></div>
      
      <div class="max-w-7xl mx-auto px-6 flex flex-col lg:flex-row gap-12 items-center relative z-10">
        <!-- Text -->
        <div class="w-full lg:w-1/2">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-bold uppercase tracking-wider mb-6">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span> Nền tảng TrustFab (Fabric v3.1)
          </div>
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-extrabold text-slate-900 mb-6 leading-[1.15] tracking-tight">
            Nâng tầm giá trị <br/>
            Nông sản sạch với <br/>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-emerald-500 to-green-600">Nhật Ký Số Bất Biến.</span>
          </h1>
          <p class="text-lg text-slate-600 mb-8 leading-relaxed max-w-xl font-medium">
            Giúp <strong class="text-slate-800">Hợp Tác Xã</strong> số hóa dễ dàng quy trình VietGAP và giúp <strong class="text-slate-800">Nhà Phân Phối</strong> kiểm soát sự minh bạch của chuỗi cung ứng chỉ bằng một chạm. Bỏ qua sổ sách giấy cồng kềnh.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 mb-10">
            <a href="#" class="px-8 py-4 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-xl shadow-lg shadow-emerald-600/30 transition-all text-center text-lg flex items-center justify-center gap-2">
              Đăng Ký Thí Điểm Miễn Phí <span class="text-xl">&rarr;</span>
            </a>
            <a href="/gapchain/qr/consumer/LH-HTX001-1775545225847" target="_blank" class="px-8 py-4 bg-white hover:bg-slate-50 text-slate-700 font-bold rounded-xl border-2 border-slate-200 transition-colors flex items-center justify-center gap-2 text-lg">
              📱 Quét Thẻ QR Thử
            </a>
          </div>
          <div class="flex items-center gap-6 mt-4 opacity-70 grayscale">
            <p class="text-xs font-bold uppercase tracking-widest text-slate-400">Đạt Tiêu Chuẩn</p>
            <span class="text-sm font-bold text-slate-700">VietGAP</span>
            <span class="text-sm font-bold text-slate-700">GlobalGAP</span>
            <span class="text-sm font-bold text-slate-700">EUDR Sẵn Sàng</span>
          </div>
        </div>

        <!-- Image / Visual -->
        <div class="w-full lg:w-1/2 relative">
          <div class="bg-white rounded-3xl border border-slate-200 shadow-2xl p-4 md:p-6 pb-0 relative z-10 mx-auto overflow-hidden">
            <!-- Mock up showing an app UI with green agriculture vibes -->
            <div class="w-full h-80 bg-slate-100 rounded-t-xl overflow-hidden relative border border-slate-200">
               <!-- Abstract network dots to farm -->
               <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCI+PHBhdGggZD0iTTAgMGg0MHY0MEgweiIgZmlsbD0ibm9uZSIvPjxjaXJjbGUgY3g9IjIwIiBjeT0iMjAiIHI9IjEiIGZpbGw9IiNjYmQ1ZTEiLz48L3N2Zz4=')] opacity-50"></div>
               <!-- Content mock -->
               <div class="absolute top-10 inset-x-8 bg-white shadow-xl rounded-t-2xl p-6 border border-slate-200 h-full">
                 <div class="flex justify-between items-center mb-6">
                   <div>
                     <p class="text-xs font-bold text-slate-400 mb-1">LÔ HÀNG</p>
                     <p class="text-lg font-extrabold text-slate-800">LH-HTX001-ST25</p>
                   </div>
                   <div class="px-3 py-1 bg-emerald-100 text-emerald-700 rounded-full text-xs font-bold">An Toàn tuyệt đối</div>
                 </div>
                 
                 <div class="space-y-4">
                   <!-- timeline item 1 -->
                   <div class="flex items-start gap-4">
                     <span class="w-8 h-8 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center font-bold text-xs ring-4 ring-white">&#10003;</span>
                     <div class="pt-1">
                       <p class="text-sm font-bold text-slate-700">Bón Phân Hữu Cơ</p>
                       <p class="text-xs text-slate-500">Đã xác nhận Blockchain bởi HTX Bình Minh</p>
                     </div>
                   </div>
                   <!-- timeline item 2 -->
                   <div class="flex items-start gap-4">
                     <span class="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold text-xs ring-4 ring-white">&#10003;</span>
                     <div class="pt-1">
                       <p class="text-sm font-bold text-slate-700">Kiểm Định VietGAP</p>
                       <p class="text-xs text-slate-500">Ký duyệt bởi Cơ quan Quản lý</p>
                     </div>
                   </div>
                 </div>
               </div>
            </div>
            
            <!-- Abstract badge overlaid -->
            <div class="absolute top-10 -right-6 lg:-right-8 bg-[#0F172A] text-white p-4 rounded-xl shadow-2xl border border-slate-700 flex items-center gap-4 animate-bounce" style="animation-duration: 3s">
              <div class="text-3xl">🛡️</div>
              <div>
                <p class="text-xs text-emerald-400 font-bold tracking-widest uppercase">Truy nguồn gốc</p>
                <p class="text-sm font-bold">Xác thực bằng 1 chạm</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 3 — NỖI ĐAU THỊ TRƯỜNG -->
    <section id="giai-phap" class="py-24 bg-slate-50">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center max-w-2xl mx-auto mb-16">
          <h2 class="text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">Tháo Gỡ Nút Thắt Niềm Tin Nông Sản</h2>
          <p class="text-lg text-slate-600 font-medium">Bất cập lớn nhất của xuất khẩu và giá trị sản phẩm nằm ở dữ liệu truyền miệng và hồ sơ trên giấy.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Item 1 -->
          <div class="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm text-center transform transition duration-300 hover:-translate-y-1">
            <div class="w-20 h-20 mx-auto bg-rose-50 rounded-full flex items-center justify-center mb-6 text-rose-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>
            </div>
            <h3 class="text-xl font-bold text-slate-800 mb-3">Tổn thất niềm tin</h3>
            <p class="text-slate-600 text-sm leading-relaxed">Người tiêu dùng đối mặt với nguy cơ hàng giả, mạo danh tiêu chuẩn. Không có công cụ sắc bén nhận biết nguồn gốc thực.</p>
          </div>
          <!-- Item 2 -->
          <div class="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm text-center transform transition duration-300 hover:-translate-y-1">
            <div class="w-20 h-20 mx-auto bg-amber-50 rounded-full flex items-center justify-center mb-6 text-amber-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m8 17 4 4 4-4"/></svg>
            </div>
            <h3 class="text-xl font-bold text-slate-800 mb-3">Nỗ Lực Không Phản Ánh Giá Tiền</h3>
            <p class="text-slate-600 text-sm leading-relaxed">Hợp tác xã nông nghiệp sản xuất rất sạch, đáp ứng VietGAP nhưng mệt mỏi với thủ tục giấy tờ cồng kềnh, bị thương lái ép giá bèo bọt do khó chứng minh "hàng ngon, hàng cực đỉnh".</p>
          </div>
          <!-- Item 3 -->
          <div class="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm text-center transform transition duration-300 hover:-translate-y-1">
             <div class="w-20 h-20 mx-auto bg-blue-50 rounded-full flex items-center justify-center mb-6 text-blue-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/><path d="M8 14h.01"/><path d="M12 14h.01"/><path d="M16 14h.01"/><path d="M8 18h.01"/><path d="M12 18h.01"/><path d="M16 18h.01"/></svg>
            </div>
            <h3 class="text-xl font-bold text-slate-800 mb-3">Rào cản Xuất khẩu (EUDR)</h3>
            <p class="text-slate-600 text-sm leading-relaxed">Luật từ Châu Âu (EUDR) và quy định truy nguồn FSMA 204 từ Hoa Kỳ đang siết thị trường, đòi hỏi hồ sơ số hóa từ điểm canh tác. Doanh nghiệp cần phải thích ứng hoặc thua lỗ.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- SECTION 4 — LỢI ÍCH THEO ĐỐI TƯỢNG -->
    <section id="loi-ich" class="py-24 bg-[#0a192f] text-slate-300 border-t border-slate-800">
      <div class="max-w-7xl mx-auto px-6">
        <h2 class="text-3xl md:text-4xl font-extrabold text-white mb-16 text-center tracking-tight">Giá Trị Liên Kết Đa Chiều</h2>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          
          <!-- Hợp tác xã -->
          <div class="p-10 rounded-3xl bg-slate-800/40 border border-slate-700 hover:border-emerald-500/50 transition-colors">
            <h3 class="text-2xl font-bold text-emerald-400 mb-6 border-b border-slate-700/50 pb-4">Nông Dân & Hợp Tác Xã</h3>
            <p class="font-bold text-white text-lg mb-4">"Đơn giản hóa sổ sách — Tối đa hóa giá trị nông sản"</p>
            <ul class="space-y-4">
              <li class="flex items-start gap-4">
                <span class="text-emerald-500 mt-1"><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg></span>
                <div>
                  <p class="font-bold text-slate-200">Không còn đống giấy tờ</p>
                  <p class="text-sm">Ghi chép Nhật ký VietGAP trên ứng dụng di động cực kỳ nhanh chóng. Tích hợp Offline, chụp hình ảnh thực địa và dán vị trí GPS tự động.</p>
                </div>
              </li>
              <li class="flex items-start gap-4">
                <span class="text-emerald-500 mt-1"><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg></span>
                <div>
                  <p class="font-bold text-slate-200">Đàm phán giá dễ dàng</p>
                  <p class="text-sm">Khi mọi lịch sử minh bạch rõ ràng không thể chối cãi, người nông dân có cớ để yêu cầu siêu thị nâng giá mua tương xứng hoặc bán thẳng cho khách hàng.</p>
                </div>
              </li>
            </ul>
          </div>

          <!-- Siêu thị / Nhà phân phối -->
          <div class="p-10 rounded-3xl bg-slate-800/40 border border-slate-700 hover:border-blue-500/50 transition-colors">
            <h3 class="text-2xl font-bold text-blue-400 mb-6 border-b border-slate-700/50 pb-4">Siêu Thị & Kênh Phân Phối</h3>
            <p class="font-bold text-white text-lg mb-4">"Kiểm soát chất lượng không phụ thuộc bên trung gian"</p>
            <ul class="space-y-4">
              <li class="flex items-start gap-4">
                <span class="text-blue-500 mt-1"><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg></span>
                <div>
                  <p class="font-bold text-slate-200">Hồ Sơ Luôn Sẵn Sàng</p>
                  <p class="text-sm">Hệ thống Data luôn luôn đáp ứng Audit theo mọi chuẩn VietGAP, EUDR mới nhất. Trình báo cho lực lượng quản lý lập tức bằng QRCode.</p>
                </div>
              </li>
              <li class="flex items-start gap-4">
                <span class="text-blue-500 mt-1"><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg></span>
                <div>
                  <p class="font-bold text-slate-200">Giảm Rủi Ro Thu Hồi Cấp Tốc</p>
                  <p class="text-sm">Nếu phát ngiện sự cố độc dược trong lô nông sản, thuật toán Blockchain khóa chặt lô đó và phát tín hiệu đỏ chỉ với 30 giây truy vết (thay vì 10 ngày như truyền thống).</p>
                </div>
              </li>
            </ul>
          </div>

        </div>

      </div>
    </section>

    <!-- SECTION 5 — CÁC TÍNH NĂNG CỐT LÕI (MVP) -->
    <section id="tinh-nang" class="py-24 bg-white border-t border-slate-200 overflow-hidden">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">Tính Năng Định Hình Niềm Tin</h2>
          <p class="text-lg text-slate-600 font-medium">Bỏ đi các từ khóa Blockchain phức tạp, đây là những gì GAPChain làm cho bạn.</p>
        </div>

        <!-- Feature Block 1 -->
        <div class="flex flex-col md:flex-row items-center gap-12 mb-20">
          <div class="w-full md:w-1/2">
            <div class="bg-gradient-to-br from-emerald-50 to-green-100 p-8 rounded-3xl border border-emerald-200 shadow-sm relative">
              <div class="bg-white p-4 rounded-xl shadow-xl border border-slate-100 absolute top-0 -right-6 transform md:translate-x-1/2 rotate-12 rotate-[-5deg] z-20">
                 <p class="text-emerald-600 font-bold text-sm mb-1">✓ Lên Lịch Gieo Hạt</p>
                 <p class="text-emerald-600 font-bold text-sm">✓ Bón Phân (Đảm bảo 4 đúng)</p>
                 <p class="text-slate-400 text-xs italic mt-2">Dấu băm số đã khóa lúc 08:31, không ai lùi giờ sửa chữa.</p>
              </div>
              <img src="/farmer_app_mockup.png" alt="Nhật Ký Số - Người Nông Dân" class="w-full rounded-2xl shadow-lg border border-white/50 relative z-10" />
            </div>
          </div>
          <div class="w-full md:w-1/2 md:pl-8">
            <div class="w-12 h-12 bg-emerald-100 text-emerald-600 rounded-xl flex items-center justify-center font-bold mb-4">1</div>
            <h3 class="text-3xl font-extrabold text-slate-900 mb-4">Nhật Ký Số Bất Biến</h3>
            <p class="text-slate-600 text-lg leading-relaxed mb-6">
              Mỗi hành động bón phân, tưới nước hay thu hoạch của người nông dân đều được số hóa, mã hóa thành một mốc thời gian không thể tẩy xóa hay chỉnh sửa lại. Đảm bảo 100% "Lời hứa trồng sạch là danh dự đã chốt sổ".
            </p>
          </div>
        </div>

        <!-- Feature Block 2 -->
        <div class="flex flex-col md:flex-row-reverse items-center gap-12 mb-10">
          <div class="w-full md:w-1/2">
            <div class="relative w-full rounded-3xl overflow-hidden shadow-2xl border border-slate-200 group">
              <!-- Background Image -->
              <img src="/saleswoman_rice.png" alt="Gạo Hữu Cơ Hạt Ngọc - Cô Cửa Hàng" class="w-full h-auto object-cover" />
              
              <!-- Injected Live QR Code absolutely positioned exactly on the white label -->
              <div class="absolute z-10" style="bottom: 15%; left: 50%; transform: translateX(-50%); width: 15%; aspect-ratio: 1/1;">
                <img src="https://api.qrserver.com/v1/create-qr-code/?margin=2&size=180x180&data=https://trustfab.github.io/gapchain/qr/consumer/LH-HTX001-1775545225847" 
                     alt="Scan ME" 
                     class="w-full h-full mix-blend-darken opacity-90 group-hover:opacity-100 transition-all cursor-pointer hover:scale-105 rounded-sm" />
              </div>

              <!-- Floating badge hinting users to scan -->
              <div class="absolute top-6 left-6 animate-bounce" style="animation-duration: 4s">
                <span class="text-xs font-bold text-white bg-blue-600/90 backdrop-blur px-4 py-2 rounded-full shadow-xl border border-white/20 flex items-center gap-2 uppercase tracking-wide">
                  📱 Quét để Truy Nguồn "1 Chạm"
                </span>
              </div>
            </div>
          </div>
          <div class="w-full md:w-1/2 lg:pr-12 text-left">
            <div class="w-12 h-12 bg-blue-100 text-blue-600 rounded-xl flex items-center justify-center font-bold mb-4">2</div>
            <h3 class="text-3xl font-extrabold text-slate-900 mb-4">Một Trình Quét — Thỏa Sự Tò Mò</h3>
            <p class="text-slate-600 text-lg leading-relaxed mb-6">
              Trao hệ thống xác thực vào tay người tiêu dùng. Thay vì dòng chữ vô thưởng vô phạt trên nhãn hiệu bao bì, quét mã QR sẽ kéo ra Timeline hiển thị đến từng lô đất nuôi trồng thực tế. Nếu tem bị lạm dụng scan nhiều nơi một lúc, app sẽ báo động cảnh báo giả mạo lập tức.
            </p>
          </div>
        </div>

        <!-- Feature Block 3 -->
        <div class="flex flex-col md:flex-row items-center gap-12 mt-20">
          <div class="w-full md:w-1/2">
            <div class="bg-gradient-to-br from-indigo-50 to-blue-100 p-8 rounded-3xl border border-indigo-200 shadow-sm relative overflow-hidden h-full min-h-[250px] flex items-center justify-center">
               <!-- Abstract AI / IoT visualization -->
               <div class="absolute inset-0 opacity-30" style="background-image: radial-gradient(#4f46e5 1px, transparent 1px); background-size: 20px 20px;"></div>
               <div class="relative z-10 bg-white/90 backdrop-blur-sm p-6 rounded-2xl border border-white shadow-xl text-center w-4/5 md:w-3/4 transform transition hover:scale-105 duration-300">
                 <div class="text-4xl mb-4">🤖 📡</div>
                 <p class="font-bold text-slate-800 text-lg mb-2">Hệ Sinh Thái Thông Minh</p>
                 <div class="flex flex-wrap justify-center gap-2 mt-3">
                   <span class="text-[10px] font-bold px-2 py-1 bg-indigo-100 text-indigo-700 rounded-lg">IoT Sensors</span>
                   <span class="text-[10px] font-bold px-2 py-1 bg-indigo-100 text-indigo-700 rounded-lg">AI Optimization</span>
                 </div>
               </div>
            </div>
          </div>
          <div class="w-full md:w-1/2 md:pl-8 text-left">
            <div class="w-12 h-12 bg-indigo-100 text-indigo-600 rounded-xl flex items-center justify-center font-bold mb-4">3</div>
            <h3 class="text-3xl font-extrabold text-slate-900 mb-4">Tối Ưu Hóa Bằng AI & IoT</h3>
            <p class="text-slate-600 text-lg leading-relaxed mb-6">
              Nền tảng của chúng tôi <strong>thiết kế bằng kiến trúc dựa trên quy tắc (rule-based)</strong> với tính linh hoạt và khả năng mở rộng cao. Đặc biệt hỗ trợ <strong>xây dựng hệ thống AI và IoT</strong> để đo lường chính xác, qua đó tối ưu hóa toàn bộ quy trình sản xuất và mạng lưới phân phối. Nhờ áp dụng <strong>phát triển phần mềm hỗ trợ bằng AI</strong>, tiến độ triển khai liên tục các tính năng phục vụ doanh nghiệp bạn sẽ được đẩy nhanh thần tốc.
            </p>
          </div>
        </div>

      </div>
    </section>

    <!-- SECTION 6 — CTA ĐĂNG KÝ -->
    <section class="py-24 bg-gradient-to-b from-slate-900 to-[#12182b] text-white overflow-hidden relative">
      <div class="absolute w-full h-full inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMSIgZmlsbD0iI2ZmZiIgZmlsbC1vcGFjaXR5PSIuMSIvPjwvc3ZnPg==')] opacity-60"></div>
      
      <div class="max-w-4xl mx-auto px-6 relative z-10 text-center">
         <h2 class="text-4xl md:text-5xl font-extrabold text-white mb-8 leading-tight">Sẵn Sàng Hệ Thống Chuyển Đổi Số <br/>Ngay Từ Hôm Nay</h2>
         <p class="text-xl text-slate-400 mb-12 max-w-2xl mx-auto font-medium leading-relaxed">
           Chúng tôi đang tiếp nhận đăng ký tham gia Thí điểm Chương Trình (Pilot) cho các HTX & Đơn vị phân phối chiến lược khu vực Tây Nam Bộ và Tây Nguyên. Hỗ trợ hệ thống không chi phí thiết lập.
         </p>
         
         <form class="bg-white/5 p-2 rounded-2xl border border-slate-700 max-w-lg mx-auto flex flex-col md:flex-row gap-2 shadow-2xl backdrop-blur-sm">
           <input type="email" placeholder="Nhập địa chỉ email của bạn..." class="w-full bg-slate-800/80 border-none text-white px-6 py-4 rounded-xl focus:ring-2 focus:ring-emerald-500 placeholder-slate-500 font-medium"/>
           <button type="button" class="w-full md:w-auto px-8 py-4 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-xl transition-colors whitespace-nowrap flex-shrink-0 shadow-lg shadow-emerald-900/30">
             🔥 Nhận Tư Vấn Mở Cổng
           </button>
         </form>
         <p class="text-sm font-medium text-slate-500 mt-6 flex items-center justify-center gap-2">
           <span>🛡️ Cam kết bảo mật dữ liệu</span> |
           <span>Khởi tạo Demo chỉ 24h</span>
         </p>
      </div>
    </section>

    <!-- SECTION 7 — FOOTER -->
    <footer class="bg-[#05080f] py-16 border-t border-slate-800/80 text-center relative z-10">
      <div class="max-w-7xl mx-auto px-6">
        <div class="flex items-center justify-center gap-3 mb-6 opacity-60">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-slate-600 to-slate-800 flex items-center justify-center font-bold text-white text-sm">G</div>
          <span class="font-extrabold text-xl text-slate-400 tracking-tight">GAPChain Business MVP</span>
        </div>
        <p class="text-sm font-medium text-slate-500 mb-6 flex flex-wrap items-center justify-center gap-3">
          <span class="px-3 py-1 bg-slate-800 rounded-md shadow-sm border border-slate-700/50">Traceability Tech</span>
          <span class="px-3 py-1 bg-slate-800 rounded-md shadow-sm border border-slate-700/50">VietGAP Digital</span>
          <span class="px-3 py-1 bg-emerald-900/40 text-emerald-600 rounded-md shadow-sm border border-emerald-900">Green Agriculture</span>
        </p>
        <div class="mt-8 flex justify-center space-x-6 text-sm font-bold">
          <a href="#" class="text-slate-500 hover:text-emerald-400 transition-colors">Tài liệu hướng dẫn APP</a>
          <a href="#" class="text-slate-500 hover:text-emerald-400 transition-colors">Chứng nhận hệ thống</a>
          <a href="/gapchain/login" class="text-slate-500 hover:text-emerald-400 transition-colors">Login Ban Quản Trị</a>          
        </div>
        <p class="text-xs text-slate-600 mt-8">&copy; 2026. Một giải pháp cho dự án Nông Nghiệp Phát Triển Bền Vững.</p>
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
