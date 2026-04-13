<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const giaoDichs = ref([])
const availableLohangs = ref([])
const loading = ref(true)
const showForm = ref(false)

const form = ref({
  ma_lo: '',
  ma_npp: '',
  san_pham: '',
  so_luong: '',
  don_vi_tinh: 'kg',
  don_gia: '',
  ty_le_hoa_hong: '5',
  ghi_chu: '',
})
const submitting = ref(false)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    // HTX xem danh sach giao dich cua minh
    const res = await api.get(`/api/v1/giaodich/htx/${auth.tenantId}`)
    const data = res.data?.data || res.data || []
    giaoDichs.value = Array.isArray(data) ? data : []

    // Fetch Danh sach Lo Hang cua HTX de lap dropdown
    const resLo = await api.get(`/api/v1/lohang/htx/${auth.tenantId}`)
    const loData = resLo.data?.data || resLo.data || []
    if (Array.isArray(loData)) {
      availableLohangs.value = loData.filter(lo => {
        const ton = typeof lo.so_luong_con_lai === 'number' ? lo.so_luong_con_lai : lo.so_luong
        return lo.trang_thai === 'san_sang_ban' && ton > 0
      })
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const selectedLohangInventory = computed(() => {
  const lo = availableLohangs.value.find(l => l.ma_lo === form.value.ma_lo)
  if (!lo) return null
  return typeof lo.so_luong_con_lai === 'number' ? lo.so_luong_con_lai : lo.so_luong
})

const isOverInventory = computed(() => {
  if (selectedLohangInventory.value === null) return false
  const sl = parseFloat(form.value.so_luong) || 0
  return sl > selectedLohangInventory.value
})

watch(() => form.value.ma_lo, (newMaLo) => {
  const lo = availableLohangs.value.find(l => l.ma_lo === newMaLo)
  if (lo) {
    form.value.san_pham = lo.ten_san_pham || form.value.san_pham
    form.value.don_vi_tinh = lo.don_vi_tinh || form.value.don_vi_tinh
  }
})

async function taoGiaoDich() {
  submitting.value = true
  try {
    const now = Date.now()
    await api.post('/api/v1/giaodich', {
      ma_giao_dich: `GD-${auth.tenantId}-${now}`,
      ma_lo: form.value.ma_lo,
      ma_htx: auth.tenantId,
      ma_npp: form.value.ma_npp,
      san_pham: form.value.san_pham,
      so_luong: parseFloat(form.value.so_luong) || 0,
      don_vi_tinh: form.value.don_vi_tinh,
      don_gia: parseFloat(form.value.don_gia) || 0,
      ty_le_hoa_hong: parseFloat(form.value.ty_le_hoa_hong) || 0,
      ghi_chu: form.value.ghi_chu,
    })
    showForm.value = false
    form.value = { ma_lo: '', ma_npp: '', san_pham: '', so_luong: '', don_vi_tinh: 'kg', don_gia: '', ty_le_hoa_hong: '5', ghi_chu: '' }
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi tạo giao dịch')
  } finally {
    submitting.value = false
  }
}

const formatStatus = (s) => {
  const map = {
    cho_duyet: { label: '⏳ Chờ duyệt', class: 'bg-yellow-100 text-yellow-700' },
    da_duyet: { label: '📝 Đã duyệt', class: 'bg-blue-100 text-blue-700' },
    dang_giao: { label: '🚚 Đang giao', class: 'bg-purple-100 text-purple-700' },
    da_giao: { label: '📦 Đã nhận hàng', class: 'bg-indigo-100 text-indigo-700' },
    cho_thanh_toan: { label: '💳 Chờ thanh toán', class: 'bg-orange-100 text-orange-700' },
    da_thanh_toan: { label: '💰 Đã thanh toán', class: 'bg-green-100 text-green-700' },
    huy_bo: { label: '❌ Đã hủy', class: 'bg-gray-100 text-gray-700' }
  }
  return map[s] || { label: s || '⏳ Chờ duyệt', class: 'bg-gray-100 text-gray-700' }
}

const stats = computed(() => {
  const total = giaoDichs.value.length
  let pending = 0
  let completed = 0
  let totalRevenue = 0
  let pendingRevenue = 0

  giaoDichs.value.forEach(gd => {
    if (gd.trang_thai === 'da_thanh_toan') {
      completed++
      totalRevenue += (gd.so_luong * gd.don_gia)
    } else if (gd.trang_thai !== 'huy_bo') {
      pending++
      pendingRevenue += (gd.so_luong * gd.don_gia)
    }
  })

  return { total, pending, completed, totalRevenue, pendingRevenue }
})
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Giao dịch thu mua</h2>
      <button @click="showForm = !showForm"
        class="px-5 py-2.5 bg-gradient-to-r from-amber-500 to-orange-600 text-white rounded-xl font-medium btn-glow flex items-center gap-2 transition-transform hover:scale-105 shadow-md">
        <span>{{ showForm ? '✕ Đóng form' : '+ Tạo Hợp Đồng' }}</span>
      </button>
    </div>

    <!-- Overview Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
      <div class="glass rounded-2xl p-4 border border-white/40 shadow-sm flex flex-col justify-center relative justify-between overflow-hidden group hover:shadow-md transition-all">
        <div class="absolute -right-4 -top-4 w-16 h-16 bg-blue-100 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-500"></div>
        <p class="text-xs text-gray-500 font-semibold mb-1 relative z-10 flex items-center gap-1.5"><span class="text-blue-500">📊</span> Tổng Hợp Đồng</p>
        <p class="text-2xl font-bold text-gray-800 relative z-10 font-mono">{{ stats.total }}</p>
      </div>
      <div class="glass rounded-2xl p-4 border border-white/40 shadow-sm flex flex-col justify-center relative justify-between overflow-hidden group hover:shadow-md transition-all">
         <div class="absolute -right-4 -top-4 w-16 h-16 bg-amber-100 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-500"></div>
        <p class="text-xs text-gray-500 font-semibold mb-1 relative z-10 flex items-center gap-1.5"><span class="text-amber-500">⏳</span> Đang Xử Lý</p>
        <p class="text-2xl font-bold text-amber-600 relative z-10 font-mono">{{ stats.pending }}</p>
      </div>
      <div class="glass rounded-2xl p-4 border border-white/40 shadow-sm flex flex-col justify-center relative justify-between overflow-hidden group hover:shadow-md transition-all">
         <div class="absolute -right-4 -top-4 w-16 h-16 bg-emerald-100 rounded-full opacity-50 group-hover:scale-150 transition-transform duration-500"></div>
        <p class="text-xs text-gray-500 font-semibold mb-1 relative z-10 flex items-center gap-1.5"><span class="text-emerald-500">✅</span> Hoàn Tất</p>
        <p class="text-2xl font-bold text-emerald-600 relative z-10 font-mono">{{ stats.completed }}</p>
      </div>
      <div class="bg-gradient-to-br from-orange-50 to-amber-100/60 rounded-2xl p-4 border border-amber-100 shadow-sm flex flex-col justify-center relative justify-between overflow-hidden">
        <p class="text-xs text-amber-700 font-semibold mb-1 relative z-10 flex items-center gap-1.5">💸 Chờ Thu Tiền (VNĐ)</p>
        <p class="text-xl font-bold text-amber-700 relative z-10 font-mono tracking-tight" :title="stats.pendingRevenue.toLocaleString()">{{ stats.pendingRevenue >= 1000000000 ? (stats.pendingRevenue / 1000000000).toFixed(2) + ' Tỷ' : stats.pendingRevenue >= 1000000 ? (stats.pendingRevenue / 1000000).toFixed(1) + ' Tr' : stats.pendingRevenue.toLocaleString() }}</p>
      </div>
      <div class="bg-gradient-to-br from-green-50 to-emerald-100/60 rounded-2xl p-4 border border-emerald-100 shadow-sm flex flex-col justify-center relative justify-between overflow-hidden">
        <p class="text-xs text-green-700 font-semibold mb-1 relative z-10 flex items-center gap-1.5">💰 Thực Thu (VNĐ)</p>
        <p class="text-xl font-bold text-green-700 relative z-10 font-mono tracking-tight" :title="stats.totalRevenue.toLocaleString()">{{ stats.totalRevenue >= 1000000000 ? (stats.totalRevenue / 1000000000).toFixed(2) + ' Tỷ' : stats.totalRevenue >= 1000000 ? (stats.totalRevenue / 1000000).toFixed(1) + ' Tr' : stats.totalRevenue.toLocaleString() }}</p>
      </div>
    </div>

    <!-- Form -->
    <transition name="slide">
      <div v-if="showForm" class="glass rounded-2xl p-6 border border-white/40 mb-6">
        <h3 class="font-semibold text-gray-800 mb-4">Tạo giao dịch mới</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select v-model="form.ma_lo" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70">
            <option value="" disabled>-- Chọn một lô hàng (đã thu hoạch) * --</option>
            <option v-for="lo in availableLohangs" :key="lo.ma_lo" :value="lo.ma_lo">
              {{ lo.ten_san_pham }} (Tồn: {{ typeof lo.so_luong_con_lai === 'number' ? lo.so_luong_con_lai : lo.so_luong }} {{ lo.don_vi_tinh }})
            </option>
          </select>
          <input v-model="form.ma_npp" placeholder="Mã NPP (VD: NPP001) *" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70" />
          <input v-model="form.san_pham" placeholder="Tên sản phẩm (VD: Táo Ninh Thuận) *" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70" />
          
          <div class="relative">
            <input v-model="form.so_luong" placeholder="Số lượng *" type="number" step="0.1" :class="isOverInventory ? 'border-red-400 focus:ring-red-200' : 'border-gray-200 focus:border-amber-400 focus:ring-amber-200'" class="px-4 py-3 rounded-xl border focus:ring-2 outline-none transition-all bg-white/70 w-full" />
            <p v-if="isOverInventory" class="text-xs text-red-500 mt-1 absolute -bottom-5 left-1">⚠️ Vượt quá Tồn kho ({{ selectedLohangInventory }})</p>
          </div>

          <select v-model="form.don_vi_tinh" class="px-4 py-3 mt-1 md:mt-0 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70">
            <option value="kg">kg</option>
            <option value="tan">tấn</option>
            <option value="thung">thùng</option>
          </select>
          <input v-model="form.don_gia" placeholder="Đơn giá (VND) *" type="number" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70" />
          <input v-model="form.ty_le_hoa_hong" placeholder="Tỷ lệ hoa hồng (%) *" type="number" step="0.5" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70" />
          <input v-model="form.ghi_chu" placeholder="Ghi chú" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-200 outline-none transition-all bg-white/70" />
        </div>
        <div class="flex gap-3 mt-4">
          <button @click="taoGiaoDich" :disabled="submitting || !form.ma_lo || !form.ma_npp || !form.san_pham || !form.so_luong || !form.don_gia || isOverInventory"
            class="px-5 py-2.5 bg-amber-600 text-white rounded-xl font-medium hover:bg-amber-700 transition-all disabled:opacity-40 disabled:cursor-not-allowed">
            {{ submitting ? 'Đang gửi...' : 'Gửi yêu cầu' }}
          </button>
          <button @click="showForm = false" class="px-5 py-2.5 bg-gray-200 text-gray-700 rounded-xl font-medium hover:bg-gray-300 transition-all">Hủy</button>
        </div>
      </div>
    </transition>

    <!-- List -->
    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <div v-else-if="giaoDichs.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Chưa có giao dịch nào.
    </div>

    <div v-else class="space-y-4">
      <div v-for="gd in giaoDichs" :key="gd.ma_giao_dich || gd.id" class="glass rounded-xl p-5 border border-white/40 hover:shadow-md transition-all">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h4 class="font-semibold text-gray-800 text-lg">{{ gd.san_pham || 'Sản phẩm' }}</h4>
              <span class="text-xs px-2.5 py-1 rounded-full font-medium" :class="formatStatus(gd.trang_thai).class">
                {{ formatStatus(gd.trang_thai).label }}
              </span>
            </div>
            
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mt-3">
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Mã GD / Lô</p>
                <p class="text-xs font-mono text-gray-600 mt-1 truncate" :title="gd.ma_giao_dich">{{ gd.ma_giao_dich }}</p>
                <p class="text-xs font-mono text-gray-500 truncate" :title="gd.ma_lo">{{ gd.ma_lo }}</p>
                <p class="text-xs text-gray-400 mt-1 italic">{{ gd.ngay_tao ? new Date(gd.ngay_tao).toLocaleDateString('vi-VN') : 'Chưa ghi nhận' }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Bên mua (NPP)</p>
                <div class="flex items-center gap-1.5 mt-1">
                  <span class="text-sm">🏬</span>
                  <p class="text-sm font-medium text-gray-700">{{ gd.ma_npp }}</p>
                </div>
              </div>
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Sản lượng</p>
                <p class="text-sm font-medium text-gray-700 mt-1">{{ gd.so_luong }} <span class="text-xs text-gray-500">{{ gd.don_vi_tinh }}</span></p>
              </div>
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Đơn giá</p>
                <p class="text-sm font-medium text-gray-700 mt-1">{{ gd.don_gia?.toLocaleString() }} <span class="text-xs text-gray-500">VND</span></p>
              </div>
            </div>

            <div v-if="gd.ghi_chu" class="mt-4 pt-3 border-t border-gray-100 flex items-start gap-2 text-sm text-gray-500">
              <span>📌</span>
              <p class="italic">{{ gd.ghi_chu }}</p>
            </div>
          </div>

          <div class="text-right ml-4 pl-4 border-l border-gray-100 flex flex-col justify-center min-w-[140px]">
            <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold mb-1">Thành tiền</p>
            <p class="text-xl font-bold text-gray-800">{{ (gd.so_luong * gd.don_gia).toLocaleString() }}</p>
            <p class="text-xs font-medium text-gray-500">VND</p>
            
            <div class="mt-3 pt-3 border-t border-gray-100">
              <p class="text-xs text-gray-400">Hoa hồng NPP</p>
              <p class="text-sm font-semibold text-emerald-600 mt-0.5">{{ gd.ty_le_hoa_hong }}%</p>
            </div>
            
            <div class="mt-3" v-if="['cho_duyet', 'da_duyet', 'dang_giao', 'da_giao', 'cho_thanh_toan'].includes(gd.trang_thai)">
               <span class="text-xs text-amber-600 italic">Đang chờ Platform/NPP xử lý...</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
