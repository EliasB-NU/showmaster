<script setup>

import { watch } from 'vue'
import { useRoute } from 'vue-router'
import Footer from '@/components/Footer.vue'

const route = useRoute();

watch (
  () => route.params.projectId,
  (newId) => {

  }
)

</script>

<template>
  <Header project=" |  {{ route.params.projctId }}" render-logout=true />

  <Footer />
</template>

<style scoped>

</style>