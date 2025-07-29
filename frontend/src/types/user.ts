export interface User {
  id: number
  username: string
  email: string
  is_active: boolean
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user: User
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  token: string
  refresh_token: string
}

export interface UpdateUserRequest {
  username?: string
  email?: string
  password?: string
  is_active?: boolean
}

export interface UsersListResponse {
  users: User[]
  total: number
  page: number
  limit: number
}