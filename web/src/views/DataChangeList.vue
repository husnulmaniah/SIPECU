<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-xl font-bold text-white">Pengajuan Perubahan Data</h1>
        <p class="text-xs text-slate-400 mt-1">Layanan mandiri (Self-Service) pegawai untuk mengajukan pembaruan data profil resmi dengan bukti berkas SK.</p>
      </div>

      <!-- Submit Request Button (Employee only) -->
      <button
        v-if="authStore.role === 'employee'"
        @click="showCreateModal = true"
        class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2"
      >
        <UserCog class="w-4 h-4" />
        <span>Buat Pengajuan Perubahan</span>
      </button>
    </div>

    <!-- DataTable container -->
    <div class="glass-panel p-4 rounded-2xl" @click="handleTableClick">
      <DataTable
        ref="tableRef"
        :columns="tableColumns"
        ajaxUrl="/api/data-change-requests"
      />
    </div>

    <!-- SUBMIT DATA CHANGE REQUEST MODAL (Employee) -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div class="glass-panel w-full max-w-xl rounded-2xl p-6 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Formulir Koreksi Data Profil</h2>
          <button @click="showCreateModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <form @submit.prevent="submitRequest" class="space-y-4">
          <div v-if="modalError" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl">
            {{ modalError }}
          </div>

          <div class="p-4 bg-slate-950/50 border border-slate-850 rounded-xl space-y-3">
            <h3 class="text-xs font-bold text-slate-300">Pilih Field Yang Ingin Diubah & Isi Nilai Barunya</h3>
            
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <!-- Select field -->
              <div>
                <label class="block text-[10px] font-semibold text-slate-450 uppercase mb-1">Pilih Field Data*</label>
                <select v-model="changeForm.field" required
                        class="w-full px-2 py-1.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                  <option value="nama">Nama Lengkap</option>
                  <option value="jabatan">Jabatan</option>
                  <option value="jenis_jabatan">Jenis Jabatan</option>
                  <option value="tempat_tugas">Tempat Tugas (Unit Kerja)</option>
                  <option value="jenis_tempat">Jenis Tempat (Dinas/Sekolah)</option>
                  <option value="tempat_lahir">Tempat Lahir</option>
                  <option value="tanggal_lahir">Tanggal Lahir</option>
                  <option value="tanggal_kgb_terakhir">Tanggal SK KGB Baru (Memicu histori)</option>
                  <option value="tanggal_kenaikan_pangkat_terakhir">Tanggal SK Naik Pangkat Baru (Memicu histori)</option>
                </select>
              </div>

              <!-- Input new value -->
              <div>
                <label class="block text-[10px] font-semibold text-slate-450 uppercase mb-1">Nilai Baru*</label>
                <!-- Show date input for dates -->
                <DatePicker
                  v-if="changeForm.field.startsWith('tanggal')"
                  v-model="changeForm.value"
                  inputClass="w-full px-2 py-1.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
                />
                <!-- Select for jenis jabatan -->
                <select v-else-if="changeForm.field === 'jenis_jabatan'" v-model="changeForm.value" required
                        class="w-full px-2 py-1.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                  <option value="Fungsional">Fungsional</option>
                  <option value="Pelaksana">Pelaksana</option>
                  <option value="Struktural">Struktural</option>
                </select>
                <!-- Select for jenis tempat -->
                <select v-else-if="changeForm.field === 'jenis_tempat'" v-model="changeForm.value" required
                        class="w-full px-2 py-1.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                  <option value="Dinas">Dinas</option>
                  <option value="Sekolah">Sekolah</option>
                </select>
                <!-- Standard text -->
                <input v-else v-model="changeForm.value" type="text" placeholder="Isi data koreksi baru" required
                       class="w-full px-2 py-1.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
              </div>
            </div>
          </div>

          <!-- Document Upload -->
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Unggah Dokumen Bukti (SK CPNS/Pangkat/KGB Terakhir) - PDF/Image Max 5MB*</label>
            <input type="file" @change="onFileChange" accept=".pdf,image/*" required class="w-full text-xs text-slate-400" />
          </div>

          <!-- Actions -->
          <div class="border-t border-slate-800 pt-3 flex justify-end gap-2">
            <button type="button" @click="showCreateModal = false" :disabled="submitting"
                    class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="submitting"
                    class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="submitting" class="w-3 animate-spin" />
              <span>Kirim Pengajuan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- TINJAU PERUBAHAN DATA MODAL (Admin) -->
    <div v-if="showReviewModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Tinjau Perubahan Data Pegawai</h2>
          <button @click="showReviewModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <div class="bg-slate-950/40 border border-slate-850 p-4 rounded-xl text-xs space-y-3">
          <div class="grid grid-cols-3 gap-1.5">
            <span class="text-slate-500 font-medium">Pegawai</span>
            <span class="col-span-2 text-slate-200 font-bold">{{ activeRequest.employee?.nama }}</span>

            <span class="text-slate-500 font-medium">NIP</span>
            <span class="col-span-2 text-slate-200">{{ activeRequest.employee?.nip }}</span>

            <span class="text-slate-500 font-medium">Koreksi Data</span>
            <span class="col-span-2 text-amber-400 font-mono">{{ parseChanges(activeRequest.data_json) }}</span>
          </div>

          <div class="border-t border-slate-850 pt-2 flex justify-between items-center">
            <span class="text-slate-400 font-medium">Dokumen Bukti SK</span>
            <a :href="activeRequest.sk_terakhir_file" target="_blank"
               class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 rounded font-bold text-[10px]">
              Buka File SK
            </a>
          </div>
        </div>

        <form @submit.prevent="submitReview" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Keputusan*</label>
            <select v-model="reviewForm.status" required
                    class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
              <option value="Disetujui">Setujui (Update otomatis profile pegawai)</option>
              <option value="Dikembalikan">Kembalikan ke Pegawai (Butuh revisi)</option>
              <option value="Ditolak">Tolak Pengajuan</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Catatan</label>
            <textarea v-model="reviewForm.catatan_admin" placeholder="Tuliskan catatan perbaikan atau alasan..." rows="3"
                      class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none resize-none"></textarea>
          </div>

          <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
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
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import DataTable from '@/components/DataTable.vue'
import DatePicker from '@/components/DatePicker.vue'
import { UserCog, X, Loader2 } from 'lucide-vue-next'
import Swal from 'sweetalert2'

const authStore = useAuthStore()
const tableRef = ref(null)

const showCreateModal = ref(false)
const showReviewModal = ref(false)
const submitting = ref(false)
const modalSubmitting = ref(false)
const modalError = ref('')

const changeForm = ref({ field: 'nama', value: '' })
const skFile = ref(null)

const activeRequest = ref({})
const reviewForm = ref({ status: 'Disetujui', catatan_admin: '' })

const tableColumns = computed(() => {
  const cols = [
    { data: 'id', title: 'ID', width: '50px' },
    {
      data: 'data_json',
      title: 'Perubahan Data',
      render: (data) => parseChanges(data)
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
        
        return `<span class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold ${badgeColor}">${data}</span>`
      }
    }
  ]

  if (authStore.role === 'admin') {
    cols.splice(1, 0, { data: 'employee.nama', title: 'Nama Pegawai' })
  }

  cols.push({
    data: null,
    title: 'Aksi / Dokumen',
    orderable: false,
    render: (data, type, row) => {
      let buttons = ''

      if (row.sk_terakhir_file) {
        buttons += `<a href="${row.sk_terakhir_file}" target="_blank" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded border border-slate-700 text-[10px] font-semibold mr-1 inline-block">File Bukti SK</a>`
      }

      if (authStore.role === 'admin' && row.status === 'Diajukan') {
        buttons += `<button data-id="${row.id}" class="btn-review px-2.5 py-1 bg-primary-600 hover:bg-primary-700 text-white rounded text-[10px] font-semibold transition-colors">Tinjau</button>`
      }

      if (row.status === 'Dikembalikan' && row.catatan_admin) {
        buttons += `<span class="text-[10px] text-amber-500 font-semibold italic">Catatan: ${row.catatan_admin}</span>`
      }
      if (row.status === 'Ditolak' && row.catatan_admin) {
        buttons += `<span class="text-[10px] text-red-500 font-semibold italic">Ditolak: ${row.catatan_admin}</span>`
      }

      return buttons || '-'
    }
  })

  return cols
})

const handleTableClick = (e) => {
  const target = e.target
  const idStr = target.getAttribute('data-id')
  if (!idStr) return
  const id = parseInt(idStr)

  if (target.classList.contains('btn-review')) {
    const request = tableRef.value.dt.row(target.closest('tr')).data()
    activeRequest.value = request
    reviewForm.value = { status: 'Disetujui', catatan_admin: '' }
    showReviewModal.value = true
  }
}

const onFileChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    skFile.value = file
  }
}

const submitRequest = async () => {
  submitting.value = true
  modalError.value = ''

  const formData = new FormData()
  // Package changeForm.field -> value inside map
  const changesMap = {}
  changesMap[changeForm.value.field] = changeForm.value.value

  formData.append('data_json', JSON.stringify(changesMap))
  formData.append('sk_terakhir_file', skFile.value)

  try {
    await api.post('/data-change-requests', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showCreateModal.value = false
    changeForm.value = { field: 'nama', value: '' }
    skFile.value = null
    if (tableRef.value) {
      tableRef.value.reload()
    }
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Pengajuan perubahan data berhasil dikirim.',
      timer: 1550,
      showConfirmButton: false,
      background: '#1E293B',
      color: '#F8FAFC'
    })
  } catch (error) {
    modalError.value = error.response?.data?.error || 'Gagal mengirim pengajuan'
  } finally {
    submitting.value = false
  }
}

const submitReview = async () => {
  modalSubmitting.value = true
  try {
    await api.put(`/data-change-requests/${activeRequest.value.id}/status`, {
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
      text: 'Keputusan perubahan data berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan keputusan',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    modalSubmitting.value = false
  }
}

const parseChanges = (jsonStr) => {
  if (!jsonStr) return ''
  try {
    const obj = JSON.parse(jsonStr)
    const entries = Object.entries(obj)
    if (entries.length === 0) return ''
    
    // Map database field keys to user-friendly labels in Indonesian
    const labelMap = {
      nama: 'Nama Lengkap',
      jabatan: 'Jabatan',
      jenis_jabatan: 'Jenis Jabatan',
      tempat_tugas: 'Tempat Tugas (Unit)',
      jenis_tempat: 'Jenis Unit',
      tempat_lahir: 'Tempat Lahir',
      tanggal_lahir: 'Tanggal Lahir',
      tanggal_kgb_terakhir: 'SK KGB Baru',
      tanggal_kenaikan_pangkat_terakhir: 'SK Pangkat Baru'
    }

    const [key, val] = entries[0]
    const label = labelMap[key] || key
    
    // Format date if val is a date string
    let displayVal = val
    if (key.startsWith('tanggal')) {
      displayVal = new Date(val).toLocaleDateString('id-ID')
    }

    return `${label} → ${displayVal}`
  } catch (e) {
    return jsonStr
  }
}
</script>
