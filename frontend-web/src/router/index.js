import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    name: 'LandingPage',
    component: () => import('@/views/LandingPageView.vue'),
    meta: { guest: true },
  },
  {
    path: '/business',
    name: 'MVPLanding',
    component: () => import('@/views/MVPLandingView.vue'),
    meta: { guest: true },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { guest: true },
  },
  // Public QR Portal
  {
    path: '/qr/consumer/:maLo',
    name: 'QRPortal',
    component: () => import('@/views/public/QRPortalView.vue'),
    meta: { guest: true },
  },
  // HTX Routes
  {
    path: '/htx',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: 'htx' },
    children: [
      { path: '', name: 'HTXDashboard', component: () => import('@/views/htx/DashboardView.vue') },
      { path: 'lohang', name: 'HTXLoHang', component: () => import('@/views/htx/LoHangView.vue') },
      { path: 'lohang/:maLo', name: 'HTXLoHangDetail', component: () => import('@/views/htx/LoHangDetailView.vue') },
      { path: 'nhatky', name: 'HTXNhatKy', component: () => import('@/views/htx/NhatKyView.vue') },
      { path: 'giaodich', name: 'HTXGiaoDich', component: () => import('@/views/htx/GiaoDichView.vue') },
    ],
  },
  // NPP Routes
  {
    path: '/npp',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: ['npp', 'nptc'] },
    children: [
      { path: '', name: 'NPPDashboard', component: () => import('@/views/npp/DashboardView.vue') },
      { path: 'giaodich', name: 'NPPGiaoDich', component: () => import('@/views/npp/GiaoDichView.vue') },
      { path: 'congno', name: 'NPPCongNo', component: () => import('@/views/npp/CongNoView.vue') },
    ],
  },
  // BVTV Routes
  {
    path: '/bvtv',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: 'bvtv' },
    children: [
      { path: '', name: 'BVTVDashboard', component: () => import('@/views/bvtv/DashboardView.vue') },
      { path: 'duyet', name: 'BVTVDuyet', component: () => import('@/views/bvtv/DuyetNhatKyView.vue') },
      { path: 'lohang', name: 'BVTVLoHang', component: () => import('@/views/bvtv/KiemDinhLoHangView.vue') },
      { path: 'lohang/:maLo', name: 'BVTVLoHangDetail', component: () => import('@/views/bvtv/KiemDinhLoHangDetailView.vue') },
    ],
  },
  // Platform Routes
  {
    path: '/platform',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: 'platform' },
    children: [
      { path: '', name: 'PlatformDashboard', component: () => import('@/views/platform/DashboardView.vue') },
      { path: 'lohang', name: 'PlatformLoHang', component: () => import('@/views/platform/QuanLyLoHangView.vue') },
      { path: 'giaodich', name: 'PlatformGiaoDich', component: () => import('@/views/platform/QuanLyGiaoDichView.vue') },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/login' },
]

const router = createRouter({
  history: createWebHistory('/gapchain/'),
  routes,
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  if (to.meta.guest) {
    if (auth.isLoggedIn && to.name === 'Login') {
      return next(auth.roleHomeRoute)
    }
    return next()
  }

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return next('/login')
  }

  if (to.meta.role) {
    const allowedRoles = Array.isArray(to.meta.role) ? to.meta.role : [to.meta.role]
    if (!allowedRoles.includes(auth.role)) {
      return next(auth.roleHomeRoute)
    }
  }

  next()
})

export default router
