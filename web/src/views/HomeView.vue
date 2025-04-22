<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'
import { onMounted, ref } from 'vue'

interface Event {

}

const events = ref<Event[]>([])

const fetchEvents = async () => {

}

onMounted(async () => {
  await fetchEvents()
})

</script>

<template>
  <div>
    <HeaderComponent />
    <h1>Hello World</h1>
  </div>
</template>

<style scoped>

</style>
