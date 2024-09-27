<script>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref } from 'vue'
import Projects from '@/components/Projects.vue'
import axios from 'axios'

export default {
  name: "HomeView",
  components: { Projects, Footer, Header },
  setup () {
    const informationPopUp = ref(false)
    const informationMessage = ref("")

    const openInformationPopUp = () => {
      informationPopUp.value = true;
    }

    const closeInformationPopUp = () => {
      informationPopUp.value = false;
    }


    const newProjectPopUp = ref(false)
    const name = ref("")

    const openNewProjectPopUp = () => {
      newProjectPopUp.value = true;
    }

    const newProject = () =>  {
      axios
        .post('/api/newproject', {
          name: name.value,
          creator: String(localStorage.getItem('email')),
        })
        .then(() => {
          newProjectPopUp.value = false;
          informationPopUp.value = true;
          informationMessage.value = "Successfully Created new Project";
        })
        .catch((error) => {
          console.log(error)
          informationPopUp.value = true;
          informationMessage.value = "Creation Failed";
        })
    }

    const closeNewProjectPopUp = () => {
      newProjectPopUp.value = false;
    }

    return {
      informationPopUp,
      informationMessage,
      openInformationPopUp,
      closeInformationPopUp,

      newProject,
      newProjectPopUp,
      openNewProjectPopUp,
      closeNewProjectPopUp,
      name,
    }
  }
}
</script>

<template>
  <div class="page-container">
    <Header project=" | Projects" render-logout=true />

    <!-- Projects -->
    <Projects />

    <!-- NewProject PopUp -->
    <div v-if="newProjectPopUp" class="modal-overlay">
      <div class="login-box">
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
    <div class="newProject-container">
      <b class="newProject-Button" @click="openNewProjectPopUp"><div class="circle"></div></b>
    </div>
    <Footer />
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

.newProject-container {
  display: flex;
  position: absolute; bottom: 20px; right: 20px;
  align-items: center;
  justify-content: center;
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