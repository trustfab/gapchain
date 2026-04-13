<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '@/plugins/axios'

import { useAuthStore } from '@/stores/auth'

const stats = ref({
  fabric_nodes: 5,
  lohang_total: 0,
  nhatky_total: 0,
  giaodich_total: 0,
})
const animatedStats = ref({ ...stats.value })
const activityFeed = ref([])
const loading = ref(true)

let feedInterval = null

onMounted(async () => {
  // Load real stats and real feed data
  try {
    const auth = useAuthStore()
    const [statRes, loRes, nkRes, gdRes] = await Promise.allSettled([
      api.get('/api/v1/nhatky/thongke'),
      api.get(`/api/v1/lohang/htx/HTX001`), // MVP HTX001
      api.get(`/api/v1/nhatky/htx/HTX001`),
      api.get(`/api/v1/giaodich/htx/HTX001`)
    ])

    let tongLoHang = 0
    let tongNhatKy = 0
    let tongGiaoDich = 0
    let feedItems = []

    if (loRes.status === 'fulfilled') {
      const data = loRes.value.data?.data || loRes.value.data || []
      const lohangs = Array.isArray(data) ? data : []
      tongLoHang = lohangs.length
      lohangs.forEach(lo => {
        feedItems.push({
          type: 'lohang',
          msg: `${lo.ma_htx} tạo lô hàng mã ${lo.ma_lo}: ${lo.ten_san_pham}`,
          time: new Date().toISOString(), // Lô hàng chưa có field ngay_tao explicit trong payload chuẩn
          timestamp: Date.now() - Math.random() * 86400000, // randomizing slightly for sorting if no timestamp
          icon: '📦'
        })
        if (lo.chung_nhan && lo.chung_nhan.length > 0) {
          lo.chung_nhan.forEach(cn => {
            feedItems.push({
               type: 'chungnhan',
               msg: `Chi cục BVTV đã cấp ${cn.loai_chung_nhan} cho lô ${lo.ma_lo}`,
               time: cn.ngay_cap,
               timestamp: new Date(cn.ngay_cap).getTime(),
               icon: '🛡️'
            })
          })
        }
      })
    }
    
    if (nkRes.status === 'fulfilled') {
      const data = nkRes.value.data?.data || nkRes.value.data || []
      const nhatkys = Array.isArray(data) ? data : []
      tongNhatKy = nhatkys.length
      nhatkys.forEach(nk => {
        const timeStr = nk.ngay || nk.ngay_ghi
        const ts = timeStr ? new Date(timeStr).getTime() : Date.now()
        feedItems.push({
          type: 'nhatky',
          msg: `${nk.ma_htx} thêm hoạt động '${nk.loai_hoat_dong?.replace(/_/g, ' ')}' cho lô ${nk.ma_lo}`,
          time: timeStr,
          timestamp: ts,
          icon: '📝'
        })
        if(nk.trang_thai === 'da_duyet') {
           feedItems.push({
             type: 'duyet',
             msg: `Cập nhật: BVTV đã phê duyệt nhật ký canh tác của ${nk.ma_htx}`,
             time: timeStr,
             timestamp: ts + 10000,
             icon: '✅'
           })
        }
      })
    }

    if (gdRes.status === 'fulfilled') {
      const data = gdRes.value.data?.data || gdRes.value.data || []
      const giaodichs = Array.isArray(data) ? data : []
      tongGiaoDich = giaodichs.length
      giaodichs.forEach(gd => {
        const timeStr = gd.ngay_tao
        const ts = timeStr ? new Date(timeStr).getTime() : Date.now()
        feedItems.push({
          type: 'giaodich',
          msg: `Phát sinh giao dịch: ${gd.ma_npp} lên đơn mua ${gd.so_luong} ${gd.don_vi_tinh || 'kg'} lô ${gd.ma_lo}`,
          time: timeStr,
          timestamp: ts,
          icon: '💼'
        })
        if (gd.trang_thai === 'da_thanh_toan') {
           feedItems.push({
             type: 'thanhtoan',
             msg: `${gd.ma_npp} đã thanh toán ${(gd.so_luong * gd.don_gia).toLocaleString()} VNĐ cho ${gd.ma_htx}`,
             time: timeStr,
             timestamp: ts + 20000,
             icon: '💰'
           })
        }
      })
    }

    // Sort descending by timestamp
    feedItems.sort((a, b) => b.timestamp - a.timestamp)
    
    // Format times into readable string
    feedItems = feedItems.map(item => {
       let displayTime = 'Chưa rõ'
       if (item.time) {
          try { displayTime = new Date(item.time).toLocaleString('vi-VN') } catch(e){}
       }
       return { ...item, time: displayTime }
    })

    stats.value = {
      fabric_nodes: 5, // 4 Peers (Platform, HTX, NPP, BVTV) + 1 Orderer
      lohang_total: tongLoHang,
      nhatky_total: tongNhatKy,
      giaodich_total: tongGiaoDich,
    }
    
    activityFeed.value = feedItems.slice(0, 10) // Display top 10

  } catch (err) {
    console.error('Error loading real feed', err)
  }
  loading.value = false

  // Animate counters
  animateCounters()
})

// Removed onUnmounted mock interval hook
onUnmounted(() => {})

function animateCounters() {
  const duration = 1500
  const start = Date.now()
  const targets = { ...stats.value }

  function tick() {
    const elapsed = Date.now() - start
    const progress = Math.min(elapsed / duration, 1)
    const eased = 1 - Math.pow(1 - progress, 3)

    animatedStats.value = {
      fabric_nodes: Math.round(targets.fabric_nodes * eased),
      lohang_total: Math.round(targets.lohang_total * eased),
      nhatky_total: Math.round(targets.nhatky_total * eased),
      giaodich_total: Math.round(targets.giaodich_total * eased),
    }

    if (progress < 1) requestAnimationFrame(tick)
  }
  requestAnimationFrame(tick)
}

const typeColor = (type) => {
  const map = {
    lohang: 'border-green-400',
    nhatky: 'border-blue-400',
    giaodich: 'border-amber-400',
    duyet: 'border-emerald-400',
    thanhtoan: 'border-purple-400',
    chungnhan: 'border-orange-400',
  }
  return map[type] || 'border-gray-300'
}
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-800 mb-2">Ecosystem Overview</h2>
    <p class="text-gray-500 mb-6">Giám sát toàn bộ hệ sinh thái GAPChain</p>

    <!-- Animated Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
      <div class="glass rounded-2xl p-6 border border-white/40 text-center">
        <div class="w-14 h-14 mx-auto rounded-2xl bg-gradient-to-br from-purple-400 to-violet-600 flex items-center justify-center text-2xl mb-3">🖧</div>
        <p class="text-4xl font-bold text-gray-800 tabular-nums">{{ animatedStats.fabric_nodes }}</p>
        <p class="text-sm text-gray-500 mt-1">Fabric Nodes</p>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40 text-center">
        <div class="w-14 h-14 mx-auto rounded-2xl bg-gradient-to-br from-green-400 to-emerald-600 flex items-center justify-center text-2xl mb-3">📦</div>
        <p class="text-4xl font-bold text-gray-800 tabular-nums">{{ animatedStats.lohang_total }}</p>
        <p class="text-sm text-gray-500 mt-1">Lô Hàng</p>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40 text-center">
        <div class="w-14 h-14 mx-auto rounded-2xl bg-gradient-to-br from-blue-400 to-indigo-600 flex items-center justify-center text-2xl mb-3">📝</div>
        <p class="text-4xl font-bold text-gray-800 tabular-nums">{{ animatedStats.nhatky_total }}</p>
        <p class="text-sm text-gray-500 mt-1">Nhật Ký</p>
      </div>
      <div class="glass rounded-2xl p-6 border border-white/40 text-center">
        <div class="w-14 h-14 mx-auto rounded-2xl bg-gradient-to-br from-amber-400 to-orange-600 flex items-center justify-center text-2xl mb-3">💰</div>
        <p class="text-4xl font-bold text-gray-800 tabular-nums">{{ animatedStats.giaodich_total }}</p>
        <p class="text-sm text-gray-500 mt-1">Giao Dịch</p>
      </div>
    </div>

    <!-- Live Activity Feed -->
    <div class="glass rounded-2xl p-6 border border-white/40">
      <div class="flex items-center gap-3 mb-5">
        <div class="w-3 h-3 rounded-full bg-green-500 animate-pulse"></div>
        <h3 class="text-lg font-semibold text-gray-800">Live Activity Feed</h3>
      </div>

      <div class="space-y-3">
        <transition-group name="slide">
          <div
            v-for="(activity, idx) in activityFeed"
            :key="activity.msg + idx"
            class="flex items-center gap-4 p-4 rounded-xl bg-white/40 border-l-4 transition-all duration-500"
            :class="typeColor(activity.type)"
          >
            <span class="text-2xl">{{ activity.icon }}</span>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-700 truncate">{{ activity.msg }}</p>
            </div>
            <span class="text-xs text-gray-400 whitespace-nowrap">{{ activity.time }}</span>
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>
