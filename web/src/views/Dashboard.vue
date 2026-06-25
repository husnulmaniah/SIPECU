<template>
  <div class="space-y-6">
    <!-- Welcome section -->
    <div class="glass-panel p-6 rounded-2xl flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white">Selamat Datang, {{ authStore.employee?.nama || 'Administrator' }}!</h1>
        <p class="text-xs text-slate-400 mt-1">Hari ini adalah {{ currentFormattedDate }}. Semoga aktivitas Anda berjalan lancar.</p>
      </div>
      <div v-if="authStore.role === 'employee'" class="flex items-center gap-2">
        <span class="px-3 py-1.5 rounded-full text-xs font-semibold"
              :class="authStore.employee?.status_kepegawaian === 'Aktif' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'">
          Status: {{ authStore.employee?.status_kepegawaian }}
        </span>
        <span class="px-3 py-1.5 rounded-full text-xs font-semibold bg-primary-500/10 text-primary-400 border border-primary-500/20">
          Jabatan: {{ authStore.employee?.jenis_jabatan }}
        </span>
      </div>
    </div>

    <!-- ADMIN DASHBOARD -->
    <div v-if="authStore.role === 'admin'" class="space-y-6">
      <!-- Statistics cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        <div @click="goToEmployees('aktif')" class="glass-panel p-5 rounded-2xl hover:border-primary-500/50 cursor-pointer transition-all duration-300 group">
          <div class="flex items-center justify-between mb-3">
            <span class="text-slate-400 text-xs font-semibold uppercase tracking-wider">Pegawai Aktif</span>
            <div class="w-8 h-8 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center group-hover:scale-110 transition-transform">
              <Users class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-black text-white">{{ stats.total_active || 0 }}</p>
          <span class="text-[10px] text-slate-500 font-medium">Klik untuk lihat daftar</span>
        </div>

        <div @click="goToEmployees('pensiun')" class="glass-panel p-5 rounded-2xl hover:border-primary-500/50 cursor-pointer transition-all duration-300 group">
          <div class="flex items-center justify-between mb-3">
            <span class="text-slate-400 text-xs font-semibold uppercase tracking-wider">Pegawai Pensiun</span>
            <div class="w-8 h-8 rounded-lg bg-slate-800 text-slate-400 flex items-center justify-center group-hover:scale-110 transition-transform">
              <UserCheck class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-black text-white">{{ stats.total_retired || 0 }}</p>
          <span class="text-[10px] text-slate-500 font-medium">Klik untuk lihat arsip</span>
        </div>

        <div @click="goToEmployees('akan_kgb')" class="glass-panel p-5 rounded-2xl hover:border-primary-500/50 cursor-pointer transition-all duration-300 group">
          <div class="flex items-center justify-between mb-3">
            <span class="text-slate-400 text-xs font-semibold uppercase tracking-wider">Akan KGB</span>
            <div class="w-8 h-8 rounded-lg bg-amber-500/10 text-amber-400 flex items-center justify-center group-hover:scale-110 transition-transform">
              <TrendingUp class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-black text-white">{{ stats.akan_kgb || 0 }}</p>
          <span class="text-[10px] text-amber-400 font-medium">&lt; {{ stats.config_months }} Bulan Ke Depan</span>
        </div>

        <div @click="goToEmployees('akan_pangkat')" class="glass-panel p-5 rounded-2xl hover:border-primary-500/50 cursor-pointer transition-all duration-300 group">
          <div class="flex items-center justify-between mb-3">
            <span class="text-slate-400 text-xs font-semibold uppercase tracking-wider">Akan Pangkat</span>
            <div class="w-8 h-8 rounded-lg bg-violet-500/10 text-violet-400 flex items-center justify-center group-hover:scale-110 transition-transform">
              <Award class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-black text-white">{{ stats.akan_pangkat || 0 }}</p>
          <span class="text-[10px] text-violet-400 font-medium">&lt; {{ stats.config_months }} Bulan Ke Depan</span>
        </div>

        <div @click="goToEmployees('akan_pensiun')" class="glass-panel p-5 rounded-2xl hover:border-primary-500/50 cursor-pointer transition-all duration-300 group">
          <div class="flex items-center justify-between mb-3">
            <span class="text-slate-400 text-xs font-semibold uppercase tracking-wider">Akan Pensiun</span>
            <div class="w-8 h-8 rounded-lg bg-red-500/10 text-red-400 flex items-center justify-center group-hover:scale-110 transition-transform">
              <CalendarX class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-black text-white">{{ stats.akan_pensiun || 0 }}</p>
          <span class="text-[10px] text-red-400 font-medium">&lt; {{ stats.config_months }} Bulan Ke Depan</span>
        </div>
      </div>

      <!-- Action Required/Pending Items Section -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Leaves pending -->
        <router-link to="/cuti" class="glass-panel p-5 rounded-2xl flex items-center justify-between hover:border-primary-500 transition-all">
          <div class="space-y-1">
            <h3 class="text-sm font-semibold text-slate-350">Antrean Pengajuan Cuti</h3>
            <p class="text-xs text-slate-500">Persetujuan & pembuatan surat rekomendasi</p>
          </div>
          <span class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm"
                :class="stats.pending_leaves > 0 ? 'bg-red-500 text-white animate-pulse' : 'bg-slate-800 text-slate-400'">
            {{ stats.pending_leaves || 0 }}
          </span>
        </router-link>

        <!-- Changes pending -->
        <router-link to="/perubahan-data" class="glass-panel p-5 rounded-2xl flex items-center justify-between hover:border-primary-500 transition-all">
          <div class="space-y-1">
            <h3 class="text-sm font-semibold text-slate-350">Ubah Profil Mandiri</h3>
            <p class="text-xs text-slate-500">Verifikasi berkas pendukung perubahan</p>
          </div>
          <span class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm"
                :class="stats.pending_changes > 0 ? 'bg-amber-500 text-white animate-pulse' : 'bg-slate-800 text-slate-400'">
            {{ stats.pending_changes || 0 }}
          </span>
        </router-link>

        <!-- BA pending -->
        <router-link to="/berita-acara" class="glass-panel p-5 rounded-2xl flex items-center justify-between hover:border-primary-500 transition-all">
          <div class="space-y-1">
            <h3 class="text-sm font-semibold text-slate-350">Verifikasi Berita Acara</h3>
            <p class="text-xs text-slate-500">Pemberkasan BA / ST / SI / SKS</p>
          </div>
          <span class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm"
                :class="stats.pending_ba > 0 ? 'bg-violet-500 text-white animate-pulse' : 'bg-slate-800 text-slate-400'">
            {{ stats.pending_ba || 0 }}
          </span>
        </router-link>
      </div>

      <!-- Chart visualizer -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="glass-panel p-6 rounded-2xl lg:col-span-2 space-y-4">
          <h2 class="text-base font-bold text-white">Visualisasi Distribusi & Alert Kepegawaian</h2>
          <div class="h-64 flex items-center justify-center bg-slate-950/20 rounded-xl border border-slate-800 p-4">
            <!-- Simulated chart representing stats visually since pure npm install is slow to boot -->
            <div class="w-full h-full flex flex-col justify-end gap-2">
              <div class="flex-1 flex items-end justify-around gap-4 px-4">
                <div class="flex flex-col items-center gap-1 w-full max-w-[50px]">
                  <div class="bg-emerald-500 w-full rounded-t-md transition-all duration-500" :style="{ height: getBarHeight(stats.total_active, 100) }"></div>
                  <span class="text-[10px] text-slate-400 text-center font-medium mt-1">Aktif</span>
                </div>
                <div class="flex flex-col items-center gap-1 w-full max-w-[50px]">
                  <div class="bg-slate-600 w-full rounded-t-md transition-all duration-500" :style="{ height: getBarHeight(stats.total_retired, 100) }"></div>
                  <span class="text-[10px] text-slate-400 text-center font-medium mt-1">Pensiun</span>
                </div>
                <div class="flex flex-col items-center gap-1 w-full max-w-[50px]">
                  <div class="bg-amber-500 w-full rounded-t-md transition-all duration-500" :style="{ height: getBarHeight(stats.akan_kgb, 100) }"></div>
                  <span class="text-[10px] text-slate-400 text-center font-medium mt-1">KGB</span>
                </div>
                <div class="flex flex-col items-center gap-1 w-full max-w-[50px]">
                  <div class="bg-violet-500 w-full rounded-t-md transition-all duration-500" :style="{ height: getBarHeight(stats.akan_pangkat, 100) }"></div>
                  <span class="text-[10px] text-slate-400 text-center font-medium mt-1">Pangkat</span>
                </div>
                <div class="flex flex-col items-center gap-1 w-full max-w-[50px]">
                  <div class="bg-red-500 w-full rounded-t-md transition-all duration-500" :style="{ height: getBarHeight(stats.akan_pensiun, 100) }"></div>
                  <span class="text-[10px] text-slate-400 text-center font-medium mt-1">Pensiun alert</span>
                </div>
              </div>
              <div class="border-t border-slate-800 pt-2 text-center">
                <span class="text-[10px] text-slate-500 font-semibold">Statistik Jumlah Pegawai per Kategori (Tinggi Bar Proporsional)</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Info Panel -->
        <div class="glass-panel p-6 rounded-2xl space-y-4">
          <h2 class="text-base font-bold text-white">Sistem Log & Notifikasi</h2>
          <div class="space-y-3">
            <div class="p-3 bg-slate-900/60 border border-slate-800 rounded-xl flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
              <div>
                <p class="text-xs font-semibold text-white">Koneksi Database Aktif</p>
                <p class="text-[10px] text-slate-500">Menggunakan SQLite Local Database (sipecut.db)</p>
              </div>
            </div>
            <div class="p-3 bg-slate-900/60 border border-slate-800 rounded-xl flex items-start gap-3">
              <CheckCircle class="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
              <div>
                <p class="text-xs font-semibold text-white">WhatsApp Gateway Simulator</p>
                <p class="text-[10px] text-slate-500">Simulator aktif, log tersimpan di notification_logs</p>
              </div>
            </div>
            <div class="p-3 bg-slate-900/60 border border-slate-800 rounded-xl flex items-start gap-3">
              <Clock class="w-5 h-5 text-violet-400 shrink-0 mt-0.5" />
              <div>
                <p class="text-xs font-semibold text-white">Auto-Retirement Checker</p>
                <p class="text-[10px] text-slate-500">Memeriksa status pensiun berkala (setiap 1 jam)</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- EMPLOYEE DASHBOARD -->
    <div v-else class="space-y-6">
      <!-- Target Countdown deadlines -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- KGB target card -->
        <div class="glass-panel p-5 rounded-2xl relative overflow-hidden group">
          <div class="absolute w-24 h-24 rounded-full bg-amber-500/5 -right-8 -bottom-8 group-hover:scale-125 transition-all"></div>
          <h3 class="text-slate-400 text-xs font-semibold uppercase tracking-wider mb-2">Target KGB Berikutnya</h3>
          <p class="text-lg font-bold text-white">{{ formatDate(stats.next_kgb) }}</p>
          <p class="text-[10px] text-amber-400 font-semibold mt-1">
            <span class="px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20">
              Siklus: 2 Tahun Sekali
            </span>
          </p>
        </div>

        <!-- Pangkat target card -->
        <div class="glass-panel p-5 rounded-2xl relative overflow-hidden group">
          <div class="absolute w-24 h-24 rounded-full bg-violet-500/5 -right-8 -bottom-8 group-hover:scale-125 transition-all"></div>
          <h3 class="text-slate-400 text-xs font-semibold uppercase tracking-wider mb-2">Target Naik Pangkat</h3>
          <p class="text-lg font-bold text-white">{{ formatDate(stats.next_pangkat) }}</p>
          <p class="text-[10px] text-violet-400 font-semibold mt-1">
            <span class="px-2 py-0.5 rounded-full bg-violet-500/10 border border-violet-500/20">
              Siklus: 4 Tahun Sekali
            </span>
          </p>
        </div>

        <!-- Pension target card -->
        <div class="glass-panel p-5 rounded-2xl relative overflow-hidden group">
          <div class="absolute w-24 h-24 rounded-full bg-red-500/5 -right-8 -bottom-8 group-hover:scale-125 transition-all"></div>
          <h3 class="text-slate-400 text-xs font-semibold uppercase tracking-wider mb-2">Target Pensiun</h3>
          <p class="text-lg font-bold text-white">{{ formatDate(stats.next_pension) }}</p>
          <p class="text-[10px] text-red-400 font-semibold mt-1">
            <span class="px-2 py-0.5 rounded-full bg-red-500/10 border border-red-500/20">
              Batas Usia: {{ authStore.employee?.jenis_jabatan === 'Fungsional' ? '60' : '58' }} Tahun
            </span>
          </p>
        </div>
      </div>

      <!-- Action Buttons Grid -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <router-link to="/cuti/tambah" class="glass-panel p-5 rounded-2xl hover:border-primary-500 hover:bg-slate-800/20 transition-all flex flex-col items-center justify-center text-center gap-3">
          <div class="w-12 h-12 rounded-xl bg-primary-500/10 text-primary-400 flex items-center justify-center">
            <CalendarDays class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-xs font-bold text-white">Ajukan Cuti</h3>
            <p class="text-[9px] text-slate-500 mt-0.5">Dapatkan surat rekomendasi</p>
          </div>
        </router-link>

        <router-link to="/perubahan-data" class="glass-panel p-5 rounded-2xl hover:border-primary-500 hover:bg-slate-800/20 transition-all flex flex-col items-center justify-center text-center gap-3">
          <div class="w-12 h-12 rounded-xl bg-amber-500/10 text-amber-400 flex items-center justify-center">
            <UserCog class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-xs font-bold text-white">Ubah Data Mandiri</h3>
            <p class="text-[9px] text-slate-500 mt-0.5">Ajukan koreksi profil</p>
          </div>
        </router-link>

        <router-link to="/berita-acara" class="glass-panel p-5 rounded-2xl hover:border-primary-500 hover:bg-slate-800/20 transition-all flex flex-col items-center justify-center text-center gap-3">
          <div class="w-12 h-12 rounded-xl bg-violet-500/10 text-violet-400 flex items-center justify-center">
            <FileText class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-xs font-bold text-white">Unggah Berita Acara</h3>
            <p class="text-[9px] text-slate-500 mt-0.5">Arsipkan berkas BA/ST/SI</p>
          </div>
        </router-link>

        <router-link to="/profil" class="glass-panel p-5 rounded-2xl hover:border-primary-500 hover:bg-slate-800/20 transition-all flex flex-col items-center justify-center text-center gap-3">
          <div class="w-12 h-12 rounded-xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center">
            <Key class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-xs font-bold text-white">Ubah Sandi Akun</h3>
            <p class="text-[9px] text-slate-500 mt-0.5">Keamanan kredensial login</p>
          </div>
        </router-link>
      </div>

      <!-- Quick status overview -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="glass-panel p-6 rounded-2xl lg:col-span-2 space-y-4">
          <h2 class="text-sm font-bold text-white">Status Pengajuan Terbaru Anda</h2>
          <div class="space-y-4">
            <div class="p-4 bg-slate-900/40 border border-slate-800 rounded-xl flex items-center justify-between gap-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-primary-500/15 text-primary-400 flex items-center justify-center">
                  <CalendarDays class="w-4 h-4" />
                </div>
                <div>
                  <h4 class="text-xs font-bold text-white">Pengajuan Cuti</h4>
                  <p class="text-[10px] text-slate-500">Mengajukan izin cuti dinamis</p>
                </div>
              </div>
              <div class="text-right">
                <span class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-slate-800 text-slate-400">
                  Total Diajukan: {{ stats.total_leaves }}
                </span>
                <p v-if="stats.pending_leaves > 0" class="text-[9px] text-amber-400 mt-1 font-semibold">{{ stats.pending_leaves }} Menunggu Review</p>
                <p v-else class="text-[9px] text-slate-500 mt-1 font-medium">Semua selesai ditinjau</p>
              </div>
            </div>

            <div class="p-4 bg-slate-900/40 border border-slate-800 rounded-xl flex items-center justify-between gap-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-amber-500/15 text-amber-400 flex items-center justify-center">
                  <UserCog class="w-4 h-4" />
                </div>
                <div>
                  <h4 class="text-xs font-bold text-white">Perubahan Profil</h4>
                  <p class="text-[10px] text-slate-500">Memperbarui KGB/Pangkat mandiri</p>
                </div>
              </div>
              <div class="text-right">
                <span class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold"
                      :class="stats.active_changes > 0 ? 'bg-amber-500/10 text-amber-400' : 'bg-slate-800 text-slate-400'">
                  {{ stats.active_changes > 0 ? 'Ada Antrean Aktif' : 'Tidak Ada Antrean' }}
                </span>
                <p class="text-[9px] text-slate-500 mt-1 font-medium">Berdasarkan dokumen SK Pendukung</p>
              </div>
            </div>
          </div>
        </div>

        <div class="glass-panel p-6 rounded-2xl space-y-4">
          <h2 class="text-sm font-bold text-white">Informasi Singkat</h2>
          <div class="space-y-3 text-xs text-slate-400 leading-relaxed">
            <p>1. Rekomendasi cuti yang telah di-**ACC oleh Admin** dapat langsung Anda unduh dalam format **PDF** melalui halaman "Pengajuan Cuti".</p>
            <p>2. Pastikan Anda mengunggah dokumen bukti yang sah (dalam format PDF/JPG/PNG max 5MB) saat melakukan pengajuan perubahan data.</p>
            <p>3. Jika pengajuan Anda berstatus **Dikembalikan**, periksa catatan admin untuk melengkapi kekurangan berkas kemudian ajukan ulang.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import {
  Users,
  UserCheck,
  TrendingUp,
  Award,
  CalendarX,
  CalendarDays,
  UserCog,
  FileText,
  Key,
  CheckCircle,
  Clock,
  AlertCircle
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const stats = ref({})

const currentFormattedDate = computed(() => {
  const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }
  return new Date().toLocaleDateString('id-ID', options)
})

const getBarHeight = (value, max) => {
  if (!value) return '0%'
  const percentage = (value / max) * 100
  return `${Math.min(percentage, 100)}%`
}

const formatDate = (dateStr) => {
  if (!dateStr || dateStr.startsWith('0001-01-01')) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const goToEmployees = (filter) => {
  if (filter === 'pensiun') {
    router.push('/pegawai/pensiun')
  } else {
    router.push({ path: '/pegawai', query: { filter } })
  }
}

onMounted(async () => {
  try {
    const response = await api.get('/dashboard/summary')
    stats.value = response.data
  } catch (error) {
    console.error('Gagal mengambil summary dashboard:', error)
  }
})
</script>
