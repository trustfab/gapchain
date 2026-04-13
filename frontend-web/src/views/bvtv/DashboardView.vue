<script setup>
import { ref, onMounted } from 'vue'
import api from '@/plugins/axios'

const stats = ref({ cho_duyet: 0, da_duyet: 0, tu_choi: 0 })
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/nhatky/thongke')
    const data = res.data?.data || res.data || {}
    stats.value.cho_duyet = data.theo_trang_thai?.cho_duyet || 0
    stats.value.da_duyet = data.theo_trang_thai?.da_duyet || 0
    stats.value.tu_choi = data.theo_trang_thai?.tu_choi || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Tổng quan BVTV</h2>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-5 mb-8">
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-yellow-400 to-amber-600 flex items-center justify-center text-2xl">⏳</div>
          <div>
            <p class="text-sm text-gray-500">Chờ duyệt</p>
            <p class="text-3xl font-bold text-amber-600">{{ stats.cho_duyet }}</p>
          </div>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-green-400 to-emerald-600 flex items-center justify-center text-2xl">✅</div>
          <div>
            <p class="text-sm text-gray-500">Đã duyệt</p>
            <p class="text-3xl font-bold text-green-600">{{ stats.da_duyet }}</p>
          </div>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-red-400 to-rose-600 flex items-center justify-center text-2xl">❌</div>
          <div>
            <p class="text-sm text-gray-500">Từ chối</p>
            <p class="text-3xl font-bold text-red-600">{{ stats.tu_choi }}</p>
          </div>
        </div>
      </div>
    </div>

    <router-link to="/bvtv/duyet"
      class="inline-flex items-center gap-2 px-5 py-3 bg-gradient-to-r from-amber-500 to-orange-600 text-white rounded-xl font-medium btn-glow">
      Xem danh sách duyệt &rarr;
    </router-link>
  </div>
</template>
