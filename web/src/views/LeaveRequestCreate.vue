<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex items-center gap-3 border-b border-slate-800 pb-3">
      <router-link to="/cuti" class="text-slate-400 hover:text-white transition-colors">
        <ArrowLeft class="w-6 h-6" />
      </router-link>
      <div>
        <h1 class="text-xl font-bold text-white">Formulir Pengajuan Cuti</h1>
        <p class="text-xs text-slate-400 mt-0.5">Silakan pilih jenis cuti dan lengkapi berkas yang dipersyaratkan.</p>
      </div>
    </div>

    <form @submit.prevent="submitRequest" class="glass-panel p-6 rounded-2xl space-y-6">
      <!-- Error Alert -->
      <div v-if="errorMsg" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl">
        {{ errorMsg }}
      </div>

      <!-- Main Fields -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="sm:col-span-2">
          <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Jenis Cuti*</label>
          <select v-model="form.jenis_cuti" required
                  class="w-full px-3 py-2.5 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
            <option value="Cuti Tahunan Biasa">Cuti Tahunan Biasa</option>
            <option value="Cuti Melahirkan">Cuti Melahirkan</option>
            <option value="Cuti Tahunan untuk Umroh">Cuti Tahunan untuk Umroh</option>
            <option value="Cuti Sakit">Cuti Sakit</option>
            <option value="Cuti Alasan Penting">Cuti Alasan Penting</option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Tanggal Mulai*</label>
          <DatePicker
            v-model="form.tanggal_mulai"
            inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Tanggal Selesai*</label>
          <DatePicker
            v-model="form.tanggal_selesai"
            inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
          />
        </div>
      </div>

      <!-- Dynamic Document Uploader Section -->
      <div class="border-t border-slate-850 pt-4 space-y-4">
        <h3 class="text-xs font-bold text-white uppercase tracking-wider">Unggah Berkas Persyaratan (PDF/JPG/PNG Max 5MB)</h3>
        
        <div class="grid grid-cols-1 gap-4">
          <!-- Required for Tahunan Biasa, Melahirkan, Umroh, Penting -->
          <div v-if="needsFile('surat_rekomendasi_kepsek')">
            <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">
              1. Surat Rekomendasi Kepala Sekolah*
            </label>
            <input type="file" @change="onFileChange($event, 'surat_rekomendasi_kepsek')" accept=".pdf,image/*" required
                   class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
          </div>

          <!-- Required for ALL cuti types -->
          <div v-if="needsFile('sk_terakhir')">
            <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">
              {{ form.jenis_cuti === 'Cuti Sakit' ? '1.' : '2.' }} SK Terakhir*
            </label>
            <input type="file" @change="onFileChange($event, 'sk_terakhir')" accept=".pdf,image/*" required
                   class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
          </div>

          <!-- Melahirkan: HPL, Buku KIA, USG (opt) -->
          <template v-if="form.jenis_cuti === 'Cuti Melahirkan'">
            <div>
              <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan HPL (Puskesmas/RS)*</label>
              <input type="file" @change="onFileChange($event, 'surat_hpl')" accept=".pdf,image/*" required
                     class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
            </div>
            <div>
              <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">4. Buku KIA*</label>
              <input type="file" @change="onFileChange($event, 'buku_kia')" accept=".pdf,image/*" required
                     class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
            </div>
            <div>
              <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">5. Hasil USG (Opsional)</label>
              <input type="file" @change="onFileChange($event, 'hasil_usg')" accept=".pdf,image/*"
                     class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
            </div>
          </template>

          <!-- Umroh: Travel -->
          <div v-if="form.jenis_cuti === 'Cuti Tahunan untuk Umroh'">
            <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan dari Travel Pemberangkatan*</label>
            <input type="file" @change="onFileChange($event, 'surat_travel')" accept=".pdf,image/*" required
                   class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
          </div>

          <!-- Sakit: Rujukan, Rawat Inap -->
          <template v-if="form.jenis_cuti === 'Cuti Sakit'">
            <div>
              <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">2. Surat Rujukan*</label>
              <input type="file" @change="onFileChange($event, 'surat_rujukan')" accept=".pdf,image/*" required
                     class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
            </div>
            <div>
              <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Surat Keterangan Rawat Inap*</label>
              <input type="file" @change="onFileChange($event, 'surat_rawat_inap')" accept=".pdf,image/*" required
                     class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
            </div>
          </template>

          <!-- Alasan Penting: Dokumen Pendukung (wajib) -->
          <div v-if="form.jenis_cuti === 'Cuti Alasan Penting'">
            <label class="block text-[10px] font-semibold text-slate-400 mb-1.5 uppercase">3. Dokumen Pendukung (Istri melahirkan/kematian/KUA/rawat inap keluarga, dll)*</label>
            <input type="file" @change="onFileChange($event, 'dokumen_pendukung')" accept=".pdf,image/*" required
                   class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
          </div>
        </div>
      </div>

      <!-- Action buttons -->
      <div class="border-t border-slate-800 pt-4 flex justify-end gap-3">
        <router-link to="/cuti" class="px-4 py-2 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-xl text-xs font-semibold transition-all">
          Batal
        </router-link>
        <button type="submit" :disabled="submitting"
                class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white rounded-xl text-xs font-bold shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2">
          <Loader2 v-if="submitting" class="w-4 h-4 animate-spin" />
          <span>{{ submitting ? 'Mengirim...' : 'Kirim Pengajuan' }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/services/api'
import { ArrowLeft, Loader2 } from 'lucide-vue-next'
import DatePicker from '@/components/DatePicker.vue'
import Swal from 'sweetalert2'

const router = useRouter()

const form = ref({
  jenis_cuti: 'Cuti Tahunan Biasa',
  tanggal_mulai: '',
  tanggal_selesai: ''
})

const files = ref({
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

const submitting = ref(false)
const errorMsg = ref('')

const needsFile = (field) => {
  const t = form.value.jenis_cuti
  if (field === 'surat_rekomendasi_kepsek') {
    return t !== 'Cuti Sakit'
  }
  if (field === 'sk_terakhir') {
    return true // Required for all
  }
  return false
}

const onFileChange = (e, field) => {
  const file = e.target.files[0]
  if (file) {
    // Validate size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      Swal.fire({
        icon: 'error',
        title: 'File Terlalu Besar',
        text: 'Ukuran file maksimal adalah 5MB.',
        confirmButtonColor: '#1A365D'
      })
      e.target.value = ''
      files.value[field] = null
      return
    }
    files.value[field] = file
  }
}

const submitRequest = async () => {
  submitting.value = true
  errorMsg.value = ''

  const formData = new FormData()
  formData.append('jenis_cuti', form.value.jenis_cuti)
  formData.append('tanggal_mulai', form.value.tanggal_mulai)
  formData.append('tanggal_selesai', form.value.tanggal_selesai)

  // Append relevant files
  for (const [field, file] of Object.entries(files.value)) {
    if (file) {
      formData.append(field, file)
    }
  }

  try {
    await api.post('/leave-requests', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    await Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Pengajuan cuti Anda telah berhasil dikirim.',
      timer: 1550,
      showConfirmButton: false
    })
    router.push('/cuti')
  } catch (error) {
    errorMsg.value = error.response?.data?.error || 'Gagal mengirim pengajuan cuti. Silakan periksa kembali berkas Anda.'
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: errorMsg.value,
      confirmButtonColor: '#1A365D'
    })
  } finally {
    submitting.value = false
  }
}
</script>
