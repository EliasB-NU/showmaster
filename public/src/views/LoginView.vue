<script>
import { ref } from 'vue'
import axios from 'axios'
import router from '@/router/index.js'
import Footer from '@/components/Footer.vue'
import Header from '@/components/Header.vue'

export default {
  name: "LoginComponent",
  components: { Header, Footer },
  setup() {
    if (localStorage.getItem('logedIn')) {
      router.push('/home');
    }

    const loginEmail = ref("");
    const loginPassword = ref("");
    const loginPopUp = ref(false);
    const loginErrorMessage = ref("");

    const login = () => {
      // Perform login logic here
      axios
        .post('/api/login', {
          email: loginEmail.value,
          password: loginPassword.value,
        })
        .then((response) => {
          switch (response.status) {
            case 202:
              console.log("Logged in successfully");
              localStorage.setItem('logedIn', JSON.stringify(true));
              localStorage.setItem('email', JSON.stringify(loginEmail.value));
              router.push('/home');
              break;
            default:
              loginPopUp.value = true;
              loginErrorMessage.value = String(response.status);
              break;
          }
        })
        .catch((error) => {
          switch (error.response.status) {
            case 404:
              console.log('User not found');
              loginPopUp.value = true;
              loginErrorMessage.value = "User not found, please register";
              break;
            case 403:
              console.log('Wrong password');
              loginPopUp.value = true;
              loginErrorMessage.value = 'Wrong Password';
              break;
            default:
              loginPopUp.value = true;
              loginErrorMessage.value = String(error.response.data);
              break;
          }
        });
    };

    const closeLoginPopUp = () => {
      loginPopUp.value = false;
    }

    const reqUsername = ref("");
    const reqEmail = ref("");
    const reqPassword = ref("");
    const reqPasswordConfirm = ref("");
    const reqPopUp = ref(false);
    const reqErrorMessage = ref("");

    const register = () => {
      // Opens registration popup
      if (reqPassword.value !== reqPasswordConfirm.value) {
        loginPopUp.value = true;
        loginErrorMessage.value = "Password does not match";
      } else {
        axios
          .post('/api/register', {
            name: reqUsername.value,
            email: reqEmail.value,
            password: reqPassword.value,
          })
          .then((response) => {
            if (response.status === 200) {
              reqPopUp.value = false;
              loginPopUp.value = true;
              loginErrorMessage.value = "Successfully registered";
            }
          })
          .catch((error) => {
            console.log(error);
            loginPopUp.value = true;
            loginErrorMessage.value = "Registration Failed";
          });
      }
    };

    const openReqPopUp = () => {
      reqPopUp.value = true;
    }

    const closeReqPopUp = () => {
      reqPopUp.value = false;
    }

    const adminSite = () => {
      // Redirect to admin page
      router.push('/admin');
    };

    return {
      login,
      loginEmail,
      loginPassword,
      loginPopUp,
      loginErrorMessage,
      closeLoginPopUp,

      register,
      openReqPopUp,
      reqUsername,
      reqEmail,
      reqPassword,
      reqPasswordConfirm,
      reqPopUp,
      reqErrorMessage,
      closeReqPopUp,

      adminSite,
    };
  },
};
</script>

<template>
  <div class="page-container">
    <Header project="" logout=false />

    <!-- Login Section -->
    <div class="login-container">
      <div class="login-box">
        <h2 class="text-center">Login</h2>
        <form @submit.prevent="login">
          <div class="form-group mb-3">
            <label for="email">Email</label>
            <input
              type="email"
              v-model="loginEmail"
              id="email"
              class="input-control"
              placeholder="Enter your email"
              required
            />
          </div>
          <div class="form-group mb-3">
            <label for="password">Password</label>
            <input
              type="password"
              v-model="loginPassword"
              id="password"
              class="input-control"
              placeholder="Enter your password"
              required
            />
          </div>
          <button type="submit" class="btn-primary w-100">Login</button>
        </form>
        <div class="footer-links">
          <button @click="adminSite" class="btn-link">Admin Site</button>
          <button @click="openReqPopUp" class="btn-link">Register</button>
        </div>
      </div>
    </div>

    <!-- Registration PopUp -->
    <div v-if="reqPopUp" class="modal-overlay">
      <div class="login-box">
        <h2 class="text-center">Login</h2>
        <form @submit.prevent="register">
          <div class="form-group mb-3">
            <label for="email">Email</label>
            <input
              type="email"
              v-model="reqEmail"
              id="email"
              class="input-control"
              placeholder="Enter your email"
              required
            />
          </div>
          <div class="form-group mb-3">
            <label for="text">Username</label>
            <input
              type="text"
              v-model="reqUsername"
              id="username"
              class="input-control"
              placeholder="Enter your username"
              required
            />
          </div>
          <div class="form-group mb-3">
            <label for="password">Password</label>
            <input
              type="password"
              v-model="reqPassword"
              id="password"
              class="input-control"
              placeholder="Enter your password"
              required
            />
          </div>
          <div class="form-group mb-3">
            <input
              type="password"
              v-model="reqPasswordConfirm"
              id="passwordConfirm"
              class="input-control"
              placeholder="Confirm your password"
              required
            />
          </div>
          <button type="submit" class="btn-primary w-100">Register</button>
        </form>
        <div class="footer-links">
          <button @click="closeReqPopUp" class="btn-link">cancel</button>
        </div>
      </div>
    </div>

    <!-- Modal for Errors -->
    <div v-if="loginPopUp" class="modal-overlay">
      <div class="modal-content">
        <p>{{ loginErrorMessage }}</p>
        <button @click="closeLoginPopUp" class="btn-primary">Confirm</button>
      </div>
    </div>

    <Footer/>
  </div>
</template>

