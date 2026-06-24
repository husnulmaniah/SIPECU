<template>
  <div class="max-w-md mx-auto space-y-6">
    <div>
      <h1 class="text-xl font-bold text-white">Profil & Keamanan Akun</h1>
      <p class="text-xs text-slate-400 mt-1">Ubah kata sandi akun login Anda untuk menjaga keamanan database.</p>
    </div>

    <!-- Password form -->
    <form @submit.prevent="handleChangePassword" class="glass-panel p-6 rounded-2xl space-y-5 text-xs">
      <h2 class="text-sm font-bold text-white border-b border-slate-800 pb-2 flex items-center gap-2">
        <Key class="w-4 h-4 text-primary-400" />
        <span>Ganti Kata Sandi</span>
      </h2>

      <!-- Feedback notifications -->
      <div v-if="successMsg" class="p-3 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl">
        {{ successMsg }}
      </div>
      <div v-if="errorMsg" class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
        {{ errorMsg }}
      </div>

      <div>
        <label class="block text-[10px] font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Kata Sandi Lama*</label>
        <input v-model="form.old_password" type="password" required placeholder="Masukkan kata sandi lama Anda"
               class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-slate-100 outline-none" />
      </div>

      <div>
        <label class="block text-[10px] font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Kata Sandi Baru*</label>
        <input v-model="form.new_password" type="password" required placeholder="Minimal 6 karakter"
               class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-slate-100 outline-none" />
      </div>

      <div>
        <label class="block text-[10px] font-semibold text-slate-350 mb-1.5 uppercase tracking-wider">Konfirmasi Kata Sandi Baru*</label>
        <input v-model="form.confirm_password" type="password" required placeholder="Ulangi kata sandi baru"
               class="w-full px-3 py-2 bg-slate-950/60 border border-slate-800 focus:border-primary-500 rounded-lg text-slate-100 outline-none" />
      </div>

      <div class="border-t border-slate-800 pt-4 flex justify-end">
        <button type="submit" :disabled="submitting"
                class="px-4 py-2 bg-gradient-to-r from-primary-600 to-violet-600 hover:from-primary-700 hover:to-violet-700 text-white rounded-xl font-bold shadow-lg shadow-primary-500/20 transition-all flex items-center gap-2">
          <Loader2 v-if="submitting" class="w-4 h-4 animate-spin" />
          <span>{{ submitting ? 'Memperbarui...' : 'Ubah Kata Sandi' }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/services/api'
import Swal from 'sweetalert2'
import { Key, Loader2 } from 'lucide-vue-next'

const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const handleChangePassword = async () => {
  errorMsg.value = ''
  successMsg.value = ''

  if (form.value.new_password !== form.value.confirm_password) {
    errorMsg.value = 'Konfirmasi kata sandi baru tidak cocok'
    Swal.fire({
      icon: 'error',
      title: 'Validasi Gagal',
      text: 'Konfirmasi kata sandi baru tidak cocok.',
      confirmButtonColor: '#1A365D',
      background: '#1E293B',
      color: '#F8FAFC'
    })
    return
  }

  if (form.value.new_password.length < 6) {
    errorMsg.value = 'Kata sandi baru minimal harus terdiri dari 6 karakter'
    Swal.fire({
      icon: 'error',
      title: 'Validasi Gagal',
      text: 'Kata sandi baru minimal harus terdiri dari 6 karakter.',
      confirmButtonColor: '#1A365D',
      background: '#1E293B',
      color: '#F8FAFC'
    })
    return
  }

  submitting.value = true

  try {
    const response = await api.put('/auth/change-password', {
      old_password: form.value.old_password,
      new_password: form.value.new_password
    })
    successMsg.value = response.data.message || 'Kata sandi berhasil diperbarui'
    Swal.fire({
      icon: 'success',
      title: 'Berhasil',
      text: successMsg.value,
      timer: 1550,
      showConfirmButton: false,
      background: '#1E293B',
      color: '#F8FAFC'
    })
    form.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (error) {
    errorMsg.value = error.response?.data?.error || 'Gagal mengubah kata sandi'
    Swal.fire({
      icon: 'error',
      title: 'Gagal',
      text: errorMsg.value,
      confirmButtonColor: '#1A365D',
      background: '#1E293B',
      color: '#F8FAFC'
    })
  } finally {
    submitting.value = false
  }
}
</script>
