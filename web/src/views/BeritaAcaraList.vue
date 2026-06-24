<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-xl font-bold text-white">Dokumen Berita Acara</h1>
        <p class="text-xs text-slate-400 mt-1">Layanan pengunggahan berkas pendukung kedinasan pegawai (BA, ST, SI, SKS) untuk ditinjau Admin.</p>
      </div>

      <!-- Upload BA Button (Employee only) -->
      <button
        v-if="authStore.role === 'employee'"
        @click="showCreateModal = true"
        class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2"
      >
        <FileUp class="w-4 h-4" />
        <span>Unggah Dokumen Baru</span>
      </button>
    </div>

    <!-- DataTable container -->
    <div class="glass-panel p-4 rounded-2xl" @click="handleTableClick">
      <DataTable
        ref="tableRef"
        :columns="tableColumns"
        ajaxUrl="/api/berita-acara"
      />
    </div>

    <!-- SUBMIT BERITA ACARA MODAL (Employee) -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Unggah Dokumen Kedinasan</h2>
          <button @click="showCreateModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <form @submit.prevent="submitRequest" class="space-y-4">
          <div v-if="modalError" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl">
            {{ modalError }}
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Kategori Dokumen*</label>
            <select v-model="form.jenis" required
                    class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
              <option value="BA">BA (Berita Acara)</option>
              <option value="ST">ST (Surat Tugas)</option>
              <option value="SI">SI (Surat Izin)</option>
              <option value="SKS">SKS (Surat Keterangan Sakit)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Berkas Dokumen (PDF/JPG/PNG Max 5MB)*</label>
            <input type="file" @change="onFileChange" accept=".pdf,image/*" required class="w-full text-xs text-slate-400" />
          </div>

          <div class="flex justify-end gap-2 border-t border-slate-800 pt-3">
            <button type="button" @click="showCreateModal = false" :disabled="submitting"
                    class="px-3 py-1.5 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-lg text-xs font-semibold">Batal</button>
            <button type="submit" :disabled="submitting"
                    class="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-xs font-semibold flex items-center gap-1">
              <Loader2 v-if="submitting" class="w-3 animate-spin" />
              <span>Unggah Berkas</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- TINJAU BERITA ACARA MODAL (Admin) -->
    <div v-if="showReviewModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-2">
          <h2 class="text-sm font-bold text-white">Tinjau Dokumen Berita Acara</h2>
          <button @click="showReviewModal = false" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <div class="bg-slate-950/40 border border-slate-850 p-4 rounded-xl text-xs space-y-2.5">
          <div class="grid grid-cols-3 gap-1.5">
            <span class="text-slate-500 font-medium">Pegawai</span>
            <span class="col-span-2 text-slate-200 font-bold">{{ activeRequest.employee?.nama }}</span>

            <span class="text-slate-500 font-medium">NIP</span>
            <span class="col-span-2 text-slate-200">{{ activeRequest.employee?.nip }}</span>

            <span class="text-slate-500 font-medium">Jenis Berkas</span>
            <span class="col-span-2 text-violet-400 font-bold">{{ activeRequest.jenis }}</span>
          </div>

          <div class="border-t border-slate-850 pt-2 flex justify-between items-center">
            <span class="text-slate-400 font-medium">File Terlampir</span>
            <a :href="activeRequest.file_path" target="_blank"
               class="px-2.5 py-0.5 bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 rounded font-bold text-[10px]">
              Unduh / Lihat File
            </a>
          </div>
        </div>

        <form @submit.prevent="submitReview" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Keputusan*</label>
            <select v-model="reviewForm.status" required
                    class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
              <option value="Disetujui">Setujui & Arsipkan</option>
              <option value="Dikembalikan">Kembalikan ke Pegawai (Kekurangan berkas)</option>
              <option value="Ditolak">Tolak Dokumen</option>
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
              <span>Simpan Keputusan</span>
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
import { FileUp, X, Loader2 } from 'lucide-vue-next'
import Swal from 'sweetalert2'

const authStore = useAuthStore()
const tableRef = ref(null)

const showCreateModal = ref(false)
const showReviewModal = ref(false)
const submitting = ref(false)
const modalSubmitting = ref(false)
const modalError = ref('')

const form = ref({ jenis: 'BA' })
const baFile = ref(null)

const activeRequest = ref({})
const reviewForm = ref({ status: 'Disetujui', catatan_admin: '' })

const tableColumns = computed(() => {
  const cols = [
    { data: 'id', title: 'ID', width: '50px' },
    { data: 'jenis', title: 'Jenis Dokumen' },
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
    title: 'Aksi / Berkas',
    orderable: false,
    render: (data, type, row) => {
      let buttons = ''

      if (row.file_path) {
        buttons += `<a href="${row.file_path}" target="_blank" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded border border-slate-700 text-[10px] font-semibold mr-1 inline-block">Buka File</a>`
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
    baFile.value = file
  }
}

const submitRequest = async () => {
  submitting.value = true
  modalError.value = ''

  const formData = new FormData()
  formData.append('jenis', form.value.jenis)
  formData.append('file_path', baFile.value)

  try {
    await api.post('/berita-acara', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    showCreateModal.value = false
    form.value = { jenis: 'BA' }
    baFile.value = null
    if (tableRef.value) {
      tableRef.value.reload()
    }
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Berkas Berita Acara berhasil diunggah.',
      timer: 1550,
      showConfirmButton: false,
      background: '#1E293B',
      color: '#F8FAFC'
    })
  } catch (error) {
    modalError.value = error.response?.data?.error || 'Gagal mengunggah berkas'
  } finally {
    submitting.value = false
  }
}

const submitReview = async () => {
  modalSubmitting.value = true
  try {
    await api.put(`/berita-acara/${activeRequest.value.id}/status`, {
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
      text: 'Keputusan berita acara berhasil disimpan.',
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
</script>
