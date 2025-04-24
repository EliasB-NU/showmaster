<script setup lang="ts">

import { onMounted, ref } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'
import Cookies from 'js-cookie'

const route = useRoute()

const emit = defineEmits(['close', 'created'])

interface Scene {
  scene_id: number,
  scene_name: string,
  scene_description: string,

  executes_in: number,

  audio: string,
  audio_midi_enabled: boolean,
  audio_midi_client: number,
  audio_midi_channel: number,
  audio_midi_note: string,
  audio_osc_enabled: boolean,
  audio_osc_client: number,
  audio_osc_channel: number,
  audio_osc_note: string,
  audio_gpio_enabled: boolean,
  audio_gpio_client: number,
  audio_gpio_channel: number,
  audio_gpio_note: string,

  light: string,
  light_midi_enabled: boolean,
  light_midi_client: number,
  light_midi_channel: number,
  light_midi_note: string,
  light_osc_enabled: boolean,
  light_osc_client: number,
  light_osc_channel: number,
  light_osc_note: string,
  light_gpio_enabled: boolean,
  light_gpio_client: number,
  light_gpio_channel: number,
  light_gpio_note: string,

  video: string,
  video_midi_enabled: boolean,
  video_midi_client: number,
  video_midi_channel: number,
  video_midi_note: string,
  video_osc_enabled: boolean,
  video_osc_client: number,
  video_osc_channel: number,
  video_osc_note: string,
  video_gpio_enabled: boolean,
  video_gpio_client: number,
  video_gpio_channel: number,
  video_gpio_note: string,
  video_intra_cast_enabled: boolean,
  video_intra_cast_client: number,
  video_intra_cast_channel: number,
  video_intra_cast_note: string,
}

const scene = ref<Scene>({
  scene_id: 0,
  scene_name: '',
  scene_description: '',

  executes_in: -1,

  audio: '',
  audio_midi_enabled: false,
  audio_midi_client: 0,
  audio_midi_channel: 0,
  audio_midi_note: '',
  audio_osc_enabled: false,
  audio_osc_client: 0,
  audio_osc_channel: 0,
  audio_osc_note: '',
  audio_gpio_enabled: false,
  audio_gpio_client: 0,
  audio_gpio_channel: 0,
  audio_gpio_note: '',

  light: '',
  light_midi_enabled: false,
  light_midi_client: 0,
  light_midi_channel: 0,
  light_midi_note: '',
  light_osc_enabled: false,
  light_osc_client: 0,
  light_osc_channel: 0,
  light_osc_note: '',
  light_gpio_enabled: false,
  light_gpio_client: 0,
  light_gpio_channel: 0,
  light_gpio_note: '',

  video: '',
  video_midi_enabled: false,
  video_midi_client: 0,
  video_midi_channel: 0,
  video_midi_note: '',
  video_osc_enabled: false,
  video_osc_client: 0,
  video_osc_channel: 0,
  video_osc_note: '',
  video_gpio_enabled: false,
  video_gpio_client: 0,
  video_gpio_channel: 0,
  video_gpio_note: '',
  video_intra_cast_enabled: false,
  video_intra_cast_client: 0,
  video_intra_cast_channel: 0,
  video_intra_cast_note: '',
})

const createScene = async () => {
  try {
    await axios
      .post(`/api/scenes/create/${route.params.id}`, {
        ...scene.value,
      }, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
    .then(_ => {
      emit('created')
    })
  } catch (error) {
    console.error(error)
    emit('close')
  }
}

// Get all the client types
interface MidiClient {
  ID: number,
  Name: string,
}

interface OSCClient {
  ID: number,
  Name: string,
}

interface GPIOClient {
  ID: number,
  Name: string,
}

interface IntraCastClient {
  ID: number,
  Name: string,
}

const midiClients = ref<MidiClient[]>([]);
const oscClients = ref<OSCClient[]>([]);
const gpioClients = ref<GPIOClient[]>([]);
const intracastClients = ref<IntraCastClient[]>([]);

async function fetchMidiClients() {
  try {
    await axios
      .get('/api/clients/byType/midi', {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
    .then(res => {
      midiClients.value = res.data;
    })
  } catch (error: any) {
    if (error.response.status === 404) {
      return;
    }
    console.error(error)
    emit('close')
  }
}

async function fetchOSCClients() {
  try {
    await axios
      .get('/api/clients/byType/osc', {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(res => {
        oscClients.value = res.data;
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      return;
    }
    console.error(error)
    emit('close')
  }
}

async function fetchGPIOClients() {
  try {
    await axios
      .get('/api/clients/byType/gpio', {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(res => {
        gpioClients.value = res.data;
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      return;
    }
    console.error(error)
    emit('close')
  }
}

async function fetchIntracastClients() {
  try {
    await axios
      .get('/api/clients/byType/intracast', {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(res => {
        intracastClients.value = res.data;
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      return;
    }
    console.error(error)
    emit('close')
  }
}


const closePopup = () => {
  emit('close')
}

onMounted(async () => {
  await fetchMidiClients()
  await fetchOSCClients()
  await fetchGPIOClients()
  await fetchIntracastClients()
})
</script>

<template>
  <div class="fixed inset-0 bg-opacity-50 flex items-center justify-center">
    <div class="bg-white p-6 rounded-2xl shadow-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
      <h3 class="text-lg font-semibold mb-4">Create Scene</h3>

      <!-- Form to edit member -->
      <div class="mb-4">
        <form @submit.prevent="createScene">
          <div class="grid grid-cols-2 gap-4">
            <!-- Left Column: Basic Information -->
            <div class="w-1/2 pr-4">
              <div class="mb-4">
                <label for="scene_name" class="block">Name:</label>
                <input v-model="scene.scene_name" id="name" type="text" class="w-full border px-3 py-2 rounded" required/>
              </div>
              <div class="mb-4">
                <label for="scene_id" class="block">ID:</label>
                <input v-model="scene.scene_id" id="scene_id" type="number" class="w-full border px-3 py-2 rounded" step="0.01" min="1" max="5000" required/>
              </div>
              <div class="mb-4">
                <label for="scene_description" class="block">Description:</label>
                <input v-model="scene.scene_description" id="scene_description" type="text" class="w-full border px-3 py-2 rounded"  />
              </div>
              <div class="mb-4">
                <label for="executes_in" class="block">Executes In:</label>
                <input v-model="scene.executes_in" id="executes_in" type="number" class="w-full border px-3 py-2 rounded" min="-1" max="3600" />
              </div>
              <!-- Audio -->
              <div>
                <h2 class="font-bold">Audio</h2>
                <div class="mb-4">
                  <input v-model="scene.audio" id="audio" type="text" class="w-full border px-3 py-2 rounded" min="-1" max="3600" />
                </div>
                <div class="mb-4">
                  <label for="audio_midi_enabled" class="block font-bold">Enable Midi:</label>
                  <select v-model="scene.audio_midi_enabled" id="audio_midi_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.audio_midi_enabled" class="mb-4">
                  <label for="audio_midi_client" class="block">Midi Client:</label>
                  <select v-model="scene.audio_midi_client" id="audio_midi_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in midiClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.audio_midi_enabled" class="mb-4">
                  <label for="audio_midi_channel" class="block">Midi Channel:</label>
                  <input v-model="scene.audio_midi_channel" id="audio_midi_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.audio_midi_enabled" class="mb-4">
                  <label for="audio_midi_note" class="block">Midi Note:</label>
                  <input v-model="scene.audio_midi_note" id="audio_midi_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div class="mb-4">
                  <label for="audio_midi_enabled" class="block font-bold">Enable OSC:</label>
                  <select v-model="scene.audio_osc_enabled" id="audio_osc_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.audio_osc_enabled" class="mb-4">
                  <label for="audio_osc_client" class="block">OSC Client:</label>
                  <select v-model="scene.audio_osc_client" id="audio_osc_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in oscClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.audio_osc_enabled" class="mb-4">
                  <label for="audio_osc_channel" class="block">OSC Channel:</label>
                  <input v-model="scene.audio_osc_channel" id="audio_osc_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.audio_osc_enabled" class="mb-4">
                  <label for="audio_osc_note" class="block">OSC Note:</label>
                  <input v-model="scene.audio_osc_note" id="audio_osc_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div class="mb-4">
                  <label for="audio_midi_enabled" class="block font-bold">Enable GPIO:</label>
                  <select v-model="scene.audio_gpio_enabled" id="audio_gpio_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.audio_gpio_enabled" class="mb-4">
                  <label for="audio_gpio_client" class="block">GPIO Client:</label>
                  <select v-model="scene.audio_gpio_client" id="audio_gpio_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in gpioClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.audio_gpio_enabled" class="mb-4">
                  <label for="audio_gpio_channel" class="block">GPIO Channel:</label>
                  <input v-model="scene.audio_gpio_channel" id="audio_gpio_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.audio_gpio_enabled" class="mb-4">
                  <label for="audio_gpio_note" class="block">GPIO Note:</label>
                  <input v-model="scene.audio_gpio_note" id="audio_gpio_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
              </div>
            </div>

            <!-- Right Column: Permissions -->
            <div class="w-1/2 pl-4">
              <!-- Video -->
              <div>
                <h2 class="font-bold">Light</h2>
                <div class="mb-4">
                  <input v-model="scene.light" id="birthday" type="text" class="w-full border px-3 py-2 rounded" min="-1" max="3600" />
                </div>
                <div class="mb-4">
                  <label for="light_midi_enabled" class="block font-bold">Enable Midi:</label>
                  <select v-model="scene.light_midi_enabled" id="light_midi_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.light_midi_enabled" class="mb-4">
                  <label for="light_midi_client" class="block">Midi Client:</label>
                  <select v-model="scene.light_midi_client" id="light_midi_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in midiClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.light_midi_enabled" class="mb-4">
                  <label for="light_midi_channel" class="block">Midi Channel:</label>
                  <input v-model="scene.light_midi_channel" id="light_midi_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.light_midi_enabled" class="mb-4">
                  <label for="light_midi_note" class="block">Midi Note:</label>
                  <input v-model="scene.light_midi_note" id="light_midi_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div class="mb-4">
                  <label for="light_midi_enabled" class="block font-bold">Enable OSC:</label>
                  <select v-model="scene.light_osc_enabled" id="light_osc_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.light_osc_enabled" class="mb-4">
                  <label for="light_osc_client" class="block">OSC Client:</label>
                  <select v-model="scene.light_osc_client" id="light_osc_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in oscClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.light_osc_enabled" class="mb-4">
                  <label for="light_osc_channel" class="block">OSC Channel:</label>
                  <input v-model="scene.light_osc_channel" id="light_osc_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.light_osc_enabled" class="mb-4">
                  <label for="light_osc_note" class="block">OSC Note:</label>
                  <input v-model="scene.light_osc_note" id="light_osc_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div v-if="scene.light_osc_enabled" class="mb-4 font-bold">
                  <label for="light_midi_enabled" class="block">Enable GPIO:</label>
                  <select v-model="scene.light_gpio_enabled" id="light_gpio_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.light_gpio_enabled" class="mb-4">
                  <label for="light_gpio_client" class="block">GPIO Client:</label>
                  <select v-model="scene.light_gpio_client" id="light_gpio_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in gpioClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.light_gpio_enabled" class="mb-4">
                  <label for="light_gpio_channel" class="block">GPIO Channel:</label>
                  <input v-model="scene.light_gpio_channel" id="light_gpio_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.light_gpio_enabled" class="mb-4">
                  <label for="light_gpio_note" class="block">GPIO Note:</label>
                  <input v-model="scene.light_gpio_note" id="light_gpio_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
              </div>
              <!-- Video -->
              <div>
                <div class="mb-4">
                  <label for="video" class="block">video:</label>
                  <input v-model="scene.video" id="birthday" type="text" class="w-full border px-3 py-2 rounded" min="-1" max="3600" />
                </div>
                <div class="mb-4">
                  <label for="video_midi_enabled" class="block font-bold">Enable Midi:</label>
                  <select v-model="scene.video_midi_enabled" id="video_midi_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.video_midi_enabled" class="mb-4">
                  <label for="video_midi_client" class="block">Midi Client:</label>
                  <select v-model="scene.video_midi_client" id="video_midi_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in midiClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.video_midi_enabled" class="mb-4">
                  <label for="video_midi_channel" class="block">Midi Channel:</label>
                  <input v-model="scene.video_midi_channel" id="video_midi_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div class="mb-4">
                  <label for="video_midi_note" class="block">Midi Note:</label>
                  <input v-model="scene.video_midi_note" id="video_midi_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div v-if="scene.video_midi_enabled" class="mb-4">
                  <label for="video_midi_enabled" class="block font-bold">Enable OSC:</label>
                  <select v-model="scene.video_osc_enabled" id="video_osc_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.video_osc_enabled" class="mb-4">
                  <label for="video_osc_client" class="block">OSC Client:</label>
                  <select v-model="scene.video_osc_client" id="video_osc_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in oscClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.video_osc_enabled" class="mb-4">
                  <label for="video_osc_channel" class="block">OSC Channel:</label>
                  <input v-model="scene.video_osc_channel" id="video_osc_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.video_osc_enabled" class="mb-4">
                  <label for="video_osc_note" class="block">OSC Note:</label>
                  <input v-model="scene.video_osc_note" id="video_osc_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div v-if="scene.video_osc_enabled" class="mb-4">
                  <label for="video_midi_enabled" class="block font-bold">Enable GPIO:</label>
                  <select v-model="scene.video_gpio_enabled" id="video_gpio_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.video_gpio_enabled" class="mb-4">
                  <label for="video_gpio_client" class="block">GPIO Client:</label>
                  <select v-model="scene.video_gpio_client" id="video_gpio_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in gpioClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.video_gpio_enabled" class="mb-4">
                  <label for="video_gpio_channel" class="block">GPIO Channel:</label>
                  <input v-model="scene.video_gpio_channel" id="video_gpio_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"  />
                </div>
                <div v-if="scene.video_gpio_enabled" class="mb-4">
                  <label for="video_gpio_note" class="block">GPIO Note:</label>
                  <input v-model="scene.video_gpio_note" id="video_gpio_note" type="text" class="w-full border px-3 py-2 rounded"  />
                </div>
                <div class="mb-4">
                  <label for="video_intra_cast_enabled" class="block font-bold">Enable Intracast:</label>
                  <select v-model="scene.video_intra_cast_enabled" id="video_intra_cast_enabled" class="w-full border px-3 py-2 rounded">
                    <option value=true>Enable</option>
                    <option value=false>Disable</option>
                  </select>
                </div>
                <div v-if="scene.video_intra_cast_enabled" class="mb-4">
                  <label for="video_intra_cast_client" class="block">Intracast Client:</label>
                  <select v-model="scene.video_intra_cast_client" id="video_intra_cast_client" class="w-full border px-3 py-2 rounded">
                    <option v-for="client in intracastClients" :key="client.ID" :value="client.ID">{{ client.Name }}</option>
                  </select>
                </div>
                <div v-if="scene.video_intra_cast_enabled" class="mb-4">
                  <label for="video_intra_cast_channel" class="block">Intracast Channel:</label>
                  <input v-model="scene.video_intra_cast_channel" id="video_intra_cast_channel" type="number" class="w-full border px-3 py-2 rounded" min="1" max="5000"/>
                </div>
                <div v-if="scene.video_intra_cast_enabled" class="mb-4">
                  <label for="video_intra_cast_note" class="block">Intracast Note:</label>
                  <input v-model="scene.video_intra_cast_note" id="video_intra_cast_note" type="text" class="w-full border px-3 py-2 rounded"  />
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
