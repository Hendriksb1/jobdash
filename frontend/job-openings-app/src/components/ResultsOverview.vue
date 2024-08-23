<template>
  <div class="results-overview-widget">
    <h3>Results Overview</h3>
    <Pie :data="chartData" />
  </div>
</template>

<script>
import axios from 'axios';
import { Pie } from 'vue-chartjs';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  CategoryScale,
} from 'chart.js';

// Register necessary Chart.js components
ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale);

export default {
  components: {
    Pie,  // Register the Pie chart component
  },
  data() {
    return {
      resultOverview: [],
      chartData: {
        labels: [],  // Labels for the pie chart
        datasets: [
          {
            label: 'Results Overview',
            backgroundColor: ['#333', '#c84138', '#f44336', '#42b983'],  // Update with your colors
            data: [],  // Data for each label
          },
        ],
      },
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

        // Update chart data based on the API response
        this.updateChartData();
      } catch (error) {
        console.error('Error fetching results overview:', error);
      }
    },
    updateChartData() {
      const labels = this.resultOverview.map(overview => overview.result_name);
      const data = this.resultOverview.map(overview => overview.count_total);

      // Use Object.assign to ensure Vue detects the change
      this.chartData = Object.assign({}, this.chartData, {
        labels: labels,
        datasets: [
          {
            ...this.chartData.datasets[0],
            data: data,
          },
        ],
      });
    }
  },
  watch: {
    resultOverview() {
      // Update chart data when resultOverview changes
      this.updateChartData();
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
  max-width: calc(33.33% - 20px);
  float: left;
  clear: both;
}

h3 {
  color: #e0e0e0;
  margin-bottom: 8px;
  font-size: 1.2rem;
}

canvas {
  width: 100% !important;
  height: auto !important;
  max-width: 400px;
  max-height: 400px;
}
</style>
