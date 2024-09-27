<script>
import { ref } from 'vue'
import axios from 'axios'
import router from '@/router/index.js'

export default {
  name: "LoginComponent",
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
    <!-- Header Section -->
    <header class="header">
      <div class="header-title">ShowMaster V3</div>
    </header>

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
        <button @click="closeLoginPopUp" class="btn-primary">Close</button>
      </div>
    </div>

    <!-- Footer -->
    <footer class="footer">
      <span>Copyright: Elias Braun 2024</span>
      <a href="https://github.com/EliasB-NU/showmaster" target="_blank" class="footer-link">
        GitHub
      </a>
    </footer>
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

/* Login Stuff */
.login-box {
  background-color: #fff;
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  min-width: 350px;
  max-width: 1000px;
  margin-bottom: 15px;
}

/* Centering the login box */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-grow: 0.5;
}

/* Styling the login form */
h2 {
  font-size: 28px;
  margin-bottom: 25px;
  text-align: center;
}

.form-group {
  margin-bottom: 25px;
}

label {
  font-size: 16px;
  color: #555;
  margin-bottom: 8px;
  display: block;
}

.input-control {
  width: 100%;
  padding: 14px;
  border-radius: 6px;
  border: 1px solid #ddd;
  font-size: 18px;
  box-sizing: border-box;
}

/* Primary Button */
.btn-primary {
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

.btn-primary:hover {
  background-color: #0056b3;
}

/* Links for password reset and register */
.footer-links {
  text-align: center;
  margin-top: 25px;
}

.btn-link {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-size: 16px;
  text-decoration: none;
}

.btn-link:hover {
  text-decoration: underline;
}

/* PopUp */
/* Modal Overlay */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5); /* Dark background with transparency */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

/* Modal Content */
.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  position: relative;
}

/* Button Inside Modal */
.modal-content button {
  margin-top: 15px;
  padding: 10px 20px;
}
</style>
