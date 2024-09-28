<script setup>
  import axios from 'axios'
  import { onMounted, ref, watch } from 'vue'
  import router from '@/router/index.js'

  const props = defineProps([
    'reload',
    'edit',
    'delete',
  ]);
  let projects = ref([]);

  const informationPopUp = ref(false);
  const informationMessage = ref("");
  const closeInformationPopUp = () => {
    informationPopUp.value = false;
  }

  onMounted(() => {
    axios
      .get('/api/getprojects')
      .then(res => {
        projects.value = res.data;
      })
      .catch((error) => {
        console.log(error);
        informationPopUp.value = true;
        informationMessage.value = "Failed to get projects";
      })
  })

  const editProjectPopUp = ref(false);
  const newName = ref("");
  const projectToEdit = ref("");

  const editProject = () => {
    axios
      .patch('/api/updateproject', {
        oldname: projectToEdit.value,
        newname: newName.value,
      })
      .then(() => {
        editProjectPopUp.value = false;
        informationPopUp.value = true;
        informationMessage.value = "Successfully Updated Project";
        getProjects();
      })
      .catch((error) => {
        informationPopUp.value = true;
        informationMessage.value = "Failed to update project";
        console.log(error);
      })
  }

  const openEditProjectPopUp = (data) => {
    projectToEdit.value = data.name;
    editProjectPopUp.value = true;
  }

  const closeEditProjectPopUp = () => {
    editProjectPopUp.value = false;
  }

  const deleteConfirmPopUp = ref(false);
  const projectToDelete = ref("")
  const startDeletion = (data) => {
    projectToDelete.value = data.name;
    deleteConfirmPopUp.value = true;
  }

  const cancelDelete = () => {
    deleteConfirmPopUp.value = false;
  }

  const deleteProject = () => {
    axios
      .delete('/api/deleteproject', {
        data: {
          name: projectToDelete.value,
        }
      })
      .then(() => {
        deleteConfirmPopUp.value = false;
        informationPopUp.value = true;
        informationMessage.value = "Successfully Deleted Project";
        getProjects();
      })
      .catch((error) => {
        console.log(error);
        informationPopUp.value = true;
        informationMessage.value = "Failed to delete Project";
      })
  }

  watch(
    () => props.reload,
    () => {
      getProjects();
    },
    { deep: true }
  )

  function getProjects() {
    axios
      .get('/api/getprojects')
      .then(res => {
        projects.value = res.data;
      })
      .catch((error) => {
        console.log(error);
        informationPopUp.value = true;
        informationMessage.value = "Failed to get projects";
      })
  }

  const projectViewRoute = (data) => {
    const firstStep = String(data.name);
    const secondStep = firstStep.replaceAll(" ", "_");
    const projectName = secondStep.toLowerCase();

    router.push('/project/'+projectName);
  }
</script>

<template>
  <div class="page-container-projects">

    <!-- Project Renderer -->
    <div class="projects-grid" id="projects-grid">
      <div v-for="project in projects" :key="project.name">
        <div class="project-container" @click="projectViewRoute(project)">
          <div class="project-box">
            <h2 class="text-center">{{ project.name }}</h2>
            <span>Creator: {{ project.creator }}</span><br>
            <span>Timer: {{ project.timer }}</span>
            <div v-if="edit">
              <button @click="openEditProjectPopUp(project)" class="btn-primary w-100">Edit Project</button>
            </div>
            <div v-if="delete">
              <button @click="startDeletion(project)" class="btn-primary w-100" style="background: red">Delete Project</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit PopUp -->
    <div v-if="editProjectPopUp" class="modal-overlay">
      <div class="modal-content">
        <h2 class="text-center">Edit Project</h2>
        <form @submit.prevent="editProject">
          <div class="form-group mb-3">
            <label for="name">New Name</label>
            <input
              type="text"
              v-model="newName"
              id="email"
              class="input-control"
              placeholder="Enter new name"
              required
            />
          </div>
          <button type="submit" class="btn-primary">Update Project</button>
        </form>
        <div class="footer-links">
          <button @click="closeEditProjectPopUp" class="btn-link">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm PopUp -->
    <div v-if="deleteConfirmPopUp" class="modal-overlay">
      <div class="modal-content">
        <h2 class="text-center">DO YOU REALLY WANT TO DELETE THIS PROJECT?</h2>
        <button @click="cancelDelete" class="btn-primary">Cancel</button>
        <button @click="deleteProject" class="btn-primary" style="background: red">Confirm</button>
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
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.project-box {
  background-color: #fff;
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  min-width: 300px;
  margin-bottom: 15px;
}

.project-container {
  display: flex;
  justify-content: center;
}

.project-container:hover {
  cursor: pointer;
}
</style>