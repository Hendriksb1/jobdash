<template>
  <div id="app">
    <AddOpening @opening-added="handleOpeningAdded" />
    <div class="widgets-row">
      <ResultsOverview :key="overviewKey" />
      <JobsOverview :key="jobKey" />
      <ApplicationsThisWeek />
      <ApplicationsPerWeek />
    </div>
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

export default {
  name: 'App',
  components: {
    AddOpening,
    OpeningsList,
    ResultsOverview,
    JobsOverview,
    ApplicationsThisWeek,
    ApplicationsPerWeek
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
  justify-content: start; /* Adjust as needed */
  margin-top: -20px;
}
</style>
