import { createPinia } from 'pinia'
import { createApp } from 'vue'
import i18n from '@/lib/locals'
import router from '@/lib/router'
import App from './App.vue'

import './styles/globals.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.mount('#app')
