import http from '../utils/http'
import { User, UsersResponse } from '../types/user'

class UserService {
  /**
   * 获取所有用户
   * @returns 用户列表
   */
  async getUsers(): Promise<User[]> {
    // 明确告诉 http.get，我们期望得到 UsersResponse 类型的响应
    const response = await http.get<UsersResponse>('/users')
    // 正确地从响应数据中返回 users 数组
    return response.data.users
  }

  /**
   * 根据ID获取用户
   * @param id 用户ID
   * @returns 用户信息
   */
  async getUserById(id: number): Promise<User> {
    const response = await http.get<User>(`/users/${id}`)
    return response.data
  }

  /**
   * 更新用户信息
   * @param id 用户ID
   * @param data 更新数据
   * @returns 更新后的用户信息
   */
  async updateUser(id: number, data: Partial<User>): Promise<User> {
    const response = await http.put<User>(`/users/${id}`, data)
    return response.data
  }

  /**
   * 删除用户
   * @param id 用户ID
   */
  async deleteUser(id: number): Promise<void> {
    await http.delete(`/users/${id}`)
  }
}

export default new UserService()