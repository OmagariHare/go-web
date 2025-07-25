import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { router } from './router'
import { store } from './store' // 假设store定义在./store文件中

import App from './App.vue'
import './assets/main.css'

const app = createApp(App)
const pinia = createPinia() // 创建Pinia实例

app.use(pinia) // 使用Pinia
app.use(router)

app.mount('#app')