<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const congNo = ref(null)
const loading = ref(true)

async function fetchData() {
  loading.value = true
  try {
    const [cnRes, hhRes] = await Promise.allSettled([
      api.get(`/api/v1/giaodich/npp/${auth.tenantId}/congno`),
      api.get(`/api/v1/giaodich/npp/${auth.tenantId}/hoahong`)
    ])
    
    let list = []
    let totalDaThanhToan = 0

    if (cnRes.status === 'fulfilled') {
      const cnData = cnRes.value.data?.data || cnRes.value.data || []
      list = Array.isArray(cnData) ? cnData : []
    }
    
    if (hhRes.status === 'fulfilled') {
      const hhData = hhRes.value.data?.data || hhRes.value.data || {}
      totalDaThanhToan = hhData.tong_doanh_thu || 0
    } else {
      console.error("Lỗi khi load Hoa Hồng:", hhRes.reason)
    }

    const totalConLai = list.reduce((acc, curr) => acc + (curr.so_luong * curr.don_gia), 0)

    congNo.value = {
      tong_cong_no: totalConLai + totalDaThanhToan,
      da_thanh_toan: totalDaThanhToan,
      con_lai: totalConLai,
      giao_dich: list
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchData()
})

async function xacNhanThanhToan(id) {
  try {
    await api.put(`/api/v1/giaodich/${id}/trangthai`, { trang_thai: 'da_thanh_toan' })
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái')
  }
}
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Công Nợ</h2>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <template v-else>
      <!-- Summary -->
      <div class="glass rounded-2xl p-6 border border-white/40 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="text-center">
            <p class="text-sm text-gray-500 mb-1">Tổng công nợ</p>
            <p class="text-3xl font-bold text-red-600">{{ congNo?.tong_cong_no?.toLocaleString() || '0' }} VND</p>
          </div>
          <div class="text-center">
            <p class="text-sm text-gray-500 mb-1">Đã thanh toán</p>
            <p class="text-3xl font-bold text-green-600">{{ congNo?.da_thanh_toan?.toLocaleString() || '0' }} VND</p>
          </div>
          <div class="text-center">
            <p class="text-sm text-gray-500 mb-1">Còn lại</p>
            <p class="text-3xl font-bold text-gray-800">{{ congNo?.con_lai?.toLocaleString() || '0' }} VND</p>
          </div>
        </div>
      </div>

      <!-- Detail table -->
      <div class="glass rounded-2xl p-6 border border-white/40">
        <h3 class="font-semibold text-gray-800 mb-4">Chi tiết giao dịch</h3>

        <div v-if="!congNo?.giao_dich?.length" class="text-center py-8 text-gray-400">
          Không có dữ liệu công nợ.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200">
                <th class="text-left py-3 px-2 text-gray-500 font-medium">Mã GD</th>
                <th class="text-left py-3 px-2 text-gray-500 font-medium">Bên bán (HTX)</th>
                <th class="text-left py-3 px-2 text-gray-500 font-medium">Sản phẩm</th>
                <th class="text-right py-3 px-2 text-gray-500 font-medium">Thành tiền</th>
                <th class="text-right py-3 px-2 text-gray-500 font-medium">Hoa hồng NPP</th>
                <th class="text-center py-3 px-2 text-gray-500 font-medium">Trạng thái</th>
                <th class="text-center py-3 px-2 text-gray-500 font-medium w-36">Xử lý</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="gd in congNo.giao_dich" :key="gd.ma_giao_dich" class="border-b border-gray-100 hover:bg-white/30">
                <td class="py-3 px-2 font-mono text-xs">{{ gd.ma_giao_dich || gd.id }}</td>
                <td class="py-3 px-2 font-medium">{{ gd.ma_htx }}</td>
                <td class="py-3 px-2">
                  <p class="font-medium text-gray-800">{{ gd.san_pham }}</p>
                  <p class="text-xs text-gray-500">{{ gd.so_luong }} {{ gd.don_vi_tinh }}</p>
                </td>
                <td class="py-3 px-2 text-right font-bold text-gray-700">{{ (gd.so_luong * gd.don_gia).toLocaleString() }} đ</td>
                <td class="py-3 px-2 text-right">
                  <p class="text-emerald-600 font-medium">{{ (gd.so_luong * gd.don_gia * gd.ty_le_hoa_hong / 100).toLocaleString() }} đ</p>
                  <p class="text-xs text-gray-400">({{ gd.ty_le_hoa_hong }}%)</p>
                </td>
                <td class="py-3 px-2 text-center">
                  <span class="text-xs px-2.5 py-1 rounded-full font-medium whitespace-nowrap"
                    :class="gd.trang_thai === 'da_thanh_toan' ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'">
                    {{ gd.trang_thai === 'cho_thanh_toan' ? 'Chờ thanh toán' : gd.trang_thai }}
                  </span>
                </td>
                <td class="py-3 px-2 text-center">
                  <button v-if="gd.trang_thai === 'cho_thanh_toan'"
                    @click="xacNhanThanhToan(gd.ma_giao_dich || gd.id)"
                    class="px-3 py-1.5 bg-emerald-50 text-emerald-600 border border-emerald-200 rounded-lg text-xs font-medium hover:bg-emerald-100 transition-all whitespace-nowrap shadow-sm hover:shadow">
                    Đã chuyển khoản
                  </button>
                  <span v-else class="text-xs text-gray-400 text-center block">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>
