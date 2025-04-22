<script setup lang="ts">
const props = defineProps({
  token: {
    type: String,
    required: true,
  },
  visible: {
    type: Boolean,
    default: false,
    required: true,
  }
})

const emit = defineEmits(['close'])

const copyToken = async () => {
  try {
    await navigator.clipboard.writeText(props.token);
    alert("Token copied to clipboard!");
  } catch (err) {
    alert("Failed to copy token.");
  }
};

const onClose = () => {
  emit('close');
}
</script>

<template>
  <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-50">
    <div class="bg-white rounded-xl shadow-lg w-full max-w-md p-6 space-y-4 relative">
      <!-- Close Button -->
      <button @click="onClose" class="absolute top-3 right-3 text-gray-500 hover:text-gray-700">
        ✕
      </button>

      <!-- Title -->
      <h2 class="text-xl font-semibold">Your One-Time Token</h2>

      <!-- Token Display -->
      <div class="bg-gray-100 text-gray-800 font-mono p-4 rounded-lg relative overflow-x-auto">
        <span class="whitespace-nowrap break-all">{{ token }}</span>
        <button
          @click="copyToken"
          class="absolute top-2 right-2 px-2 py-1 text-sm bg-gray-300 hover:bg-gray-400 rounded"
        >
          Copy
        </button>
      </div>

      <!-- Note -->
      <p class="text-sm text-red-500">This token will only be shown once. Make sure to save it!</p>
    </div>
  </div>
</template>

<style scoped>

</style>
