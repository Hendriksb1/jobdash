<template>
    <div>
      <h1>User Management</h1>
  
      <!-- Register User -->
      <div>
        <h2>Register User</h2>
        <input v-model="registerName" placeholder="Name" />
        <input v-model="registerEmail" placeholder="Email" />
        <button @click="registerUser">Register</button>
      </div>
  
      <!-- Login User -->
      <div>
        <h2>Login User</h2>
        <input v-model="loginEmail" placeholder="Email" />
        <!-- <input v-model="loginPassword" type="password" placeholder="Password" /> -->
        <button @click="loginUser">Login</button>
      </div>
  
      <!-- Current User -->
      <div v-if="user">
        <h2>Current User</h2>
        <p>Logged in as: {{ user.name }} ({{ user.email }})</p>
        <button @click="logoutUser">Logout</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        registerName: '',
        registerEmail: '',
        loginEmail: '',
        // loginPassword: ''
      };
    },
    computed: {
      user() {
        return this.$store.getters.user;
      }
    },
    methods: {
      async registerUser() {
        if (!this.registerName || !this.registerEmail) {
          alert('Please enter both name and email');
          return;
        }
  
        const newUser = {
          name: this.registerName,
          email: this.registerEmail
        };
  
        // Register user
        await this.$store.dispatch('registerUser', newUser);
  
        this.registerName = '';
        this.registerEmail = '';
      },
  
      async loginUser() {
        if (!this.loginEmail) {
          alert('Please enter email');
          return;
        }
  
        // Log in user
        await this.$store.dispatch('loginUser', { email: this.loginEmail, password: this.loginPassword });
  
        this.loginEmail = '';
        // this.loginPassword = '';
      },
  
      logoutUser() {
        this.$store.dispatch('logoutUser');
      }
    }
  };
  </script>
  
  <style scoped>
  input {
    margin-right: 8px;
  }
  button {
    margin-bottom: 16px;
  }
  </style>