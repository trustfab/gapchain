<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const navItems = computed(() => {
  const map = {
    htx: [
      { label: 'Tổng quan', icon: '📊', to: '/htx' },
      { label: 'Lô Hàng', icon: '📦', to: '/htx/lohang' },
      { label: 'Nhật Ký', icon: '📝', to: '/htx/nhatky' },
      { label: 'Giao Dịch', icon: '💰', to: '/htx/giaodich' },
    ],
    npp: [
      { label: 'Tổng quan', icon: '📊', to: '/npp' },
      { label: 'Giao Dịch', icon: '📋', to: '/npp/giaodich' },
      { label: 'Công Nợ', icon: '💳', to: '/npp/congno' },
    ],
    nptc: [
      { label: 'Tổng quan', icon: '📊', to: '/npp' },
      { label: 'Giao Dịch', icon: '📋', to: '/npp/giaodich' },
      { label: 'Công Nợ', icon: '💳', to: '/npp/congno' },
    ],
    bvtv: [
      { label: 'Tổng quan', icon: '📊', to: '/bvtv' },
      { label: 'Kiểm định Lô Hàng', icon: '🏅', to: '/bvtv/lohang' },
      { label: 'Duyệt Nhật Ký', icon: '✅', to: '/bvtv/duyet' },
    ],
    platform: [
      { label: 'Ecosystem', icon: '🌐', to: '/platform' },
      { label: 'Q.lý Lô Hàng', icon: '📦', to: '/platform/lohang' },
      { label: 'Q.lý Giao Dịch', icon: '💰', to: '/platform/giaodich' },
    ],
  }
  return map[auth.role] || []
})

const roleLabel = computed(() => {
  const map = { htx: 'Hợp tác xã', npp: 'Nhà phân phối', bvtv: 'Chi cục BVTV', platform: 'Platform Admin' }
  return map[auth.role] || ''
})

const roleBg = computed(() => {
  const map = { htx: 'from-green-600 to-emerald-700', npp: 'from-blue-600 to-indigo-700', bvtv: 'from-amber-600 to-orange-700', platform: 'from-purple-600 to-violet-700' }
  return map[auth.role] || 'from-gray-600 to-gray-700'
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-gradient-to-b min-h-screen p-4 flex flex-col shadow-xl" :class="roleBg">
      <div class="mb-8">
        <h1 class="text-white text-xl font-bold tracking-wide">GAPChain</h1>
        <p class="text-white/60 text-sm mt-1">{{ roleLabel }}</p>
      </div>

      <nav class="flex-1 space-y-1">
        <router-link
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-white/80 hover:bg-white/15 transition-all duration-200"
          active-class="!bg-white/20 !text-white font-semibold"
        >
          <span class="text-lg">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="mt-auto pt-4 border-t border-white/20">
        <div class="flex items-center gap-3 px-4 py-2">
          <div class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center text-white text-sm font-bold">
            {{ auth.username?.charAt(0)?.toUpperCase() }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-white text-sm font-medium truncate">{{ auth.username }}</p>
            <p class="text-white/50 text-xs">{{ auth.tenantId }}</p>
          </div>
        </div>
        <button
          @click="handleLogout"
          class="w-full mt-2 px-4 py-2 text-sm text-white/70 hover:text-white hover:bg-white/10 rounded-lg transition-all"
        >
          Đăng xuất
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-6 overflow-auto">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>
