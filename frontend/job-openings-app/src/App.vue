<template>
  <div id="app">
    <AddOpening @opening-added="handleOpeningAdded" />
    <div class="widgets-row">
      <ResultsOverview :key="overviewKey" />
      <JobsOverview :key="jobKey" />
      <ApplicationsThisWeek />
      <ApplicationsPerWeek />
    </div>
    <UserOperations />
  
    <OpeningsList :newOpening="newOpening" @data-changed="handleDataChanged" />
  </div>
</template>

<script>
import AddOpening from './components/AddOpening.vue';
import OpeningsList from './components/OpeningsList.vue';
import ResultsOverview from './components/ResultsOverview.vue';
import JobsOverview from './components/JobsOverview.vue';
import ApplicationsThisWeek from './components/ApplicationsThisWeek.vue';
import ApplicationsPerWeek from './components/ApplicationsPerWeek.vue';
import UserOperations from './components/UserOperations.vue';

export default {
  name: 'App',
  components: {
    AddOpening,
    OpeningsList,
    ResultsOverview,
    JobsOverview,
    ApplicationsThisWeek,
    ApplicationsPerWeek,
    UserOperations
  },
  data() {
    return {
      newOpening: null,
      overviewKey: 0 // Key to force ResultsOverview to remount on data change
    };
  },
  methods: {
    handleOpeningAdded(opening) {
      this.newOpening = opening;
    },
    handleDataChanged() {
      // Incrementing key will force ResultsOverview to remount, triggering fetch
      this.overviewKey++;
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
  justify-content: flex-start; /* Adjust as needed */
  margin-top: -20px;
}

.widgets-row > div {
  max-height: 340px; /* Set a max height */
  overflow-y: auto; /* Enable vertical scrolling if content exceeds max height */
  margin: 10px; /* Space between widgets */
  background-color: #1e1e1e; /* Background color for the widgets */
  padding: 15px; /* Padding inside the widgets */
  border-radius: 10px; /* Rounded corners for the widgets */
  max-width: calc(20% - 20px);
}

.widgets-row > div.applications-per-week-widget {
  max-width: calc(40% - 20px);

}


</style>./components/UserOperations.vue
