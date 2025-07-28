export interface User {
  id: number
  username: string
  email: string
  role_id: number
  role: string
}

// 新增：定义获取用户列表的响应体结构
export interface UsersResponse {
  users: User[];
  // 如果后端还返回了其他字段，可以加在这里
  // total?: number;
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