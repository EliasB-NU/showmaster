<script>
import router from '@/router/index.js'

let project = "";
let renderLogout = false;
export default {
  name: 'Header',
  props: {
    project,
    renderLogout,
  },
  setup() {
    const logout = () => {
      localStorage.clear();
      router.push("/")
    }

    return {
      logout,
    }
  },
}

</script>

<template>
  <!-- Header Section -->
  <header class="header">
    <div class="header-title">ShowMaster V3 {{ project }}</div>
  </header>
  <div class="logout-container" v-if="renderLogout">
    <button class="logoutButton" @click="logout">Logout</button>
  </div>
</template>

<style scoped>
.logout-container {
  display: flex;
  position: absolute; top: 20px; right: 20px;
}

.logoutButton {
  width: 100%;
  padding: 14px;
  background-color: #007bff;
  border: none;
  color: white;
  font-size: 18px;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 15px;
}
</style>