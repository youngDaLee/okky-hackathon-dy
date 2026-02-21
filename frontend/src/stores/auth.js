import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const accessToken = ref(localStorage.getItem('access_token') || null)
  const refreshToken = ref(localStorage.getItem('refresh_token') || null)

  const isAuthenticated = computed(() => !!accessToken.value)

  const login = async (email, password) => {
    try {
      const { data } = await api.post('/auth/login', { email, password })
      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      localStorage.setItem('access_token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)
      await fetchUser()
      return data
    } catch (error) {
      throw error
    }
  }

  const signup = async (email, password, nickname) => {
    try {
      const { data } = await api.post('/auth/signup', {
        email,
        password,
        nickname
      })
      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      localStorage.setItem('access_token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)
      await fetchUser()
      return data
    } catch (error) {
      throw error
    }
  }

  const fetchUser = async () => {
    try {
      const { data } = await api.get('/users/me')
      user.value = data
    } catch (error) {
      console.error('Failed to fetch user:', error)
    }
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout', {
        refresh_token: refreshToken.value
      })
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    login,
    signup,
    fetchUser,
    logout
  }
})
