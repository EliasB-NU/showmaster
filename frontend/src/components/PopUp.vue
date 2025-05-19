<script setup lang="ts">
import { ref, defineExpose } from "vue";

const visible = ref(false);
const message = ref("");

const show = (msg: string, duration = 3000) => {
  message.value = msg;
  visible.value = true;
  setTimeout(() => {
    visible.value = false;
  }, duration);
};

const close = () => {
  visible.value = false;
};

defineExpose({ show });
</script>

<template>
  <transition
    enter-active-class="transform transition duration-300 ease-out"
    enter-from-class="translate-y-5 opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transform transition duration-300 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-5 opacity-0"
  >
    <div
      v-if="visible"
      class="fixed bottom-5 right-5 bg-gray-800 text-white py-3 px-5 rounded-lg shadow-lg flex items-center space-x-3"
    >
      <span>{{ message }}</span>
      <button @click="close" class="text-gray-400 hover:text-white">&times;</button>
    </div>
  </transition>
</template>

<style scoped>

</style>
