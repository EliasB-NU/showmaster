<script setup>
  import router from '@/router/index.js'

  defineProps([
    'renderLogout',
    'project'
  ]);

  const logout = () => {
    localStorage.clear();
    router.push('/');
  };
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

<style>
.header {
  background-color: #fff;
  padding: 20px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  max-width: 450px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
  text-align: center;
  position: fixed; top: 20px; left: 20px;
}

.header-title {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  color: #007bff;
}

.logout-container {
  display: flex;
  position: fixed; top: 20px; right: 20px;
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