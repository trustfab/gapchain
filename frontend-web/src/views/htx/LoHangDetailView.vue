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

onMounted(async () => {
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
})

const showTachLoModal = ref(false)
const submittingTachLo = ref(false)
const formTachLo = ref({
  ma_lo_moi: '',
  so_luong_tach: ''
})

function openTachLo() {
  const now = Date.now()
  formTachLo.value.ma_lo_moi = `LH-CON-${auth.tenantId}-${now}`
  formTachLo.value.so_luong_tach = ''
  showTachLoModal.value = true
}

async function submitTachLo() {
  if (!formTachLo.value.ma_lo_moi || !formTachLo.value.so_luong_tach) return
  submittingTachLo.value = true
  try {
    await api.post(`/api/v1/lohang/${maLo}/tach`, {
      ma_lo_moi: formTachLo.value.ma_lo_moi,
      so_luong_tach: parseFloat(formTachLo.value.so_luong_tach)
    })
    showTachLoModal.value = false
    // Tải lại dữ liệu lô hàng
    const res = await api.get(`/api/v1/lohang/${maLo}`)
    lohang.value = res.data?.data || res.data
    alert('Tách lô thành công! Lô con đã được tạo.')
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi khi tách lô')
  } finally {
    submittingTachLo.value = false
  }
}

async function capNhatTrangThai(newStatus) {
  if (!confirm(`Bạn có chắc muốn chuyển trạng thái này không?`)) return
  try {
    await api.put(`/api/v1/lohang/${maLo}/trangthai`, { trang_thai: newStatus })
    const res = await api.get(`/api/v1/lohang/${maLo}`)
    lohang.value = res.data?.data || res.data
    alert('Cập nhật trạng thái thành công!')
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi khi cập nhật trạng thái')
  }
}


const activityIcon = (loai) => {
  const map = {
    gieo_hat: '🌱', bon_phan: '🧪', tuoi_nuoc: '💧',
    phun_thuoc: '🔬', thu_hoach: '🌾', kiem_tra: '🔍',
  }
  return map[loai] || '📝'
}

const formatLoHangStatus = (s) => {
  const map = {
    dang_trong: { label: '🌱 Đang trồng', class: 'bg-green-100 text-green-700' },
    da_thu_hoach: { label: '🌾 Đã thu hoạch', class: 'bg-yellow-100 text-yellow-700' },
    cho_chung_nhan: { label: '🔍 Chờ chứng nhận', class: 'bg-orange-100 text-orange-700' },
    san_sang_ban: { label: '✅ Sẵn sàng bán', class: 'bg-blue-100 text-blue-700' },
    het_hang: { label: '📦 Hết hàng', class: 'bg-gray-100 text-gray-700' },
    dinh_chi: { label: '⛔ Đình chỉ', class: 'bg-red-100 text-red-700' }
  }
  return map[s] || { label: s || 'N/A', class: 'bg-gray-100 text-gray-700' }
}

const getLoaiHoatDongLabel = (value) => {
  const map = {
    gieo_hat: "Gieo hạt", bon_phan: "Bón phân", tuoi_nuoc: "Tưới nước",
    phun_thuoc: "Phun thuốc", kiem_tra: "Kiểm tra", thu_hoach: "Thu hoạch",
    dong_goi: "Đóng gói", van_chuyen: "Vận chuyển", khac: "Khác"
  }
  return map[value] || (value?.replace(/_/g, ' ') || 'Hoạt động')
}
</script>

<template>
  <div>
    <router-link to="/htx/lohang" class="text-green-600 hover:text-green-700 text-sm font-medium mb-4 inline-block">
      &larr; Quay lại danh sách
    </router-link>

    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <template v-else>
      <!-- Lo Hang Info -->
      <div class="glass rounded-2xl p-6 border border-white/40 mb-6">
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
            <button v-if="lohang?.trang_thai === 'dang_trong'" 
                    @click="capNhatTrangThai('da_thu_hoach')"
                    class="px-4 py-2 bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-xl font-medium btn-glow text-sm hover:opacity-90">
              🌾 Báo Thu Hoạch
            </button>
            <button v-if="lohang?.trang_thai === 'da_thu_hoach'" 
                    @click="capNhatTrangThai('cho_chung_nhan')"
                    class="px-4 py-2 bg-gradient-to-r from-indigo-500 to-purple-600 text-white rounded-xl font-medium btn-glow text-sm hover:opacity-90">
              🔍 Gửi Kiểm Định
            </button>
            <button v-if="lohang?.trang_thai === 'cho_chung_nhan'" 
                    @click="capNhatTrangThai('san_sang_ban')"
                    class="px-4 py-2 bg-gradient-to-r from-emerald-500 to-green-600 text-white rounded-xl font-medium btn-glow text-sm hover:opacity-90">
              ✅ Sẵn Sàng Bán
            </button>
            <button v-if="lohang?.trang_thai === 'san_sang_ban'" 
                    @click="openTachLo"
                    class="px-4 py-2 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-xl font-medium btn-glow text-sm hover:opacity-90">
              ✂️ Tách Lô
            </button>
            <span class="px-3 py-1.5 rounded-full text-xs font-medium" :class="formatLoHangStatus(lohang?.trang_thai).class">
              {{ formatLoHangStatus(lohang?.trang_thai).label }}
            </span>
          </div>
        </div>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-5">
          <div>
            <p class="text-xs text-gray-400">Loại sản phẩm</p>
            <p class="font-medium text-gray-700 capitalize">{{ lohang?.loai_san_pham?.replace('_', ' ') || '-' }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400">Số lượng</p>
            <p class="font-medium text-gray-700">{{ lohang?.so_luong || '-' }} {{ lohang?.don_vi_tinh || '' }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400">Còn lại (Tồn)</p>
            <p class="font-medium" :class="(lohang?.so_luong_con_lai === 0) ? 'text-red-500' : 'text-gray-700'">
              {{ typeof lohang?.so_luong_con_lai === 'number' ? lohang.so_luong_con_lai : lohang?.so_luong || '-' }} {{ lohang?.don_vi_tinh || '' }}
            </p>
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

        <!-- Chung nhan -->
        <div v-if="lohang?.chung_nhan?.length" class="mt-5 pt-4 border-t border-gray-100">
          <p class="text-sm font-semibold text-gray-600 mb-2">Chứng nhận</p>
          <div class="flex flex-wrap gap-2">
            <span v-for="cn in lohang.chung_nhan" :key="cn.ma_chung_nhan"
              class="px-3 py-1 text-xs rounded-full bg-emerald-100 text-emerald-700 font-medium">
              {{ cn.ten_chung_nhan || cn.ma_chung_nhan }}
            </span>
          </div>
        </div>
      </div>

      <!-- Nhat Ky Timeline -->
      <div class="glass rounded-2xl p-6 border border-white/40">
        <h3 class="text-lg font-semibold text-gray-800 mb-6">Nhật ký canh tác</h3>

        <div v-if="nhatkys.length === 0" class="text-center py-8 text-gray-400">
          Chưa có nhật ký nào cho lô hàng này.
        </div>

        <div v-else class="relative">
          <!-- Timeline line -->
          <div class="absolute left-6 top-0 bottom-0 w-0.5 bg-gradient-to-b from-green-400 to-emerald-200"></div>

          <div v-for="(nk, idx) in nhatkys" :key="nk.id || idx" class="relative pl-16 pb-8 last:pb-0">
            <!-- Timeline dot -->
            <div class="absolute left-4 w-5 h-5 rounded-full bg-white border-2 border-green-500 flex items-center justify-center text-xs z-10">
              {{ activityIcon(nk.loai_hoat_dong) }}
            </div>

            <div class="glass rounded-xl p-4 border border-white/30 hover:shadow-md transition-all">
              <div class="flex items-start justify-between">
                <div>
                  <h4 class="font-semibold text-gray-800 capitalize">{{ getLoaiHoatDongLabel(nk.loai_hoat_dong) }}</h4>
                  <p class="text-sm text-gray-600 mt-1">{{ nk.chi_tiet || nk.mo_ta || nk.noi_dung || '' }}</p>
                  
                  <div class="flex flex-wrap items-center gap-4 mt-3">
                    <div class="flex items-center gap-1.5 text-xs text-gray-500">
                      <span>👤</span>
                      <span class="font-medium">{{ nk.nguoi_thuc_hien || 'Không rõ' }}</span>
                    </div>
                    <div class="flex items-center gap-1.5 text-xs text-gray-500">
                      <span>📍</span>
                      <span>{{ nk.vi_tri || 'Chưa định vị' }}</span>
                    </div>
                    <div v-if="nk.minh_chung_hash" class="flex items-center gap-1.5 text-xs text-gray-500" title="Mã chứng nhận (Hash)">
                      <span>🔗</span>
                      <span class="font-mono text-[10px] bg-gray-100 border border-gray-200 px-1.5 py-0.5 rounded">{{ nk.minh_chung_hash.substring(0, 10) }}...</span>
                    </div>
                  </div>
                </div>
                
                <div class="text-right ml-4 flex-shrink-0">
                  <span class="text-xs font-medium bg-gray-100/80 px-2.5 py-1 rounded-lg text-gray-600 block mb-2">{{ nk.ngay_ghi || nk.ngay || '' }}</span>
                  <div v-if="nk.trang_thai">
                    <span class="text-xs px-2 py-0.5 rounded-full inline-block"
                      :class="
                        nk.trang_thai === 'da_duyet' ? 'bg-green-100 text-green-700' : 
                        nk.trang_thai === 'tu_choi' ? 'bg-red-100 text-red-700' : 
                        nk.trang_thai === 'cho_nong_dan_xac_nhan' ? 'bg-purple-100 text-purple-700' :
                        'bg-yellow-100 text-yellow-700'">
                      {{ nk.trang_thai === 'da_duyet' ? '✅ Đã duyệt' : nk.trang_thai === 'tu_choi' ? '❌ Từ chối' : nk.trang_thai === 'cho_nong_dan_xac_nhan' ? '⏳ Chờ cập nhật' : '⏳ Chờ duyệt' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tach Lo Modal -->
      <transition name="slide">
        <div v-if="showTachLoModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
          <div class="bg-white rounded-2xl w-full max-w-md overflow-hidden shadow-2xl">
            <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
              <h3 class="font-bold text-gray-800 text-lg">Tách Lô Hàng (Partial Batch)</h3>
              <button @click="showTachLoModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
                <span class="text-2xl">&times;</span>
              </button>
            </div>
            
            <div class="p-6">
              <div class="bg-blue-50 text-blue-800 text-sm p-3 rounded-lg border border-blue-100 mb-5">
                <p>Số lượng tồn kho lô mẹ hiện tại: <strong>{{ lohang?.so_luong_con_lai }} {{ lohang?.don_vi_tinh }}</strong></p>
              </div>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Mã Lô Mới (Kế thừa truy vết)</label>
                  <input v-model="formTachLo.ma_lo_moi" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 outline-none transition-all" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Số lượng tách <span class="text-gray-400 font-normal">({{ lohang?.don_vi_tinh }})</span> *</label>
                  <input v-model="formTachLo.so_luong_tach" type="number" step="0.1" :max="lohang?.so_luong_con_lai" placeholder="VD: 100" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 outline-none transition-all" />
                </div>
              </div>
            </div>

            <div class="px-6 py-4 bg-gray-50 flex justify-end gap-3 rounded-b-2xl">
              <button @click="showTachLoModal = false" class="px-5 py-2 rounded-xl text-gray-600 font-medium hover:bg-gray-200 transition-colors">
                Hủy
              </button>
              <button @click="submitTachLo" :disabled="submittingTachLo || !formTachLo.ma_lo_moi || !formTachLo.so_luong_tach"
                class="px-5 py-2 bg-indigo-600 text-white rounded-xl font-medium hover:bg-indigo-700 disabled:opacity-50 transition-all flex items-center gap-2">
                <span v-if="submittingTachLo" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                Xác nhận tách
              </button>
            </div>
          </div>
        </div>
      </transition>
    </template>
  </div>
</template>
