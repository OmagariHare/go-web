<template>
  <el-container class="h-screen">
    <el-aside width="200px" class="bg-gray-800 text-white">
      <div class="p-4 text-2xl font-bold">管理后台</div>
      <el-menu
        default-active="1"
        class="el-menu-vertical-demo bg-gray-800 text-white"
        active-text-color="#ffd04b"
        background-color="#545c64"
        text-color="#fff"
      >
        <el-menu-item index="1">
          <el-icon><i-ep-menu /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="2">
          <el-icon><i-ep-user /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="flex justify-between items-center bg-white border-b">
        <div></div>
        <el-dropdown>
          <span class="el-dropdown-link">
            欢迎, {{ authStore.user?.username }}
            <el-icon class="el-icon--right"><i-ep-arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="bg-gray-100 p-8">
        <el-card>
          <template #header>
            <div class="text-xl font-semibold">欢迎回来！</div>
          </template>
          <div v-if="authStore.user">
            <p><strong>用户名:</strong> {{ authStore.user.username }}</p>
            <p><strong>邮箱:</strong> {{ authStore.user.email }}</p>
            <p><strong>角色:</strong> {{ authStore.user.role }}</p>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

onMounted(() => {
  authStore.checkAuth()
  if (!authStore.isAuthenticated) {
    router.push('/')
  }
})

const handleLogout = () => {
  ElMessageBox.confirm('您确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    authStore.logout()
    router.push('/')
    ElMessage({
      type: 'success',
      message: '退出成功！'
    })
  }).catch(() => {
    // Cancel
  })
}
</script>

<style scoped>
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}
.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}
</style>
