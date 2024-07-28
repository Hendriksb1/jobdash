<template>
  <div class="applications-per-week-widget">
    <h3>Applications Per Week</h3>
    <canvas id="applicationsPerWeekChart"></canvas>
  </div>
</template>

<script>
import axios from 'axios';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  LineController,
  CategoryScale,
  LinearScale,
  PointElement
} from 'chart.js';
import { getWeek, parseISO } from 'date-fns';

ChartJS.register(Title, Tooltip, Legend, LineElement, LineController, CategoryScale, LinearScale, PointElement);

export default {
  data() {
    return {
      chartData: {
        labels: [],
        datasets: [
          {
            label: 'Applications Count',
            data: [],
            backgroundColor: 'rgba(66, 185, 131, 0.2)',
            borderColor: '#42b983',
            borderWidth: 1,
            fill: true,
          }
        ]
      },
      chartOptions: {
        responsive: true,
        maintainAspectRatio: true,
        scales: {
          x: {
            title: {
              display: true,
              text: 'Week Number'
            }
          },
          y: {
            title: {
              display: true,
              text: 'Applications Count'
            },
            beginAtZero: true
          }
        }
      },
      chart: null
    };
  },
  created() {
    this.fetchApplicationsPerWeekCount();
  },
  methods: {
    async fetchApplicationsPerWeekCount() {
      try {
        const response = await axios.get('http://localhost:8080/getCountsPerWeek');
        const data = response.data;

        // Process the data to extract labels (week numbers) and counts
        const labels = data.map(entry => `Week ${getWeek(parseISO(entry.week_end_date))}`);
        const counts = data.map(entry => entry.entry_count);

        this.chartData.labels = labels;
        this.chartData.datasets[0].data = counts;

        this.renderChart();
      } catch (error) {
        console.error('Error fetching applications count per week:', error);
      }
    },
    renderChart() {
      const ctx = document.getElementById('applicationsPerWeekChart').getContext('2d');
      if (this.chart) {
        this.chart.destroy(); // Destroy existing chart instance to avoid duplication
      }
      this.chart = new ChartJS(ctx, {
        type: 'line',
        data: this.chartData,
        options: this.chartOptions
      });
    }
  },
  beforeUnmount() {
    if (this.chart) {
      this.chart.destroy();
    }
  }
};
</script>

<style scoped>
.applications-per-week-widget {
  background-color: #1e1e1e;
  padding: 20px;
  border-radius: 10px;
  margin: 20px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  max-width: calc(33.33% - 20px); /* One-third of the maximum width minus margins */
  float: left; /* Ensures it occupies only the necessary width */
}

h3 {
  color: #e0e0e0;
  margin-bottom: 8px;
  font-size: 1.2rem;
}

canvas {
  max-width: 100%;
  height: 400px; /* Set a fixed height */
}
</style>
