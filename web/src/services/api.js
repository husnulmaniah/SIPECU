import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request Interceptor: Attach JWT Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('sipecut_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor: Auto-Refresh Token on 401
let isRefreshing = false
let failedQueue = []

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  
  failedQueue = []
}

api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return api(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      const refreshToken = localStorage.getItem('sipecut_refresh_token')
      if (!refreshToken) {
        // No refresh token, trigger logout/redirect
        handleSessionExpired()
        return Promise.reject(error)
      }

      try {
        const response = await axios.post('/api/auth/refresh', { refresh_token: refreshToken })
        const { token, refresh_token: newRefreshToken } = response.data
        
        localStorage.setItem('sipecut_token', token)
        localStorage.setItem('sipecut_refresh_token', newRefreshToken)
        
        api.defaults.headers.common.Authorization = `Bearer ${token}`
        originalRequest.headers.Authorization = `Bearer ${token}`
        
        processQueue(null, token)
        isRefreshing = false
        
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        isRefreshing = false
        handleSessionExpired()
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

function handleSessionExpired() {
  localStorage.removeItem('sipecut_token')
  localStorage.removeItem('sipecut_refresh_token')
  localStorage.removeItem('sipecut_role')
  localStorage.removeItem('sipecut_nip')
  
  // Only redirect if we are not already on login page
  if (window.location.pathname !== '/login') {
    window.location.href = '/login?expired=1'
  }
}

export default api
