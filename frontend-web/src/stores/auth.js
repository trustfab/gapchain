import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/plugins/axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => user.value?.role || '')
  const tenantId = computed(() => user.value?.tenant_id || '')
  const mspId = computed(() => user.value?.msp_id || '')
  const username = computed(() => user.value?.username || '')

  const isPlatform = computed(() => role.value === 'platform')
  const isHTX = computed(() => role.value === 'htx')
  const isNPP = computed(() => role.value === 'npp' || role.value === 'nptc')
  const isBVTV = computed(() => role.value === 'bvtv')

  const roleHomeRoute = computed(() => {
    const map = {
      platform: '/platform',
      admin: '/platform',
      htx: '/htx',
      npp: '/npp',
      nptc: '/npp',
      bvtv: '/bvtv',
    }
    return map[role.value] || '/login'
  })

  // Infer multi-tenant fields from username if backend doesn't return them
  function inferTenant(username, serverUser) {
    const defaults = {
      htx: { role: 'htx', msp_id: 'HTXNongSanOrgMSP', tenant_id: 'HTX001' },
      npp: { role: 'npp', msp_id: 'NPPXanhOrgMSP', tenant_id: 'NPP001' },
      nptc: { role: 'nptc', msp_id: 'NPPTieuChuanOrgMSP', tenant_id: 'NPTC001' },
      bvtv: { role: 'bvtv', msp_id: 'ChiCucBVTVOrgMSP', tenant_id: 'BVTV_CT' },
      platform: { role: 'platform', msp_id: 'PlatformOrgMSP', tenant_id: 'PLATFORM' },
    }
    const prefix3 = username?.substring(0, 3)
    const prefix4 = username?.substring(0, 4)
    let fallback = null
    if (username === 'platform') fallback = defaults.platform
    else if (username === 'bvtv') fallback = defaults.bvtv
    else if (prefix4 === 'nptc') fallback = defaults.nptc
    else if (prefix3 === 'htx') fallback = defaults.htx
    else if (prefix3 === 'npp') fallback = defaults.npp

    return {
      username: serverUser.username || username,
      role: serverUser.tenant_id ? serverUser.role : (fallback?.role || serverUser.role),
      msp_id: serverUser.msp_id || fallback?.msp_id || '',
      tenant_id: serverUser.tenant_id || fallback?.tenant_id || '',
    }
  }

  async function login(uname, password) {
    const res = await api.post('/api/v1/auth/login', { username: uname, password })
    token.value = res.data.token
    user.value = inferTenant(uname, res.data.user || {})
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('user', JSON.stringify(user.value))
    return res.data
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token, user,
    isLoggedIn, role, tenantId, mspId, username,
    isPlatform, isHTX, isNPP, isBVTV,
    roleHomeRoute,
    login, logout,
  }
})
