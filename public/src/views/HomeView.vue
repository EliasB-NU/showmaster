<script setup>
import Projects from '@/components/Projects.vue'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { onMounted, ref } from 'vue'
import axios from 'axios'

// Message PopUp
const informationPopUp = ref(false);
const informationMessage = ref("");

const closeInformationPopUp = () => {
  informationPopUp.value = false;
}

// Initial checks
const loadContent = ref(false);
const loadNewProjects = ref(false);
const loadEditProjects = ref(false);

onMounted(() => {
  switch (localStorage.getItem('permlvl')) {
    case '0':
      informationPopUp.value = true;
      informationMessage.value = "You currently don't have any permissions, please contact site admin.";
      break;
    case '1':
      loadContent.value = true;
      break;
    case '2':
      loadContent.value = true;
      loadNewProjects.value = true;
      loadEditProjects.value = true;
      break;
    case '3':
      loadContent.value = true;
      loadNewProjects.value = true;
      loadEditProjects.value = true;
      break;
    default:
      informationPopUp.value = true;
      informationMessage.value = "Error loading permissions, please login again.";
      break;
  }
})

// Logic for new project
const newProjectPopUp = ref(false);
const name = ref("")
const rl = ref(false)

const newProject = () => {
  axios
    .post('/api/newproject', {
      name: name.value,
      creator: String(localStorage.getItem('email')),
    })
    .then(() => {
      newProjectPopUp.value = false;
      informationPopUp.value = true;
      informationMessage.value = "Successfully Created new Project";
      rl.value = !rl.value;
      name.value = "";
    })
    .catch((error) => {
      console.log(error)
      informationPopUp.value = true;
      informationMessage.value = "Creation Failed";
    })
}

const openNewProjectPopUp = () => {
  newProjectPopUp.value = true;
}

const closeNewProjectPopUp = () => {
  newProjectPopUp.value = false;
}

</script>

<template>
  <div class="page-container">
    <div class="header-container">
      <Header project=" |  Projects" render-logout=true />
    </div>

    <!-- Projects -->
    <div v-if="loadContent">
      <Projects :reload="rl" :render="loadEditProjects" />
    </div>

    <!-- NewProject PopUp -->
    <div v-if="newProjectPopUp" class="modal-overlay">
      <div class="modal-content">
        <h2 class="text-center">Login</h2>
        <form @submit.prevent="newProject">
          <div class="form-group mb-3">
            <label for="name">Name</label>
            <input
              type="text"
              v-model="name"
              id="email"
              class="input-control"
              placeholder="Enter a name"
              required
            />
          </div>
          <button type="submit" class="btn-primary w-100">Create Project</button>
        </form>
        <div class="footer-links">
          <button @click="closeNewProjectPopUp" class="btn-link">Cancel</button>
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

    <!-- New Project -->
    <div class="newProject-container" v-if="loadNewProjects">
      <b class="newProject-Button" @click="openNewProjectPopUp"><div class="circle"></div></b>
    </div>

    <div class="footer-container">
      <Footer />
    </div>
  </div>
</template>

<style scoped>
/* Overall page container */
.page-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f4f7fa;
  padding: 20px;
  box-sizing: border-box;
}

.header-container {
  padding: 40px;
  z-index: 1000;
}

.footer-container {
  padding: 40px;
  z-index: 1000;
}

.newProject-container {
  display: flex;
  position: fixed;
  bottom: 20px;
  right: 20px;
}

.circle {
  border-radius: 40%;
  width: 60px;
  height: 60px;
  background-color: #007bff;
  cursor: pointer;
}

.circle::before {
  content: "+";
  height: 60px;
  width: 60px;
  font-size: 60px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-family: courier, serif;
  color: white;
}
</style>