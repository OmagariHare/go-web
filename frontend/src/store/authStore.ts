import { defineStore } from 'pinia'
import authService from '../services/authService'
import { User } from '../types/user'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token')
  }),

  actions: {
    /**
     * 用户登录
     * @param username 用户名
     * @param password 密码
     */
    async login(username: string, password: string) {
      try {
        const response = await authService.login({ username, password })
        this.user = response.user
        this.token = response.token
        this.isAuthenticated = true
        
        // 保存到localStorage
        localStorage.setItem('token', response.token)
        localStorage.setItem('user', JSON.stringify(response.user))
        
        return response
      } catch (error) {
        throw error
      }
    },

    /**
     * 用户登出
     */
    logout() {
      authService.logout()
      this.user = null
      this.token = null
      this.isAuthenticated = false
    },

    /**
     * 检查认证状态
     */
    checkAuth() {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('user')
      
      if (token && user) {
        this.token = token
        this.user = JSON.parse(user)
        this.isAuthenticated = true
      } else {
        this.logout()
      }
    }
  }
})