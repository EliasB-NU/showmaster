import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import ProjectView from '@/views/ProjectView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import LoginView from '@/views/LoginView.vue'

const routes = [
  { path: '/login', name: 'login', component: LoginView, alias: '/' },

  {path: '/projects', name: 'projects', component: HomeView, alias: '/home'},

  {path: '/project/:projectId', component: ProjectView},

  { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView, alias: '/404' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router