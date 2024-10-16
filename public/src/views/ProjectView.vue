<script setup>

import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Footer from '@/components/Footer.vue'
import axios from 'axios'
import InformationPopUp from '@/components/InformationPopUp.vue'

const route = useRoute();
const infoPopUp = ref(false);
const message = ref("");

watch (
  () => route.params.projectId,
  (newId) => {
    getData();
  }
)

function getData() {
  axios
    .get('/api/getrows'+route.params.projectId+'table')
    .then(res => {

    })
    .catch((error) => {

    })
}

</script>

<template>
  <Header project=" |  {{ route.params.projctId }}" render-logout=true />

  <InformationPopUp :shown="infoPopUp" :message="message"/>
  <Footer />
</template>

<style scoped>

</style>