<script setup>
import { ref, onMounted } from 'vue'
import api from '@/plugins/axios'

const giaoDichs = ref([])
const loading = ref(true)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    // MVP: Lấy dữ liệu của HTX001 làm mẫu do Platform chưa có endpoint GetAll
    const res = await api.get(`/api/v1/giaodich/htx/HTX001`)
    const data = res.data?.data || res.data || []
    giaoDichs.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function thucHienHanhDong(id, action, trangThaiMoi) {
  try {
    if (action === 'duyet') {
      await api.put(`/api/v1/giaodich/${id}/duyet`)
    } else {
      await api.put(`/api/v1/giaodich/${id}/trangthai`, { trang_thai: trangThaiMoi })
    }
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || `Lỗi cập nhật trạng thái đơn ${id}`)
  }
}

async function xacNhanHuy(id) {
  if (confirm(`Bạn có chắc chắn muốn HỦY BỎ giao dịch ${id} không?`)) {
    await thucHienHanhDong(id, 'trangthai', 'huy_bo')
  }
}

const statusBadge = (s) => {
  const map = {
    cho_duyet: { class: 'bg-yellow-100 text-yellow-700 border-yellow-200', text: 'Chờ duyệt', icon: '⏳' },
    da_duyet: { class: 'bg-blue-100 text-blue-700 border-blue-200', text: 'Đã duyệt', icon: '📝' },
    dang_giao: { class: 'bg-orange-100 text-orange-700 border-orange-200', text: 'Đang giao', icon: '🚚' },
    da_giao: { class: 'bg-indigo-100 text-indigo-700 border-indigo-200', text: 'Đã nhận hàng', icon: '📦' },
    cho_thanh_toan: { class: 'bg-red-100 text-red-700 border-red-200', text: 'Chờ thanh toán', icon: '💳' },
    da_thanh_toan: { class: 'bg-green-100 text-green-700 border-green-200', text: 'Đã thanh toán', icon: '✅' },
    huy_bo: { class: 'bg-gray-100 text-gray-700 border-gray-200', text: 'Đã huỷ', icon: '❌' },
  }
  return map[s] || { class: 'bg-gray-100 text-gray-700 border-gray-200', text: s }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Điều phối Giao Dịch toàn mạng</h2>
    </div>

    <!-- Alert for MVP Data Source -->
    <div class="mb-6 p-4 rounded-xl bg-purple-50 text-purple-700 border border-purple-200 text-sm flex items-start gap-3">
      <span class="text-xl">ℹ️</span>
      <div>
        <p class="font-semibold">Bộ điều phối uỷ quyền (Platform Role)</p>
        <p>Hệ thống hỗ trợ thao tác chuyển giai đoạn (Duyệt, Giao hàng, Ấn định Công nợ) để phục vụ việc demo nhanh vòng đời giao dịch.</p>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang đồng bộ giao dịch từ sổ cái...</div>

    <div v-else-if="giaoDichs.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Không có luồng giao dịch điện tử nào.
    </div>

    <div v-else class="space-y-4">
      <div v-for="gd in giaoDichs" :key="gd.ma_giao_dich || gd.id" class="glass rounded-2xl p-6 border border-white/40 shadow-sm relative hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <span class="text-lg font-bold text-gray-800 font-mono">{{ gd.ma_giao_dich || gd.id }}</span>
              <span class="text-xs px-2.5 py-1 rounded-full border shadow-sm flex items-center gap-1 font-medium" :class="statusBadge(gd.trang_thai).class">
                <span>{{ statusBadge(gd.trang_thai).icon }}</span>
                {{ statusBadge(gd.trang_thai).text }}
              </span>
            </div>
            
            <div class="text-sm text-gray-500 flex flex-wrap items-center gap-4 mt-3">
              <p class="flex items-center gap-1.5"><span class="bg-gray-100 p-1 rounded">📅</span> <span class="font-medium text-gray-700">{{ gd.ngay_tao ? new Date(gd.ngay_tao).toLocaleString('vi-VN') : 'Chưa ghi nhận' }}</span></p>
              <p class="flex items-center gap-1.5"><span class="bg-gray-100 p-1 rounded">🏠</span> HTX: <span class="font-medium text-gray-700">{{ gd.ma_htx }}</span></p>
              <p class="flex items-center gap-1.5"><span class="bg-gray-100 p-1 rounded">🏭</span> NPP: <span class="font-medium text-gray-700">{{ gd.ma_npp }}</span></p>
              <p class="flex items-center gap-1.5"><span class="bg-gray-100 p-1 rounded">📦</span> Lô: <span class="font-mono text-gray-700">{{ gd.ma_lo }}</span></p>
            </div>
            
            <div class="mt-4 p-3 bg-white/50 border border-gray-100 rounded-xl inline-block min-w-[300px]">
              <div class="flex justify-between items-center text-gray-700">
                <span class="text-xs text-gray-500 uppercase tracking-wider">Khối lượng</span>
                <span class="font-medium">{{ gd.so_luong }} {{ gd.don_vi_tinh || 'kg' }}</span>
              </div>
              <div class="flex justify-between items-center text-gray-700 mt-2">
                <span class="text-xs text-gray-500 uppercase tracking-wider">Đơn giá</span>
                <span class="font-medium">{{ gd.don_gia?.toLocaleString() }} ₫</span>
              </div>
              <div class="flex justify-between items-center text-gray-800 mt-3 pt-2 border-t border-gray-200">
                <span class="font-bold text-green-700">Giá trị hợp đồng</span>
                <span class="font-bold text-lg text-green-700">{{ (gd.so_luong * gd.don_gia).toLocaleString() }} ₫</span>
              </div>
            </div>
          </div>

          <div class="ml-8 flex flex-col gap-2 min-w-[140px]">
            <button v-if="gd.trang_thai === 'cho_duyet'"
              @click="thucHienHanhDong(gd.ma_giao_dich || gd.id, 'duyet')"
              class="px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-xl text-sm font-medium hover:shadow-lg transition-all text-center">
              Duyệt Giao Dịch
            </button>
            <button v-if="gd.trang_thai === 'da_duyet'"
              @click="thucHienHanhDong(gd.ma_giao_dich || gd.id, 'trangthai', 'dang_giao')"
              class="px-4 py-2 bg-gradient-to-r from-orange-400 to-orange-500 text-white rounded-xl text-sm font-medium hover:shadow-lg transition-all text-center">
              Xác nhận Giao
            </button>
            <button v-if="gd.trang_thai === 'da_giao' || gd.trang_thai === 'da_nhan_hang'"
              @click="thucHienHanhDong(gd.ma_giao_dich || gd.id, 'trangthai', 'cho_thanh_toan')"
              class="px-4 py-2 bg-gradient-to-r from-red-400 to-red-500 text-white rounded-xl text-sm font-medium hover:shadow-lg transition-all text-center">
              Chốt Công Nợ
            </button>
            
            <button v-if="gd.trang_thai !== 'da_thanh_toan' && gd.trang_thai !== 'huy_bo'"
              @click="xacNhanHuy(gd.ma_giao_dich || gd.id)"
              class="px-4 py-2 bg-gray-100 text-gray-600 border border-gray-200 rounded-xl text-sm font-medium hover:bg-gray-200 hover:text-red-500 transition-all text-center mt-auto">
              Hủy hợp đồng
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
