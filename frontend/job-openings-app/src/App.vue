<template>
  <div id="app">
    <!-- Conditionally render based on login status -->
    <UserLoginLogout v-if="!loggedInUser" @login-success="handleLoginSuccess" />

    <div v-else>
      <div v-if="!isInRouterView">hello</div>

      <!-- Display main application components when not in a router view -->
      <AddOpening @opening-added="handleOpeningAdded" v-if="!isInRouterView" />
      <div class="widgets-row" v-if="!isInRouterView">
        <JobsOverview :key="jobKey" />
        <ResultsOverview :key="overviewKey" />
        <ApplicationsPerWeek />
      </div>
      <div class="widgets-row" v-if="!isInRouterView">
        <ApplicationsThisWeek />
      </div>

      <OpeningsList 
        :newOpening="newOpening"
        @data-changed="handleDataChanged"
        v-if="!isInRouterView"
        @show-router-view="handleShowRouterView" 
      />

      <button @click="logout" v-if="!isInRouterView">Logout</button>

      <!-- Router view for displaying route-specific components -->
      <router-view v-if="isInRouterView" @route-change="handleShowRouterView" />
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
      isInRouterView: false, // Track if we are in a router view
    };
  },
  computed: {
    ...mapGetters({
      loggedInUser: 'user' // Using 'user' getter to check logged-in status
    })
  },
  methods: {
    ...mapActions(['loginUser', 'logoutUser']),
    
    handleLoginSuccess(email) {
      this.loginUser({ email });
    },
    
    handleOpeningAdded(opening) {
      this.newOpening = opening;
    },
    
    handleDataChanged() {
      this.overviewKey++;
    },
    
    logout() {
      this.logoutUser();
    },
    
    handleShowRouterView() {
      console.log("Navigating to a new route");
      this.isInRouterView = true; // Set to true when navigating to a new route
    }
  },
  watch: {
    $route(to) {
      // Check if the new route is a detail route and set visibility accordingly
      this.isInRouterView = to.name === 'job-detail'; // Update this condition as needed for your routes
    }
  },
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