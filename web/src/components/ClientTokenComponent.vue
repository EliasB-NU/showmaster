<script setup lang="ts">
const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  content: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close'])

const closePopUp = () => {
  emit('close')
}
</script>

<template>
  <div v-if="visible" class="fixed inset-0 bg-opacity-50 flex items-center justify-center">
    <h2 class="text-lg font-semibold mb-4">Client Token</h2>
    <p>Please copy the client token, you cannot view it again</p>
    <p>{{ content }}</p>

    <div class="p-3 flex space-x-2 mt-4 justify-end">
      <button type="button" @click="closePopUp" class="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700" >Close</button>
    </div>
  </div>
</template>

<style scoped>

</style>
