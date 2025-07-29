import axios, { type AxiosInstance } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'
import type { 
  LoginRequest, 
  RegisterRequest, 
  LoginResponse, 
  RefreshTokenResponse,
  UpdateUserRequest,
  UsersListResponse,
  User
} from '@/types/user'

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  async error => {
    const userStore = useUserStore()
    
    if (error.response?.status === 401) {
      // Token过期，尝试刷新
      if (userStore.refreshToken) {
        try {
          await userStore.refreshTokens()
          // 重新发送原请求
          return api(error.config)
        } catch (refreshError) {
          // 刷新失败，跳转到登录页
          userStore.logout()
          router.push('/login')
          ElMessage.error('登录已过期，请重新登录')
        }
      } else {
        userStore.logout()
        router.push('/login')
      }
    } else if (error.response?.status >= 500) {
      ElMessage.error('服务器错误，请稍后重试')
    } else if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    }
    
    return Promise.reject(error)
  }
)

// 认证相关API
export const authAPI = {
  register: (data: RegisterRequest) => api.post<User>('/auth/register', data),
  login: (data: LoginRequest) => api.post<LoginResponse>('/auth/login', data),
  logout: () => api.post<{ message: string }>('/auth/logout'),
  refreshToken: (refreshToken: string) => api.post<RefreshTokenResponse>('/auth/refresh', { refresh_token: refreshToken })
}

// 用户相关API
export const userAPI = {
  getUsers: (params?: { page?: number; limit?: number }) => api.get<UsersListResponse>('/users', { params }),
  getUser: (id: number) => api.get<User>(`/users/${id}`),
  updateUser: (id: number, data: UpdateUserRequest) => api.put<User>(`/users/${id}`, data),
  deleteUser: (id: number) => api.delete<{ message: string }>(`/users/${id}`),
  getProfile: () => api.get<User>('/users/profile'),
  updateProfile: (data: UpdateUserRequest) => api.put<User>('/users/profile', data)
}

export default api