<script>
import { ref } from 'vue'
import axios from 'axios'

export default {
  name: "Projects",
  setup() {
    const informationPopUp = ref(false)
    const informationMessage = ref("")

    const openInformationPopUp = () => {
      informationPopUp.value = true;
    }

    const closeInformationPopUp = () => {
      informationPopUp.value = false;
    }

    const projects = ref([])

    axios
      .get('/api/getprojects')
      .then((res) => {
        console.log(res.status);
        projects.value = res.data;
      })
      .catch((err) => {
        informationPopUp.value = true;
        informationMessage.value = "Something went wrong";
        console.log(err)
      });

    return {
      informationPopUp,
      informationMessage,
      openInformationPopUp,
      closeInformationPopUp,

      projects,
    }
  },
}
</script>

<template>
  <div class="page-container-projects">

    <!-- Project Renderer -->
    <div class="projects-grid">
      <div v-for="project in projects" :key="project.id">
        <div class="project-container">
          <div class="project-box">
            <h2 class="text-center">{{ project.name }}</h2>
            <span>Creator: {{ project.creator }}</span><br>
            <span>Timer: {{ project.timer }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Information PopUp -->
    <div v-if="informationPopUp" class="modal-overlay">
      <div class="modal-content">
        <p>{{ informationMessage }}</p>
        <button @click="closeInformationPopUp" class="btn-primary">Confirm</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-container-projects {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f4f7fa;
  padding: 20px;
  box-sizing: border-box;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); /* Responsive grid */
  gap: 20px; /* Space between grid items */
  width: 100%; /* Full width grid */
  max-width: 1200px; /* Optional max-width to control the grid size */
  margin: 0 auto; /* Center grid horizontally */
}

.project-box {
  background-color: #fff;
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  min-width: 300px; /* Ensure a minimum width */
  margin-bottom: 15px;
}

.project-container {
  display: flex;
  justify-content: center;
}
</style>