import type { AxiosError } from 'axios'

export interface ApiResponse<T = any> {
  data: T
  message?: string
}

export interface ApiError {
  error: string
  message?: string
}

export type ApiErrorResponse = AxiosError<ApiError>