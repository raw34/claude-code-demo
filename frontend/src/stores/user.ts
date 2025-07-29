import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI, userAPI } from '@/api'
import type { User, LoginRequest, RegisterRequest, UpdateUserRequest } from '@/types/user'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')

  const isAuthenticated = computed(() => !!token.value)

  const login = async (credentials: LoginRequest) => {
    const response = await authAPI.login(credentials)
    const { token: newToken, refresh_token, user: userData } = response.data
    
    token.value = newToken
    refreshToken.value = refresh_token
    user.value = userData
    
    localStorage.setItem('token', newToken)
    localStorage.setItem('refreshToken', refresh_token)
  }

  const register = async (userData: RegisterRequest) => {
    const response = await authAPI.register(userData)
    return response.data
  }

  const logout = async () => {
    try {
      await authAPI.logout()
    } finally {
      token.value = ''
      refreshToken.value = ''
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
    }
  }

  const fetchProfile = async () => {
    const response = await userAPI.getProfile()
    user.value = response.data
  }

  const updateProfile = async (updates: UpdateUserRequest) => {
    const response = await userAPI.updateProfile(updates)
    user.value = response.data
  }

  const refreshTokens = async () => {
    const response = await authAPI.refreshToken(refreshToken.value)
    const { token: newToken, refresh_token } = response.data
    
    token.value = newToken
    refreshToken.value = refresh_token
    
    localStorage.setItem('token', newToken)
    localStorage.setItem('refreshToken', refresh_token)
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    logout,
    fetchProfile,
    updateProfile,
    refreshTokens
  }
})