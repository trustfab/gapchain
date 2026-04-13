<script setup>
import { ref, onMounted } from 'vue'
import api from '@/plugins/axios'

const lohangs = ref([])
const loading = ref(true)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    // Fetch lots from HTX001 for MVP
    const res = await api.get(`/api/v1/lohang/htx/HTX001`)
    const data = res.data?.data || res.data || []
    lohangs.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function capNhatTrangThai(id, trangThaiMoi) {
  if (!confirm(`Bạn có chắc muốn chuyển trạng thái này không?`)) return
  try {
    await api.put(`/api/v1/lohang/${id}/trangthai`, { trang_thai: trangThaiMoi })
    alert('Cập nhật trạng thái thành công!')
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái lô hàng')
  }
}

const getLohangBadge = (s) => {
  const map = {
    dang_trong: { color: 'bg-green-100 text-green-700', text: 'Đang trồng', icon: '🌱' },
    da_thu_hoach: { color: 'bg-yellow-100 text-yellow-700', text: 'Đã thu hoạch', icon: '🌾' },
    cho_chung_nhan: { color: 'bg-orange-100 text-orange-700', text: 'Chờ chứng nhận', icon: '🔍' },
    san_sang_ban: { color: 'bg-blue-100 text-blue-700', text: 'Sẵn sàng bán', icon: '✅' },
    het_hang: { color: 'bg-gray-100 text-gray-700', text: 'Hết hàng', icon: '📦' },
    dinh_chi: { color: 'bg-red-100 text-red-700', text: 'Đình chỉ', icon: '⛔' }
  }
  return map[s] || { color: 'bg-gray-100 text-gray-700', text: s || 'N/A', icon: '📌' }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Kiểm định & Cấp Chứng Nhận Lô Hàng</h2>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải biểu dữ liệu lô hàng...</div>

    <div v-else-if="lohangs.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Cơ sở dữ liệu Blockchain hiện không có lô hàng nào.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div v-for="lo in lohangs" :key="lo.ma_lo" class="glass rounded-2xl p-5 border border-white/40 shadow-sm relative overflow-hidden flex flex-col h-full"
           :class="lo.trang_thai === 'cho_chung_nhan' ? 'ring-2 ring-orange-400/50' : ''">
        <div class="flex items-start justify-between mb-3 border-b border-gray-100 pb-3">
          <div>
            <h3 class="font-semibold text-gray-800 text-lg">{{ lo.ten_san_pham || 'Không rõ tên' }}</h3>
            <p class="text-xs font-mono text-gray-500 mt-0.5">Mã: {{ lo.ma_lo }}</p>
          </div>
          <span class="text-xs px-2.5 py-1 rounded-full font-medium whitespace-nowrap flex items-center gap-1 shadow-sm" :class="getLohangBadge(lo.trang_thai).color">
            <span>{{ getLohangBadge(lo.trang_thai).icon }}</span>
            <span>{{ getLohangBadge(lo.trang_thai).text }}</span>
          </span>
        </div>
        
        <div class="space-y-2 text-sm text-gray-600 flex-1">
          <div class="flex justify-between">
            <span class="text-gray-400">Hợp tác xã</span>
            <span class="font-medium text-gray-700">{{ lo.ma_htx }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Số lượng tồn</span>
            <span class="font-medium text-gray-700">{{ lo.so_luong_con_lai || lo.so_luong }} {{ lo.don_vi_tinh }}</span>
          </div>
          <div v-if="lo.chung_nhan && lo.chung_nhan.length" class="mt-2 pt-2 border-t border-gray-100">
            <span class="text-xs text-gray-400 block mb-1">Chứng nhận đã cấp:</span>
            <div class="flex flex-wrap gap-1">
              <span v-for="cn in lo.chung_nhan" :key="cn.ma_chung_nhan" class="text-[10px] px-2 py-0.5 bg-green-100 text-green-700 rounded-sm font-medium">
                {{ cn.loai_chung_nhan }} ({{ cn.ma_chung_nhan }})
              </span>
            </div>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-gray-100 flex gap-2 justify-end">
          <router-link :to="`/bvtv/lohang/${lo.ma_lo}`"
            class="px-4 py-2 bg-gradient-to-r from-orange-400 to-amber-500 text-white rounded-lg text-sm font-medium hover:opacity-90 transition-all flex items-center gap-1 flex-1 justify-center shadow-md">
            <span>🔍</span> Xem & Kiểm định
          </router-link>
          <button v-if="lo.trang_thai !== 'dinh_chi' && lo.trang_thai !== 'het_hang'"
            @click="capNhatTrangThai(lo.ma_lo, 'dinh_chi')"
            class="px-3 py-2 bg-red-50 text-red-600 rounded-lg text-sm font-medium hover:bg-red-100 border border-red-200 transition-all flex items-center gap-1">
            <span>⛔</span> Đình chỉ ngay
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
