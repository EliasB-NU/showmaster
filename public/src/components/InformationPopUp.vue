<script setup>
const props = defineProps([
  'shown',
  'message',
])

const closePopUp = async () => {
  props.shown = false;
}
</script>

<template>
  <div v-if="shown" class="modal-overlay">
    <div class="modal-content">
      <p>{{ message }}</p>
      <button @click="closePopUp" class="btn-primary">Confirm</button>
    </div>
  </div>
</template>