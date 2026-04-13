<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const lohangs = ref([])
const loading = ref(true)
const showForm = ref(false)

const form = ref({
  ten_san_pham: '',
  loai_san_pham: '',
  so_luong: '',
  don_vi_tinh: 'kg',
  vu_mua: '',
  dia_diem: '',
})
const submitting = ref(false)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/lohang/htx/${auth.tenantId}`)
    const data = res.data?.data || res.data || []
    lohangs.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function taoLoHang() {
  submitting.value = true
  try {
    const now = Date.now()
    await api.post('/api/v1/lohang', {
      ma_lo: `LH-${auth.tenantId}-${now}`,
      ma_htx: auth.tenantId,
      ten_san_pham: form.value.ten_san_pham,
      loai_san_pham: form.value.loai_san_pham,
      so_luong: parseFloat(form.value.so_luong) || 0,
      don_vi_tinh: form.value.don_vi_tinh,
      vu_mua: form.value.vu_mua,
      dia_diem: form.value.dia_diem,
    })
    showForm.value = false
    form.value = { ten_san_pham: '', loai_san_pham: '', so_luong: '', don_vi_tinh: 'kg', vu_mua: '', dia_diem: '' }
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Loi tao lo hang')
  } finally {
    submitting.value = false
  }
}

const formatStatus = (s) => {
  const map = {
    dang_trong: { label: '🌱 Đang trồng', class: 'bg-green-100 text-green-700' },
    da_thu_hoach: { label: '🌾 Đã thu hoạch', class: 'bg-yellow-100 text-yellow-700' },
    cho_chung_nhan: { label: '🔍 Chờ chứng nhận', class: 'bg-orange-100 text-orange-700' },
    san_sang_ban: { label: '✅ Sẵn sàng bán', class: 'bg-blue-100 text-blue-700' },
    het_hang: { label: '📦 Hết hàng', class: 'bg-gray-100 text-gray-700' },
    dinh_chi: { label: '⛔ Đình chỉ', class: 'bg-red-100 text-red-700' }
  }
  return map[s] || { label: s || 'N/A', class: 'bg-gray-100 text-gray-700' }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Quản lý Lô Hàng</h2>
      <button
        @click="showForm = !showForm"
        class="px-5 py-2.5 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl font-medium btn-glow"
      >
        + Tạo Lô Hàng
      </button>
    </div>

    <!-- Create Form -->
    <transition name="slide">
      <div v-if="showForm" class="glass rounded-2xl p-6 border border-white/40 mb-6">
        <h3 class="font-semibold text-gray-800 mb-4">Tạo Lô Hàng mới</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input v-model="form.ten_san_pham" placeholder="Tên sản phẩm (VD: Táo Ninh Thuận) *" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70" />
          <select v-model="form.loai_san_pham" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70">
            <option value="">-- Loại sản phẩm * --</option>
            <option value="rau_cu">Rau củ</option>
            <option value="trai_cay">Trái cây</option>
            <option value="lua_gao">Lúa gạo</option>
            <option value="thuy_san">Thủy sản</option>
            <option value="khac">Khác</option>
          </select>
          <input v-model="form.so_luong" placeholder="Số lượng *" type="number" step="0.1" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70" />
          <select v-model="form.don_vi_tinh" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70">
            <option value="kg">kg</option>
            <option value="tan">tấn</option>
            <option value="thung">thùng</option>
            <option value="bao">bao</option>
          </select>
          <input v-model="form.vu_mua" placeholder="Vụ mùa (VD: Đông Xuân 2025)" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70" />
          <input v-model="form.dia_diem" placeholder="Địa điểm (VD: Ninh Thuận)" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-green-400 focus:ring-2 focus:ring-green-200 outline-none transition-all bg-white/70" />
        </div>
        <div class="flex gap-3 mt-4">
          <button @click="taoLoHang" :disabled="submitting || !form.ten_san_pham || !form.loai_san_pham || !form.so_luong"
            class="px-5 py-2.5 bg-green-600 text-white rounded-xl font-medium hover:bg-green-700 transition-all disabled:opacity-40 disabled:cursor-not-allowed">
            {{ submitting ? 'Đang lưu...' : 'Lưu' }}
          </button>
          <button @click="showForm = false" class="px-5 py-2.5 bg-gray-200 text-gray-700 rounded-xl font-medium hover:bg-gray-300 transition-all">Hủy</button>
        </div>
      </div>
    </transition>

    <!-- List -->
    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <div v-else-if="lohangs.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Chưa có lô hàng nào.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <router-link
        v-for="lo in lohangs"
        :key="lo.ma_lo"
        :to="`/htx/lohang/${lo.ma_lo}`"
        class="glass rounded-2xl p-5 border border-white/40 hover:shadow-lg hover:scale-[1.02] transition-all duration-200 block"
      >
        <div class="flex items-start justify-between mb-3">
          <h3 class="font-semibold text-gray-800">{{ lo.ten_san_pham || 'Chưa đặt tên' }}</h3>
          <span class="text-xs px-2.5 py-1 rounded-full font-medium whitespace-nowrap" :class="formatStatus(lo.trang_thai).class">
            {{ formatStatus(lo.trang_thai).label }}
          </span>
        </div>
        <div class="space-y-1 text-sm text-gray-500">
          <p>Mã lô: <span class="font-mono text-gray-700">{{ lo.ma_lo }}</span></p>
          <p v-if="lo.loai_san_pham">Phân loại: <span class="capitalize text-gray-700">{{ lo.loai_san_pham?.replace('_', ' ') }}</span></p>
          <p v-if="lo.so_luong">Sản lượng / Tồn: <span class="font-medium" :class="lo.so_luong_con_lai === 0 ? 'text-red-500' : 'text-gray-700'">{{ typeof lo.so_luong_con_lai === 'number' ? lo.so_luong_con_lai : lo.so_luong }} / {{ lo.so_luong }}</span> {{ lo.don_vi_tinh }}</p>
        </div>
      </router-link>
    </div>
  </div>
</template>
