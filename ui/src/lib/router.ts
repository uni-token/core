import { createRouter, createWebHistory } from 'vue-router'
import About from '@/views/About.vue'
import ActionHandler from '@/views/ActionHandler.vue'
import AppManagement from '@/views/AppManagement.vue'
import Models from '@/views/Models.vue'
import NotFound from '@/views/NotFound.vue'
import Settings from '@/views/Settings.vue'
import Usage from '@/views/Usage.vue'

const routes = [
  {
    path: '/',
    redirect: '/apps',
  },
  {
    path: '/apps',
    name: 'AppManagement',
    component: AppManagement,
  },
  {
    path: '/models',
    name: 'Models',
    component: Models,
  },
  {
    path: '/usage',
    name: 'Usage',
    component: Usage,
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
  },
  {
    path: '/about',
    name: 'About',
    component: About,
  },
  {
    path: '/action/:actionType',
    name: 'ActionHandler',
    component: ActionHandler,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
