<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const loadingKey = ref('')
const error = ref('')

const roles = [
  {
    key: 'htx',
    label: 'Hợp tác xã',
    desc: 'Quản lý lô hàng, nhật ký canh tác',
    username: 'htx001',
    icon: '🧑‍🌾',
    gradient: 'from-green-400 to-emerald-600',
    bg: 'bg-green-50',
  },
  {
    key: 'platform',
    label: 'Platform Admin',
    desc: 'Giám sát toàn bộ hệ sinh thái',
    username: 'platform',
    icon: '🌐',
    gradient: 'from-purple-400 to-violet-600',
    bg: 'bg-purple-50',
  },
  {
    key: 'npp',
    label: 'Nhà phân phối Xanh',
    desc: 'Đơn hàng, công nợ, hoa hồng',
    username: 'npp001',
    icon: '🚚',
    gradient: 'from-blue-400 to-indigo-600',
    bg: 'bg-blue-50',
  },
  {
    key: 'bvtv',
    label: 'Chi cục BVTV An Giang',
    desc: 'Kiểm định, phê duyệt nhật ký',
    username: 'bvtv',
    icon: '🛡️',
    gradient: 'from-amber-400 to-orange-600',
    bg: 'bg-amber-50',
  },
  {
    key: 'bvtv',
    label: 'Chi cục BVTV Vĩnh Long',
    desc: 'Kiểm định, phê duyệt nhật ký',
    username: 'bvtv',
    icon: '🛡️',
    gradient: 'from-amber-400 to-orange-600',
    bg: 'bg-amber-50',
  },
  {
    key: 'nptc',
    label: 'Nhà phân phối Tiêu Chuẩn',
    desc: 'Đơn hàng, công nợ, hoa hồng',
    username: 'nptc001',
    icon: '🚚',
    gradient: 'from-cyan-400 to-blue-600',
    bg: 'bg-cyan-50',
  },
]

async function quickLogin(role) {
  loading.value = true
  loadingKey.value = role.key
  error.value = ''
  try {
    await auth.login(role.username, '123456')
    const target = auth.roleHomeRoute
    if (target === '/login') {
      error.value = `Role "${auth.role}" không được hỗ trợ. Vui lòng restart backend.`
    } else {
      router.push(target)
    }
  } catch (e) {
    error.value = e.response?.data?.error || e.message || 'Không thể kết nối đến server'
  } finally {
    loading.value = false
    loadingKey.value = ''
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <!-- Background decorations -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-96 h-96 bg-green-300/30 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-yellow-300/30 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-emerald-200/20 rounded-full blur-3xl"></div>
    </div>

    <div class="relative z-10 w-full max-w-4xl pt-12">
      <!-- Back Link -->
      <div class="absolute top-0 left-0">
        <router-link to="/" class="inline-flex items-center gap-1.5 text-gray-500 hover:text-emerald-600 transition-colors text-sm font-medium bg-white/50 px-3 py-1.5 rounded-full hover:bg-white shadow-sm border border-transparent hover:border-emerald-100">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
          Quay lại Trang chủ
        </router-link>
      </div>

      <!-- Header -->
      <div class="text-center mb-10 mt-6">
        <div class="inline-flex items-center gap-3 mb-4">
          <div class="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-700 rounded-2xl flex items-center justify-center shadow-lg">
            <span class="text-white text-2xl font-bold">G</span>
          </div>
          <h1 class="text-4xl font-bold bg-gradient-to-r from-green-700 to-emerald-500 bg-clip-text text-transparent">
            GAPChain
          </h1>
        </div>
        <p class="text-gray-600 text-lg">Hệ thống truy xuất nguồn gốc nông sản trên Blockchain</p>
        <p class="text-gray-400 text-sm mt-2">Chọn vai trò để truy cập Demo</p>
      </div>

      <!-- Error -->
      <div v-if="error" class="mb-6 p-4 glass rounded-xl border border-red-200 text-red-600 text-center text-sm">
        {{ error }}
      </div>

      <!-- Role Cards Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
        <button
          v-for="role in roles"
          :key="role.key"
          @click="quickLogin(role)"
          :disabled="loading"
          class="group relative glass rounded-2xl p-6 text-left transition-all duration-300 hover:scale-[1.03] hover:shadow-2xl cursor-pointer disabled:opacity-50 disabled:cursor-wait border border-white/40 overflow-hidden"
        >
          <!-- Gradient overlay on hover -->
          <div
            class="absolute inset-0 bg-gradient-to-br opacity-0 group-hover:opacity-10 transition-opacity duration-300 rounded-2xl"
            :class="role.gradient"
          ></div>

          <div class="relative z-10 flex items-start gap-4">
            <div
              class="w-16 h-16 rounded-2xl bg-gradient-to-br flex items-center justify-center text-3xl shadow-lg group-hover:scale-110 transition-transform duration-300"
              :class="role.gradient"
            >
              <span v-if="loadingKey === role.key" class="w-6 h-6 border-3 border-white/30 border-t-white rounded-full animate-spin"></span>
              <span v-else>{{ role.icon }}</span>
            </div>
            <div class="flex-1">
              <h3 class="text-lg font-bold text-gray-800 group-hover:text-gray-900">{{ role.label }}</h3>
              <p class="text-sm text-gray-500 mt-1">{{ role.desc }}</p>
              <div class="mt-3 flex items-center gap-2">
                <span class="text-xs px-2 py-1 rounded-full bg-gray-100 text-gray-500 font-mono">{{ role.username }}</span>
                <span class="text-xs text-gray-400">Click để đăng nhập</span>
              </div>
            </div>
          </div>

          <!-- Glow effect -->
          <div
            class="absolute -bottom-1 left-1/2 -translate-x-1/2 w-3/4 h-1 rounded-full bg-gradient-to-r opacity-0 group-hover:opacity-100 transition-opacity duration-300 blur-sm"
            :class="role.gradient"
          ></div>
        </button>
      </div>

      <!-- Footer -->
      <p class="text-center text-gray-400 text-xs mt-8">
        Powered by Hyperledger Fabric v3.1.1 &middot; MVP Demo
      </p>
    </div>
  </div>
</template>
