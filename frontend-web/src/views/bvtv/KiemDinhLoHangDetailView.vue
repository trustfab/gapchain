<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const route = useRoute()
const auth = useAuthStore()
const maLo = route.params.maLo

const lohang = ref(null)
const nhatkys = ref([])
const loading = ref(true)

const showChungNhanModal = ref(false)
const submittingChungNhan = ref(false)
const formChungNhan = ref({
  loai_chung_nhan: 'VietGAP',
  ma_chung_nhan: '',
  co_quan_cap: 'Chi Cục BVTV',
  ngay_cap: '',
  ngay_het_han: '',
  ghi_chu: ''
})

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const [loRes, nkRes] = await Promise.allSettled([
      api.get(`/api/v1/lohang/${maLo}`),
      api.get(`/api/v1/nhatky/lo/${maLo}`),
    ])
    if (loRes.status === 'fulfilled') {
      lohang.value = loRes.value.data?.data || loRes.value.data
    }
    if (nkRes.status === 'fulfilled') {
      const data = nkRes.value.data?.data || nkRes.value.data || []
      nhatkys.value = Array.isArray(data) ? data : []
    }
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

async function capNhatTrangThaiLo(trangThaiMoi) {
  if (!confirm(`Bạn có chắc muốn chuyển trạng thái này không?`)) return
  try {
    await api.put(`/api/v1/lohang/${maLo}/trangthai`, { trang_thai: trangThaiMoi })
    alert('Cập nhật trạng thái thành công!')
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi cập nhật trạng thái lô hàng')
  }
}

function openChungNhanModal() {
  const today = new Date().toISOString().split('T')[0]
  const nextYear = new Date()
  nextYear.setFullYear(nextYear.getFullYear() + 1)
  
  formChungNhan.value = {
    loai_chung_nhan: 'VietGAP',
    ma_chung_nhan: `VG-${Date.now().toString().slice(-6)}`,
    co_quan_cap: 'Chi Cục BVTV',
    ngay_cap: today,
    ngay_het_han: nextYear.toISOString().split('T')[0],
    ghi_chu: ''
  }
  showChungNhanModal.value = true
}

async function submitChungNhan() {
  if (!formChungNhan.value.ma_chung_nhan) return
  submittingChungNhan.value = true
  try {
    await api.post(`/api/v1/lohang/${maLo}/chungnhan`, formChungNhan.value)
    // Tự động cập nhật san_sang_ban nếu chứng nhận OK
    await api.put(`/api/v1/lohang/${maLo}/trangthai`, { trang_thai: 'san_sang_ban' })
    showChungNhanModal.value = false
    alert('Đã cấp chứng nhận và chuyển trạng thái Sẵn sàng bán thành công!')
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi khi cấp chứng nhận')
  } finally {
    submittingChungNhan.value = false
  }
}

const getLohangBadge = (s) => {
  const map = {
    dang_trong: { color: 'bg-green-100 text-green-700', text: 'Đang trồng', icon: '🌱' },
    da_thu_hoach: { color: 'bg-yellow-100 text-yellow-700', text: 'Đã thu hoạch', icon: '🌾' },
    cho_chung_nhan: { color: 'bg-orange-100 text-orange-700', text: 'Chờ chứng nhận', icon: '🔍' },
    san_sang_ban: { color: 'bg-blue-100 text-blue-700', text: 'Sẵn sàng bán', icon: '✅' },
    het_hang: { color: 'bg-gray-100 text-gray-700', text: 'Hết hàng', icon: '📦' },
    dinh_chi: { color: 'bg-red-100 text-red-700', text: 'Đình chỉ', icon: '⛔' }
  }
  return map[s] || { color: 'bg-gray-100 text-gray-700', text: s || 'N/A', icon: '📌' }
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
    <router-link to="/bvtv/lohang" class="text-emerald-600 hover:text-emerald-700 text-sm font-medium mb-4 inline-block">
      &larr; Quay lại danh sách Lô Hàng
    </router-link>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải chi tiết lô hàng...</div>

    <template v-else>
      <!-- Thông tin Lô Hàng -->
      <div class="glass rounded-2xl p-6 border border-white/40 mb-6 relative overflow-hidden" 
           :class="lohang?.trang_thai === 'cho_chung_nhan' ? 'ring-2 ring-orange-400/50' : ''">
        <div class="flex items-start justify-between">
          <div>
            <h2 class="text-2xl font-bold text-gray-800">{{ lohang?.ten_san_pham || maLo }}</h2>
            <p class="text-gray-500 font-mono text-sm mt-1">
              {{ maLo }} 
              <span v-if="lohang?.ma_lo_me" class="ml-2 px-2 py-0.5 bg-indigo-50 text-indigo-600 rounded text-xs border border-indigo-100">
                Lô gốc: {{ lohang.ma_lo_me }}
              </span>
            </p>
          </div>
          <div class="flex items-center gap-3">
            <button v-if="lohang?.trang_thai === 'cho_chung_nhan'" 
                    @click="openChungNhanModal"
                    class="px-4 py-2 bg-gradient-to-r from-orange-400 to-amber-500 text-white rounded-xl font-medium btn-glow text-sm hover:opacity-90 shadow-md">
              🏅 Cấp Chứng Nhận
            </button>
            <button v-if="lohang?.trang_thai !== 'dinh_chi' && lohang?.trang_thai !== 'het_hang'"
                    @click="capNhatTrangThaiLo('dinh_chi')"
                    class="px-4 py-2 bg-red-50 text-red-600 rounded-xl text-sm font-medium hover:bg-red-100 border border-red-200 transition-all flex items-center gap-1">
              ⛔ Đình chỉ
            </button>
            <span class="px-3 py-1.5 rounded-full text-xs font-medium flex items-center gap-1" :class="getLohangBadge(lohang?.trang_thai).color">
              <span>{{ getLohangBadge(lohang?.trang_thai).icon }}</span>
              <span>{{ getLohangBadge(lohang?.trang_thai).text }}</span>
            </span>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-5">
          <div>
            <p class="text-xs text-gray-400">Hợp tác xã</p>
            <p class="font-medium text-gray-700">{{ lohang?.ma_htx || '-' }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400">Số lượng xuất / Tồn kho</p>
            <p class="font-medium text-gray-700">{{ lohang?.so_luong_con_lai }} / {{ lohang?.so_luong }} {{ lohang?.don_vi_tinh }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400">Vụ mùa</p>
            <p class="font-medium text-gray-700">{{ lohang?.vu_mua || '-' }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400">Địa điểm</p>
            <p class="font-medium text-gray-700">{{ lohang?.dia_diem || '-' }}</p>
          </div>
        </div>

        <div v-if="lohang?.chung_nhan && lohang.chung_nhan.length" class="mt-5 pt-4 border-t border-gray-100">
          <p class="text-sm font-semibold text-gray-600 mb-2">Chứng nhận đã cấp</p>
          <div class="flex flex-wrap gap-2">
            <span v-for="cn in lohang.chung_nhan" :key="cn.ma_chung_nhan"
              class="px-3 py-1.5 text-xs rounded-lg bg-emerald-50 text-emerald-700 font-medium border border-emerald-100 flex flex-col">
              <span>{{ cn.loai_chung_nhan }} - {{ cn.ma_chung_nhan }}</span>
              <span class="text-emerald-500 text-[10px]">{{ cn.ngay_cap }} &rarr; {{ cn.ngay_het_han }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- Danh sách Nhật Ký -->
      <div class="glass rounded-2xl p-6 border border-white/40">
        <h3 class="text-xl font-bold text-gray-800 mb-6 flex items-center justify-between">
          <span>Nhật ký canh tác & Xác nhận</span>
          <span class="text-sm font-normal text-gray-500 bg-gray-100 px-3 py-1 rounded-full">{{ nhatkys.length }} bản ghi</span>
        </h3>

        <div v-if="nhatkys.length === 0" class="text-center py-8 text-gray-400 border border-dashed rounded-xl border-gray-200">
          Lô hàng này chưa có nhật ký canh tác nào.
        </div>

        <div v-else class="space-y-4">
          <div v-for="nk in nhatkys" :key="nk.ma_nhat_ky" class="bg-white/60 rounded-xl p-4 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <span class="font-semibold text-gray-800 capitalize">{{ nk.loai_hoat_dong?.replace(/_/g, ' ') }}</span>
                  <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="statusStyle(nk.trang_thai)">
                    {{ nk.trang_thai === 'da_duyet' ? 'Đã duyệt' : nk.trang_thai === 'tu_choi' ? 'Từ chối' : 'Chờ duyệt' }}
                  </span>
                </div>
                <p class="text-sm text-gray-600">{{ nk.chi_tiet || nk.noi_dung || nk.mo_ta }}</p>
                <div class="flex items-center gap-4 mt-3 text-xs text-gray-500">
                  <span>👤 {{ nk.nguoi_thuc_hien }}</span>
                  <span>📍 {{ nk.vi_tri || 'Chưa định vị' }}</span>
                  <span class="bg-gray-100 px-2 py-0.5 rounded">Ngày: {{ nk.ngay || nk.ngay_ghi }}</span>
                  <span v-if="nk.minh_chung_hash" class="font-mono bg-gray-100 px-1.5 py-0.5 border rounded" title="Minh chứng Hash">
                    🔗 {{ nk.minh_chung_hash.substring(0, 8) }}...
                  </span>
                </div>
              </div>

              <!-- Buttons for Duyệt Nhật Ký -->
              <div v-if="nk.trang_thai === 'cho_duyet'" class="flex gap-2 ml-4 flex-col sm:flex-row">
                <button @click="duyetNhatKy(nk.ma_nhat_ky, true)"
                  class="px-4 py-2 bg-green-50 text-green-600 border border-green-200 rounded-lg text-sm font-medium hover:bg-green-500 hover:text-white transition-all shadow-sm">
                  Phê duyệt
                </button>
                <button @click="duyetNhatKy(nk.ma_nhat_ky, false)"
                  class="px-4 py-2 bg-red-50 text-red-600 border border-red-200 rounded-lg text-sm font-medium hover:bg-red-500 hover:text-white transition-all shadow-sm">
                  Bác bỏ
                </button>
              </div>
              <div v-else-if="nk.quyet_dinh || nk.nguoi_duyet" class="ml-4 text-right">
                <p class="text-[10px] text-gray-400">Người duyệt:</p>
                <p class="text-xs font-medium text-gray-600">{{ nk.nguoi_duyet }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Modal Cấp Chứng Nhận -->
    <transition name="slide">
      <div v-if="showChungNhanModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
        <div class="bg-white rounded-2xl w-full max-w-md overflow-hidden shadow-2xl">
          <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
            <h3 class="font-bold text-gray-800 text-lg">Cấp Chứng Nhận Cho Lô Hàng</h3>
            <button @click="showChungNhanModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
              <span class="text-2xl">&times;</span>
            </button>
          </div>
          
          <div class="p-6">
            <p class="mb-4 text-sm text-gray-600">Đang cấp chứng nhận cho lô: <strong class="text-indigo-600">{{ maLo }}</strong></p>
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Loại chứng nhận *</label>
                <select v-model="formChungNhan.loai_chung_nhan" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all">
                  <option value="VietGAP">VietGAP</option>
                  <option value="GlobalGAP">GlobalGAP</option>
                  <option value="Organic">Hữu cơ (Organic)</option>
                  <option value="TCVN">TCVN cơ bản</option>
                  <option value="Khac">Khác</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Mã chứng nhận *</label>
                <input v-model="formChungNhan.ma_chung_nhan" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all" />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Ngày cấp</label>
                  <input v-model="formChungNhan.ngay_cap" type="date" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Ngày hết hạn</label>
                  <input v-model="formChungNhan.ngay_het_han" type="date" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all" />
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Ghi chú thêm</label>
                <input v-model="formChungNhan.ghi_chu" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all" />
              </div>
            </div>
          </div>

          <div class="px-6 py-4 bg-gray-50 flex justify-end gap-3 rounded-b-2xl">
            <button @click="showChungNhanModal = false" class="px-5 py-2 rounded-xl text-gray-600 font-medium hover:bg-gray-200 transition-colors">
              Hủy
            </button>
            <button @click="submitChungNhan" :disabled="submittingChungNhan || !formChungNhan.ma_chung_nhan"
              class="px-5 py-2 bg-gradient-to-r from-amber-500 to-orange-600 text-white rounded-xl font-medium hover:opacity-90 disabled:opacity-50 transition-all flex items-center gap-2">
              <span v-if="submittingChungNhan" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              Cấp chứng nhận ngay
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>
