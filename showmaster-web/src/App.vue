<template>
  <div v-if="showPOPUP">
    <POPUP title="Test" text="helloWorld" @close="togglePopUp"/>
  </div>
  <button @click="togglePopUp">HUHN</button>
</template>

<script>
import POPUP from './components/popup.vue'

export default {
  name: "showmaster",
  components: { POPUP },
  data() {
    return {
      showPOPUP: false,
    }
  },
  methods: {
    togglePopUp() {
      this.showPOPUP = !this.showPOPUP
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

h1 {
  border-bottom: 1px solid #ddd;
  display: inline-block;
  padding-bottom: 10px;
}

.active {
  color: #2c3e50;
}
</style>
