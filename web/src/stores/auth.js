import { defineStore } from 'pinia'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('sipecut_token') || null,
    refreshToken: localStorage.getItem('sipecut_refresh_token') || null,
    role: localStorage.getItem('sipecut_role') || null,
    nip: localStorage.getItem('sipecut_nip') || null,
    employee: JSON.parse(localStorage.getItem('sipecut_employee') || 'null')
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.role === 'admin'
  },
  actions: {
    async login(nip, password) {
      try {
        const response = await api.post('/auth/login', { nip, password })
        const data = response.data
        
        this.token = data.token
        this.refreshToken = data.refresh_token
        this.role = data.role
        this.nip = data.nip
        this.employee = data.employee

        localStorage.setItem('sipecut_token', data.token)
        localStorage.setItem('sipecut_refresh_token', data.refresh_token)
        localStorage.setItem('sipecut_role', data.role)
        localStorage.setItem('sipecut_nip', data.nip)
        localStorage.setItem('sipecut_employee', JSON.stringify(data.employee))

        return true
      } catch (error) {
        throw error.response?.data?.error || 'Koneksi ke server gagal'
      }
    },
    logout() {
      this.token = null
      this.refreshToken = null
      this.role = null
      this.nip = null
      this.employee = null

      localStorage.removeItem('sipecut_token')
      localStorage.removeItem('sipecut_refresh_token')
      localStorage.removeItem('sipecut_role')
      localStorage.removeItem('sipecut_nip')
      localStorage.removeItem('sipecut_employee')
    },
    async updateSelfProfile() {
      if (!this.nip || this.role === 'admin' && this.nip === 'admin') return
      try {
        const response = await api.get(`/employees/${this.nip}`)
        const data = response.data
        this.employee = data.employee
        localStorage.setItem('sipecut_employee', JSON.stringify(data.employee))
      } catch (error) {
        console.error('Gagal memperbarui profil store:', error)
      }
    }
  }
})
