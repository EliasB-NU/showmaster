<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'
import PopUp from '@/components/PopUp.vue'

import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const popup = ref<InstanceType<typeof PopUp>| null>(null)

interface Scene {
  ID: Number,
  SceneID: Number,
  SceneName: String,
  SceneDescription: String,

  // Audio
  Audio: String,
  AudioMidiEnabled: Boolean,
  AudioMidiClient: Number,
  AudioMidiChannel: Number,
  AudioMidiNote: String,
  AudioOSCEnable: Boolean,
  AudioOSCClient: Number,
  AudioOSCChannel: Number,
  AudioOSCNote: String,
  AudioGPIOEnable: Boolean,
  AudioGPIOClient: Number,
  AudioGPIOChannel: Number,
  AudioGPIONote: String,
  // Light
  Light: String,
  LightMidiEnabled: Boolean,
  LightMidiClient: Number,
  LightMidiChannel: Number,
  LightMidiNote: String,
  LightOSCEnable: Boolean,
  LightOSCClient: Number,
  LightOSCChannel: Number,
  LightOSCNote: String,
  LightGPIOEnable: Boolean,
  LightGPIOClient: Number,
  LightGPIOChannel: Number,
  LightGPIONote: String,
  // Video
  Video: String,
  VideoMidiEnabled: Boolean,
  VideoMidiClient: Number,
  VideoMidiChannel: Number,
  VideoMidiNote: String,
  VideoOSCEnable: Boolean,
  VideoOSCClient: Number,
  VideoOSCChannel: Number,
  VideoOSCNote: String,
  VideoGPIOEnable: Boolean,
  VideoGPIOClient: Number,
  VideoGPIOChannel: Number,
  VideoGPIONote: String,
  VideoIntraCastEnabled: Boolean,
  VideoIntraCastClient: Number,
  VideoIntraCastChannel: Number,
  VideoIntraCastNote: String,
}

const scenes = ref<Scene[]>([])
const activeScene = ref<Number>(0)

const fetchScenes = async () => {
  try {
    await axios
      .get(`/api/scenes/${route.params.id}`, {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        }
      })
      .then((res) => {
        scenes.value = res.data.scenes
        activeScene.value = res.data.activeScene
      })
  } catch (error: any) {
    if (error.response.status === 404) {
      popup.value?.show("No scenes found")
    }
    popup.value?.show("Error fetching Scenes.")
    console.log(error)
  }
}

onMounted(async () => {
  await fetchScenes()
})

</script>

<template>
  <div class="flex flex-col min-h-screen">
    <HeaderComponent :eventID="Number(route.params.id)" />
    <div class="bg-white py-12 px-4 md:px-16 text-gray-800">

    </div>
    <PopUp ref="popup" />
  </div>
</template>

<style scoped>

</style>
