<template>
    <div>
      <h1>User Management</h1>
  
      <div>
        <h2>Register User</h2>
        <input v-model="registerName" placeholder="Name" />
        <input v-model="registerEmail" placeholder="Email" />
        <button @click="registerUser">Register</button>
      </div>
  
      <div>
        <h2>Get User</h2>
        <input v-model="getUserId" placeholder="User ID" />
        <button @click="getUser">Get User</button>
        <div v-if="user">
          <p>ID: {{ user.id }}</p>
          <p>Name: {{ user.name }}</p>
          <p>Email: {{ user.email }}</p>
        </div>
      </div>
  
      <div>
        <h2>Unregister User</h2>
        <input v-model="unregisterUserId" placeholder="User ID" />
        <button @click="unregisterUser">Unregister</button>
      </div>
  
      <div>
        <h2>Change User</h2>
        <input v-model="changeUserId" placeholder="User ID" />
        <input v-model="changeName" placeholder="New Name" />
        <input v-model="changeEmail" placeholder="New Email" />
        <button @click="changeUser">Change User</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        registerName: '',
        registerEmail: '',
        getUserId: '',
        user: null,
        unregisterUserId: '',
        changeUserId: '',
        changeName: '',
        changeEmail: '',
      };
    },
    methods: {
      async registerUser() {
        try {
        const response = await fetch('http://localhost:8080/registerUser', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              name: this.registerName,
              email: this.registerEmail,
            }),
          });
  
          if (!response.ok) {
            throw new Error('Failed to register user');
          }
  
          this.registerName = '';
          this.registerEmail = '';
          alert('User registered successfully');
        } catch (error) {
          alert(error.message);
        }
      },
      async getUser() {
        try {
          const response = await fetch(`http://localhost:8080/getUser?id=${this.getUserId}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
          });

          if (!response.ok) {
            throw new Error('Failed to get user');
          }

          this.user = await response.json();
        } catch (error) {
          alert(error.message);
        }
      },
      async unregisterUser() {
        try {
          const response = await fetch('http://localhost:8080/unRegisterUser', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: this.unregisterUserId }),
          });
  
          if (!response.ok) {
            throw new Error('Failed to unregister user');
          }
  
          this.unregisterUserId = '';
          alert('User unregistered successfully');
        } catch (error) {
          alert(error.message);
        }
      },
      async changeUser() {
        try {
          const response = await fetch('http://localhost:8080/changeUser', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              id: this.changeUserId,
              name: this.changeName,
              email: this.changeEmail,
            }),
          });
  
          if (!response.ok) {
            throw new Error('Failed to change user');
          }
  
          this.changeUserId = '';
          this.changeName = '';
          this.changeEmail = '';
          alert('User updated successfully');
        } catch (error) {
          alert(error.message);
        }
      },
    },
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
  