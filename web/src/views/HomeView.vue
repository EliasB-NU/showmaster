<script setup lang="ts">
import PopUp from '@/components/PopUp.vue'
import EventCreateComponent from '@/components/EventCreateComponent.vue'
import HeaderComponent from '@/components/HeaderComponent.vue'

import { onMounted, ref } from 'vue'
import axios from 'axios'
import Cookies from 'js-cookie'
import router from '@/router'
import EventEditComponent from '@/components/EventEditComponent.vue'

const popup = ref<InstanceType<typeof PopUp>| null>(null)

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
    popup.value?.show('Error fetching events.')
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

const deleteEvent = async (id: number) => {
  try {
    await axios
      .delete(`/api/event/delete/${id}`, {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
    .then(_ => {
      popup.value?.show('Successfully deleted event.')
      fetchEvents()
    })
  } catch (error) {
    console.log(error)
    popup.value?.show('Error deleting event.')
  }
}

const activateEvent = async (id: number) => {
  if (loadedEvent.value === id) {
    popup.value?.show("Already activated")
    return
  }
  try {
    await axios
      .post(`/api/event/activate/${id}`, {}, {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(_ => {
        popup.value?.show("Event activated.")
        loadedEvent.value = id
    })
  } catch (error) {
    console.log(error)
    popup.value?.show("Error deleting event.")
  }
}

const showCreateModal = ref<boolean>(false)

const eventCreated = async () => {
  showCreateModal.value = false
  await fetchEvents()
}

const showEditModal = ref<boolean>(false)
const toEditEvent = ref<Event>({} as Event)

const editEvent = (event: Event) => {
  toEditEvent.value = event
  showEditModal.value = true
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
            class="bg-gray-800 hover:bg-gray-700 text-white font-medium px-4 py-2 rounded"
          >
            Create Event
          </button>
        </div>

        <!-- Loading State -->
        <p v-if="loading" class="text-gray-500">Loading...</p>
        <!-- All events -->
        <div v-else class="grid gap-10 md:grid-cols-2">
          <div
            v-for="event in events"
            :key="event.ID"
            class="bg-white rounded-xl shadow-sm hover:shadow-md transition overflow-hidden"
          >
            <div class="grid gap-15 grid-cols-2 p-2">
              <div>
                <h3 class="text-xl font-semibold mb-1 p-1">{{ event.Name }}</h3>
                <p class="text-sm text-gray-600 mb-2 p-1">{{ event.Description }}</p>
                <p class="text-sm p-1"><span class="font-bold">Location:</span> {{ event.Location }}</p>
                <p class="text-sm p-1"><span class="font-bold">Start:</span> {{ formatDate(event.StartDate) }}</p>
                <p class="text-sm p-1"><span class="font-bold">End:</span> {{ formatDate(event.EndDate) }}</p>
                <p class="text-sm p-1"><span class="font-bold">Time Elapsed:</span> {{ formatDuration(event.TimeElapsed) }}</p>
              </div>
              <div class="grid gap-5 grid-cols-1">
                <button @click="router.push(`/event/${event.ID}`)" class="px-4 py-2 bg-gray-800 text-white text-sm rounded hover:bg-gray-700">Open</button>
                <button @click="editEvent(event)" v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 2" class="px-4 py-2 bg-gray-800 text-white text-sm rounded hover:bg-gray-700">Edit</button>
                <button @click="deleteEvent(event.ID)" v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 3" class="px-4 py-2 bg-red-500 text-white text-sm rounded hover:bg-red-600">Delete</button>
                <button
                  v-if="Cookies.get('admin') || Number(Cookies.get('events')) >= 2"
                  class="px-4 py-2 text-sm rounded"
                  @click="activateEvent(event.ID)"
                  :class="event.ID === loadedEvent ? 'bg-green-300 text-green-900' : 'bg-green-600 text-white hover:bg-green-700'"
                >
                  {{ event.ID === loadedEvent ? 'Active' : 'Activate' }}
                </button>
              </div>
            </div>

          </div>
        </div>
        <PopUp ref="popup" />
      </div>

      <EventCreateComponent
        v-if="showCreateModal"
        @close="showCreateModal = false"
        @created="eventCreated"
      />

      <EventEditComponent
        v-if="showEditModal"
        :event="toEditEvent"
        @close="showEditModal = false"
      />
    </div>
  </div>
</template>

<style scoped>

</style>
