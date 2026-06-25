import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/pegawai',
    name: 'Pegawai',
    component: () => import('@/views/EmployeeList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/pegawai/pensiun',
    name: 'PegawaiPensiun',
    component: () => import('@/views/EmployeeRetiredList.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/pegawai/:nip',
    name: 'PegawaiDetail',
    component: () => import('@/views/EmployeeDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/cuti',
    name: 'LeaveRequestList',
    component: () => import('@/views/LeaveRequestList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/cuti/tambah',
    name: 'LeaveRequestCreate',
    component: () => import('@/views/LeaveRequestCreate.vue'),
    meta: { requiresAuth: true, role: 'employee' }
  },
  {
    path: '/perubahan-data',
    name: 'DataChangeList',
    component: () => import('@/views/DataChangeList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/berita-acara',
    name: 'BeritaAcara',
    component: () => import('@/views/BeritaAcaraList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/master-aturan',
    name: 'MasterRules',
    component: () => import('@/views/MasterRules.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/profil',
    name: 'SelfProfile',
    component: () => import('@/views/SelfProfile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Route Guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Guest checking (redirect away from login if already authenticated)
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    return next('/dashboard')
  }

  // Auth checking
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next('/login')
  }

  // Role checking (Admin restriction)
  if (to.meta.role && authStore.role !== to.meta.role) {
    // If employee, redirect /pegawai directly to self details page
    if (to.path === '/pegawai' && authStore.role === 'employee') {
      return next(`/pegawai/${authStore.nip}`)
    }
    // Otherwise return to dashboard
    return next('/dashboard')
  }

  // If going to /pegawai and is employee, redirect to /pegawai/:nip
  if (to.path === '/pegawai' && authStore.role === 'employee') {
    return next(`/pegawai/${authStore.nip}`)
  }

  next()
})

export default router
