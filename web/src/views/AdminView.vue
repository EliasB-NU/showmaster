<script setup lang="ts">
import HeaderComponent from '@/components/HeaderComponent.vue'
import UserEditComponent from '@/components/UserEditComponent.vue'
import UserCreateComponent from '@/components/UserCreateComponent.vue'
import PopUp from '@/components/PopUp.vue'
import { onMounted, ref } from 'vue'
import Cookies from 'js-cookie'
import axios from 'axios'
import ClientEditComponent from '@/components/ClientEditComponent.vue'
import ClientCreateComponent from '@/components/ClientCreateComponent.vue'

const popup = ref<InstanceType<typeof PopUp> | null>(null);

interface perms {
  login: boolean
  admin: boolean
  events: number
}
interface User {
  id: number
  email: string
  name: string
  permissions: perms

}

const users = ref<User[]>([])
const showCreateUser = ref<boolean>(false)
const showEditUser = ref<boolean>(false)
const userToEdit = ref<User>({} as User)
const loadingUsers = ref<boolean>(false)

async function fetchUsers() {
  try {
    await axios
      .get('/api/getUsers', {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${Cookies.get('token')}`,
        },
      })
      .then((response) => {
        users.value = response.data
        loadingUsers.value = false
      })
      .catch((error) => {
        console.error('Error fetching users:', error)
        popup.value?.show('Error fetching users')
        loadingUsers.value = false
      })
  } catch (error) {
    popup.value?.show('Error fetching users')
    console.error('Error fetching users:', error)
    loadingUsers.value = false
  }
}

const editUser = (user: User) => {
  userToEdit.value = user
  showEditUser.value = true
}

const deleteUser = async (id: Number) => {
  try {
    await axios
      .delete(`/api/deleteUser/${id}`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${Cookies.get('token')}`,
        },
      })
      .then((response) => {
        if (response.status === 200) {
          popup.value?.show('User deleted successfully')
          fetchUsers()
        } else {
          popup.value?.show('Error deleting user')
        }
      })
  } catch (error) {
    popup.value?.show('Error deleting user')
    console.error('Error deleting user:', error)
  }
}


interface Client {
  id: number
  name: string
  type: string
  ip: string
  port: number
}

const clients = ref<Client[]>([])
const loadingClients = ref<boolean>(false)
const showCreateClient = ref<boolean>(false)
const showEditClient = ref<boolean>(false)
const clientToEdit = ref<Client>({} as Client)

async function fetchClients() {
  try {
    await axios
      .get('/api/getClients', {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${Cookies.get('token')}`,
        }
      })
    .then((response) => {
      clients.value = response.data
    })
  } catch (error) {
    popup.value?.show('Error getting clients')
  }
}

const editClient = (client: Client) => {
  showEditClient.value = true
  clientToEdit.value = client
}

const deleteClient = async (id: number) => {
  try {
    await axios
    .delete(`/api/deleteClient/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${Cookies.get('token')}`,
      }
    })
    .then((response) => {
      if (response.status === 200) {
        popup.value?.show('Successfully deleted client')
      }
    })
  } catch (error) {
    console.error('Error deleting client:', error)
    popup.value?.show('Error deleting client')
  }
}

onMounted(async () => {
  await fetchClients()
  await fetchUsers()

})
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <HeaderComponent />
    <div class="grid gap-10 md:grid-cols-1 lg:grid-cols-2">
      <div class="max-w-4xl mx-auto p-6 space-y-2">
        <div class="max-w-7xl mx-auto">
          <h2 class="text-2xl font-semibold mb-4">Users</h2>
          <button
            @click="showCreateUser = true" class="bg-gray-800 text-white px-4 py-2 rounded-lg hover:bg-gray-700">
            Create User
          </button>
        </div>

        <!-- Loading State -->
        <p v-if="loadingUsers" class="text-gray-500">Loading...</p>

        <!-- Members Table -->
        <div v-else class="bg-white shadow rounded-lg">
          <table class="w-full border-collapse">
            <thead>
            <tr class="bg-gray-200">
              <th class="p-3 text-left">Name</th>
              <th class="p-3 text-left">Email</th>
              <th class="p-3">Actions</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="user in users" :key="user.id" class="border-t">
              <td class="p-3">{{ user.name }}</td>
              <td class="p-3">{{ user.email }}</td>
              <td class="p-3 flex space-x-2">
                <button
                  v-if="Cookies.get('admin') === 'true'"
                  @click="editUser(user)"
                  class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700"
                >
                  Edit
                </button>
                <button
                  v-if="Cookies.get('admin') === 'true'"
                  @click="deleteUser(user.id)"
                  class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
                >
                  Delete
                </button>
              </td>
            </tr>
            </tbody>
          </table>
        </div>

        <!-- MembersEditComponent Popup -->
        <UserEditComponent
          :visible="showEditUser"
          :user="userToEdit"
          @close="showEditUser = false"
          @updatedUser="fetchUsers"
        />

        <UserCreateComponent
          :visible="showCreateUser"
          @close="showCreateUser = false"
          @created-user="fetchUsers"
        />
      </div>
      <div class="max-w-4xl mx-auto p-6 space-y-2">
        <div>
          <h2 class="text-2xl font-semibold p-6 space-y-2">Clients</h2>
          <button
            @click="showCreateClient = true" class="bg-gray-800 text-white px-4 py-2 rounded-lg hover:bg-gray-700">
            Create Client
          </button>
        </div>

        <!-- Loading State -->
        <p v-if="loadingClients" class="text-gray-500">Loading...</p>

        <!-- Clients Table -->
        <div v-else class="bg-white shadow rounded-lg">
          <table class="w-full border-collapse">
            <thead>
            <tr class="bg-gray-200">
              <th class="p-3 text-left">Name</th>
              <th class="p-3 text-left">Type</th>
              <th class="p-3">Actions</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="client in clients" :key="client.id" class="border-t">
              <td class="p-3">{{ client.name }}</td>
              <td class="p-3">{{ client.type }}</td>
              <td class="p-3 flex space-x-2">
                <button
                  v-if="Cookies.get('admin') === 'true'"
                  @click="editClient(client)"
                  class="bg-gray-800 text-white px-3 py-1 rounded hover:bg-gray-700"
                >
                  Edit
                </button>
                <button
                  v-if="Cookies.get('admin') === 'true'"
                  @click="deleteClient(client.id)"
                  class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
                >
                  Delete
                </button>
              </td>
            </tr>
            </tbody>
          </table>
        </div>

        <ClientEditComponent
          :user="clientToEdit"
          :visible="showEditUser"
          @close="showEditUser = false"
          @clientEdited="fetchClients"
        />

        <ClientCreateComponent
          :visible="showCreateClient"
          @close="showCreateClient = false"
          @clientCreated="fetchClients"
        />
      </div>
    </div>
    <PopUp ref="popup" />
  </div>
</template>

<style scoped>

</style>
