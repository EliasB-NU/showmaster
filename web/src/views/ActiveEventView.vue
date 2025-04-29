<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'

let connection = new WebSocket('/api/ws', )

connection.onopen = () => {
  console.log('connection opened')
  connection.send("Test")
}
</script>

<template>
  <div>
    <HeaderComponent />
    <h1>!!!Active Event!!!</h1>
  </div>
</template>

<style scoped>

</style>
