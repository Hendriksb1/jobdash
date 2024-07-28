<template>
  <div class="results-overview-widget">
    <h3>Results Overview</h3>
    <table>
      <tbody>
        <tr v-for="overview in resultOverview" :key="overview.result_name">
          <td>{{ overview.result_name }}</td>
          <td>{{ overview.count_total }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      resultOverview: []
    };
  },
  created() {
    this.fetchResultsOverview();
  },
  methods: {
    async fetchResultsOverview() {
      try {
        const response = await axios.get('http://localhost:8080/getOpeningsByResult');
        this.resultOverview = response.data;
      } catch (error) {
        console.error('Error fetching results overview:', error);
      }
    }
  },
  watch: {
    // Watch for changes in key prop to fetch data again
    '$route': {
      handler() {
        this.fetchResultsOverview();
      },
      immediate: true
    }
  }
};
</script>

<style scoped>
.results-overview-widget {
  background-color: #1e1e1e;
  padding: 20px;
  border-radius: 10px;
  margin: 20px;
  text-align: left;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  max-width: calc(33.33% - 20px); /* One-third of the maximum width minus margins */
  float: left; /* Ensures it occupies only the necessary width */
  clear: both; /* Clears any previous floats */
}

h3 {
  color: #e0e0e0;
  margin-bottom: 8px;
  font-size: 1.2rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 8px;
  color: #e0e0e0;
}

td {
  padding: 6px;
  text-align: left;
  border-bottom: 1px solid #444;
}

tr:last-child td {
  border-bottom: none;
}
</style>
