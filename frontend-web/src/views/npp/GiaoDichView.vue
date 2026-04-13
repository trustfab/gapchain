<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const giaoDichs = ref([])
const loading = ref(true)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/giaodich/npp/${auth.tenantId}`)
    const data = res.data?.data || res.data || []
    giaoDichs.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function tiepNhanHang(id) {
  try {
    await api.put(`/api/v1/giaodich/${id}/trangthai`, { trang_thai: 'da_giao' })
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái')
  }
}

async function xacNhanThanhToan(id) {
  try {
    await api.put(`/api/v1/giaodich/${id}/trangthai`, { trang_thai: 'da_thanh_toan' })
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái')
  }
}

const statusColor = (s) => {
  const map = {
    cho_duyet: 'bg-yellow-100 text-yellow-700 border-yellow-200',
    da_duyet: 'bg-blue-100 text-blue-700 border-blue-200',
    dang_giao: 'bg-purple-100 text-purple-700 border-purple-200',
    da_giao: 'bg-indigo-100 text-indigo-700 border-indigo-200',
    cho_thanh_toan: 'bg-orange-100 text-orange-700 border-orange-200',
    da_thanh_toan: 'bg-emerald-100 text-emerald-700 border-emerald-200',
    huy_bo: 'bg-gray-100 text-gray-700 border-gray-200',
  }
  return map[s] || 'bg-gray-100 text-gray-700 border-gray-200'
}

const statusText = (s) => {
  const map = {
    cho_duyet: 'Chờ duyệt',
    da_duyet: 'Đã duyệt',
    dang_giao: 'Đang giao',
    da_giao: 'Đã nhận hàng',
    cho_thanh_toan: 'Chờ thanh toán',
    da_thanh_toan: 'Đã thanh toán',
    huy_bo: 'Đã hủy',
  }
  return map[s] || s
}
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Quản lý Đơn hàng</h2>

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
              <span class="text-xs px-2.5 py-1.5 rounded-full border shadow-sm font-medium" :class="statusColor(gd.trang_thai)">{{ statusText(gd.trang_thai) }}</span>
            </div>
            
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mt-3">
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Mã GD / Lô</p>
                <p class="text-xs font-mono text-gray-600 mt-1 truncate" :title="gd.ma_giao_dich">{{ gd.ma_giao_dich || gd.id }}</p>
                <p class="text-xs font-mono text-gray-500 truncate" :title="gd.ma_lo">{{ gd.ma_lo }}</p>
                <p class="text-xs text-gray-400 mt-1 italic">{{ gd.ngay_tao ? new Date(gd.ngay_tao).toLocaleDateString('vi-VN') : 'Chưa ghi nhận' }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Bên bán (HTX)</p>
                <div class="flex items-center gap-1.5 mt-1">
                  <span class="text-sm">🏠</span>
                  <p class="text-sm font-medium text-gray-700">{{ gd.ma_htx }}</p>
                </div>
              </div>
              <div>
                <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold">Sản lượng</p>
                <p class="text-sm font-medium text-gray-700 mt-1">{{ gd.so_luong }} <span class="text-xs text-gray-500">{{ gd.don_vi_tinh || 'kg' }}</span></p>
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

          <div class="text-right ml-4 pl-4 border-l border-gray-100 flex flex-col justify-center min-w-[140px] gap-4">
            <div>
              <p class="text-[11px] text-gray-400 uppercase tracking-wider font-semibold mb-1">Thành tiền</p>
              <p class="text-xl font-bold text-gray-800 text-green-700">{{ (gd.so_luong * gd.don_gia).toLocaleString() }}</p>
              <p class="text-xs font-medium text-gray-500">VND</p>
            </div>
            
            <div class="flex flex-col gap-2 mt-2">
              <button v-if="gd.trang_thai === 'dang_giao'"
                @click="tiepNhanHang(gd.ma_giao_dich || gd.id)"
                class="px-4 py-2 bg-indigo-500 text-white rounded-lg text-sm font-medium hover:bg-indigo-600 transition-all shadow-sm w-full text-center border border-indigo-600">
                Nhận hàng
              </button>
              <button v-if="gd.trang_thai === 'cho_thanh_toan'"
                @click="xacNhanThanhToan(gd.ma_giao_dich || gd.id)"
                class="px-4 py-2 bg-emerald-500 text-white rounded-lg text-sm font-medium hover:bg-emerald-600 transition-all shadow-sm w-full text-center border border-emerald-600">
                Đã thanh toán
              </button>
              <span v-if="['cho_duyet', 'da_duyet', 'da_giao'].includes(gd.trang_thai)" class="text-[11px] text-center text-amber-600 mt-1 italic">
                Đợi Platform điều phối...
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
