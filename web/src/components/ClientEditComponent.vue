<script setup lang="ts">

defineProps({
  visible: {
    type: Boolean,
    default: false,
    required: true
  },
  user: {
    type: Object,
    default: null,
    required: true
  }
})

const emit = defineEmits(['close', 'clientEdited'])
</script>

<template>

</template>

<style scoped>

</style>
