<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <div>
        <h1 class="text-xl font-bold text-white">Data Kepegawaian Aktif</h1>
        <p class="text-xs text-slate-400 mt-1">Daftar seluruh pegawai aktif yang terdaftar di sistem.</p>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-wrap gap-2">
        <!-- Download Template -->
        <a :href="templateUrl" target="_blank"
           class="px-3 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 font-bold text-xs rounded-xl border border-slate-700 flex items-center gap-2 transition-all">
          <FileDown class="w-4 h-4 text-emerald-400" />
          <span>Template</span>
        </a>

        <!-- Import Data -->
        <button @click="showImportModal = true"
                class="px-3 py-2 bg-emerald-600/20 hover:bg-emerald-600/30 text-emerald-400 font-bold text-xs rounded-xl border border-emerald-500/30 flex items-center gap-2 transition-all">
          <FileUp class="w-4 h-4" />
          <span>Import Excel</span>
        </button>

        <!-- Export Data -->
        <a :href="exportUrl" target="_blank"
           class="px-3 py-2 bg-violet-600/20 hover:bg-violet-600/30 text-violet-400 font-bold text-xs rounded-xl border border-violet-500/30 flex items-center gap-2 transition-all">
          <Download class="w-4 h-4" />
          <span>Export Excel</span>
        </a>

        <!-- Tambah Pegawai -->
        <button @click="openCreateModal"
                class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-primary-500/20 hover:shadow-primary-500/30 transition-all flex items-center gap-2">
          <UserPlus class="w-4 h-4" />
          <span>Tambah Pegawai</span>
        </button>

        <!-- Hapus Semua -->
        <button @click="confirmDeleteAll"
                class="px-3 py-2 bg-rose-600/20 hover:bg-rose-600/30 text-rose-400 font-bold text-xs rounded-xl border border-rose-500/30 flex items-center gap-2 transition-all">
          <Trash2 class="w-4 h-4" />
          <span>Hapus Semua</span>
        </button>
      </div>
    </div>

    <!-- Active filter alert -->
    <div v-if="activeFilter" class="p-3 bg-amber-500/10 border border-amber-500/20 text-amber-400 text-xs rounded-xl flex items-center justify-between">
      <div class="flex items-center gap-2">
        <AlertTriangle class="w-4 h-4 shrink-0" />
        <span>Menampilkan filter: <strong>{{ getFilterLabel(activeFilter) }}</strong> (Rentang target 3 bulan ke depan).</span>
      </div>
      <button @click="clearFilter" class="text-xs font-bold underline hover:text-white">Bersihkan Filter</button>
    </div>

    <!-- DataTable container -->
    <div class="glass-panel p-4 rounded-2xl" @click="handleTableClick">
      <DataTable
        ref="tableRef"
        :columns="tableColumns"
        ajaxUrl="/api/employees"
        :extraParams="extraQueryParams"
      />
    </div>

    <!-- ===== IMPORT MODAL ===== -->
    <div v-if="showImportModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-2xl rounded-2xl p-6 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div>
            <h2 class="text-base font-bold text-white">Import Data Pegawai</h2>
            <p class="text-xs text-slate-400 mt-1">Upload file Excel (.xlsx) untuk import bulk data pegawai.</p>
          </div>
          <button @click="closeImportModal" class="text-slate-400 hover:text-white"><X class="w-5 h-5" /></button>
        </div>

        <!-- Upload Area -->
        <div v-if="!importResult"
             @dragover.prevent="dragOver = true"
             @dragleave.prevent="dragOver = false"
             @drop.prevent="onImportDrop"
             :class="['border-2 border-dashed rounded-xl p-8 text-center transition-all cursor-pointer',
                      dragOver ? 'border-primary-500 bg-primary-500/10' : 'border-slate-700 hover:border-slate-600']"
             @click="$refs.importFileInput.click()">
          <FileUp :class="['w-10 h-10 mx-auto mb-3', importFile ? 'text-emerald-400' : 'text-slate-500']" />
          <p class="text-sm font-semibold text-slate-300">
            {{ importFile ? importFile.name : 'Klik atau seret file Excel ke sini' }}
          </p>
          <p class="text-xs text-slate-500 mt-1">Format: .xlsx · Maks. 10MB</p>
          <input ref="importFileInput" type="file" accept=".xlsx" class="hidden" @change="onImportFileSelect" />
        </div>

        <!-- Import Result -->
        <div v-if="importResult" class="space-y-3">
          <div :class="['p-3 rounded-xl text-xs font-semibold',
                        importResult.fail_count === 0
                          ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                          : 'bg-amber-500/10 border border-amber-500/20 text-amber-400']">
            {{ importResult.message }}
          </div>
          <!-- Detail table -->
          <div class="max-h-64 overflow-y-auto rounded-xl border border-slate-800">
            <table class="w-full text-xs">
              <thead class="bg-slate-900 sticky top-0">
                <tr>
                  <th class="p-2 text-left text-slate-400">Baris</th>
                  <th class="p-2 text-left text-slate-400">NIP</th>
                  <th class="p-2 text-left text-slate-400">Nama</th>
                  <th class="p-2 text-left text-slate-400">Status</th>
                  <th class="p-2 text-left text-slate-400">Keterangan</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="detail in importResult.details" :key="detail.row"
                    :class="['border-t border-slate-800', detail.status === 'Gagal' ? 'bg-red-500/5' : '']">
                  <td class="p-2 text-slate-400">{{ detail.row }}</td>
                  <td class="p-2 text-slate-300">{{ detail.nip }}</td>
                  <td class="p-2 text-slate-300">{{ detail.nama }}</td>
                  <td class="p-2">
                    <span :class="['px-1.5 py-0.5 rounded text-[10px] font-bold',
                                   detail.status === 'Gagal' ? 'bg-red-500/20 text-red-400' :
                                   detail.status === 'Skip' ? 'bg-slate-500/20 text-slate-400' :
                                   'bg-emerald-500/20 text-emerald-400']">
                      {{ detail.status }}
                    </span>
                  </td>
                  <td class="p-2 text-slate-500">{{ detail.message }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
          <button @click="closeImportModal" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-semibold">
            {{ importResult ? 'Tutup' : 'Batal' }}
          </button>
          <button v-if="!importResult" @click="submitImport" :disabled="!importFile || importing"
                  class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white rounded-lg text-xs font-semibold flex items-center gap-1.5">
            <Loader2 v-if="importing" class="w-3.5 h-3.5 animate-spin" />
            <span>{{ importing ? 'Mengimport...' : 'Mulai Import' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ===== CREATE EMPLOYEE MODAL ===== -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div class="glass-panel w-full max-w-3xl rounded-2xl shadow-2xl p-6 relative max-h-[90vh] overflow-y-auto space-y-6">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <h2 class="text-base font-bold text-white">Formulir Tambah Pegawai Baru</h2>
          <button @click="showCreateModal = false" class="text-slate-400 hover:text-white">
            <X class="w-6 h-6" />
          </button>
        </div>

        <form @submit.prevent="submitCreateEmployee" class="space-y-4">
          <!-- Error feedback -->
          <div v-if="modalError" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 text-xs rounded-xl">
            {{ modalError }}
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- NIP -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">NIP (Username)*</label>
              <input v-model="form.nip" type="text" placeholder="Contoh: 19900101001" required
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
            </div>

            <!-- Nama -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Nama Lengkap*</label>
              <input v-model="form.nama" type="text" placeholder="Nama Beserta Gelar" required
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
            </div>

            <!-- Jenis Jabatan -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Jenis Jabatan*</label>
              <select v-model="form.jenis_jabatan" required @change="onJenisJabatanChange"
                      class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                <option value="Fungsional">Fungsional</option>
                <option value="Pelaksana">Pelaksana</option>
                <option value="Struktural">Struktural</option>
              </select>
            </div>

            <!-- Jabatan with Autocomplete -->
            <div class="relative">
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Jabatan*</label>
              <input v-model="form.jabatan" type="text" placeholder="Ketik atau pilih jabatan..." required
                     @input="onJabatanInput" @focus="showJabatanDropdown = true" @blur="onJabatanBlur"
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
              <div v-if="showJabatanDropdown && filteredJabatan.length > 0"
                   class="absolute z-50 top-full left-0 right-0 mt-1 bg-slate-900 border border-slate-700 rounded-lg shadow-xl max-h-48 overflow-y-auto">
                <button v-for="j in filteredJabatan" :key="j.id" type="button"
                        @mousedown.prevent="selectJabatan(j)"
                        class="w-full text-left px-3 py-2 text-xs text-slate-300 hover:bg-slate-800 transition-colors">
                  {{ j.nama_jabatan }}
                  <span class="text-slate-500 ml-1">({{ j.jenis_jabatan }})</span>
                </button>
              </div>
            </div>

            <!-- Tempat Lahir -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tempat Lahir</label>
              <input v-model="form.tempat_lahir" type="text" placeholder="Kota Lahir"
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
            </div>

            <!-- Tanggal Lahir -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tanggal Lahir</label>
              <DatePicker
                v-model="form.tanggal_lahir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>

            <!-- Tempat Tugas with Autocomplete -->
            <div class="relative">
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tempat Tugas</label>
              <input v-model="form.tempat_tugas" type="text" placeholder="Ketik atau pilih unit kerja..."
                     @input="onTempatInput" @focus="showTempatDropdown = true" @blur="onTempatBlur"
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
              <div v-if="showTempatDropdown && filteredTempat.length > 0"
                   class="absolute z-50 top-full left-0 right-0 mt-1 bg-slate-900 border border-slate-700 rounded-lg shadow-xl max-h-48 overflow-y-auto">
                <button v-for="t in filteredTempat" :key="t.id" type="button"
                        @mousedown.prevent="selectTempat(t)"
                        class="w-full text-left px-3 py-2 text-xs text-slate-300 hover:bg-slate-800 transition-colors">
                  {{ t.nama_tempat }}
                  <span class="text-slate-500 ml-1">({{ t.jenis_tempat }})</span>
                </button>
              </div>
            </div>

            <!-- Jenis Tempat -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Jenis Unit Kerja</label>
              <select v-model="form.jenis_tempat"
                      class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                <option value="Dinas">Dinas</option>
                <option value="Sekolah">Sekolah</option>
              </select>
            </div>

            <!-- Pengangkatan -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Pengangkatan</label>
              <input v-model="form.pengangkatan" type="text" placeholder="Contoh: CPNS 2019"
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
            </div>

            <!-- Jenis Pengangkatan -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Jenis Kepegawaian</label>
              <select v-model="form.jenis_pengangkatan"
                      class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-300 outline-none">
                <option value="PNS">PNS</option>
                <option value="PPPK">PPPK</option>
              </select>
            </div>

            <!-- No HP -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">No. WhatsApp Pegawai</label>
              <input v-model="form.no_hp" type="text" placeholder="Contoh: 0812345678"
                     class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none" />
            </div>

            <!-- KGB Terakhir -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tanggal KGB Terakhir</label>
              <DatePicker
                v-model="form.tanggal_kgb_terakhir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>

            <!-- Kenaikan Pangkat Terakhir -->
            <div>
              <label class="block text-xs font-semibold text-slate-350 mb-1.5 uppercase">Tanggal Pangkat Terakhir</label>
              <DatePicker
                v-model="form.tanggal_kenaikan_pangkat_terakhir"
                inputClass="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-xs text-slate-100 outline-none"
              />
            </div>
          </div>

          <!-- Document Uploads -->
          <div class="border-t border-slate-800 pt-4 space-y-4">
            <h3 class="text-xs font-bold text-white uppercase tracking-wider">Unggah File Lampiran SK (PDF/JPG/PNG Max 5MB)</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-[10px] font-semibold text-slate-400 mb-1">Foto Profil</label>
                <input type="file" @change="onFileChange($event, 'foto_profil')" accept="image/*"
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
              <div>
                <label class="block text-[10px] font-semibold text-slate-400 mb-1">SK CPNS / PPPK</label>
                <input type="file" @change="onFileChange($event, 'sk_cpns_pppk_file')" accept=".pdf,image/*"
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
              <div>
                <label class="block text-[10px] font-semibold text-slate-400 mb-1">SK PNS (Jika Ada)</label>
                <input type="file" @change="onFileChange($event, 'sk_pns_file')" accept=".pdf,image/*"
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
              <div>
                <label class="block text-[10px] font-semibold text-slate-400 mb-1">SK KGB Terbaru</label>
                <input type="file" @change="onFileChange($event, 'sk_kgb_file')" accept=".pdf,image/*"
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
              <div>
                <label class="block text-[10px] font-semibold text-slate-400 mb-1">SK Pangkat Terbaru</label>
                <input type="file" @change="onFileChange($event, 'sk_pangkat_file')" accept=".pdf,image/*"
                       class="w-full text-xs text-slate-400 file:mr-3 file:py-1 file:px-2 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-slate-800 file:text-slate-300 hover:file:bg-slate-700 cursor-pointer" />
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="border-t border-slate-800 pt-4 flex justify-end gap-3">
            <button type="button" @click="showCreateModal = false" :disabled="submitting"
                    class="px-4 py-2 bg-slate-800 hover:bg-slate-750 text-slate-300 rounded-xl text-xs font-semibold transition-all">
              Batal
            </button>
            <button type="submit" :disabled="submitting"
                    class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white rounded-xl text-xs font-semibold shadow-lg shadow-primary-500/20 transition-all flex items-center gap-2">
              <Loader2 v-if="submitting" class="w-3.5 h-3.5 animate-spin" />
              <span>{{ submitting ? 'Menyimpan...' : 'Simpan Pegawai' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ===== CONFIRM DELETE ALL MODAL ===== -->
    <div v-if="showDeleteAllModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="glass-panel w-full max-w-md rounded-2xl p-6 space-y-4">
        <div class="flex items-center gap-3 text-rose-500">
          <AlertTriangle class="w-8 h-8 shrink-0" />
          <h2 class="text-base font-bold text-white">Hapus Seluruh Data Pegawai?</h2>
        </div>
        <p class="text-xs text-slate-400 leading-relaxed">
          Tindakan ini akan <strong>menghapus secara permanen</strong> seluruh data pegawai (aktif & pensiun), beserta seluruh riwayat KGB, Pangkat, pengajuan cuti, berita acara, dan akun login terkait pegawai.
        </p>
        <div class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-3 text-[11px] text-rose-400">
          <strong>PENTING:</strong> Tindakan ini tidak dapat dibatalkan. Silakan ketik <strong>HAPUS</strong> untuk mengonfirmasi.
        </div>
        <input v-model="confirmText" type="text" placeholder="Ketik HAPUS di sini..."
               class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-rose-500 rounded-lg text-xs text-slate-100 outline-none" />
        <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
          <button @click="closeDeleteAllModal" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-semibold">
            Batal
          </button>
          <button @click="submitDeleteAll" :disabled="confirmText !== 'HAPUS' || deletingAll"
                  class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 disabled:opacity-40 disabled:cursor-not-allowed text-white rounded-lg text-xs font-semibold flex items-center gap-1.5">
            <Loader2 v-if="deletingAll" class="w-3.5 h-3.5 animate-spin" />
            <span>{{ deletingAll ? 'Menghapus...' : 'Hapus Semua' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'
import DataTable from '@/components/DataTable.vue'
import DatePicker from '@/components/DatePicker.vue'
import { UserPlus, X, Loader2, AlertTriangle, FileDown, FileUp, Download, Trash2 } from 'lucide-vue-next'
import Swal from 'sweetalert2'

const route = useRoute()
const router = useRouter()

const tableRef = ref(null)
const showCreateModal = ref(false)
const showImportModal = ref(false)
const submitting = ref(false)
const importing = ref(false)
const modalError = ref('')
const dragOver = ref(false)
const importFile = ref(null)
const importResult = ref(null)
const importFileInput = ref(null)

// Master data for autocomplete
const masterJabatan = ref([])
const masterTempat = ref([])
const showJabatanDropdown = ref(false)
const showTempatDropdown = ref(false)

// API base URL for direct downloads
const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const token = computed(() => localStorage.getItem('sipecut_token') || '')
const templateUrl = computed(() => `${apiBase}/api/employees/import-template?token=${token.value}`)
const exportUrl = computed(() => `${apiBase}/api/employees/export?token=${token.value}`)

const files = ref({
  foto_profil: null,
  sk_cpns_pppk_file: null,
  sk_pns_file: null,
  sk_kgb_file: null,
  sk_pangkat_file: null
})

const form = ref({
  nip: '',
  nama: '',
  jenis_jabatan: 'Fungsional',
  jabatan: '',
  tempat_lahir: '',
  tanggal_lahir: '',
  tempat_tugas: '',
  jenis_tempat: 'Dinas',
  pengangkatan: '',
  jenis_pengangkatan: 'PNS',
  no_hp: '',
  tanggal_kgb_terakhir: '',
  tanggal_kenaikan_pangkat_terakhir: ''
})

const activeFilter = computed(() => route.query.filter || '')
const extraQueryParams = computed(() => {
  if (activeFilter.value) return { dashboard_filter: activeFilter.value }
  return {}
})

watch(() => route.query.filter, () => {
  if (tableRef.value) tableRef.value.reload()
})

// ---- Master Autocomplete ----
const filteredJabatan = computed(() => {
  const q = (form.value.jabatan || '').toLowerCase()
  return masterJabatan.value.filter(j =>
    j.nama_jabatan.toLowerCase().includes(q) &&
    (!form.value.jenis_jabatan || j.jenis_jabatan === form.value.jenis_jabatan)
  ).slice(0, 10)
})

const filteredTempat = computed(() => {
  const q = (form.value.tempat_tugas || '').toLowerCase()
  return masterTempat.value.filter(t =>
    t.nama_tempat.toLowerCase().includes(q) &&
    (!form.value.jenis_tempat || t.jenis_tempat === form.value.jenis_tempat)
  ).slice(0, 10)
})

const onJenisJabatanChange = () => {
  form.value.jabatan = ''
  form.value.jenis_tempat = form.value.jenis_jabatan === 'Fungsional' ? 'Dinas' : 'Dinas'
}

const onJabatanInput = () => { showJabatanDropdown.value = true }
const onJabatanBlur = () => setTimeout(() => { showJabatanDropdown.value = false }, 150)
const selectJabatan = (j) => {
  form.value.jabatan = j.nama_jabatan
  showJabatanDropdown.value = false
}

const onTempatInput = () => { showTempatDropdown.value = true }
const onTempatBlur = () => setTimeout(() => { showTempatDropdown.value = false }, 150)
const selectTempat = (t) => {
  form.value.tempat_tugas = t.nama_tempat
  form.value.jenis_tempat = t.jenis_tempat
  showTempatDropdown.value = false
}

const fetchMasterData = async () => {
  try {
    const [jabRes, tempatRes] = await Promise.all([
      api.get('/master/jabatan'),
      api.get('/master/tempat-tugas')
    ])
    masterJabatan.value = jabRes.data.data || []
    masterTempat.value = tempatRes.data.data || []
  } catch (e) {
    console.warn('Gagal memuat master data:', e)
  }
}

// ---- Table ----
const getFilterLabel = (filter) => {
  if (filter === 'akan_kgb') return 'Pegawai yang akan kenaikan gaji berkala (KGB)'
  if (filter === 'akan_pangkat') return 'Pegawai yang akan kenaikan pangkat'
  if (filter === 'akan_pensiun') return 'Pegawai yang akan memasuki batas usia pensiun'
  return filter
}
const clearFilter = () => router.push('/pegawai')

const tableColumns = [
  { data: 'nip', title: 'NIP' },
  { data: 'nama', title: 'Nama Pegawai' },
  { data: 'jenis_jabatan', title: 'Jenis Jabatan' },
  { data: 'jabatan', title: 'Jabatan' },
  { data: 'tempat_tugas', title: 'Tempat Tugas' },
  { data: 'umur', title: 'Umur', orderable: false },
  {
    data: null,
    title: 'Aksi',
    orderable: false,
    render: (data, type, row) => {
      return `<button data-nip="${row.nip}" class="btn-detail px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white font-bold rounded-lg text-xs transition-colors">
        Detail Profil
      </button>`
    }
  }
]

const handleTableClick = (e) => {
  const target = e.target
  if (target.classList.contains('btn-detail')) {
    const nip = target.getAttribute('data-nip')
    router.push(`/pegawai/${nip}`)
  }
}

const onFileChange = (e, field) => {
  const file = e.target.files[0]
  if (file) files.value[field] = file
}

const openCreateModal = () => {
  modalError.value = ''
  showCreateModal.value = true
}

const submitCreateEmployee = async () => {
  submitting.value = true
  modalError.value = ''

  const formData = new FormData()
  for (const [key, value] of Object.entries(form.value)) {
    formData.append(key, value)
  }
  for (const [key, file] of Object.entries(files.value)) {
    if (file) formData.append(key, file)
  }

  try {
    await api.post('/employees', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    showCreateModal.value = false
    form.value = { nip: '', nama: '', jenis_jabatan: 'Fungsional', jabatan: '', tempat_lahir: '', tanggal_lahir: '', tempat_tugas: '', jenis_tempat: 'Dinas', pengangkatan: '', jenis_pengangkatan: 'PNS', no_hp: '', tanggal_kgb_terakhir: '', tanggal_kenaikan_pangkat_terakhir: '' }
    files.value = { foto_profil: null, sk_cpns_pppk_file: null, sk_pns_file: null, sk_kgb_file: null, sk_pangkat_file: null }
    if (tableRef.value) tableRef.value.reload()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Pegawai baru berhasil ditambahkan.',
      timer: 1550,
      showConfirmButton: false,
      background: '#1E293B',
      color: '#F8FAFC'
    })
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Gagal menyimpan data pegawai'
  } finally {
    submitting.value = false
  }
}

// ---- Import ----
const closeImportModal = () => {
  showImportModal.value = false
  importFile.value = null
  importResult.value = null
  dragOver.value = false
  if (importResult.value?.success_count > 0 && tableRef.value) tableRef.value.reload()
}

const onImportFileSelect = (e) => {
  importFile.value = e.target.files[0] || null
}

const onImportDrop = (e) => {
  dragOver.value = false
  const file = e.dataTransfer.files[0]
  if (file && file.name.endsWith('.xlsx')) {
    importFile.value = file
  } else {
    Swal.fire({
      icon: 'error',
      title: 'Format Salah',
      text: 'Hanya file .xlsx yang diizinkan',
      confirmButtonColor: '#1A365D'
    })
  }
}

const submitImport = async () => {
  if (!importFile.value) return
  importing.value = true
  const formData = new FormData()
  formData.append('file', importFile.value)
  try {
    const res = await api.post('/employees/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    importResult.value = res.data
    if (tableRef.value) tableRef.value.reload()
  } catch (err) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal Import',
      text: err.response?.data?.error || 'Gagal mengimport data',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    importing.value = false
  }
}
// ---- Delete All ----
const showDeleteAllModal = ref(false)
const confirmText = ref('')
const deletingAll = ref(false)

const confirmDeleteAll = () => {
  confirmText.value = ''
  showDeleteAllModal.value = true
}

const closeDeleteAllModal = () => {
  showDeleteAllModal.value = false
  confirmText.value = ''
}

const submitDeleteAll = async () => {
  if (confirmText.value !== 'HAPUS') return
  deletingAll.value = true
  try {
    const res = await api.delete('/employees/all')
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: res.data.message || 'Semua data pegawai berhasil dihapus',
      confirmButtonColor: '#1A365D'
    })
    showDeleteAllModal.value = false
    if (tableRef.value) tableRef.value.reload()
  } catch (err) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: err.response?.data?.error || 'Gagal menghapus semua data pegawai',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    deletingAll.value = false
  }
}

onMounted(() => {
  fetchMasterData()
})
</script>
