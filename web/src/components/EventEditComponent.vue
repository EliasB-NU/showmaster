<script setup lang="ts">
import axios from 'axios'
import Cookies from 'js-cookie'

const props = defineProps({
  event: {
    type: Object,
    required: true
  },
})

const emit = defineEmits(["close"]);

const updateEvent = async () => {
  try {
    await axios
      .post(`/api/event/update`, {
        id: props.event.ID,
        name: props.event.Name,
        description: props.event.Description,
        location: props.event.Location,
        startDate: new Date(props.event.StartDate),
        endDate: new Date(props.event.EndDate),
      }, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(_ => {
        emit("close")
      })
  } catch (error) {
    console.log(error)
    emit('close')
  }
}
</script>

<template>
  <div class="fixed inset-0 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-2xl shadow-lg w-full max-w-xl p-6 relative">
      <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 text-2xl">&times;</button>

      <h2 class="text-xl font-bold mb-4">Update Event</h2>

      <form @submit.prevent="updateEvent" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Name</label>
          <input v-model="event.Name" type="text" class="w-full border px-3 py-2 rounded" required />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <textarea v-model="event.Description" class="w-full border px-3 py-2 rounded" required></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Location</label>
          <input v-model="event.Location" type="text" class="w-full border px-3 py-2 rounded" required />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium mb-1">Start Date</label>
            <input v-model="event.StartDate" type="text" class="w-full border px-3 py-2 rounded" required />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">End Date</label>
            <input v-model="event.EndDate" type="text" class="w-full border px-3 py-2 rounded" required />
          </div>
        </div>

        <button
          type="submit"
          class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          Update Event
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>
