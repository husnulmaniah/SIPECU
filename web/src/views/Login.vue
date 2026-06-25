<template>
  <!-- Full-page login: split layout navy kiri + offwhite kanan -->
  <div class="min-h-screen flex font-sans">

    <!-- Kiri: Branding Navy Panel -->
    <div class="hidden lg:flex flex-col justify-between w-1/2 p-12 relative overflow-hidden"
         style="background: linear-gradient(160deg, #1A365D 0%, #0F2040 100%);">
      <!-- Decorative circles -->
      <div class="absolute -top-32 -left-32 w-80 h-80 rounded-full opacity-10" style="background: #D69E2E;"></div>
      <div class="absolute -bottom-20 -right-20 w-64 h-64 rounded-full opacity-10" style="background: #D69E2E;"></div>
      <div class="absolute top-1/3 right-0 w-48 h-48 rounded-full opacity-5" style="background: #FFFFFF;"></div>

      <!-- Brand -->
      <div class="relative z-10 flex items-center gap-4">
        <div class="w-14 h-14 rounded-2xl p-2 flex items-center justify-center" style="background: rgba(255,255,255,0.12);">
          <img src="/logo-morowali-utara.png" alt="Logo" class="w-full h-full object-contain" />
        </div>
        <div>
          <h1 class="text-2xl font-extrabold tracking-[0.15em] text-white">SIPECUT</h1>
          <p class="text-xs font-medium" style="color: rgba(191,219,254,0.7);">Sistem Kepegawaian Digital</p>
        </div>
      </div>

      <!-- Center Quote -->
      <div class="relative z-10 space-y-6">
        <div class="w-12 h-1 rounded-full" style="background: #D69E2E;"></div>
        <h2 class="text-3xl font-bold text-white leading-snug">
          Kelola Data<br/>Kepegawaian<br/><span style="color: #ECC94B;">Lebih Cerdas</span>
        </h2>
        <p class="text-sm leading-relaxed" style="color: rgba(191,219,254,0.7);">
          Platform digital terintegrasi untuk manajemen kepegawaian Dinas Pendidikan dan Kebudayaan Daerah Kabupaten Morowali Utara.
        </p>
        <!-- Feature Badges -->
        <div class="flex flex-wrap gap-2 pt-2">
          <span v-for="f in ['Data Pegawai', 'Pengajuan Cuti', 'KGB & Pangkat', 'Berita Acara']" :key="f"
                class="px-3 py-1 text-xs font-semibold rounded-full"
                style="background: rgba(255,255,255,0.1); color: rgba(254,243,199,0.85); border: 1px solid rgba(255,255,255,0.12);">
            {{ f }}
          </span>
        </div>
      </div>

      <!-- Footer Branding -->
      <div class="relative z-10">
        <p class="text-[11px]" style="color: rgba(148,163,184,0.6);">
          Hak Cipta &copy; 2026 · Dinas Pendidikan dan Kebudayaan Daerah Morowali Utara
        </p>
      </div>
    </div>

    <!-- Kanan: Login Form Panel -->
    <div class="flex-1 flex items-center justify-center p-8" style="background-color: #F7FAFC;">
      <div class="w-full max-w-sm">
        <!-- Mobile Logo -->
        <div class="lg:hidden text-center mb-8">
          <img src="/logo-morowali-utara.png" alt="Logo" class="w-16 h-16 mx-auto object-contain mb-3" />
          <h1 class="text-xl font-extrabold tracking-widest" style="color: #1A365D;">SIPECUT</h1>
          <p class="text-xs mt-1" style="color: #718096;">Sistem Kepegawaian Digital</p>
        </div>

        <!-- Form Header -->
        <div class="mb-8">
          <h2 class="text-2xl font-extrabold" style="color: #1A365D;">Selamat Datang</h2>
          <p class="text-sm mt-1" style="color: #718096;">Masuk ke akun Anda untuk melanjutkan</p>
        </div>

        <!-- Session Expired -->
        <div v-if="sessionExpired" class="mb-5 p-3 rounded-xl flex items-center gap-2 text-xs font-medium"
             style="background: #FFFBEB; border: 1px solid #FCD34D; color: #92400E;">
          <AlertTriangle class="w-4 h-4 shrink-0" />
          <span>Sesi Anda telah berakhir. Silakan login kembali.</span>
        </div>

        <!-- Error -->
        <div v-if="errorMsg" class="mb-5 p-3 rounded-xl flex items-center gap-2 text-xs font-medium"
             style="background: #FEF2F2; border: 1px solid #FCA5A5; color: #991B1B;">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- NIP -->
          <div>
            <label for="nip" class="block text-xs font-bold uppercase tracking-wider mb-2" style="color: #2D3748;">
              NIP / Username
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center" style="color: #A0AEC0;">
                <User class="w-4 h-4" />
              </span>
              <input
                v-model="nip"
                type="text"
                id="nip"
                placeholder="Masukkan NIP Anda"
                required
                class="w-full pl-10 pr-4 py-3 rounded-xl text-sm outline-none transition-all"
                style="background: #FFFFFF; border: 1.5px solid #CBD5E0; color: #2D3748;"
                @focus="e => e.target.style.borderColor='#1A365D'"
                @blur="e => e.target.style.borderColor='#CBD5E0'"
              />
            </div>
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="block text-xs font-bold uppercase tracking-wider mb-2" style="color: #2D3748;">
              Kata Sandi
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center" style="color: #A0AEC0;">
                <Lock class="w-4 h-4" />
              </span>
              <input
                v-model="password"
                type="password"
                id="password"
                placeholder="••••••••"
                required
                class="w-full pl-10 pr-4 py-3 rounded-xl text-sm outline-none transition-all"
                style="background: #FFFFFF; border: 1.5px solid #CBD5E0; color: #2D3748;"
                @focus="e => e.target.style.borderColor='#1A365D'"
                @blur="e => e.target.style.borderColor='#CBD5E0'"
              />
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 rounded-xl text-sm font-bold flex items-center justify-center gap-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed mt-2"
            style="background: linear-gradient(135deg, #1A365D 0%, #2A4A7F 100%); color: #FFFFFF; box-shadow: 0 4px 14px rgba(26,54,93,0.3);"
            @mouseover="e => !loading && (e.target.style.boxShadow='0 6px 20px rgba(26,54,93,0.45)')"
            @mouseout="e => e.target.style.boxShadow='0 4px 14px rgba(26,54,93,0.3)'"
          >
            <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
            <span>{{ loading ? 'Memverifikasi...' : 'Masuk ke Aplikasi' }}</span>
          </button>

          <!-- Gold Accent Line -->
          <div class="flex items-center gap-3 pt-1">
            <div class="flex-1 h-px" style="background: #E2E8F0;"></div>
            <div class="w-6 h-1 rounded-full" style="background: #D69E2E;"></div>
            <div class="flex-1 h-px" style="background: #E2E8F0;"></div>
          </div>
        </form>

        <!-- Footer -->
        <p class="text-center text-[10px] mt-6" style="color: #A0AEC0;">
          Hak Cipta &copy; 2026 SIPECUT · Dinas Pendidikan dan Kebudayaan Daerah Morowali Utara
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, Lock, AlertCircle, AlertTriangle, Loader2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const nip = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')
const sessionExpired = ref(false)

onMounted(() => {
  if (route.query.expired) {
    sessionExpired.value = true
  }
})

const handleLogin = async () => {
  loading.value = true
  errorMsg.value = ''
  sessionExpired.value = false

  try {
    const success = await authStore.login(nip.value, password.value)
    if (success) {
      router.push('/dashboard')
    }
  } catch (err) {
    errorMsg.value = err
  } finally {
    loading.value = false
  }
}
</script>
