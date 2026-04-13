<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()

const nhatkys = ref([])
const loading = ref(true)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    // Fetch all nhat ky from a known HTX for BVTV to review
    const res = await api.get('/api/v1/nhatky/htx/HTX001')
    const data = res.data?.data || res.data || []
    nhatkys.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function duyetNhatKy(id, approved) {
  const ly_do_tu_choi = approved ? '' : prompt('Lý do từ chối:')
  if (!approved && !ly_do_tu_choi) return

  try {
    await api.put(`/api/v1/nhatky/${id}/duyet`, {
      quyet_dinh: approved ? 'da_duyet' : 'tu_choi',
      ly_do_tu_choi: ly_do_tu_choi,
      nguoi_duyet: auth.username || 'bvtv',
    })
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi duyệt nhật ký')
  }
}

const statusStyle = (s) => {
  const map = {
    da_duyet: 'bg-green-100 text-green-700',
    tu_choi: 'bg-red-100 text-red-700',
    cho_duyet: 'bg-yellow-100 text-yellow-700',
  }
  return map[s] || 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Duyệt Nhật Ký canh tác</h2>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <div v-else-if="nhatkys.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Không có nhật ký nào cần duyệt.
    </div>

    <div v-else class="space-y-4">
      <div v-for="nk in nhatkys" :key="nk.ma_nhat_ky" class="glass rounded-2xl p-5 border border-white/40">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <span class="font-semibold text-gray-800">{{ nk.loai_hoat_dong }}</span>
              <span class="text-xs font-mono text-gray-400 bg-gray-100 px-2 py-0.5 rounded">{{ nk.ma_lo }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full" :class="statusStyle(nk.trang_thai)">
                {{ nk.trang_thai || 'cho_duyet' }}
              </span>
            </div>
            <p class="text-sm text-gray-600">{{ nk.chi_tiet || nk.mo_ta || nk.noi_dung }}</p>
            <p class="text-xs text-gray-400 mt-2">HTX: {{ nk.ma_htx }} | Ngày: {{ nk.ngay || nk.ngay_ghi }}</p>
          </div>

          <div v-if="nk.trang_thai === 'cho_duyet'" class="flex flex-col gap-2 ml-4">
            <button @click="duyetNhatKy(nk.ma_nhat_ky, true)"
              class="px-4 py-2 bg-green-500 text-white rounded-lg text-sm font-medium hover:bg-green-600 transition-all text-center">
              Phê duyệt
            </button>
            <button @click="duyetNhatKy(nk.ma_nhat_ky, false)"
              class="px-4 py-2 bg-red-500 text-white rounded-lg text-sm font-medium hover:bg-red-600 transition-all text-center">
              Bác bỏ
            </button>
            <router-link :to="`/bvtv/lohang/${nk.ma_lo}`" class="mt-2 text-xs text-indigo-600 hover:text-indigo-800 text-center font-medium">
              Xem Lô Hàng &rarr;
            </router-link>
          </div>
          <div v-else class="flex flex-col items-end gap-2 ml-4">
            <span class="text-sm font-medium text-gray-500">Người duyệt: {{ nk.nguoi_duyet || 'Không rõ' }}</span>
            <router-link :to="`/bvtv/lohang/${nk.ma_lo}`" class="mt-1 text-xs text-indigo-600 hover:text-indigo-800 font-medium">
              Vào Lô Hàng &rarr;
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
