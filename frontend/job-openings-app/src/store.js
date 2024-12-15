import { createStore } from 'vuex';
import axios from 'axios';
import createPersistedState from 'vuex-persistedstate'; 
import Cookies from 'js-cookie'; // Import js-cookie

const appName = "jobDash"

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
        let encodedEmail = encodeURIComponent(email);
        // Make a request to your getUser endpoint
        const response = await axios.get(`http://localhost:8080/getUserByEmail?email=${encodedEmail}`);
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
        console.log("register user in store.js");
        // Make a request to register the user (assuming you have an endpoint for this)
        await axios.post('http://localhost:8080/registerUser', user);
        commit('ADD_USER', user);
      } catch (error) {
        console.error(error);
        alert('Registration failed');
      }
    },
    logoutUser({ commit }) {
      commit('LOGOUT_USER');
      Cookies.remove(appName);  // Clear persisted cookie data
    }
  },
  plugins: [
    createPersistedState({      
      key: appName, // The key to store in cookies
      paths: ['loggedInUser'], // Persist only the loggedInUser state
      storage: {
        // Custom storage object to use cookies
        getItem: (key) => Cookies.get(key),  // Retrieve from cookie
        setItem: (key, value) => Cookies.set(key, value, { expires: 3, secure: false }), // Store in cookie with expiration and secure flag
        removeItem: (key) => Cookies.remove(key) // Remove from cookie
      }
    })
  ]
});