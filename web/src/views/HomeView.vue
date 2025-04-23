<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'

import { onMounted, ref } from 'vue'
import axios from 'axios'
import Cookies from 'js-cookie'
import CreateEventComponent from '@/components/CreateEventComponent.vue'

interface Event {
  ID: number
  Name: string
  Description: string
  Location: string
  StartDate: string
  EndDate: string

  CurrentScene: number
  TimeElapsed: number
}

const events = ref<Event[]>([])
const loadedEvent = ref<Number>(0)
const loading = ref<boolean>(true)

const fetchEvents = async () => {
  try {
    await axios
      .get('/api/events', {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(response => {
        events.value = response.data.events
        loadedEvent.value = response.data.active
        loading.value = false
      })
  } catch (error) {
    console.log(error)

  }
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDuration = (seconds: number): string => {
  const hrs = Math.floor(seconds / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  return `${hrs}h ${mins}m ${secs}s`
}

const showCreateModal = ref<boolean>(false)

const eventCreated = async () => {
  showCreateModal.value = false
  await fetchEvents()
}

onMounted(async () => {
  await fetchEvents()
})

</script>

<template>
  <div class="flex flex-col min-h-screen">
    <HeaderComponent />
    <div class="bg-white py-12 px-4 md:px-16 text-gray-800">
      <div class="max-w-6xl mx-auto">
        <h2 class="text-3xl font-bold text-center mb-10">All Events</h2>

        <div class="flex justify-end mb-6">
          <button
            @click="showCreateModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded"
          >
            + Create Event
          </button>
        </div>

        <!-- Loading State -->
        <p v-if="loading" class="text-gray-500">Loading...</p>
        <!-- All events -->
        <div v-else class="grid gap-10 md:grid-cols-2 lg:grid-cols-5">
          <div
            v-for="event in events"
            :key="event.ID"
            class="bg-white rounded-xl shadow-sm hover:shadow-md transition overflow-hidden"
          >
            <div class="grid gap-15 grid-cols-2">
              <div>
                <h3 class="text-xl font-semibold mb-1">{{ event.Name }}</h3>
                <p class="text-sm text-gray-600 mb-2">{{ event.Description }}</p>
                <p class="text-sm"><span class="font-medium">Location:</span> {{ event.Location }}</p>
                <p class="text-sm"><span class="font-medium">Start:</span> {{ formatDate(event.StartDate) }}</p>
                <p class="text-sm"><span class="font-medium">End:</span> {{ formatDate(event.EndDate) }}</p>
                <p class="text-sm"><span class="font-medium">Time Elapsed:</span> {{ formatDuration(event.TimeElapsed) }}</p>
              </div>
              <div>
                <button class="px-4 py-2 bg-blue-600 text-white text-sm rounded hover:bg-blue-700">Edit</button>
                <button class="px-4 py-2 bg-red-500 text-white text-sm rounded hover:bg-red-600">Delete</button>
                <button
                  class="px-4 py-2 text-sm rounded"
                  :class="event.ID === loadedEvent ? 'bg-green-300 text-green-900' : 'bg-green-600 text-white hover:bg-green-700'"
                >
                  {{ event.ID === loadedEvent ? 'Active' : 'Activate' }}
                </button>
              </div>
            </div>

          </div>
        </div>
      </div>

      <CreateEventComponent
        @close="showCreateModal = false"
        @created="eventCreated"
      />
    </div>
  </div>
</template>

<style scoped>

</style>
