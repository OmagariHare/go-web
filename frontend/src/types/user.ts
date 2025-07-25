export interface User {
  id: number
  username: string
  email: string
  role_id: number
  role: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}