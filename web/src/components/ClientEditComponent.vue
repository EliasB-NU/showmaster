<script setup lang="ts">
import axios from 'axios'
import Cookies from 'js-cookie'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
    required: true
  },
  client: {
    type: Object,
    default: null,
    required: true
  }
})

const emit = defineEmits(['close', 'clientEdited'])

const updateClient = async () => {
  try {
    await axios
      .post('/api/updateClient', {
        ...props.client,
      }, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`
        }
      })
      .then(() => {
        emit('close')
      })
  } catch (error) {
    console.error(error)
    emit('close')
  }
}
const closePopup = () => {
  emit('close')
}
</script>

<template>
  <div v-if="visible" class="fixed inset-0 bg-opacity-50 flex items-center justify-center">
    <div class="bg-white p-6 rounded-2xl shadow-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
      <h3 class="text-lg font-semibold mb-4">Edit client</h3>

      <!-- Form to edit client -->
      <div class="mb-4">
        <form @submit.prevent="updateClient">
          <div class="grid grid-cols-2 gap-4">
            <!-- Left Column: Basic Information -->
            <div class="w-1/2 pr-4">
              <div class="mb-4">
                <label for="name" class="block">Name:</label>
                <input v-model="client.Name" id="name" type="text" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="type" class="block">Type:</label>
                <select v-model="client.Type" id="type" class="input-field border-2 border-b-black rounded-sm">
                  <option value="osc">OSC</option>
                  <option value="midi">MIDI</option>
                  <option value="gpio">GPIO</option>
                  <option value="llls">LLLS</option>
                </select>
              </div>
              <div class="mb-4">
                <label for="ip" class="block">IP:</label>
                <input v-model="client.IP" id="ip" type="text" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="port" class="block">Port:</label>
                <input v-model="client.Port" id="port" type="number" min="1" max="65535" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
            </div>
          </div>

          <div class="p-3 flex space-x-2 mt-4 justify-end">
            <button type="submit" class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700">Update Client</button>
            <button type="button" @click="closePopup" class="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700" >Cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
