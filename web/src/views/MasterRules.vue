<template>
  <div class="space-y-6">
    <!-- Tab Navigation & Action Buttons -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-slate-200 pb-2 gap-4">
      <div class="flex flex-wrap gap-2">
        <button v-for="tab in tabs" :key="tab.id" @click="currentTab = tab.id"
                class="pb-3 px-1 text-xs font-bold uppercase tracking-wider relative transition-all"
                :class="currentTab === tab.id ? 'text-navy-600 font-extrabold' : 'text-slate-500 hover:text-slate-800'">
          {{ tab.name }}
          <div v-if="currentTab === tab.id" class="absolute bottom-0 left-0 right-0 h-0.5 bg-navy-600"></div>
        </button>
      </div>

      <!-- Master Bulk Actions (Admin only) -->
      <div v-if="['jabatan', 'tempat_tugas'].includes(currentTab)" class="flex items-center gap-2 pb-2">
        <a :href="masterTemplateUrl" target="_blank"
           class="px-3 py-1.5 bg-white hover:bg-slate-100 text-slate-700 hover:text-navy-900 border border-slate-300 rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all shadow-sm">
          <Download class="w-3.5 h-3.5 text-navy-600" />
          <span>Template Master</span>
        </a>
        <button @click="triggerMasterImport"
                class="px-3 py-1.5 bg-navy-600 hover:bg-navy-700 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all shadow-md shadow-navy-500/10">
          <Upload class="w-3.5 h-3.5" />
          <span>Import Master</span>
        </button>
        <button @click="confirmClearAllMaster"
                class="px-3 py-1.5 bg-rose-600/20 hover:bg-rose-600/30 text-rose-400 border border-rose-500/30 rounded-xl text-xs font-bold flex items-center gap-1.5 transition-all">
          <Trash2 class="w-3.5 h-3.5" />
          <span>Hapus Semua</span>
        </button>
        <input type="file" ref="masterFileInput" class="hidden" accept=".xlsx" @change="handleMasterImport" />
      </div>
    </div>

    <!-- TAB: Aturan Pensiun / KGB / Pangkat -->
    <div v-if="['pension', 'kgb', 'pangkat'].includes(currentTab)" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Rules Form -->
      <div class="glass-panel p-5 rounded-2xl h-fit space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">
          {{ ruleForm.id ? 'Edit Aturan' : 'Tambah Aturan Baru' }}
        </h2>
        <form @submit.prevent="submitRule" class="space-y-4 text-xs">
          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Jenis Jabatan*</label>
            <select v-model="ruleForm.jenis_jabatan" required
                    class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none">
              <option value="Fungsional">Fungsional</option>
              <option value="Pelaksana">Pelaksana</option>
              <option value="Struktural">Struktural</option>
            </select>
          </div>

          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Nama Jabatan (Isi * untuk Semua)*</label>
            <input v-model="ruleForm.jabatan" type="text" placeholder="Contoh: Guru Utama, Staf TU" required
                   class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none" />
          </div>

          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">
              {{ currentTab === 'pension' ? 'Batas Usia Pensiun (Tahun)*' : 'Siklus Periode (Tahun)*' }}
            </label>
            <input v-model="ruleForm.value" type="number" min="1" max="100" required
                   class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none" />
          </div>

          <div class="flex justify-end gap-2 pt-2 border-t border-slate-200">
            <button v-if="ruleForm.id" type="button" @click="resetForm"
                    class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg">Batal</button>
            <button type="submit" class="px-3 py-1.5 bg-navy-600 hover:bg-navy-700 text-white rounded-lg font-semibold">
              Simpan Aturan
            </button>
          </div>
        </form>
      </div>

      <!-- Rules list table -->
      <div class="glass-panel p-6 rounded-2xl lg:col-span-2 space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">
          Daftar Konfigurasi Aturan Master
        </h2>
        <div class="overflow-x-auto">
          <table class="w-full text-xs text-left text-slate-700 border-collapse">
            <thead>
              <tr class="border-b border-slate-200 text-slate-500">
                <th class="py-2.5 px-3">Jenis Jabatan</th>
                <th class="py-2.5 px-3">Jabatan</th>
                <th class="py-2.5 px-3">
                  {{ currentTab === 'pension' ? 'Batas Usia' : 'Siklus' }} (Tahun)
                </th>
                <th class="py-2.5 px-3 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="rulesList.length === 0">
                <td colspan="4" class="py-4 text-center text-slate-400 italic">Belum ada aturan khusus dikonfigurasi. Menggunakan fallback bawaan.</td>
              </tr>
              <tr v-for="item in rulesList" :key="item.id" class="border-b border-slate-100 hover:bg-slate-50">
                <td class="py-2.5 px-3 font-semibold text-navy-900">{{ item.jenis_jabatan }}</td>
                <td class="py-2.5 px-3 text-slate-800">{{ item.jabatan === '*' ? 'Semua Jabatan' : item.jabatan }}</td>
                <td class="py-2.5 px-3 font-bold text-navy-700">
                  {{ currentTab === 'pension' ? item.batas_usia_pensiun : item.siklus_tahun }} Tahun
                </td>
                <td class="py-2.5 px-3 text-right space-x-2">
                  <button @click="editRule(item)" class="text-amber-600 hover:underline font-semibold">Edit</button>
                  <button @click="deleteRule(item.id)" class="text-rose-600 hover:underline font-semibold">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB: Notifikasi & Kriteria Pensiun -->
    <div v-if="currentTab === 'settings'" class="glass-panel p-6 rounded-2xl max-w-3xl mx-auto space-y-6">
      <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">Pengaturan Umum & Notifikasi WhatsApp</h2>
      
      <form @submit.prevent="submitSettings" class="space-y-5 text-xs">
        <!-- Pension Criteria -->
        <div class="p-4 bg-slate-50 border border-slate-200 rounded-xl space-y-2">
          <label class="block font-bold text-slate-700 uppercase tracking-wider">Kriteria Kelengkapan Pensiun (Opsi B)*</label>
          <div class="flex flex-col gap-2 mt-2">
            <label class="flex items-center gap-2 text-slate-600 cursor-pointer">
              <input type="radio" v-model="settings.kriteria_pensiun_lengkap" value="dokumen_pemberhentian" class="accent-navy-600" />
              <span>Tanggal Pensiun Lewat <strong>DAN</strong> Dokumen Pemberhentian Pembayaran Telah Diunggah (Direkomendasikan)</span>
            </label>
            <label class="flex items-center gap-2 text-slate-600 cursor-pointer">
              <input type="radio" v-model="settings.kriteria_pensiun_lengkap" value="tanggal_saja" class="accent-navy-600" />
              <span>Tanggal Pensiun Lewat Saja (Beralih pensiun otomatis tanpa syarat upload berkas)</span>
            </label>
          </div>
        </div>

        <!-- Dashboard range alert -->
        <div>
          <label class="block font-semibold text-slate-700 mb-1.5 uppercase">Rentang Notifikasi Dashboard (Bulan)*</label>
          <input v-model="settings.dashboard_alert_months" type="number" required
                 class="w-full max-w-[150px] px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none" />
          <p class="text-[10px] text-slate-500 mt-1">Mengatur jarak waktu ke depan untuk mendeteksi pegawai yang "Akan KGB/Kenaikan Pangkat/Pensiun" di widget dashboard.</p>
        </div>

        <!-- WhatsApp Templates -->
        <div class="space-y-4">
          <h3 class="font-bold text-slate-700 uppercase border-b border-slate-200 pb-1">Format Template Pesan WhatsApp</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-slate-600 mb-1">Cuti Baru (Untuk Admin)</label>
              <textarea v-model="settings.whatsapp_template_leave_new" rows="3"
                        class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 outline-none resize-none font-mono text-[10px]"></textarea>
            </div>
            <div>
              <label class="block text-slate-600 mb-1">Cuti Update (Untuk Pegawai)</label>
              <textarea v-model="settings.whatsapp_template_leave_update" rows="3"
                        class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 outline-none resize-none font-mono text-[10px]"></textarea>
            </div>
            <div>
              <label class="block text-slate-600 mb-1">Ubah Profil Baru (Untuk Admin)</label>
              <textarea v-model="settings.whatsapp_template_change_new" rows="3"
                        class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 outline-none resize-none font-mono text-[10px]"></textarea>
            </div>
            <div>
              <label class="block text-slate-600 mb-1">Ubah Profil Update (Untuk Pegawai)</label>
              <textarea v-model="settings.whatsapp_template_change_update" rows="3"
                        class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-slate-800 outline-none resize-none font-mono text-[10px]"></textarea>
            </div>
          </div>
        </div>

        <!-- Submit settings button -->
        <div class="border-t border-slate-200 pt-4 flex justify-end">
          <button type="submit" :disabled="settingsSubmitting"
                  class="px-4 py-2 bg-gradient-to-r from-navy-600 to-indigo-700 text-white rounded-xl font-bold shadow-lg shadow-navy-500/20 hover:shadow-navy-500/35 transition-all flex items-center gap-2">
            <Loader2 v-if="settingsSubmitting" class="w-4 h-4 animate-spin" />
            <span>{{ settingsSubmitting ? 'Menyimpan...' : 'Simpan Pengaturan' }}</span>
          </button>
        </div>
      </form>
    </div>

    <!-- TAB: Master Jabatan -->
    <div v-if="currentTab === 'jabatan'" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Form -->
      <div class="glass-panel p-5 rounded-2xl h-fit space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">
          {{ jabatanForm.id ? 'Edit Jabatan' : 'Tambah Jabatan Baru' }}
        </h2>
        <form @submit.prevent="submitJabatan" class="space-y-4 text-xs">
          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Jenis Jabatan*</label>
            <select v-model="jabatanForm.jenis_jabatan" required
                    class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none">
              <option value="Fungsional">Fungsional</option>
              <option value="Pelaksana">Pelaksana</option>
              <option value="Struktural">Struktural</option>
            </select>
          </div>
          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Nama Jabatan*</label>
            <input v-model="jabatanForm.nama_jabatan" type="text" placeholder="Contoh: Guru Ahli Pertama" required
                   class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none" />
          </div>
          <div class="flex justify-end gap-2 pt-2 border-t border-slate-200">
            <button v-if="jabatanForm.id" type="button"
                    @click="jabatanForm = { id: null, nama_jabatan: '', jenis_jabatan: 'Fungsional' }"
                    class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-755 rounded-lg text-xs">Batal</button>
            <button type="submit" :disabled="jabatanSubmitting"
                    class="px-3 py-1.5 bg-navy-600 hover:bg-navy-700 text-white rounded-lg font-semibold text-xs flex items-center gap-1">
              <Loader2 v-if="jabatanSubmitting" class="w-3 h-3 animate-spin" />
              {{ jabatanForm.id ? 'Perbarui' : 'Simpan Jabatan' }}
            </button>
          </div>
        </form>
      </div>
      <!-- List Table -->
      <div class="glass-panel p-5 rounded-2xl lg:col-span-2 space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">Daftar Master Jabatan</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-xs text-left text-slate-700 border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200">
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Jenis</th>
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Nama Jabatan</th>
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="masterJabatanList.length === 0">
                <td colspan="3" class="p-4 text-center text-slate-400 italic">Belum ada data</td>
              </tr>
              <tr v-for="j in masterJabatanList" :key="j.id"
                  class="border-t border-slate-100 hover:bg-slate-50 transition-colors">
                <td class="p-3">
                  <span :class="j.jenis_jabatan === 'Fungsional' ? 'text-navy-600 font-semibold' : (j.jenis_jabatan === 'Struktural' ? 'text-indigo-600 font-semibold' : 'text-amber-600 font-semibold')">{{ j.jenis_jabatan }}</span>
                </td>
                <td class="p-3 text-slate-800 font-medium">{{ j.nama_jabatan }}</td>
                <td class="p-3">
                  <button @click="jabatanForm = { id: j.id, nama_jabatan: j.nama_jabatan, jenis_jabatan: j.jenis_jabatan }"
                          class="px-2 py-1 bg-slate-100 hover:bg-slate-200 text-slate-775 rounded text-[10px] mr-1">Edit</button>
                  <button @click="deleteJabatan(j.id)"
                          class="px-2 py-1 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded text-[10px]">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB: Master Tempat Tugas -->
    <div v-if="currentTab === 'tempat_tugas'" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Form -->
      <div class="glass-panel p-5 rounded-2xl h-fit space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">
          {{ tempatForm.id ? 'Edit Tempat Tugas' : 'Tambah Tempat Tugas Baru' }}
        </h2>
        <form @submit.prevent="submitTempat" class="space-y-4 text-xs">
          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Jenis Tempat*</label>
            <select v-model="tempatForm.jenis_tempat" required
                    class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none">
              <option value="Dinas">Dinas</option>
              <option value="Sekolah">Sekolah</option>
            </select>
          </div>
          <div>
            <label class="block text-[10px] font-semibold text-slate-600 mb-1 uppercase">Nama Tempat Tugas*</label>
            <input v-model="tempatForm.nama_tempat" type="text" placeholder="Contoh: SDN 01 Lembo" required
                   class="w-full px-3 py-2 bg-white border border-slate-300 focus:border-navy-500 rounded-lg text-slate-800 outline-none" />
          </div>
          <div class="flex justify-end gap-2 pt-2 border-t border-slate-200">
            <button v-if="tempatForm.id" type="button"
                    @click="tempatForm = { id: null, nama_tempat: '', jenis_tempat: 'Dinas' }"
                    class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-755 rounded-lg text-xs">Batal</button>
            <button type="submit" :disabled="tempatSubmitting"
                    class="px-3 py-1.5 bg-navy-600 hover:bg-navy-700 text-white rounded-lg font-semibold text-xs flex items-center gap-1">
              <Loader2 v-if="tempatSubmitting" class="w-3 h-3 animate-spin" />
              {{ tempatForm.id ? 'Perbarui' : 'Simpan Tempat Tugas' }}
            </button>
          </div>
        </form>
      </div>
      <!-- List Table -->
      <div class="glass-panel p-5 rounded-2xl lg:col-span-2 space-y-4">
        <h2 class="text-sm font-bold text-navy-800 border-b border-slate-200 pb-2">Daftar Master Tempat Tugas</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-xs text-left text-slate-700 border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200">
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Jenis</th>
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Nama Tempat</th>
                <th class="p-3 text-[10px] uppercase font-bold text-slate-500">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="masterTempatList.length === 0">
                <td colspan="3" class="p-4 text-center text-slate-400 italic">Belum ada data</td>
              </tr>
              <tr v-for="t in masterTempatList" :key="t.id"
                  class="border-t border-slate-100 hover:bg-slate-50 transition-colors">
                <td class="p-3">
                  <span :class="t.jenis_tempat === 'Sekolah' ? 'text-emerald-600 font-semibold' : 'text-violet-600 font-semibold'">{{ t.jenis_tempat }}</span>
                </td>
                <td class="p-3 text-slate-800 font-medium">{{ t.nama_tempat }}</td>
                <td class="p-3">
                  <button @click="tempatForm = { id: t.id, nama_tempat: t.nama_tempat, jenis_tempat: t.jenis_tempat }"
                          class="px-2 py-1 bg-slate-100 hover:bg-slate-200 text-slate-775 rounded text-[10px] mr-1">Edit</button>
                  <button @click="deleteTempat(t.id)"
                          class="px-2 py-1 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded text-[10px]">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB: Log Audit Trail -->
    <div v-if="currentTab === 'audit'" class="glass-panel p-4 rounded-2xl">
      <DataTable
        ref="auditTableRef"
        :columns="auditColumns"
        ajaxUrl="/api/dashboard/audit-history"
      />
    </div>

    <!-- TAB: Log WhatsApp -->
    <div v-if="currentTab === 'whatsapp'" class="glass-panel p-4 rounded-2xl">
      <DataTable
        ref="waTableRef"
        :columns="waColumns"
        ajaxUrl="/api/dashboard/notification-logs"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import api from '@/services/api'
import DataTable from '@/components/DataTable.vue'
import { Loader2, Download, Upload, Trash2 } from 'lucide-vue-next'
import Swal from 'sweetalert2'

const currentTab = ref('pension')
const tabs = [
  { id: 'pension', name: 'Aturan Pensiun' },
  { id: 'kgb', name: 'Aturan KGB' },
  { id: 'pangkat', name: 'Aturan Pangkat' },
  { id: 'jabatan', name: 'Master Jabatan' },
  { id: 'tempat_tugas', name: 'Master Tempat Tugas' },
  { id: 'settings', name: 'Notifikasi & Kriteria' },
  { id: 'audit', name: 'Log Audit Trail' },
  { id: 'whatsapp', name: 'Log WhatsApp' }
]

const rulesList = ref([])
const ruleForm = ref({ id: null, jenis_jabatan: 'Fungsional', jabatan: '', value: 2 })

// Settings Map
const settings = ref({
  kriteria_pensiun_lengkap: 'dokumen_pemberhentian',
  dashboard_alert_months: '3',
  whatsapp_template_leave_new: '',
  whatsapp_template_leave_update: '',
  whatsapp_template_change_new: '',
  whatsapp_template_change_update: '',
  whatsapp_template_ba_new: '',
  whatsapp_template_ba_update: ''
})
const settingsSubmitting = ref(false)

// Master Bulk Actions Elements
const masterFileInput = ref(null)
const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const token = computed(() => localStorage.getItem('sipecut_token') || '')
const masterTemplateUrl = computed(() => `${apiBase}/api/master/import-template?token=${token.value}`)

const triggerMasterImport = () => {
  if (masterFileInput.value) {
    masterFileInput.value.click()
  }
}

const handleMasterImport = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await api.post('/master/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: res.data.message || 'Master data berhasil diimpor!',
      confirmButtonColor: '#1A365D'
    })
    // Reload active tab data
    if (['pension', 'kgb', 'pangkat'].includes(currentTab.value)) {
      await fetchRules()
    } else if (currentTab.value === 'jabatan') {
      await fetchMasterJabatan()
    } else if (currentTab.value === 'tempat_tugas') {
      await fetchMasterTempat()
    }
  } catch (e) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: e.response?.data?.error || 'Gagal mengimpor master data',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    if (event.target) event.target.value = ''
  }
}

const confirmClearAllMaster = async () => {
  const isJabatan = currentTab.value === 'jabatan'
  const title = isJabatan ? 'Hapus Semua Master Jabatan?' : 'Hapus Semua Master Tempat Tugas?'
  const text = isJabatan
    ? 'Apakah Anda yakin ingin menghapus seluruh data master jabatan? Tindakan ini tidak dapat dibatalkan.'
    : 'Apakah Anda yakin ingin menghapus seluruh data master tempat tugas? Tindakan ini tidak dapat dibatalkan.'

  const result = await Swal.fire({
    title: title,
    text: text,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#E11D48',
    cancelButtonColor: '#1A365D',
    confirmButtonText: 'Ya, Hapus Semua',
    cancelButtonText: 'Batal',
    background: '#1E293B',
    color: '#F8FAFC'
  })

  if (result.isConfirmed) {
    try {
      const endpoint = isJabatan ? '/master/jabatan/all' : '/master/tempat-tugas/all'
      const res = await api.delete(endpoint)
      
      Swal.fire({
        icon: 'success',
        title: 'Berhasil',
        text: res.data.message || 'Data master berhasil dihapus.',
        timer: 1550,
        showConfirmButton: false,
        background: '#1E293B',
        color: '#F8FAFC'
      })

      if (isJabatan) {
        await fetchMasterJabatan()
      } else {
        await fetchMasterTempat()
      }
    } catch (e) {
      Swal.fire({
        icon: 'error',
        title: 'Gagal',
        text: e.response?.data?.error || 'Gagal menghapus data master.',
        confirmButtonColor: '#1A365D',
        background: '#1E293B',
        color: '#F8FAFC'
      })
    }
  }
}

// Watch tab change to fetch data
watch(currentTab, (newTab) => {
  if (['pension', 'kgb', 'pangkat'].includes(newTab)) {
    fetchRules()
    resetForm()
  } else if (newTab === 'settings') {
    fetchSettings()
  } else if (newTab === 'jabatan') {
    fetchMasterJabatan()
    jabatanForm.value = { id: null, nama_jabatan: '', jenis_jabatan: 'Fungsional' }
  } else if (newTab === 'tempat_tugas') {
    fetchMasterTempat()
    tempatForm.value = { id: null, nama_tempat: '', jenis_tempat: 'Dinas' }
  }
})

const fetchRules = async () => {
  try {
    let endpoint = ''
    if (currentTab.value === 'pension') endpoint = '/master/pension-rules'
    if (currentTab.value === 'kgb') endpoint = '/master/kgb-cycle-rules'
    if (currentTab.value === 'pangkat') endpoint = '/master/pangkat-cycle-rules'

    const response = await api.get(endpoint)
    rulesList.value = response.data.data || []
  } catch (error) {
    console.error('Gagal mengambil aturan master:', error)
  }
}

const fetchSettings = async () => {
  try {
    const response = await api.get('/master/settings')
    // Populate settings mapping
    for (const [k, v] of Object.entries(response.data.data)) {
      settings.value[k] = v
    }
  } catch (error) {
    console.error('Gagal mengambil settings:', error)
  }
}

const resetForm = () => {
  ruleForm.value = {
    id: null,
    jenis_jabatan: 'Fungsional',
    jabatan: '',
    value: currentTab.value === 'pension' ? 58 : 2
  }
}

const editRule = (item) => {
  ruleForm.value = {
    id: item.id,
    jenis_jabatan: item.jenis_jabatan,
    jabatan: item.jabatan,
    value: currentTab.value === 'pension' ? item.batas_usia_pensiun : item.siklus_tahun
  }
}

const submitRule = async () => {
  try {
    let endpoint = ''
    let payload = {
      jenis_jabatan: ruleForm.value.jenis_jabatan,
      jabatan: ruleForm.value.jabatan
    }

    if (ruleForm.value.id) payload.id = ruleForm.value.id

    if (currentTab.value === 'pension') {
      endpoint = '/master/pension-rules'
      payload.batas_usia_pensiun = parseInt(ruleForm.value.value)
    } else if (currentTab.value === 'kgb') {
      endpoint = '/master/kgb-cycle-rules'
      payload.siklus_tahun = parseInt(ruleForm.value.value)
    } else if (currentTab.value === 'pangkat') {
      endpoint = '/master/pangkat-cycle-rules'
      payload.siklus_tahun = parseInt(ruleForm.value.value)
    }

    await api.put(endpoint, payload)
    await fetchRules()
    resetForm()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Aturan master berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan aturan',
      confirmButtonColor: '#1A365D'
    })
  }
}

const deleteRule = async (id) => {
  const result = await Swal.fire({
    title: 'Apakah Anda yakin?',
    text: 'Aturan master ini akan dihapus secara permanen!',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#E11D48',
    cancelButtonColor: '#1A365D',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: 'Batal',
    background: '#1E293B',
    color: '#F8FAFC'
  })
  if (!result.isConfirmed) return

  try {
    let endpoint = ''
    if (currentTab.value === 'pension') endpoint = `/master/pension-rules/${id}`
    if (currentTab.value === 'kgb') endpoint = `/master/kgb-cycle-rules/${id}`
    if (currentTab.value === 'pangkat') endpoint = `/master/pangkat-cycle-rules/${id}`

    await api.delete(endpoint)
    await fetchRules()
    Swal.fire({
      icon: 'success',
      title: 'Dihapus!',
      text: 'Aturan master berhasil dihapus.',
      timer: 1550,
      showConfirmButton: false,
      background: '#1E293B',
      color: '#F8FAFC'
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menghapus aturan',
      confirmButtonColor: '#1A365D',
      background: '#1E293B',
      color: '#F8FAFC'
    })
  }
}

const submitSettings = async () => {
  settingsSubmitting.value = true
  try {
    await api.put('/master/settings', settings.value)
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Pengaturan umum dan template WhatsApp berhasil diperbarui.',
      confirmButtonColor: '#1A365D'
    })
  } catch (error) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: error.response?.data?.error || 'Gagal menyimpan pengaturan',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    settingsSubmitting.value = false
  }
}

// Columns definition for Logs
const auditColumns = [
  { data: 'id', title: 'ID', width: '50px' },
  { data: 'request_type', title: 'Jenis Log' },
  { data: 'status_lama', title: 'Status Lama' },
  { data: 'status_baru', title: 'Status Baru' },
  { data: 'catatan', title: 'Aktivitas / Catatan' },
  { data: 'changed_by', title: 'Oleh (NIP)' },
  {
    data: 'changed_at',
    title: 'Tanggal Log',
    render: (data) => {
      const d = new Date(data)
      return `${d.toLocaleDateString('id-ID')} ${d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })}`
    }
  }
]

const waColumns = [
  { data: 'id', title: 'ID', width: '50px' },
  { data: 'channel', title: 'Saluran' },
  { data: 'message', title: 'Isi Pesan Notifikasi' },
  {
    data: 'status',
    title: 'Status Kirim',
    render: (data) => {
      let color = 'text-slate-400'
      if (data === 'Sent') color = 'text-emerald-600 font-bold'
      if (data === 'Simulated') color = 'text-blue-600 font-semibold'
      if (data === 'Failed') color = 'text-red-600 font-bold'
      return `<span class="${color}">${data}</span>`
    }
  },
  {
    data: 'sent_at',
    title: 'Tanggal Kirim',
    render: (data) => {
      const d = new Date(data)
      return `${d.toLocaleDateString('id-ID')} ${d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })}`
    }
  }
]

// ---- Master Jabatan ----
const masterJabatanList = ref([])
const jabatanForm = ref({ id: null, nama_jabatan: '', jenis_jabatan: 'Fungsional' })
const jabatanSubmitting = ref(false)

const fetchMasterJabatan = async () => {
  try {
    const res = await api.get('/master/jabatan')
    masterJabatanList.value = res.data.data || []
  } catch (e) {
    console.error('Gagal memuat jabatan:', e)
  }
}

const submitJabatan = async () => {
  jabatanSubmitting.value = true
  try {
    if (jabatanForm.value.id) {
      await api.put(`/master/jabatan/${jabatanForm.value.id}`, jabatanForm.value)
    } else {
      await api.post('/master/jabatan', jabatanForm.value)
    }
    jabatanForm.value = { id: null, nama_jabatan: '', jenis_jabatan: 'Fungsional' }
    await fetchMasterJabatan()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Data jabatan berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (e) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: e.response?.data?.error || 'Gagal menyimpan jabatan',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    jabatanSubmitting.value = false
  }
}

const deleteJabatan = async (id) => {
  const result = await Swal.fire({
    title: 'Hapus Jabatan?',
    text: 'Apakah Anda yakin ingin menghapus jabatan ini?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#E11D48',
    cancelButtonColor: '#64748B',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: 'Batal'
  })
  if (!result.isConfirmed) return

  try {
    await api.delete(`/master/jabatan/${id}`)
    await fetchMasterJabatan()
    Swal.fire({
      icon: 'success',
      title: 'Dihapus!',
      text: 'Jabatan berhasil dihapus.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (e) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: e.response?.data?.error || 'Gagal menghapus jabatan',
      confirmButtonColor: '#1A365D'
    })
  }
}

// ---- Master Tempat Tugas ----
const masterTempatList = ref([])
const tempatForm = ref({ id: null, nama_tempat: '', jenis_tempat: 'Dinas' })
const tempatSubmitting = ref(false)

const fetchMasterTempat = async () => {
  try {
    const res = await api.get('/master/tempat-tugas')
    masterTempatList.value = res.data.data || []
  } catch (e) {
    console.error('Gagal memuat tempat tugas:', e)
  }
}

const submitTempat = async () => {
  tempatSubmitting.value = true
  try {
    if (tempatForm.value.id) {
      await api.put(`/master/tempat-tugas/${tempatForm.value.id}`, tempatForm.value)
    } else {
      await api.post('/master/tempat-tugas', tempatForm.value)
    }
    tempatForm.value = { id: null, nama_tempat: '', jenis_tempat: 'Dinas' }
    await fetchMasterTempat()
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: 'Data tempat tugas berhasil disimpan.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (e) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: e.response?.data?.error || 'Gagal menyimpan tempat tugas',
      confirmButtonColor: '#1A365D'
    })
  } finally {
    tempatSubmitting.value = false
  }
}

const deleteTempat = async (id) => {
  const result = await Swal.fire({
    title: 'Hapus Tempat Tugas?',
    text: 'Apakah Anda yakin ingin menghapus tempat tugas ini?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#E11D48',
    cancelButtonColor: '#64748B',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: 'Batal'
  })
  if (!result.isConfirmed) return

  try {
    await api.delete(`/master/tempat-tugas/${id}`)
    await fetchMasterTempat()
    Swal.fire({
      icon: 'success',
      title: 'Dihapus!',
      text: 'Tempat tugas berhasil dihapus.',
      timer: 1500,
      showConfirmButton: false
    })
  } catch (e) {
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: e.response?.data?.error || 'Gagal menghapus tempat tugas',
      confirmButtonColor: '#1A365D'
    })
  }
}

onMounted(() => {
  fetchRules()
})
</script>
