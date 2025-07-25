import http from '../utils/http'
import { LoginRequest, AuthResponse, RegisterRequest } from '../types/user'

class AuthService {
  /**
   * 用户登录
   * @param data 登录信息
   * @returns 认证响应
   */
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await http.post<AuthResponse>('/auth/login', data)
    return response.data
  }

  /**
   * 用户注册
   * @param data 注册信息
   * @returns 认证响应
   */
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await http.post<AuthResponse>('/auth/register', data)
    return response.data
  }

  /**
   * 用户登出
   */
  logout(): void {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  /**
   * 检查用户是否已认证
   * @returns 是否已认证
   */
  isAuthenticated(): boolean {
    const token = localStorage.getItem('token')
    return !!token
  }

  /**
   * 获取存储的用户信息
   * @returns 用户信息或null
   */
  getCurrentUser(): any {
    const user = localStorage.getItem('user')
    return user ? JSON.parse(user) : null
  }
}

export default new AuthService()