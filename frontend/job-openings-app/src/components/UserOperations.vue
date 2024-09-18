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

    <div v-if="loading" class="loading-spinner">Loading...</div>
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
      loading: false, // For tracking loading state
    };
  },
  methods: {
    async apiRequest(url, method, body) {
      this.loading = true;
      try {
        const response = await fetch(url, {
          method: method,
          headers: { 'Content-Type': 'application/json' },
          body: body ? JSON.stringify(body) : null,
        });

        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }

        return await response.json();
      } catch (error) {
        alert(error.message);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async registerUser() {
      if (!this.registerName || !this.registerEmail) {
        alert('Please enter both name and email');
        return;
      }

      try {
        await this.apiRequest('http://localhost:8080/registerUser', 'POST', {
          name: this.registerName,
          email: this.registerEmail,
        });

        this.registerName = '';
        this.registerEmail = '';
        alert('User registered successfully');
      } catch (error) {
        console.error('Registration failed');
      }
    },

    async getUser() {
      if (!this.getUserId) {
        alert('Please enter a User ID');
        return;
      }

      try {
        const data = await this.apiRequest(
          `http://localhost:8080/getUser?id=${this.getUserId}`,
          'GET'
        );
        this.user = data;
      } catch (error) {
        console.error('Failed to get user');
      }
    },

    async unregisterUser() {
      if (!this.unregisterUserId) {
        alert('Please enter a User ID to unregister');
        return;
      }

      try {
        await this.apiRequest('http://localhost:8080/unRegisterUser', 'POST', {
          id: this.unregisterUserId,
        });

        this.unregisterUserId = '';
        alert('User unregistered successfully');
      } catch (error) {
        console.error('Failed to unregister user');
      }
    },

    async changeUser() {
      if (!this.changeUserId || !this.changeName || !this.changeEmail) {
        alert('Please fill in all fields to update the user');
        return;
      }

      try {
        await this.apiRequest('http://localhost:8080/changeUser', 'POST', {
          id: this.changeUserId,
          name: this.changeName,
          email: this.changeEmail,
        });

        this.changeUserId = '';
        this.changeName = '';
        this.changeEmail = '';
        alert('User updated successfully');
      } catch (error) {
        console.error('Failed to update user');
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
.loading-spinner {
  color: #2b2b2b;
  font-weight: bold;
}
</style>