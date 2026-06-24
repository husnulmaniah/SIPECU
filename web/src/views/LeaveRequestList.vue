<template>
  <div class="space-y-6">
  <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-xl font-bold text-white">Daftar Pengajuan Cuti</h1>
        <p class="text-xs text-slate-400 mt-1">Daftar pengajuan surat rekomendasi cuti dan status persetujuannya.</p>
      </div>

      <div class="flex gap-2">
        <!-- Admin: Ajukan Cuti untuk Pegawai -->
        <button
          v-if="authStore.role === 'admin'"
          @click="openAdminLeaveModal"
          class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2"
        >
          <CalendarDays class="w-4 h-4" />
          <span>Ajukan Cuti untuk Pegawai</span>
        </button>

        <!-- Employee: Buat Pengajuan Cuti -->
        <router-link
          v-if="authStore.role === 'employee'"
          to="/cuti/tambah"
          class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2"
        >
          <CalendarDays class="w-4 h-4" />
          <span>Buat Pengajuan Cuti</span>
        </router-link>
      </div>
    </div>

    <!-- DataTable container -->
    <div class="glass-panel p-4 rounded-2xl" @click="handleTableClick">
      <DataTable
        ref="tableRef"
        :columns="tableColumns"
        ajaxUrl="/api/leave-requests"
      />
    </div>

    <!-- TINJAU CUTI MODAL (Admin) -->
    <div v-if="showReviewModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div class="glass-panel w-full max-w-xl rounded-2xl p-6 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Tinjau Pengajuan Cuti</h2>
          <button @click="showReviewModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <!-- Request Details -->
        <div class="bg-slate-950/30 border border-slate-850 p-4 rounded-xl text-xs space-y-3">
          <div class="grid grid-cols-3 gap-2">
            <span class="text-slate-500 font-semibold">Nama Pegawai</span>
            <span class="col-span-2 text-slate-200">{{ activeRequest.employee?.nama }}</span>

            <span class="text-slate-500 font-semibold">NIP Pegawai</span>
            <span class="col-span-2 text-slate-200">{{ activeRequest.employee?.nip }}</span>

            <span class="text-slate-500 font-semibold">Jenis Cuti</span>
            <span class="col-span-2 text-primary-400 font-bold">{{ activeRequest.jenis_cuti }}</span>

            <span class="text-slate-500 font-semibold">Rentang Cuti</span>
            <span class="col-span-2 text-slate-250 font-semibold">{{ formatDate(activeRequest.tanggal_mulai) }} s/d {{ formatDate(activeRequest.tanggal_selesai) }}</span>
          </div>

          <!-- Attachments list -->
          <div class="border-t border-slate-850 pt-2">
            <p class="font-semibold text-slate-400 mb-2">Dokumen Lampiran Pendukung:</p>
            <div v-if="activeRequest.attachments?.length === 0" class="text-slate-600">Tidak ada lampiran terunggah.</div>
            <div v-else class="space-y-1.5">
              <div v-for="att in activeRequest.attachments" :key="att.id" class="flex justify-between items-center bg-slate-900/60 p-2 rounded border border-slate-850">
                <span class="text-slate-350">{{ att.jenis_dokumen }}</span>
                <a :href="att.file_path" target="_blank" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[10px] font-bold">Buka</a>
              </div>
            </div>
          </div>
        </div>

        <!-- Action form -->
        <form @submit.prevent="submitReview" class="space-y-4 pt-2">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Keputusan Tinjauan*</label>
            <select v-model="reviewForm.status" required
                    class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
              <option value="Disetujui">Setujui (Generate rekomendasi otomatis)</option>
              <option value="Dikembalikan">Kembalikan ke Pegawai (Butuh revisi berkas)</option>
              <option value="Ditolak">Tolak Pengajuan</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase font-medium">Catatan Review / Alasan</label>
            <textarea v-model="reviewForm.catatan_admin" placeholder="Tuliskan catatan perbaikan atau detail keputusan..." rows="3"
                      class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none resize-none"></textarea>
          </div>

          <div class="flex justify-end gap-2 border-t border-slate-800 pt-3">
            <button type="button" @click="showReviewModal = false" :disabled="modalSubmitting"
                    class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting"
                    class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>Simpan Tinjauan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- UPLOAD ACC MODAL (Admin) -->
    <div v-if="showUploadAccModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Unggah Surat Rekomendasi Hasil ACC</h2>
          <button @click="showUploadAccModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <form @submit.prevent="submitUploadAcc" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1">Pilih File Hasil ACC (PDF saja, Max 5MB)*</label>
            <input type="file" @change="onAccFileChange" accept=".pdf" required class="w-full text-xs text-slate-400" />
          </div>
          <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
            <button type="button" @click="showUploadAccModal = false" :disabled="modalSubmitting"
                    class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting"
                    class="px-3 py-1.5 bg-violet-600 hover:bg-violet-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>Unggah & Selesaikan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ADMIN AJUKAN CUTI MODAL -->
    <div v-if="showAdminLeaveModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div class="glass-panel w-full max-w-lg rounded-2xl p-6 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <div>
            <h2 class="text-sm font-bold text-white">Ajukan Cuti untuk Pegawai</h2>
            <p class="text-xs text-slate-400 mt-0.5">Admin mengajukan cuti atas nama pegawai yang bersangkutan.</p>
          </div>
          <button @click="showAdminLeaveModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <div v-if="adminLeaveError" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl">
          {{ adminLeaveError }}
        </div>

        <form @submit.prevent="submitAdminLeave" class="space-y-4">
          <!-- Pilih Pegawai (Searchable) -->
          <div class="relative" ref="dropdownContainer">
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Pilih Pegawai*</label>
            <div class="relative">
              <input
                type="text"
                v-model="employeeSearchQuery"
                @focus="showEmployeeDropdown = true"
                @input="onSearchInput"
                placeholder="Cari berdasarkan nama atau NIP..."
                class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none"
                required
              />
              <button
                type="button"
                @click.stop="showEmployeeDropdown = !showEmployeeDropdown"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
              >
                <ChevronDown class="w-4 h-4" />
              </button>
            </div>
            
            <!-- Dropdown Options -->
            <div
              v-if="showEmployeeDropdown"
              class="absolute z-50 w-full mt-1 max-h-60 overflow-y-auto bg-slate-900/95 border border-slate-850 rounded-lg shadow-xl backdrop-blur-md"
            >
              <div
                v-if="filteredEmployees.length === 0"
                class="px-3 py-2 text-xs text-slate-500"
              >
                Pegawai tidak ditemukan
              </div>
              <div
                v-else
                v-for="emp in filteredEmployees"
                :key="emp.nip"
                @click="selectEmployee(emp)"
                class="px-3 py-2 text-xs text-slate-300 hover:bg-primary-600/30 hover:text-white cursor-pointer border-b border-slate-850 last:border-0"
              >
                <span class="font-semibold">{{ emp.nama }}</span>
                <span class="text-slate-500 ml-1">({{ emp.nip }})</span>
              </div>
            </div>
          </div>

          <!-- Jenis Cuti -->
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Jenis Cuti*</label>
            <select v-model="adminLeaveForm.jenis_cuti" required
                    class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
              <option value="Cuti Tahunan Biasa">Cuti Tahunan Biasa</option>
              <option value="Cuti Melahirkan">Cuti Melahirkan</option>
              <option value="Cuti Tahunan untuk Umroh">Cuti Tahunan untuk Umroh</option>
              <option value="Cuti Sakit">Cuti Sakit</option>
              <option value="Cuti Alasan Penting">Cuti Alasan Penting</option>
            </select>
          </div>

          <!-- Rentang Tanggal -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tanggal Mulai*</label>
              <DatePicker
                v-model="adminLeaveForm.tanggal_mulai"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tanggal Selesai*</label>
              <DatePicker
                v-model="adminLeaveForm.tanggal_selesai"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
          </div>

          <!-- Dynamic Document Uploader Section (Admin) -->
          <div class="border-t border-slate-850 pt-4 space-y-4">
            <h3 class="text-xs font-bold text-white uppercase tracking-wider">Unggah Berkas Persyaratan (PDF/JPG/PNG Max 5MB)</h3>
            
            <div class="grid grid-cols-1 gap-4">
              <!-- Required for Tahunan Biasa, Melahirkan, Umroh, Penting -->
              <div v-if="adminNeedsFile('surat_rekomendasi_kepsek')">
                <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">
                  1. Surat Rekomendasi Kepala Sekolah*
                </label>
                <input type="file" @change="onAdminFileChange($event, 'surat_rekomendasi_kepsek')" accept=".pdf,image/*" required
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>

              <!-- Required for ALL cuti types -->
              <div v-if="adminNeedsFile('sk_terakhir')">
                <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">
                  {{ adminLeaveForm.jenis_cuti === 'Cuti Sakit' ? '1.' : '2.' }} SK Terakhir*
                </label>
                <input type="file" @change="onAdminFileChange($event, 'sk_terakhir')" accept=".pdf,image/*" required
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>

              <!-- Melahirkan: HPL, Buku KIA, USG (opt) -->
              <template v-if="adminLeaveForm.jenis_cuti === 'Cuti Melahirkan'">
                <div>
                  <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan HPL (Puskesmas/RS)*</label>
                  <input type="file" @change="onAdminFileChange($event, 'surat_hpl')" accept=".pdf,image/*" required
                         class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
                </div>
                <div>
                  <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">4. Buku KIA*</label>
                  <input type="file" @change="onAdminFileChange($event, 'buku_kia')" accept=".pdf,image/*" required
                         class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
                </div>
                <div>
                  <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">5. Hasil USG (Opsional)</label>
                  <input type="file" @change="onAdminFileChange($event, 'hasil_usg')" accept=".pdf,image/*"
                         class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
                </div>
              </template>

              <!-- Umroh: Travel -->
              <div v-if="adminLeaveForm.jenis_cuti === 'Cuti Tahunan untuk Umroh'">
                <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan dari Travel Pemberangkatan*</label>
                <input type="file" @change="onAdminFileChange($event, 'surat_travel')" accept=".pdf,image/*" required
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>

              <!-- Sakit: Rujukan, Rawat Inap -->
              <template v-if="adminLeaveForm.jenis_cuti === 'Cuti Sakit'">
                <div>
                  <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">2. Surat Rujukan*</label>
                  <input type="file" @change="onAdminFileChange($event, 'surat_rujukan')" accept=".pdf,image/*" required
                         class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
                </div>
                <div>
                  <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan Rawat Inap*</label>
                  <input type="file" @change="onAdminFileChange($event, 'surat_rawat_inap')" accept=".pdf,image/*" required
                         class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
                </div>
              </template>

              <!-- Alasan Penting: Dokumen Pendukung (wajib) -->
              <div v-if="adminLeaveForm.jenis_cuti === 'Cuti Alasan Penting'">
                <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Dokumen Pendukung (Istri melahirkan/kematian/KUA/rawat inap keluarga, dll)*</label>
                <input type="file" @change="onAdminFileChange($event, 'dokumen_pendukung')" accept=".pdf,image/*" required
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
            <button type="button" @click="showAdminLeaveModal = false" :disabled="modalSubmitting"
                    class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="modalSubmitting"
                    class="px-4 py-1.5 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="modalSubmitting" class="w-3 animate-spin" />
              <span>{{ modalSubmitting ? 'Menyimpan...' : 'Ajukan Cuti' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import DataTable from '@/components/DataTable.vue'
import DatePicker from '@/components/DatePicker.vue'
import { CalendarDays, X, Loader2, ChevronDown } from 'lucide-vue-next'
import Swal from 'sweetalert2'

const authStore = useAuthStore()
const tableRef = ref(null)

const showReviewModal = ref(false)
const showUploadAccModal = ref(false)
const showAdminLeaveModal = ref(false)
const modalSubmitting = ref(false)

const activeRequest = ref({})
const reviewForm = ref({ status: 'Disetujui', catatan_admin: '' })
const accFile = ref(null)

// Admin ajukan cuti untuk pegawai
const employeeList = ref([])
const adminLeaveError = ref('')
const adminLeaveForm = ref({
  employee_nip: '',
  jenis_cuti: 'Cuti Tahunan Biasa',
  tanggal_mulai: '',
  tanggal_selesai: ''
})

const employeeSearchQuery = ref('')
const showEmployeeDropdown = ref(false)
const dropdownContainer = ref(null)

const adminFiles = ref({
  surat_rekomendasi_kepsek: null,
  sk_terakhir: null,
  surat_hpl: null,
  buku_kia: null,
  hasil_usg: null,
  surat_travel: null,
  surat_rujukan: null,
  surat_rawat_inap: null,
  dokumen_pendukung: null
})

const filteredEmployees = computed(() => {
  const query = employeeSearchQuery.value.toLowerCase().trim()
  if (!query) return employeeList.value
  return employeeList.value.filter(emp =>
    emp.nama.toLowerCase().includes(query) ||
    emp.nip.toLowerCase().includes(query)
  )
})

const selectEmployee = (emp) => {
  adminLeaveForm.value.employee_nip = emp.nip
  employeeSearchQuery.value = `${emp.nama} (${emp.nip})`
  showEmployeeDropdown.value = false
}

const onSearchInput = () => {
  const currentNip = adminLeaveForm.value.employee_nip
  const selectedEmp = employeeList.value.find(emp => emp.nip === currentNip)
  if (selectedEmp) {
    const expectedQuery = `${selectedEmp.nama} (${selectedEmp.nip})`
    if (employeeSearchQuery.value !== expectedQuery) {
      adminLeaveForm.value.employee_nip = ''
    }
  }
  showEmployeeDropdown.value = true
}

const adminNeedsFile = (field) => {
  const t = adminLeaveForm.value.jenis_cuti
  if (field === 'surat_rekomendasi_kepsek') {
    return t !== 'Cuti Sakit'
  }
  if (field === 'sk_terakhir') {
    return true
  }
  return false
}

const onAdminFileChange = (e, field) => {
  const file = e.target.files[0]
  if (file) {
    if (file.size > 5 * 1024 * 1024) {
      Swal.fire({
        icon: 'error',
        title: 'File Terlalu Besar',
        text: 'Ukuran file maksimal adalah 5MB.',
        confirmButtonColor: '#1A365D'
      })
      e.target.value = ''
      adminFiles.value[field] = null
      return
    }
    adminFiles.value[field] = file
  }
}

const handleClickOutside = (e) => {
  if (dropdownContainer.value && !dropdownContainer.value.contains(e.target)) {
    showEmployeeDropdown.value = false
  }
}

const fetchEmployeeList = async () => {
  try {
    const res = await api.get('/employees?all=true')
    employeeList.value = res.data.data || []
  } catch (e) {
    console.warn('Gagal memuat daftar pegawai:', e)
  }
}

const openAdminLeaveModal = async () => {
  adminLeaveError.value = ''
  adminLeaveForm.value = { employee_nip: '', jenis_cuti: 'Cuti Tahunan Biasa', tanggal_mulai: '', tanggal_selesai: '' }
  employeeSearchQuery.value = ''
  showEmployeeDropdown.value = false
  for (const key in adminFiles.value) {
    adminFiles.value[key] = null
  }
  if (employeeList.value.length === 0) await fetchEmployeeList()
  showAdminLeaveModal.value = true
}

const submitAdminLeave = async () => {
  adminLeaveError.value = ''
  if (!adminLeaveForm.value.employee_nip) {
    adminLeaveError.value = 'Pegawai wajib dipilih'
    return
  }
  modalSubmitting.value = true
  const formData = new FormData()
  formData.append('employee_nip', adminLeaveForm.value.employee_nip)
  formData.append('jenis_cuti', adminLeaveForm.value.jenis_cuti)
  formData.append('tanggal_mulai', adminLeaveForm.value.tanggal_mulai)
  formData.append('tanggal_selesai', adminLeaveForm.value.tanggal_selesai)

  for (const [field, file] of Object.entries(adminFiles.value)) {
    if (file) {
      formData.append(field, file)
    }
  }

  try {
    await api.post('/leave-requests', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    showAdminLeaveModal.value = false
    if (tableRef.value) tableRef.value.reload()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Pengajuan cuti pegawai berhasil dibuat.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (err) {
    adminLeaveError.value = err.response?.data?.error || 'Gagal mengajukan cuti'
  } finally {
    modalSubmitting.value = false
  }
}

onMounted(() => {
  if (authStore.role === 'admin') fetchEmployeeList()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const tableColumns = computed(() => {
  const cols = [
    { data: 'id', title: 'ID', width: '50px' },
    { data: 'jenis_cuti', title: 'Jenis Cuti' },
    {
      data: 'tanggal_mulai',
      title: 'Mulai',
      render: (data) => new Date(data).toLocaleDateString('id-ID')
    },
    {
      data: 'tanggal_selesai',
      title: 'Selesai',
      render: (data) => new Date(data).toLocaleDateString('id-ID')
    },
    {
      data: 'status',
      title: 'Status',
      render: (data) => {
        let badgeColor = 'bg-slate-800 text-slate-400'
        if (data === 'Diajukan') badgeColor = 'bg-blue-500/10 text-blue-400 border border-blue-500/20'
        if (data === 'Disetujui') badgeColor = 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
        if (data === 'Dikembalikan') badgeColor = 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
        if (data === 'Ditolak') badgeColor = 'bg-red-500/10 text-red-400 border border-red-500/20'
        if (data === 'Surat Terunggah') badgeColor = 'bg-violet-500/10 text-violet-400 border border-violet-500/20 font-bold'
        
        return `<span class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold ${badgeColor}">${data}</span>`
      }
    }
  ]

  // Add Employee name if Admin
  if (authStore.role === 'admin') {
    cols.splice(1, 0, { data: 'employee.nama', title: 'Nama Pegawai' })
  }

  // Action column
  cols.push({
    data: null,
    title: 'Aksi / Dokumen',
    orderable: false,
    render: (data, type, row) => {
      let buttons = ''

      // 1. Admin Actions
      if (authStore.role === 'admin') {
        if (row.status === 'Diajukan') {
          const rowJson = JSON.stringify(row).replace(/"/g, '&quot;')
          buttons += `<button data-id="${row.id}" data-row="${rowJson}" class="btn-review px-2.5 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-[10px] font-semibold transition-colors mr-1">Tinjau</button>`
        }
        if (row.status === 'Disetujui') {
          const rowJson = JSON.stringify(row).replace(/"/g, '&quot;')
          buttons += `<button data-id="${row.id}" data-row="${rowJson}" class="btn-upload px-2.5 py-1 bg-violet-600 hover:bg-violet-700 text-white rounded text-[10px] font-semibold transition-colors mr-1">Upload ACC</button>`
        }
        if (['Disetujui', 'Surat Terunggah', 'Selesai'].includes(row.status)) {
          buttons += `<button data-id="${row.id}" class="btn-regenerate px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white rounded text-[10px] font-semibold transition-colors mr-1">🔄 Regenerasi</button>`
        }
      }

      // 2. Document Downloads (Available to both Admin and Employee based on status)
      if (row.letters && row.letters.length > 0) {
        row.letters.forEach(letter => {
          const label = letter.jenis_surat === 'Formulir' ? 'Formulir' : 'Rekomendasi'
          const colorPdf = letter.jenis_surat === 'Formulir' ? 'bg-sky-700 hover:bg-sky-600' : 'bg-indigo-800 hover:bg-indigo-700'

          // Formulir: hanya PDF (tidak ada Word)
          // Rekomendasi: Word + PDF
          if (letter.jenis_surat === 'Rekomendasi' && letter.file_docx) {
            buttons += `<a href="${letter.file_docx}" target="_blank" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded border border-slate-700 text-[10px] font-medium mr-1 inline-block">Rekomendasi Word</a>`
          }
          if (letter.file_pdf) {
            buttons += `<a href="${letter.file_pdf}" target="_blank" class="px-2 py-1 ${colorPdf} text-slate-200 rounded border border-slate-700 text-[10px] font-medium mr-1 inline-block">${label} PDF</a>`
          }
          if (letter.jenis_surat === 'Rekomendasi' && letter.file_signed_pdf) {
            buttons += `<a href="${letter.file_signed_pdf}" target="_blank" class="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded text-[10px] font-bold mr-1 inline-block shadow-md">Rekomendasi ACC</a>`
          }
        })
      }

      if (row.status === 'Dikembalikan' && row.catatan_admin) {
        buttons += `<span class="text-[10px] text-amber-500 font-semibold italic">Catatan: ${row.catatan_admin}</span>`
      }
      if (row.status === 'Ditolak') {
        buttons += `<span class="text-[10px] text-red-500 font-semibold italic">Ditolak</span>`
      }

      if (buttons === '') {
        buttons = '<span class="text-[10px] text-slate-600">-</span>'
      }

      return buttons
    }
  })

  return cols
})

// Handle action button delegation clicks
const handleTableClick = (e) => {
  const target = e.target
  const idStr = target.getAttribute('data-id')
  if (!idStr) return

  // Parse the row data embedded in the button as JSON
  const parseRowData = (el) => {
    try {
      const raw = el.getAttribute('data-row')
      return raw ? JSON.parse(raw.replace(/&quot;/g, '"')) : null
    } catch (err) {
      console.error('Failed to parse row data:', err)
      return null
    }
  }

  if (target.classList.contains('btn-review')) {
    const request = parseRowData(target)
    if (!request) return
    activeRequest.value = request
    reviewForm.value = { status: 'Disetujui', catatan_admin: '' }
    showReviewModal.value = true
  }

  if (target.classList.contains('btn-upload')) {
    const request = parseRowData(target)
    if (!request) return
    activeRequest.value = request
    accFile.value = null
    showUploadAccModal.value = true
  }

  if (target.classList.contains('btn-regenerate')) {
    const id = parseInt(idStr)
    Swal.fire({
      title: 'Regenerasi Dokumen?',
      text: 'Dokumen Formulir Cuti dan Surat Rekomendasi akan dibuat ulang dengan format terbaru.',
      icon: 'question',
      showCancelButton: true,
      confirmButtonText: 'Ya, Regenerasi',
      cancelButtonText: 'Batal',
      confirmButtonColor: '#d97706',
    }).then(async (result) => {
      if (!result.isConfirmed) return
      try {
        await api.post(`/leave-requests/${id}/regenerate-docs`)
        Swal.fire({ icon: 'success', title: 'Berhasil', text: 'Dokumen berhasil di-regenerasi. Halaman akan diperbarui.', timer: 2000, showConfirmButton: false })
        if (tableRef.value) tableRef.value.reload()
      } catch (err) {
        Swal.fire({ icon: 'error', title: 'Gagal', text: err?.response?.data?.error || 'Terjadi kesalahan saat regenerasi dokumen.' })
      }
    })
  }
}

const onAccFileChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    accFile.value = file
  }
}

const submitReview = async () => {
  modalSubmitting.value = true
  try {
    await api.put(`/leave-requests/${activeRequest.value.id}/status`, {
      status: reviewForm.value.status,
      catatan_admin: reviewForm.value.catatan_admin
    })
    showReviewModal.value = false
    if (tableRef.value) {
      tableRef.value.reload()
    }
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Tinjauan pengajuan cuti berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan tinjauan',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const submitUploadAcc = async () => {
  if (!accFile.value) return
  modalSubmitting.value = true

  const formData = new FormData()
  formData.append('file_signed_pdf', accFile.value)

  try {
    await api.post(`/leave-requests/${activeRequest.value.id}/letters`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showUploadAccModal.value = false
    if (tableRef.value) {
      tableRef.value.reload()
    }
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'File ACC surat cuti berhasil diunggah.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal mengunggah file ACC',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>
