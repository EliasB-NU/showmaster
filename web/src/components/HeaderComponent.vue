<script setup lang="ts">
import Cookies from 'js-cookie'
import { ref } from 'vue'
import axios from 'axios'
import router from '@/router'

const mobileMenuOpen = ref<boolean>(false)

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
}

const logout = async () => {
  try {
    await axios
      .delete('/api/logout', {
        data: {
          token: Cookies.get("token"),
          deviceId: Cookies.get("deviceId"),
        }
      })
      .then(() => {
        Cookies.remove("token")
        Cookies.remove("admin")
        Cookies.remove("events")
        router.push({ name: "login" })
      });
  } catch (error) {
    console.error("Logout failed:", error);
  }
};
</script>

<template>
  <header class="w-full bg-gray-800 text-white shadow-md">
    <div class="max-w-7xl mx-auto px-6 flex justify-between items-center h-16">
      <!-- Logo -->
      <router-link to="/" class="text-xl font-semibold">
        Showmaster V3
      </router-link>

      <!-- Navigation (Desktop) -->
      <nav class="hidden md:flex space-x-6">
        <router-link to="/" :v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 1" class="hover:text-gray-300" >Home</router-link>
        <router-link to="/events/active" :v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 1" class="hover:text-gray-300" >Active Event</router-link>
        <router-link to="/admin" :v-if="Cookies.get('admin')" class="hover:text-gray-300" >Admin</router-link>
      </nav>

      <!-- Logout Button -->
      <button @click="logout" class="bg-red-500 px-4 py-2 rounded-lg hover:bg-red-600 transition">
        Logout
      </button>

      <!-- Mobile Menu Button -->
      <button @click="toggleMobileMenu" class="md:hidden text-white text-2xl">
        ☰
      </button>
    </div>

    <!-- Mobile Navigation -->
    <div v-if="mobileMenuOpen" class="md:hidden bg-gray-800 text-white px-4 py-6 flex flex-col space-y-4 text-lg">
      <router-link to="/" :v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 1" class="hover:text-gray-300" >Home</router-link>
      <router-link to="/events/active" :v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 1" class="hover:text-gray-300" >Active Event</router-link>
      <router-link to="/admin" :v-if="Cookies.get('admin')" class="hover:text-gray-300" >Admin</router-link>
    </div>
  </header>
</template>

<style scoped>

</style>
