<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'

interface LoginFormData {
  username: string
  password: string
}

const router = useRouter()
const authStore = useAuthStore()

const formData = ref<LoginFormData>({
  username: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(formData.value.username, formData.value.password)
    // 跳转到仪表板
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-form">
      <h2>用户登录</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="username">用户名:</label>
          <input 
            id="username"
            v-model="formData.username" 
            type="text" 
            required 
            placeholder="请输入用户名"
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码:</label>
          <input 
            id="password"
            v-model="formData.password" 
            type="password" 
            required 
            placeholder="请输入密码"
          />
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <button 
          type="submit" 
          :disabled="loading"
          class="login-button"
        >
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

input:focus {
  outline: none;
  border-color: #409eff;
}

.login-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-button:hover:not(:disabled) {
  background-color: #66b1ff;
}

.login-button:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}

.error-message {
  color: #f56c6c;
  margin-bottom: 1rem;
  text-align: center;
}
</style>