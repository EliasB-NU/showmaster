<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import Cookies from 'js-cookie'
import ClientTokenComponent from '@/components/ClientTokenComponent.vue'


defineProps({
  visible: {
    type: Boolean,
    default: false,
    required: true
  }
})

interface Client {
  name: string
  type: string
  ip: string
  port: number
}

const client = ref<Client>({} as Client)
const token = ref<string>('')
const tokenVisible = ref<boolean>(false)

const emit = defineEmits(['close', 'clientCreated'])

const createClient = async () => {
  try {
    await axios
      .post('/api/createClient', {
        ...client,
      }, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookies.get('token'),
        }
      })
      .then(response => {
        tokenVisible.value = true
        token.value = response.data.token
        emit('clientCreated')
      })
  } catch (error) {
    console.error('Error creating client:', error)
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
      <h3 class="text-lg font-semibold mb-4">Edit Member</h3>

      <!-- Form to edit member -->
      <div class="mb-4">
        <form @submit.prevent="createClient">
          <div class="grid grid-cols-2 gap-4">
            <!-- Left Column: Basic Information -->
            <div class="w-1/2 pr-4">
              <div class="mb-4">
                <label for="name" class="block">Name:</label>
                <input v-model="client.name" id="name" type="text" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="type" class="block">Type:</label>
                <select v-model="client.type" id="type" class="input-field border-2 border-b-black rounded-sm">
                  <option value="osc">OSC</option>
                  <option value="midi">MIDI</option>
                  <option value="gpio">GPIO</option>
                  <option value="llls">LLLS</option>
                </select>
              </div>
              <div class="mb-4">
                <label for="ip" class="block">IP:</label>
                <input v-model="client.ip" id="ip" type="text" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
              <div class="mb-4">
                <label for="port" class="block">Port:</label>
                <input v-model="client.port" id="port" type="number" class="input-field border-2 border-b-black rounded-sm" required />
              </div>
            </div>
          </div>

          <div class="p-3 flex space-x-2 mt-4 justify-end">
            <button type="submit" class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700">Create Client</button>
            <button type="button" @click="closePopup" class="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700" >Cancel</button>
          </div>
        </form>
      </div>
    </div>

    <ClientTokenComponent
      :content="token"
      :visible="tokenVisible"
      @close="tokenVisible = false"
    />
  </div>
</template>

<style scoped>

</style>
