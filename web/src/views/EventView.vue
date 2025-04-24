<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'
import PopUp from '@/components/PopUp.vue'

import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import Cookies from 'js-cookie'
import SceneCreateComponent from '@/components/SceneCreateComponent.vue'
import SceneEditComponent from '@/components/SceneEditComponent.vue'

const route = useRoute()
const popup = ref<InstanceType<typeof PopUp>| null>(null)

interface Scene {
  ID: number,
  SceneID: number,
  SceneName: string,
  SceneDescription: string,

  // Audio
  Audio: string,
  AudioMidiEnabled: boolean,
  AudioMidiClient: number,
  AudioMidiChannel: number,
  AudioMidiNote: string,
  AudioOSCEnable: boolean,
  AudioOSCClient: number,
  AudioOSCChannel: number,
  AudioOSCNote: string,
  AudioGPIOEnable: boolean,
  AudioGPIOClient: number,
  AudioGPIOChannel: number,
  AudioGPIONote: string,
  // Light
  Light: string,
  LightMidiEnabled: boolean,
  LightMidiClient: number,
  LightMidiChannel: number,
  LightMidiNote: string,
  LightOSCEnable: boolean,
  LightOSCClient: number,
  LightOSCChannel: number,
  LightOSCNote: string,
  LightGPIOEnable: boolean,
  LightGPIOClient: number,
  LightGPIOChannel: number,
  LightGPIONote: string,
  // Video
  Video: string,
  VideoMidiEnabled: boolean,
  VideoMidiClient: number,
  VideoMidiChannel: number,
  VideoMidiNote: string,
  VideoOSCEnable: boolean,
  VideoOSCClient: number,
  VideoOSCChannel: number,
  VideoOSCNote: string,
  VideoGPIOEnable: boolean,
  VideoGPIOClient: number,
  VideoGPIOChannel: number,
  VideoGPIONote: string,
  VideoIntraCastEnabled: boolean,
  VideoIntraCastClient: number,
  VideoIntraCastChannel: number,
  VideoIntraCastNote: string,
}

const scenes = ref<Scene[]>([])
const activeScene = ref<number>(0)
const loading = ref<boolean>(true)
const sceneToEdit = ref<Scene>({} as Scene)

const fetchScenes = async () => {
  try {
    await axios
      .get(`/api/scenes/${route.params.id}`, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then((res) => {
        scenes.value = res.data.scenes
        scenes.value.sort((a, b) => a.SceneID - b.SceneID)
        activeScene.value = res.data.activeScene
        loading.value = false
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      popup.value?.show("No scenes found")
    }
    loading.value = false
    popup.value?.show("Error fetching Scenes.")
    console.log(error)
  }
}

const deleteScene = async (sceneId: number) => {
  try {
    await axios
      .delete(`/api/scenes/delete/${sceneId}`, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
    .then(_ => {
      popup.value?.show("Successfully deleted scene.")
      fetchScenes()
    })
  } catch (error: any) {
    if (error.response.status === 404) {
      popup.value?.show("Scene not found")
    }
    popup.value?.show("Error deleting scene.")
    console.log(error)
  }
}

const selectScene = async (sceneId: number) => {
  try {
    await axios
      .post(`/api/scenes/updateActive/${route.params.id}/${sceneId}`, {}, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(_ => {
        popup.value?.show("Successfully selected scene.")
        fetchScenes()
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      popup.value?.show("Scene not found")
    }
    popup.value?.show("Error selecting scene.")
    console.log(error.response)
  }
}

const testScene = async (sceneId: number) => {
  try {
    await axios
      .post(`/api/scenes/test/${sceneId}`, {}, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${Cookies.get('token')}`,
        }
      })
      .then(_ => {
        popup.value?.show("Successfully tested scene.")
        fetchScenes()
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      popup.value?.show("Scene not found")
    }
    popup.value?.show("Error testing scene.")
    console.log(error.response)
  }
}

const showCreateModal = ref<boolean>(false)
const showEditModal = ref<boolean>(false)

onMounted(async () => {
  await fetchScenes()
})

</script>

<template>
  <div class="flex flex-col min-h-screen">
    <HeaderComponent :eventID="Number(route.params.id)" />
    <div class="bg-white py-12 px-4 md:px-16 text-gray-800">
      <div class="flex justify-end mb-6">
        <button
          @click="showCreateModal = true"
          class="bg-gray-800 hover:bg-gray-700 text-white font-medium px-4 py-2 rounded"
        >
          Create Scene
        </button>
      </div>

      <!-- Loading State -->
      <p v-if="loading" class="text-gray-500">Loading...</p>

      <!-- Scenes Table -->
      <div v-else class="bg-white shadow rounded-lg">
        <table class="w-full border-collapse">
          <thead>
          <tr class="bg-gray-200">
            <th class="p-3 text-left">ID</th>
            <th class="p-3 text-left">Name</th>
            <th class="p-3 text-left">Audio</th>
            <th class="p-3 text-left">Light</th>
            <th class="p-3 text-left">Video</th>
            <th class="p-3">Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="scene in scenes" :key="scene.SceneID" class="border-t">
            <td class="p-3">{{ scene.SceneID }}</td>
            <td class="p-3">{{ scene.SceneName }}</td>
            <td class="p-3">{{ scene.Audio }}</td>
            <td class="p-3">{{ scene.Light }}</td>
            <td class="p-3">{{ scene.Video }}</td>

            <td class="p-3 flex space-x-2">
              <button
                v-if="Cookies.get('admin') === 'true' || Number(Cookies.get('events')) >= 2"
                @click="sceneToEdit = scene; showEditModal = true"
                class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700"
              >
                Edit
              </button>
              <button
                v-if="Cookies.get('admin') === 'true' || Number(Cookies.get('events')) >= 2"
                @click="selectScene(scene.ID)"
                class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700"
              >
                Select
              </button>
              <button
                v-if="Cookies.get('admin') === 'true' || Number(Cookies.get('events')) >= 2"
                @click="testScene(scene.ID)"
                class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700"
              >
                Test
              </button>
              <button
                v-if="Cookies.get('admin') === 'true' || Number(Cookies.get('events')) >= 3"
                @click="deleteScene(scene.ID)"
                class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
              >
                Delete
              </button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>

    <SceneCreateComponent
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="fetchScenes(); showCreateModal = false"
    />

    <SceneEditComponent
      v-if="showEditModal"
      :scene="sceneToEdit"
      @close="showEditModal = false"
      @updated="fetchScenes(); showEditModal = false"
    />

    <PopUp ref="popup" />
  </div>
</template>

<style scoped>

</style>
