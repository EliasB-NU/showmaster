import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
    },
    {
      path: '/event/:id',
      name: 'events',
      component: () => import('@/views/EventView.vue'),
    },
    {
      path: '/activeEvent',
      name: 'Active Event',
      component: () => import('@/views/ActiveEventView.vue'),
    }
  ],
})

export default router
