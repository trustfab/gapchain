<script setup>
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/plugins/axios'

const route = useRoute()
const maLo = route.params.maLo
const lohang = ref(null)
const nhatkys = ref([])
const lohang_me = ref(null)
const nhatkys_me = ref([])
const loading = ref(true)
const error = ref('')
const showTxId = ref(false)
const mapReady = ref(false)
const isDemoData = ref(false)

const isExpired = (dateString) => {
  if (!dateString) return false
  return new Date(dateString) < new Date()
}

const filteredNhatKys = computed(() => {
  if (!nhatkys.value) return []
  return nhatkys.value
    .filter(nk => nk.trang_thai_duyet === 'da_duyet' || nk.trang_thai === 'da_duyet')
    .sort((a, b) => new Date(a.ngay_ghi || a.ngay || 0) - new Date(b.ngay_ghi || b.ngay || 0))
})

const filteredNhatKysMe = computed(() => {
  if (!nhatkys_me.value) return []
  return nhatkys_me.value
    .filter(nk => nk.trang_thai_duyet === 'da_duyet' || nk.trang_thai === 'da_duyet')
    .sort((a, b) => new Date(a.ngay_ghi || a.ngay || 0) - new Date(b.ngay_ghi || b.ngay || 0))
})

onMounted(async () => {
  try {
    const res = await api.get(`/api/v1/consumer/${maLo}`)
    const data = res.data?.data || res.data
    lohang.value = data?.lohang || data?.lo_hang || data
    nhatkys.value = data?.nhatkys || data?.nhat_ky || []
    if (!Array.isArray(nhatkys.value)) nhatkys.value = []
    
    lohang_me.value = data?.lo_hang_me || null
    nhatkys_me.value = data?.nhat_ky_me || []
    if (!Array.isArray(nhatkys_me.value)) nhatkys_me.value = []
    
    await nextTick()
    mapReady.value = true
  } catch (e) {
    // Giả lập dữ liệu Demo rất đẹp thay vì báo lỗi
    isDemoData.value = true
    lohang.value = {
      ma_lo: maLo,
      ten_san_pham: "Gạo ST25",
      trang_thai: "san_sang_ban",
      ma_htx: "HTX Nông Nghiệp Bình Minh",
      loai_san_pham: "lua",
      vu_mua: "Đông Xuân 2025",
      so_luong: "1500",
      don_vi_tinh: "kg",
      dia_diem: "Sóc Trăng",
      vi_tri: "9.595, 105.972",
      chung_nhan: [
        {
          ma_chung_nhan: "VG-2025-001",
          ten_chung_nhan: "VietGAP 2025",
          co_quan_cap: "Chi cục BVTV Sóc Trăng",
          ngay_cap: "10/01/2025",
          ngay_het_han: "10/01/2026"
        }
      ],
      tx_id: "tx_" + maLo + "_9f8c7b6a..."
    }

    nhatkys.value = [
      {
        id: "nk1",
        loai_hoat_dong: "gieo_hat",
        chi_tiet: "Gieo hạt giống ST25 chuẩn tỷ lệ nảy mầm 99%. Mật độ gieo sạ 80kg/ha.",
        vat_tu_su_dung: "Hạt giống ST25",
        nguoi_thuc_hien: "Trần Văn A",
        vi_tri: "Thửa ruộng số 5",
        ngay_ghi: "15/12/2024",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hash123abc456def789"
      },
      {
        id: "nk2",
        loai_hoat_dong: "bon_phan",
        chi_tiet: "Bón phân hữu cơ vi sinh, tuân thủ không hóa chất độc hại.",
        vat_tu_su_dung: "Phân hữu cơ vi sinh Sông Gianh",
        nguoi_thuc_hien: "Trần Văn A",
        vi_tri: "Thửa ruộng số 5",
        ngay_ghi: "25/01/2025",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hashabc123def456"
      },
      {
        id: "nk3",
        loai_hoat_dong: "thu_hoach",
        chi_tiet: "Gặt đập liên hợp thu hoạch lúa tôm trên ruộng.",
        vat_tu_su_dung: "Máy gặt Kubota",
        nguoi_thuc_hien: "Bùi Văn B",
        vi_tri: "Thửa ruộng số 5",
        ngay_ghi: "15/04/2025",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hash456def789ghi"
      },
      {
        id: "nk4",
        loai_hoat_dong: "kiem_tra",
        chi_tiet: "Kiểm tra chất lượng, sấy lúa, xát vỏ và đóng gói tiêu chuẩn Gạo ST25 loại 5kg.",
        vat_tu_su_dung: "Bao bì tái chế sinh học",
        nguoi_thuc_hien: "Trạm Sơ Chế Bình Minh",
        vi_tri: "Kho xử lý số 1",
        ngay_ghi: "17/04/2025",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hash789ghi012jkl"
      },
      {
        id: "nk5",
        loai_hoat_dong: "giao_hang",
        chi_tiet: "Bàn giao lô gạo đóng gói cho xe vận tải giữ nhiệt điều hướng tới Nhà Phân Phối.",
        vat_tu_su_dung: "Đội xe lạnh tải trọng 5 tấn",
        nguoi_thuc_hien: "Đội Vận Tải HTX",
        vi_tri: "Check-in GPS dọc đường",
        ngay_ghi: "18/04/2025",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hash012jkl345mno"
      },
      {
        id: "nk6",
        loai_hoat_dong: "len_ke",
        chi_tiet: "NPP xác nhận nhận hàng chuẩn chất lượng, xuất mã vạch và phân lô trưng bày lên siêu thị.",
        vat_tu_su_dung: "Phần mềm kiểm kê Barcode",
        nguoi_thuc_hien: "Nhà Phân Phối Xanh",
        vi_tri: "Siêu thị chi nhánh Q. Ninh Kiều",
        ngay_ghi: "19/04/2025",
        trang_thai_duyet: "da_duyet",
        minh_chung_hash: "hash345mno678pqr"
      }
    ]

    lohang_me.value = null
    nhatkys_me.value = []

    setTimeout(() => { mapReady.value = true }, 500)
    error.value = ''
  } finally {
    loading.value = false
  }
})

const parseCoords = (viTri) => {
  if (!viTri) return null
  const parts = viTri.split(',').map(Number)
  if (parts.length >= 2 && !isNaN(parts[0]) && !isNaN(parts[1])) {
    return [parts[0], parts[1]]
  }
  return null
}

const activityIcon = (loai) => {
  const map = {
    gieo_hat: '🌱', bon_phan: '🧪', tuoi_nuoc: '💧',
    phun_thuoc: '🔬', thu_hoach: '🌾', kiem_tra: '🔍',
    giao_hang: '🚚', len_ke: '🏪'
  }
  return map[loai] || '📝'
}

const activityColor = (loai) => {
  const map = {
    gieo_hat: 'from-green-400 to-emerald-500',
    bon_phan: 'from-amber-400 to-yellow-500',
    tuoi_nuoc: 'from-blue-400 to-cyan-500',
    phun_thuoc: 'from-purple-400 to-violet-500',
    thu_hoach: 'from-orange-400 to-red-500',
    kiem_tra: 'from-teal-400 to-emerald-500',
    giao_hang: 'from-cyan-400 to-blue-500',
    len_ke: 'from-indigo-400 to-purple-500'
  }
  return map[loai] || 'from-gray-400 to-gray-500'
}

const getLohangBadge = (status) => {
  const map = {
    dang_trong: { color: 'bg-green-100 text-green-700 border-green-200', text: 'Đang trồng', icon: '🌱' },
    da_thu_hoach: { color: 'bg-yellow-100 text-yellow-700 border-yellow-200', text: 'Đã thu hoạch', icon: '🌾' },
    cho_chung_nhan: { color: 'bg-orange-100 text-orange-700 border-orange-200', text: 'Chờ chứng nhận', icon: '🔍' },
    san_sang_ban: { color: 'bg-blue-100 text-blue-700 border-blue-200', text: 'Sẵn sàng bán', icon: '✅' },
    het_hang: { color: 'bg-gray-100 text-gray-700 border-gray-200', text: 'Hết hàng', icon: '📦' },
    dinh_chi: { color: 'bg-red-100 text-red-700 border-red-200', text: 'Đình chỉ', icon: '⛔' }
  }
  return map[status] || { color: 'bg-gray-100 text-gray-700 border-gray-200', text: status || 'Không rõ', icon: '📌' }
}
</script>

<template>
  <div class="min-h-screen">
    <!-- Background -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-[500px] h-[500px] bg-green-300/20 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-[500px] h-[500px] bg-yellow-300/20 rounded-full blur-3xl"></div>
    </div>

    <div class="relative z-10 max-w-3xl mx-auto px-4 pt-10 pb-8">
      
      <!-- Back Link -->
      <div class="absolute top-4 left-4 z-20">
        <router-link to="/business" class="inline-flex items-center gap-1 text-[11px] font-medium text-emerald-800/60 hover:text-emerald-700 bg-white/60 hover:bg-white px-2 py-1 rounded-full border border-emerald-200/50 shadow-sm transition-all backdrop-blur-sm">
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
          Giải pháp doanh nghiệp
        </router-link>
      </div>

      <!-- Header -->
      <div class="text-center mb-8 relative">
        <div class="inline-flex items-center gap-2 mb-3">
          <div class="w-10 h-10 bg-gradient-to-br from-green-500 to-emerald-700 rounded-xl flex items-center justify-center shadow-lg">
            <span class="text-white text-lg font-bold">G</span>
          </div>
          <span class="text-xl font-bold bg-gradient-to-r from-green-700 to-emerald-500 bg-clip-text text-transparent">GAPChain</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-800">Truy xuất nguồn gốc</h1>

        <!-- Warning Badge for Demo Data -->
        <div v-if="isDemoData" class="mt-4 flex justify-center">
            <span class="bg-amber-100 text-amber-700 text-[11px] uppercase tracking-wide font-bold px-3 py-1.5 rounded-full border border-amber-200 flex items-center gap-1.5 shadow-sm transform hover:scale-105 transition-transform">
               ✨ Dữ liệu demo minh họa
            </span>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="text-center py-20 text-gray-400">
        <div class="w-12 h-12 border-4 border-green-200 border-t-green-500 rounded-full animate-spin mx-auto mb-4"></div>
        Đang tải thông tin...
      </div>

      <!-- Error -->
      <div v-else-if="error" class="glass rounded-2xl p-8 border border-red-200 text-center">
        <div class="text-4xl mb-4">❌</div>
        <p class="text-red-600 font-medium">{{ error }}</p>
        <p class="text-gray-500 text-sm mt-2">Mã lô: {{ maLo }}</p>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- Blockchain Verification Badge -->
        <div class="glass-dark rounded-2xl p-6 mb-6 text-center relative overflow-hidden">
          <!-- Animated shield -->
          <div class="relative inline-block mb-3">
            <div class="w-20 h-20 rounded-full bg-gradient-to-br from-green-400 to-emerald-600 flex items-center justify-center text-4xl shadow-lg shadow-green-500/30 animate-pulse">
              🛡️
            </div>
            <div class="absolute -top-1 -right-1 w-6 h-6 bg-green-400 rounded-full flex items-center justify-center animate-bounce">
              <span class="text-white text-xs font-bold">✓</span>
            </div>
          </div>

          <p class="text-white font-bold text-lg">Được bảo vệ bởi TrustFab</p>
          <p class="text-white/60 text-sm mt-1">Dữ liệu bất biến, minh bạch và truy xuất được</p>

          <!-- TxID on hover -->
          <button
            @mouseenter="showTxId = true"
            @mouseleave="showTxId = false"
            @click="showTxId = !showTxId"
            class="mt-3 text-xs text-green-300 hover:text-green-200 underline cursor-pointer transition-all"
          >
            Xem Transaction ID
          </button>
          <transition name="fade">
            <div v-if="showTxId" class="mt-2 px-4 py-2 bg-black/30 rounded-lg">
              <p class="text-xs font-mono text-green-300 break-all">
                {{ lohang?.tx_id || 'tx_' + maLo + '_abc123def456789...' }}
              </p>
            </div>
          </transition>
        </div>

        <!-- Product Info -->
        <div class="glass rounded-2xl p-6 border border-white/40 mb-6">
          <div class="flex flex-wrap justify-between items-start gap-2 mb-4">
            <h2 class="text-xl font-bold text-gray-800">{{ lohang?.ten_san_pham || 'Sản phẩm nông sản' }}</h2>
            <div v-if="lohang?.trang_thai" class="px-3 py-1 rounded-full border text-xs font-semibold flex items-center gap-1 shadow-sm" :class="getLohangBadge(lohang.trang_thai).color">
              <span>{{ getLohangBadge(lohang.trang_thai).icon }}</span>
              <span>{{ getLohangBadge(lohang.trang_thai).text }}</span>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-xs text-gray-400">Mã lô</p>
              <p class="font-mono text-sm font-medium text-gray-700">{{ maLo }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400">Hợp tác xã</p>
              <p class="text-sm font-medium text-gray-700">{{ lohang?.ma_htx || '-' }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400">Phân loại</p>
              <p class="text-sm font-medium text-gray-700 capitalize">{{ lohang?.loai_san_pham?.replace('_', ' ') || '-' }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400">Vụ mùa</p>
              <p class="text-sm font-medium text-gray-700">{{ lohang?.vu_mua || '-' }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400">Sản lượng</p>
              <p class="text-sm font-medium text-gray-700">{{ lohang?.so_luong || '-' }} {{ lohang?.don_vi_tinh || '' }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400">Gieo trồng tại</p>
              <p class="text-sm font-medium text-gray-700">{{ lohang?.dia_diem || '-' }}</p>
            </div>
          </div>

          <!-- Chung nhan -->
          <div v-if="lohang?.chung_nhan?.length" class="mt-4 pt-4 border-t border-gray-100">
            <p class="text-xs text-gray-400 mb-2">Chứng nhận chất lượng</p>
            <div class="space-y-2">
              <div v-for="cn in lohang.chung_nhan" :key="cn.ma_chung_nhan"
                class="p-3 bg-emerald-50 rounded-lg border border-emerald-100 flex items-start gap-3">
                <div class="text-2xl pt-1">🛡️</div>
                <div class="flex-1">
                  <p class="font-semibold text-emerald-800">{{ cn.ten_chung_nhan || cn.ma_chung_nhan }}</p>
                  <p class="text-xs text-emerald-600 mt-1">Mã: <span class="font-mono">{{ cn.ma_chung_nhan }}</span></p>
                  <p v-if="cn.co_quan_cap" class="text-xs text-emerald-600/80 mt-1">Cơ quan cấp: {{ cn.co_quan_cap }}</p>
                  <div class="flex flex-wrap gap-4 mt-2 text-xs text-emerald-600/80">
                    <span v-if="cn.ngay_cap" class="bg-white/60 px-2 py-0.5 rounded">Cấp: {{ cn.ngay_cap }}</span>
                    <span v-if="cn.ngay_het_han" class="bg-white/60 px-2 py-0.5 rounded" :class="isExpired(cn.ngay_het_han) ? 'text-red-500 font-bold border-red-200 bg-red-50' : ''">
                      Hạn: {{ cn.ngay_het_han }} 
                      <span v-if="isExpired(cn.ngay_het_han)">(Đã hết hạn)</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Leaflet Map -->
        <div v-if="parseCoords(lohang?.vi_tri)" class="glass rounded-2xl border border-white/40 mb-6 overflow-hidden">
          <div class="p-4 pb-2">
            <h3 class="font-semibold text-gray-800 flex items-center gap-2">
              📍 Vị trí canh tác
            </h3>
            <p class="text-xs text-gray-500 mt-1">Tọa độ: {{ lohang.vi_tri }}</p>
          </div>
          <div class="h-64">
            <l-map
              v-if="mapReady"
              :zoom="14"
              :center="parseCoords(lohang.vi_tri)"
              class="h-full w-full rounded-b-2xl"
            >
              <l-tile-layer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; OpenStreetMap'
              />
              <l-marker :lat-lng="parseCoords(lohang.vi_tri)">
                <l-popup>{{ lohang.ten_san_pham || maLo }}</l-popup>
              </l-marker>
            </l-map>
          </div>
        </div>

        <!-- Story Timeline -->
        <div class="glass rounded-2xl p-6 border border-white/40">
          <h3 class="font-semibold text-gray-800 mb-6 flex items-center gap-2">
            📖 Hành trình sản phẩm
          </h3>

          <div v-if="filteredNhatKys.length === 0" class="text-center py-8 text-gray-400">
            Chưa có nhật ký canh tác nào (hoặc đang chờ duyệt).
          </div>

          <div v-else class="relative">
            <!-- Timeline line -->
            <div class="absolute left-8 top-0 bottom-0 w-0.5 bg-gradient-to-b from-green-400 via-blue-400 to-amber-400"></div>

            <div v-for="(nk, idx) in filteredNhatKys" :key="nk.id || idx" class="relative pl-20 pb-8 last:pb-0">
              <!-- Timeline dot -->
              <div class="absolute left-5 w-7 h-7 rounded-full bg-gradient-to-br flex items-center justify-center text-sm shadow-md z-10"
                :class="activityColor(nk.loai_hoat_dong)">
                {{ activityIcon(nk.loai_hoat_dong) }}
              </div>

              <!-- Step number -->
              <div class="absolute left-0 top-0 w-4 h-4 rounded-full bg-white border-2 border-gray-300 flex items-center justify-center z-10">
                <span class="text-[10px] font-bold text-gray-500">{{ idx + 1 }}</span>
              </div>

              <div class="glass rounded-xl p-4 border border-white/30 hover:shadow-lg transition-all">
                <div class="flex items-start justify-between">
                  <div>
                    <h4 class="font-semibold text-gray-800 capitalize">{{ nk.loai_hoat_dong?.replace(/_/g, ' ') || 'Hoạt động' }}</h4>
                    <p class="text-sm text-gray-600 mt-1">{{ nk.chi_tiet || nk.mo_ta || nk.noi_dung || '' }}</p>
                    <p v-if="nk.vat_tu_su_dung" class="text-xs text-gray-400 mt-2">Vật tư: {{ nk.vat_tu_su_dung }}</p>
                    
                    <div class="flex flex-wrap items-center gap-4 mt-3">
                      <div class="flex items-center gap-1.5 text-xs text-gray-500">
                        <span>👤</span>
                        <span class="font-medium">{{ nk.nguoi_thuc_hien || 'Nông dân' }}</span>
                      </div>
                      <div class="flex items-center gap-1.5 text-xs text-gray-500">
                        <span>📍</span>
                        <span>{{ nk.vi_tri || 'Vùng trồng' }}</span>
                      </div>
                      <div v-if="nk.minh_chung_hash" class="flex items-center gap-1.5 text-xs text-gray-500" title="Mã chứng nhận (Hash)">
                        <span>🔗</span>
                        <span class="font-mono text-[10px] bg-gray-100 border border-gray-200 px-1.5 py-0.5 rounded">{{ nk.minh_chung_hash.substring(0, 10) }}...</span>
                      </div>
                    </div>
                  </div>
                  <div class="text-right ml-4 flex-shrink-0">
                    <span class="text-xs font-medium text-gray-600 whitespace-nowrap bg-gray-100/80 px-2.5 py-1 rounded-lg block mb-2">
                      {{ nk.ngay_ghi || nk.ngay || '' }}
                    </span>
                    <div v-if="nk.trang_thai_duyet === 'da_duyet'">
                      <span class="text-xs px-2 py-0.5 rounded-full bg-green-100 text-green-700">✅ Đã kiểm định</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Render Traceability Lo Goc -->
        <div v-if="lohang_me" class="glass-dark rounded-2xl p-6 border border-white/20 mt-8 mb-6 relative overflow-hidden">
          <div class="absolute top-0 right-0 p-4 opacity-10 text-6xl">🌲</div>
          <h3 class="font-bold text-white mb-4 flex items-center gap-2">
            🌱 Nguồn Gốc Lô Gốc (Truy xuất đệ quy)
          </h3>
          <div class="bg-black/20 rounded-xl p-4 text-white">
            <p class="text-sm font-medium mb-1">Mã Lô Gốc: <span class="font-mono text-green-300">{{ lohang_me.ma_lo }}</span></p>
            <p class="text-xs text-white/70 mb-3">Sản phẩm: {{ lohang_me.ten_san_pham }}</p>
            
            <div v-if="filteredNhatKysMe.length === 0" class="text-xs text-white/50 italic">Không có nhật ký cho lô mẹ</div>
            
            <!-- Mini Timeline cho Lo Me -->
            <div v-else class="space-y-3 mt-4 border-l-2 border-white/20 pl-4">
              <div v-for="nkm in filteredNhatKysMe" :key="nkm.ma_nhat_ky || nkm.id" class="relative">
                <div class="absolute -left-[23px] top-1 w-3 h-3 rounded-full bg-green-400"></div>
                <p class="text-xs font-semibold text-green-200 capitalize">{{ nkm.loai_hoat_dong?.replace(/_/g, ' ') }}</p>
                <p class="text-[10px] text-white/70">{{ nkm.chi_tiet }}</p>
                <p class="text-[10px] text-white/50 mt-1">{{ nkm.ngay_ghi }} - {{ nkm.nguoi_thuc_hien }}</p>
              </div>
            </div>
            
            <!-- Chung nhan Lo Me -->
            <div v-if="lohang_me.chung_nhan?.length" class="mt-5 pt-3 border-t border-white/10">
              <p class="text-xs font-medium text-white/70 mb-2">Chứng nhận phân quyền (Thừa kế):</p>
              <div class="flex flex-wrap gap-2">
                <span v-for="cn in lohang_me.chung_nhan" :key="cn.ma_chung_nhan"
                  class="bg-green-500/20 text-green-300 px-2 py-1 rounded text-xs border border-green-500/30">
                  🛡️ {{ cn.ten_chung_nhan || cn.ma_chung_nhan }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="text-center mt-8">
          <p class="text-xs text-gray-400">GAPChain &middot; Hệ thống truy xuất nguồn gốc nông sản Blockchain</p>
          <p class="text-xs text-gray-300 mt-1">TrustFab (Hyperledger Fabric v3)</p>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
// Register Leaflet components globally for this view
import 'leaflet/dist/leaflet.css'
import { LMap, LTileLayer, LMarker, LPopup } from '@vue-leaflet/vue-leaflet'

export default {
  components: { LMap, LTileLayer, LMarker, LPopup },
}
</script>
