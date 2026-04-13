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
    // Platform MVP: Fetch lots from HTX001 as there is no GetAll API yet
    const res = await api.get(`/api/v1/lohang/htx/HTX001`)
    const data = res.data?.data || res.data || []
    lohangs.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function capNhatTrangThai(id, trangThaiPhucHoi) {
  try {
    await api.put(`/api/v1/lohang/${id}/trangthai`, { trang_thai: trangThaiPhucHoi })
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái lô hàng')
  }
}

async function handlePhucHoi(id) {
  const rs = prompt("Nhập trạng thái muốn phục hồi (ví dụ: san_sang_ban, dang_trong):", "san_sang_ban")
  if (rs) {
    await capNhatTrangThai(id, rs)
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
      <h2 class="text-2xl font-bold text-gray-800">Quản lý Lô Hàng toàn hệ thống</h2>
    </div>

    <!-- Alert for MVP Data Source -->
    <div class="mb-6 p-4 rounded-xl bg-purple-50 text-purple-700 border border-purple-200 text-sm flex items-start gap-3">
      <span class="text-xl">ℹ️</span>
      <div>
        <p class="font-semibold">Lưu ý mô hình hoạt động (MVP)</p>
        <p>Phiên bản MVP hiển thị mock/liên kết từ dữ liệu của "HTX001". Platform có quyền giám sát, đình chỉ đối với các lô hàng vi phạm chất lượng.</p>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải biểu dữ liệu lô hàng...</div>

    <div v-else-if="lohangs.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Cơ sở dữ liệu Blockchain hiện không có lô hàng nào.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div v-for="lo in lohangs" :key="lo.ma_lo" class="glass rounded-2xl p-5 border border-white/40 shadow-sm relative overflow-hidden flex flex-col h-full">
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
            <span class="text-gray-400">Sản lượng</span>
            <span class="font-medium text-gray-700">{{ lo.so_luong }} {{ lo.don_vi_tinh }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Thời vụ</span>
            <span class="text-gray-700">{{ lo.vu_mua || '-' }}</span>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-gray-100 flex gap-2 justify-end">
          <button v-if="lo.trang_thai !== 'dinh_chi' && lo.trang_thai !== 'het_hang'"
            @click="capNhatTrangThai(lo.ma_lo, 'dinh_chi')"
            class="px-4 py-2 bg-red-50 text-red-600 rounded-lg text-sm font-medium hover:bg-red-100 border border-red-200 transition-all flex items-center gap-1">
            <span>⛔</span> Đình chỉ ngay
          </button>
          <button v-if="lo.trang_thai === 'dinh_chi'"
            @click="handlePhucHoi(lo.ma_lo)"
            class="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg text-sm font-medium hover:bg-blue-100 border border-blue-200 transition-all flex items-center gap-1">
            <span>🛠️</span> Phục hồi
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
