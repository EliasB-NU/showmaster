<script setup lang="ts">
import { ref } from 'vue'
import PopUp from '@/components/PopUp.vue'
import axios from 'axios'
import Cookies from 'js-cookie'
import router from '@/router'

const email = ref('')
const password = ref('')
const handleLogin = async () => {
  try {
    await axios
      .post('/api/login', {
        email: email.value,
        password: password.value,
        deviceId: Cookies.get('deviceId'),
      }, {
        headers: {
          'Content-Type': 'application/json',
        },
      })
      .then((res) => {
        if (res.status === 200) {
          Cookies.set('token', res.data.token)
          Cookies.set('admin', res.data.perms.admin);
          Cookies.set('events', res.data.perms.events);
          router.push({ name: 'home' })
        } else {
          popUp.value?.show("Login failed")
          console.error(res)
        }
      })
  } catch (error) {
    console.error('Login error:', error)
    popUp.value?.show("Login failed | Wrong email or password")
    password.value = ''
  }
}

const popUp = ref<InstanceType<typeof PopUp> | null>(null);
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <div class="w-full max-w-md bg-white p-8 rounded-2xl shadow-lg">
      <h2 class="text-2xl font-semibold text-gray-800 text-center mb-6">Login</h2>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-gray-600 text-sm mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div>
          <label class="block text-gray-600 text-sm mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <button
          type="submit"
          class="w-full bg-gray-700 text-white p-3 rounded-lg hover:bg-gray-900 transition">
          Login
        </button>
      </form>
    </div>
    <PopUp ref='popUp' />
  </div>
</template>

<style scoped>
</style>
