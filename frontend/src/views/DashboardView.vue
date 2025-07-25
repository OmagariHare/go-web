import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token'))
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const isAuthenticated = ref(!!token.value)

  function setAuth(newToken: string, newUser: any) {
    token.value = newToken
    user.value = newUser
    isAuthenticated.value = true
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function logout() {
    token.value = null
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function checkAuth() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    token.value = storedToken
    user.value = storedUser ? JSON.parse(storedUser) : null
    isAuthenticated.value = !!storedToken
  }

  return { token, user, isAuthenticated, setAuth, logout, checkAuth }
})
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(true)
const error = ref('')

// 检查认证状态
const checkAuth = () => {
  authStore.checkAuth()
  if (!authStore.isAuthenticated) {
    router.push('/')
  }
}

onMounted(() => {
  checkAuth()
  loading.value = false
})

// 登出功能
const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="dashboard-container">
    <div class="header">
      <h1>仪表板</h1>
      <button @click="handleLogout" class="logout-button">登出</button>
    </div>
    
    <div v-if="loading" class="loading">
      加载中...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else-if="authStore.user" class="user-info">
      <h2>欢迎, {{ authStore.user.username }}!</h2>
      <div class="info-card">
        <h3>用户信息</h3>
        <p><strong>用户名:</strong> {{ authStore.user.username }}</p>
        <p><strong>邮箱:</strong> {{ authStore.user.email }}</p>
        <p><strong>角色:</strong> {{ authStore.user.role }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.logout-button {
  padding: 0.5rem 1rem;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.logout-button:hover {
  background-color: #f78989;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  font-size: 1.2rem;
}

.error {
  color: #f56c6c;
}

.user-info h2 {
  color: #333;
  margin-bottom: 1.5rem;
}

.info-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.info-card h3 {
  margin-top: 0;
  color: #333;
}

.info-card p {
  margin: 0.5rem 0;
  color: #666;
}
</style>