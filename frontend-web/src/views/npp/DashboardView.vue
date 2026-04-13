<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const congNo = ref(null)
const hoaHong = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const [cnRes, hhRes] = await Promise.allSettled([
      api.get(`/api/v1/giaodich/npp/${auth.tenantId}/congno`),
      api.get(`/api/v1/giaodich/npp/${auth.tenantId}/hoahong`),
    ])
    if (cnRes.status === 'fulfilled') {
      const cnData = cnRes.value.data?.data || cnRes.value.data || []
      const list = Array.isArray(cnData) ? cnData : []
      congNo.value = {
        tong_cong_no: list.reduce((acc, curr) => acc + (curr.so_luong * curr.don_gia), 0)
      }
    }
    if (hhRes.status === 'fulfilled') {
      hoaHong.value = hhRes.value.data?.data || hhRes.value.data || {}
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-6">Tổng quan Nhà Phân Phối</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-red-400 to-rose-600 flex items-center justify-center text-2xl">💳</div>
          <div>
            <p class="text-sm text-gray-500">Công nợ</p>
            <p class="text-2xl font-bold text-gray-800">{{ congNo?.tong_cong_no?.toLocaleString() || '0' }} VND</p>
          </div>
        </div>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-400 to-green-600 flex items-center justify-center text-2xl">💰</div>
          <div>
            <p class="text-sm text-gray-500">Hoa hồng</p>
            <p class="text-2xl font-bold text-gray-800">{{ hoaHong?.tong_tien_hoa_hong?.toLocaleString() || '0' }} VND</p>
          </div>
        </div>
      </div>
    </div>

    <div class="flex gap-3">
      <router-link to="/npp/giaodich"
        class="px-5 py-3 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-xl font-medium btn-glow">
        Quản lý Giao dịch
      </router-link>
      <router-link to="/npp/congno"
        class="px-5 py-3 bg-gradient-to-r from-purple-500 to-violet-600 text-white rounded-xl font-medium btn-glow">
        Xem Công Nợ
      </router-link>
    </div>
  </div>
</template>
