<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const stats = ref({ lohang: 0, nhatky: 0, giaodich: 0 })
const recentLohang = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const [loRes, nkRes, gdRes] = await Promise.allSettled([
      api.get(`/api/v1/lohang/htx/${auth.tenantId}`),
      api.get(`/api/v1/nhatky/htx/${auth.tenantId}`),
      api.get(`/api/v1/giaodich/htx/${auth.tenantId}`)
    ])
    if (loRes.status === 'fulfilled') {
      const data = loRes.value.data?.data || loRes.value.data || []
      recentLohang.value = Array.isArray(data) ? data.slice(0, 5) : []
      stats.value.lohang = Array.isArray(data) ? data.length : 0
    }
    if (nkRes.status === 'fulfilled') {
      const data = nkRes.value.data?.data || nkRes.value.data || []
      stats.value.nhatky = Array.isArray(data) ? data.length : 0
    }
    if (gdRes.status === 'fulfilled') {
      const data = gdRes.value.data?.data || gdRes.value.data || []
      stats.value.giaodich = Array.isArray(data) ? data.length : 0
    }
  } catch (e) {
    console.error('Dashboard load error:', e)
  } finally {
    loading.value = false
  }
})

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
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Tổng quan HTX</h2>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5 mb-8">
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-green-400 to-emerald-600 flex items-center justify-center text-2xl">📦</div>
          <div>
            <p class="text-sm text-gray-500">Lô hàng</p>
            <p class="text-3xl font-bold text-gray-800">{{ stats.lohang }}</p>
          </div>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-400 to-indigo-600 flex items-center justify-center text-2xl">📝</div>
          <div>
            <p class="text-sm text-gray-500">Nhật ký</p>
            <p class="text-3xl font-bold text-gray-800">{{ stats.nhatky }}</p>
          </div>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-amber-400 to-orange-600 flex items-center justify-center text-2xl">💰</div>
          <div>
            <p class="text-sm text-gray-500">Giao dịch</p>
            <p class="text-3xl font-bold text-gray-800">{{ stats.giaodich }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Lo Hang -->
    <div class="glass rounded-2xl p-6 border border-white/40">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-800">Lô hàng gần đây</h3>
        <router-link to="/htx/lohang" class="text-sm text-green-600 hover:text-green-700 font-medium">Xem tất cả</router-link>
      </div>

      <div v-if="loading" class="text-center py-8 text-gray-400">Đang tải...</div>

      <div v-else-if="recentLohang.length === 0" class="text-center py-8 text-gray-400">
        Chưa có lô hàng nào. Bấm "Xem tất cả" để tạo mới.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <router-link
          v-for="lo in recentLohang"
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
  </div>
</template>
