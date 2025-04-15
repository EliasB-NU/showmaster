<script setup lang="ts">
import FooterComponent from '@/components/FooterComponent.vue'
import { RouterView } from 'vue-router'
</script>

<template>

</template>
  <div class="flex flex-col min-h-screen">
    <main ref="mainContent" tabindex="-1" class="flex-grow outline-none focus:outline-none">
      <RouterView />
    </main>
    <FooterComponent />
  </div>
<style scoped>
</style>
