import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import HomeView from '@/views/HomeView.vue'
import ProjectView from '@/views/ProjectView.vue'

const routes = [
  {path: '/', component: LoginView, alias: '/login', name: 'login'},

  {path: '/projects', component: HomeView, alias: '/home', name: 'projects'},

  {path: '/project/:projectId', component: ProjectView},
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router