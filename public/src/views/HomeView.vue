<script>
import router from '@/router/index.js'

export default {
  name: "HomeView",
  setup () {
    const logout = () => {
      localStorage.setItem('logedIn', JSON.stringify(false));
      router.push("/");
    }

    return {
      logout,
    }
  }
}
</script>

<template>
  <h1>Hello World</h1>
  <button @click="logout">Logout</button>

</template>
