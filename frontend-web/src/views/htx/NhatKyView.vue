<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/plugins/axios'

const auth = useAuthStore()
const nhatkys = ref([])
const availableLohangs = ref([])
const loading = ref(true)
const showForm = ref(false)

const form = ref({
  ma_lo: '',
  loai_hoat_dong: '',
  chi_tiet: '',
  vi_tri: '',
  nguoi_thuc_hien: '',
  ngay_ghi: new Date().toISOString().split('T')[0],
})
const submitting = ref(false)

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/nhatky/htx/${auth.tenantId}`)
    const data = res.data?.data || res.data || []
    nhatkys.value = Array.isArray(data) ? data : []

    // Fetch Danh sach Lo Hang cua HTX de lap dropdown
    const resLo = await api.get(`/api/v1/lohang/htx/${auth.tenantId}`)
    const loData = resLo.data?.data || resLo.data || []
    availableLohangs.value = Array.isArray(loData) ? loData : []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function ghiNhatKy() {
  submitting.value = true
  try {
    const now = Date.now()
    await api.post('/api/v1/nhatky', {
      ma_nhat_ky: `NK-${auth.tenantId}-${now}`,
      ma_lo: form.value.ma_lo,
      ma_htx: auth.tenantId,
      loai_hoat_dong: form.value.loai_hoat_dong,
      chi_tiet: form.value.chi_tiet,
      vi_tri: form.value.vi_tri,
      nguoi_thuc_hien: form.value.nguoi_thuc_hien || auth.username,
      ngay_ghi: form.value.ngay_ghi,
    })
    showForm.value = false
    form.value = { ma_lo: '', loai_hoat_dong: '', chi_tiet: '', vi_tri: '', nguoi_thuc_hien: '', ngay_ghi: new Date().toISOString().split('T')[0] }
    await fetchData()
  } catch (e) {
    alert(e.response?.data?.error || 'Loi ghi nhat ky')
  } finally {
    submitting.value = false
  }
}

async function xacNhanFarmer(nk) {
  const hash = window.prompt(`Nông dân xác nhận: Nhập mã Hash OTP / Tệp tin kiểm định cho ${nk.loai_hoat_dong}`, '');
  if (!hash) return;
  
  try {
    const id = nk.ma_nhat_ky || nk.id;
    await api.put(`/api/v1/nhatky/${id}/xacnhan`, { 
      minh_chung_hash: hash 
    });
    alert('Đã nộp minh chứng thành công! Hồ sơ chờ chi cục duyệt.');
    await fetchData();
  } catch (e) {
    alert(e.response?.data?.error || 'Lỗi khi xác nhận bằng chứng');
  }
}

const DANH_MUC_HOAT_DONG = [
  {"value": "gieo_hat",   "label": "Gieo hạt"},
  {"value": "bon_phan",   "label": "Bón phân"},
  {"value": "tuoi_nuoc",  "label": "Tưới nước"},
  {"value": "phun_thuoc", "label": "Phun thuốc"},
  {"value": "kiem_tra",   "label": "Kiểm tra"},
  {"value": "thu_hoach",  "label": "Thu hoạch"},
  {"value": "dong_goi",   "label": "Đóng gói"},
  {"value": "van_chuyen", "label": "Vận chuyển"},
  {"value": "khac",       "label": "Khác"}
];

const loaiOptionsToState = {
  dang_trong: ['gieo_hat', 'bon_phan', 'tuoi_nuoc', 'phun_thuoc', 'thu_hoach', 'kiem_tra', 'khac'],
  da_thu_hoach: ['thu_hoach', 'kiem_tra', 'dong_goi', 'van_chuyen', 'khac'],
  san_sang_ban: ['van_chuyen', 'khac']
}

const filteredLoaiOptions = computed(() => {
  if (!form.value.ma_lo) return []
  const lo = availableLohangs.value.find(l => l.ma_lo === form.value.ma_lo)
  const allowedValues = lo ? (loaiOptionsToState[lo.trang_thai] || ['khac']) : ['khac']
  return DANH_MUC_HOAT_DONG.filter(item => allowedValues.includes(item.value))
})

const getLoaiHoatDongLabel = (value) => {
  const item = DANH_MUC_HOAT_DONG.find(i => i.value === value)
  return item ? item.label : (value?.replace(/_/g, ' ') || 'Hoạt động')
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Nhật ký canh tác</h2>
      <button @click="showForm = !showForm"
        class="px-5 py-2.5 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-xl font-medium btn-glow">
        + Ghi Nhật Ký
      </button>
    </div>

    <!-- Form -->
    <transition name="slide">
      <div v-if="showForm" class="glass rounded-2xl p-6 border border-white/40 mb-6">
        <h3 class="font-semibold text-gray-800 mb-4">Ghi nhật ký mới</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select v-model="form.ma_lo" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70">
            <option value="" disabled>-- Chọn một lô hàng * --</option>
            <option v-for="lo in availableLohangs" :key="lo.ma_lo" :value="lo.ma_lo">
              {{ lo.ten_san_pham }} ({{ lo.ma_lo.split('-').pop() }})
            </option>
          </select>
          <select v-model="form.loai_hoat_dong" :disabled="!form.ma_lo" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70 disabled:opacity-50">
            <option value="">-- Loại hoạt động * --</option>
            <option v-for="opt in filteredLoaiOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
          <input v-model="form.chi_tiet" placeholder="Chi tiết hoạt động *" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70 md:col-span-2" />
          <input v-model="form.vi_tri" placeholder="Vị trí (VD: 11.58,108.99) *" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70" />
          <input v-model="form.nguoi_thuc_hien" :placeholder="`Người thực hiện (mặc định: ${auth.username})`" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70" />
          <div>
            <label class="text-sm text-gray-500 mb-1 block">Ngày ghi *</label>
            <input v-model="form.ngay_ghi" type="date" class="px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-white/70 w-full" />
          </div>
        </div>
        <div class="flex gap-3 mt-4">
          <button @click="ghiNhatKy" :disabled="submitting || !form.ma_lo || !form.loai_hoat_dong || !form.chi_tiet || !form.vi_tri || !form.ngay_ghi"
            class="px-5 py-2.5 bg-blue-600 text-white rounded-xl font-medium hover:bg-blue-700 transition-all disabled:opacity-40 disabled:cursor-not-allowed">
            {{ submitting ? 'Đang lưu...' : 'Lưu' }}
          </button>
          <button @click="showForm = false" class="px-5 py-2.5 bg-gray-200 text-gray-700 rounded-xl font-medium hover:bg-gray-300 transition-all">Hủy</button>
        </div>
      </div>
    </transition>

    <!-- List -->
    <div v-if="loading" class="text-center py-12 text-gray-400">Đang tải...</div>

    <div v-else-if="nhatkys.length === 0" class="glass rounded-2xl p-12 border border-white/40 text-center text-gray-400">
      Chưa có nhật ký nào.
    </div>

    <div v-else class="space-y-3">
      <div v-for="nk in nhatkys" :key="nk.ma_nhat_ky || nk.id" class="glass rounded-xl p-5 border border-white/30 hover:shadow-md transition-all">
        <div class="flex items-start justify-between">
          <div>
            <div class="flex items-center gap-2">
              <h4 class="font-semibold text-gray-800 capitalize">{{ getLoaiHoatDongLabel(nk.loai_hoat_dong) }}</h4>
              <span class="text-xs font-mono text-gray-500 bg-gray-100 px-2 py-0.5 rounded border border-gray-200">{{ nk.ma_lo }}</span>
            </div>
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
              <span class="text-xs px-2 py-0.5 rounded-full inline-block mb-2"
                :class="
                  nk.trang_thai === 'da_duyet' ? 'bg-green-100 text-green-700' : 
                  nk.trang_thai === 'tu_choi' ? 'bg-red-100 text-red-700' : 
                  nk.trang_thai === 'cho_nong_dan_xac_nhan' ? 'bg-purple-100 text-purple-700' :
                  'bg-yellow-100 text-yellow-700'
                ">
                {{ nk.trang_thai === 'da_duyet' ? '✅ Đã duyệt' : nk.trang_thai === 'tu_choi' ? '❌ Từ chối' : nk.trang_thai === 'cho_nong_dan_xac_nhan' ? '⏳ Chờ dân xác nhận' : '⏳ Chờ duyệt' }}
              </span>
              <div v-if="nk.trang_thai === 'cho_nong_dan_xac_nhan'">
                <button @click="xacNhanFarmer(nk)" class="text-xs text-purple-600 border border-purple-300 rounded px-2 py-1 hover:bg-purple-50">
                  + Xác nhận OTP
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
