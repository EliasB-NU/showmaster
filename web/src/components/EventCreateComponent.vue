<script setup lang="ts">

import axios from 'axios'
import Cookies from 'js-cookie'
import { ref } from 'vue'

const emit = defineEmits(['close', 'created'])

interface Event {
  name: string
  description: string
  location: string
  startDate: string
  endDate: string
}

const event = ref<Event>({
  name: '',
  description: '',
  location: '',
  startDate: '',
  endDate: '',
})

const createEvent = async () => {
  try {
    await axios
      .post('/api/event/create', {
        name: event.value.name,
        description: event.value.description,
        location: event.value.location,
        startDate: new Date(event.value.startDate),
        endDate: new Date(event.value.endDate),
      }, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(() => {
        event.value = {} as Event
        emit('created')
      })
  } catch (error) {
    console.log(error)
    event.value = {} as Event
    emit('close')
  }
}
</script>

<template>
  <div class="fixed inset-0 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-2xl shadow-lg w-full max-w-xl p-6 relative">
      <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 text-2xl">&times;</button>

      <h2 class="text-xl font-bold mb-4">Create New Event</h2>

      <form @submit.prevent="createEvent" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Name</label>
          <input v-model="event.name" type="text" class="w-full border px-3 py-2 rounded" required />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <textarea v-model="event.description" class="w-full border px-3 py-2 rounded" required></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Location</label>
          <input v-model="event.location" type="text" class="w-full border px-3 py-2 rounded" required />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium mb-1">Start Date</label>
            <input v-model="event.startDate" type="text" class="w-full border px-3 py-2 rounded" required />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">End Date</label>
            <input v-model="event.endDate" type="text" class="w-full border px-3 py-2 rounded" required />
          </div>
        </div>

        <button
          type="submit"
          class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          Create Event
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>
