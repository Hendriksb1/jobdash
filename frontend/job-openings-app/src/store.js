import { createStore } from 'vuex';
import axios from 'axios';
import createPersistedState from 'vuex-persistedstate'; // Import vuex-persistedstate

export default createStore({
  state: {
    users: [], // List of registered users
    loggedInUser: null // Currently logged-in user
  },
  getters: {
    user: (state) => state.loggedInUser,
    getUserByEmail: (state) => (email) => {
      return state.users.find(user => user.email === email);
    }
  },
  mutations: {
    SET_USER(state, user) {
      state.loggedInUser = user;
    },
    ADD_USER(state, user) {
      state.users.push(user);
    },
    LOGOUT_USER(state) {
      state.loggedInUser = null;
    }
  },
  actions: {
    async loginUser({ commit }, { email }) {
      try {
        // Make a request to your getUser endpoint
        const response = await axios.get(`http://localhost:8080/getUserByEmail?email=${email}`);
        const user = response.data;

        if (user) {
          commit('SET_USER', user);
        } else {
          throw new Error('User not found');
        }
      } catch (error) {
        alert(error.message);
      }
    },
    async registerUser({ commit }, user) {
      try {
        // Make a request to register the user (assuming you have an endpoint for this)
        await axios.post('/api/registerUser', user);
        commit('ADD_USER', user);
      } catch (error) {
        alert('Registration failed');
      }
    },
    logoutUser({ commit }) {
      commit('LOGOUT_USER');
    }
  },
  plugins: [
    createPersistedState({
      key: 'my-app', // The key to store in localStorage
      paths: ['loggedInUser'], // Persist only the loggedInUser state
    })
  ]
});