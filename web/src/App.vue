<template>
  <div class="min-h-screen font-sans selection:bg-gold-200 selection:text-navy-800" style="background-color: #F7FAFC; color: #2D3748;">
    <!-- Unauthenticated Layout (Login Page) -->
    <template v-if="isLoginPage">
      <router-view />
    </template>

    <!-- Authenticated Dashboard Layout -->
    <template v-else>
      <div class="flex flex-1 relative overflow-hidden min-h-screen">

        <!-- ── Sidebar Desktop ───────────────────────────────── -->
        <aside class="hidden md:flex flex-col w-64 shrink-0 z-30 border-r border-navy-600/30"
               style="background: linear-gradient(180deg, #1A365D 0%, #152D4E 60%, #0F2040 100%); min-height: 100vh;">
          <!-- Logo & Brand -->
          <div class="p-6 border-b border-white/10 flex items-center gap-3">
            <div class="w-11 h-11 shrink-0 rounded-xl overflow-hidden bg-white/10 p-1">
              <img src="/logo-morowali-utara.png" alt="Logo Morowali Utara" class="w-full h-full object-contain" />
            </div>
            <div>
              <h1 class="font-extrabold text-lg tracking-widest text-white">SIPECUT</h1>
              <p class="text-[11px] text-blue-200/70 font-medium leading-tight">Sistem Kepegawaian</p>
            </div>
          </div>

          <!-- Nav Links -->
          <nav class="flex-1 px-3 py-5 space-y-0.5 overflow-y-auto">
            <router-link
              v-for="item in menuItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-200"
              :class="activeRoute === item.path
                ? 'bg-gold-500 text-navy-900 shadow-lg shadow-gold-500/20 font-bold'
                : 'text-blue-100/70 hover:bg-white/10 hover:text-white'"
            >
              <component :is="item.icon" class="w-4 h-4 shrink-0" />
              <span>{{ item.name }}</span>
            </router-link>
          </nav>

          <!-- User Footer -->
          <div class="p-4 border-t border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-9 h-9 rounded-full bg-gold-500/20 border border-gold-400/40 flex items-center justify-center font-bold text-gold-300 overflow-hidden text-sm">
                <img v-if="authStore.employee?.foto_profil" :src="authStore.employee.foto_profil" class="w-full h-full object-cover" />
                <span v-else>{{ authStore.employee?.nama?.charAt(0) || authStore.role?.charAt(0)?.toUpperCase() }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-white truncate">{{ authStore.employee?.nama || 'Administrator' }}</p>
                <p class="text-[10px] text-blue-200/60 truncate">{{ authStore.role === 'admin' ? 'Admin' : authStore.nip }}</p>
              </div>
            </div>
            <button
              @click="logout"
              class="w-full flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-semibold transition-all duration-200 border border-white/10 text-blue-100/60 hover:bg-red-900/30 hover:text-red-300 hover:border-red-500/30"
            >
              <LogOut class="w-4 h-4" />
              <span>Keluar Sesi</span>
            </button>
          </div>
        </aside>

        <!-- ── Sidebar Mobile Overlay ────────────────────────── -->
        <div v-if="mobileSidebarOpen"
             class="md:hidden fixed inset-0 z-40 backdrop-blur-sm"
             style="background: rgba(10,20,40,0.55);"
             @click="mobileSidebarOpen = false">
        </div>
        <aside
          class="md:hidden fixed top-0 bottom-0 left-0 w-64 z-50 flex flex-col border-r border-white/10 transition-transform duration-300 transform"
          style="background: linear-gradient(180deg, #1A365D 0%, #152D4E 60%, #0F2040 100%);"
          :class="mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
        >
          <div class="p-5 border-b border-white/10 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 shrink-0 rounded-lg overflow-hidden bg-white/10 p-0.5">
                <img src="/logo-morowali-utara.png" alt="Logo" class="w-full h-full object-contain" />
              </div>
              <h1 class="font-extrabold text-base tracking-widest text-white">SIPECUT</h1>
            </div>
            <button @click="mobileSidebarOpen = false" class="text-blue-200/60 hover:text-white">
              <X class="w-6 h-6" />
            </button>
          </div>
          <nav class="flex-1 px-3 py-5 space-y-0.5 overflow-y-auto" @click="mobileSidebarOpen = false">
            <router-link
              v-for="item in menuItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-200"
              :class="activeRoute === item.path
                ? 'bg-gold-500 text-navy-900 font-bold shadow-lg shadow-gold-500/20'
                : 'text-blue-100/70 hover:bg-white/10 hover:text-white'"
            >
              <component :is="item.icon" class="w-4 h-4 shrink-0" />
              <span>{{ item.name }}</span>
            </router-link>
          </nav>
          <div class="p-4 border-t border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-9 h-9 rounded-full bg-gold-500/20 border border-gold-400/40 flex items-center justify-center font-bold text-gold-300 overflow-hidden text-sm">
                <img v-if="authStore.employee?.foto_profil" :src="authStore.employee.foto_profil" class="w-full h-full object-cover" />
                <span v-else>{{ authStore.employee?.nama?.charAt(0) || authStore.role?.charAt(0)?.toUpperCase() }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-white truncate">{{ authStore.employee?.nama || 'Administrator' }}</p>
                <p class="text-[10px] text-blue-200/60 truncate">{{ authStore.role === 'admin' ? 'Admin' : authStore.nip }}</p>
              </div>
            </div>
            <button @click="logout"
              class="w-full flex items-center justify-center gap-2 py-2 px-3 rounded-lg text-xs font-semibold transition-all border border-white/10 text-blue-100/60 hover:bg-red-900/30 hover:text-red-300">
              <LogOut class="w-4 h-4" /><span>Keluar Sesi</span>
            </button>
          </div>
        </aside>

        <!-- ── Main Content ───────────────────────────────────── -->
        <div class="flex-1 flex flex-col min-w-0 overflow-y-auto">
          <!-- Top Header -->
          <header class="h-16 sticky top-0 z-20 flex items-center justify-between px-6 border-b"
                  style="background: rgba(247,250,252,0.95); backdrop-filter: blur(10px); border-color: #CBD5E0; box-shadow: 0 1px 8px rgba(26,54,93,0.07);">
            <div class="flex items-center gap-4">
              <button @click="mobileSidebarOpen = true" class="md:hidden" style="color: #4A5568;">
                <Menu class="w-6 h-6" />
              </button>
              <!-- Accent bar kiri judul -->
              <div class="flex items-center gap-3">
                <div class="hidden sm:block w-1 h-6 rounded-full" style="background: #D69E2E;"></div>
                <h2 class="text-base font-bold capitalize" style="color: #1A365D;">{{ pageTitle }}</h2>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <router-link to="/profil" class="flex items-center gap-2.5 hover:opacity-80 transition-opacity">
                <div class="hidden sm:block text-right">
                  <p class="text-xs font-semibold" style="color: #1A365D;">{{ authStore.employee?.nama || 'Admin' }}</p>
                  <p class="text-[10px] capitalize" style="color: #718096;">{{ authStore.role }}</p>
                </div>
                <div class="w-8 h-8 rounded-full border-2 flex items-center justify-center text-sm font-bold overflow-hidden"
                     style="border-color: #D69E2E; background: #EBF0F8; color: #1A365D;">
                  <img v-if="authStore.employee?.foto_profil" :src="authStore.employee.foto_profil" class="w-full h-full object-cover" />
                  <span v-else>{{ authStore.employee?.nama?.charAt(0) || authStore.role?.charAt(0)?.toUpperCase() }}</span>
                </div>
              </router-link>
            </div>
          </header>

          <!-- Page Content -->
          <main class="flex-1 p-6" style="background-color: #F7FAFC;">
            <router-view />
          </main>

          <!-- Footer -->
          <footer class="px-6 py-3 text-center text-[11px] border-t" style="color: #718096; border-color: #E2E8F0; background: #EDF2F7;">
            Hak Cipta &copy; 2026 SIPECUT Dinas Pendidikan dan Kebudayaan Daerah Morowali Utara. All rights reserved.
          </footer>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Swal from 'sweetalert2'
import {
  LayoutDashboard,
  Users,
  UserCheck,
  CalendarDays,
  FileCheck,
  FileText,
  Settings,
  User,
  LogOut,
  Menu,
  X
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const mobileSidebarOpen = ref(false)

const isLoginPage = computed(() => route.name === 'Login')
const activeRoute = computed(() => route.path)

const pageTitle = computed(() => {
  if (route.name === 'Dashboard') return 'Dashboard'
  if (route.path.startsWith('/pegawai/pensiun')) return 'Arsip Pegawai Pensiun'
  if (route.path.startsWith('/pegawai')) return 'Data Kepegawaian'
  if (route.path.startsWith('/cuti')) return 'Pengajuan Cuti'
  if (route.path.startsWith('/perubahan-data')) return 'Pengajuan Perubahan Data'
  if (route.path.startsWith('/berita-acara')) return 'Dokumen Berita Acara'
  if (route.path.startsWith('/master-aturan')) return 'Pengaturan Aturan Master'
  if (route.path.startsWith('/profil')) return 'Profil Saya'
  return ''
})

const menuItems = computed(() => {
  const items = [
    { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard }
  ]

  if (authStore.role === 'admin') {
    items.push(
      { name: 'Data Pegawai Aktif', path: '/pegawai', icon: Users },
      { name: 'Arsip Pegawai Pensiun', path: '/pegawai/pensiun', icon: UserCheck },
      { name: 'Pengajuan Cuti', path: '/cuti', icon: CalendarDays },
      { name: 'Perubahan Data', path: '/perubahan-data', icon: FileCheck },
      { name: 'Berita Acara', path: '/berita-acara', icon: FileText },
      { name: 'Aturan Master', path: '/master-aturan', icon: Settings }
    )
  } else {
    items.push(
      { name: 'Profil Pegawai', path: `/pegawai/${authStore.nip}`, icon: User },
      { name: 'Pengajuan Cuti', path: '/cuti', icon: CalendarDays },
      { name: 'Perubahan Data', path: '/perubahan-data', icon: FileCheck },
      { name: 'Berita Acara', path: '/berita-acara', icon: FileText }
    )
  }

  items.push({ name: 'Profil & Akun', path: '/profil', icon: User })

  return items
})

const logout = async () => {
  const result = await Swal.fire({
    title: 'Keluar Sesi?',
    text: 'Apakah Anda yakin ingin keluar dari aplikasi?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#1A365D',
    cancelButtonColor: '#E11D48',
    confirmButtonText: 'Ya, Keluar',
    cancelButtonText: 'Batal',
    background: '#1E293B',
    color: '#F8FAFC'
  })
  if (result.isConfirmed) {
    authStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  authStore.updateSelfProfile()
})
</script>
