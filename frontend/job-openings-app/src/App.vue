<template>
  <div id="app">
    <!-- Conditionally render based on login status -->
    <UserLoginLogout v-if="!loggedInUser" @login-success="handleLoginSuccess" />
    
    <div v-else>
      <AddOpening @opening-added="handleOpeningAdded" />
      <div class="widgets-row">
        <JobsOverview :key="jobKey" />
        <ResultsOverview :key="overviewKey" />
        <ApplicationsPerWeek />
      </div>
      <div class="widgets-row">
        <ApplicationsThisWeek />
      </div>
      <OpeningsList :newOpening="newOpening" @data-changed="handleDataChanged" />
      <button @click="logout">Logout</button>
    </div>
  </div>
</template>

<script>
import AddOpening from './components/AddOpening.vue';
import OpeningsList from './components/OpeningsList.vue';
import ResultsOverview from './components/ResultsOverview.vue';
import JobsOverview from './components/JobsOverview.vue';
import ApplicationsThisWeek from './components/ApplicationsThisWeek.vue';
import ApplicationsPerWeek from './components/ApplicationsPerWeek.vue';
import UserLoginLogout from './components/UserLoginLogout.vue';
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'App',
  components: {
    AddOpening,
    OpeningsList,
    ResultsOverview,
    JobsOverview,
    ApplicationsThisWeek,
    ApplicationsPerWeek,
    UserLoginLogout
  },
  data() {
    return {
      newOpening: null,
      overviewKey: 0, // Key to force ResultsOverview to remount on data change
    };
  },
  computed: {
    // Map Vuex getter to check if user is logged in
    ...mapGetters({
      loggedInUser: 'user' // Using 'user' getter to check logged-in status
    })
  },
  methods: {
    // Map Vuex actions for login and logout
    ...mapActions(['loginUser', 'logoutUser']),
    
    handleLoginSuccess(email) {
      // Pass the email to login action in Vuex
      this.loginUser({ email });
    },
    
    handleOpeningAdded(opening) {
      this.newOpening = opening;
    },
    
    handleDataChanged() {
      // Incrementing key will force ResultsOverview to remount, triggering fetch
      this.overviewKey++;
    },
    
    logout() {
      // Call Vuex action to log out the user
      this.logoutUser();
    }
  }
};
</script>

<style>
html {
  background-color: #000;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

.widgets-row {
  display: flex;
  justify-content: flex-start;
  margin-top: -20px;
}

.widgets-row > div {
  max-height: 460px;
  overflow-y: auto;
  margin: 10px;
  background-color: #1e1e1e;
  padding: 15px;
  border-radius: 10px;
  max-width: calc(20% - 20px);
}

.widgets-row > div.applications-per-week-widget {
  max-width: calc(40% - 20px);
}
</style>