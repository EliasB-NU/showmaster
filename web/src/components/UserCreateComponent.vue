<script setup lang="ts">

import { ref } from 'vue'
import axios from 'axios'
import Cookies from 'js-cookie'

defineProps({
  visible: {
    type: Boolean,
    required: true,
  },
})

const emit = defineEmits(['close', 'createdUser'])

interface User {
  email: string
  name: string
  password: string
  permissions: {
    login: boolean
    admin: boolean
    events: number
  }
}

const user = ref<User>({
  email: '',
  name: '',
  password: '',
  permissions: {
    login: false,
    admin: false,
    events: 0,
  },
})

const createUser = async () => {
  try {
    await axios
      .post('/api/createUser', {
        ...user.value,
      },{
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        },
      })
      .then((response) => {
        if (response.status === 200) {
          emit('createdUser', response.data)
          closePopup()
        } else {
          console.error('Error creating user:', response.data)
        }
      })
  } catch (error) {
    console.error('Error creating user:', error)
    emit('createdUser', null)
  }
}

const closePopup = () => {
  emit('close')
}

</script>

<template>
  <div v-if="visible" class="fixed inset-0 bg-opacity-50 flex items-center justify-center">
    <div class="bg-white p-6 rounded-2xl shadow-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
      <h3 class="text-lg font-semibold mb-4">Create User</h3>

      <!-- Form to create user -->
      <div class="mb-4">
        <form @submit.prevent="createUser">
          <div class="grid grid-cols-2 gap-4">
            <!-- Left Column: Basic Information -->
            <div class="w-1/2 pr-4">
              <div class="mb-4">
                <label for="name" class="block">Name:</label>
                <input v-model="user.name" id="name" type="text" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="email" class="block">Email:</label>
                <input v-model="user.email" id="email" type="email" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="password" class="block">Password:</label>
                <input v-model="user.password" id="password" type="password" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
            </div>

            <!-- Right Column: Permissions -->
            <div class="w-1/2 pl-4">
              <label class="block mb-2">Permissions: </label>
              <div class="permissions-grid">
                <div>
                  <label for="login">Login: </label>
                  <input v-model="user.permissions.login" type="checkbox" id="login" />
                </div>
                <div>
                  <label for="login">Admin: </label>
                  <input v-model="user.permissions.admin" type="checkbox" id="admin" />
                </div>
                <div>
                  <label for="events">Events: </label>
                  <input v-model="user.permissions.events" type="number" id="events" min="0" max="3" />
                </div>
              </div>
            </div>
          </div>

          <div class="p-3 flex space-x-2 mt-4 justify-end">
            <button type="submit" class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700">Create</button>
            <button type="button" @click="closePopup" class="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700" >Cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
