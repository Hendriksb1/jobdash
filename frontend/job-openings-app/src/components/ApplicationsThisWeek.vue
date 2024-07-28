<template>
    <div class="applications-this-week-widget">
      <h3>Applications This Week</h3>
      <div class="count">{{ applicationsCount }}</div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    data() {
      return {
        applicationsCount: 0
      };
    },
    created() {
      this.fetchApplicationsCount();
    },
    methods: {
      async fetchApplicationsCount() {
        try {
          const response = await axios.get('http://localhost:8080/getCountThisWeek');
          this.applicationsCount = response.data.count;
        } catch (error) {
          console.error('Error fetching applications count this week:', error);
        }
      }
    },
    watch: {
      // Watch for changes in key prop to fetch data again
      '$route': {
        handler() {
          this.fetchApplicationsCount();
        },
        immediate: true
      }
    }
  };
  </script>
  
  <style scoped>
  .applications-this-week-widget {
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
  
  .count {
    font-size: 2rem;
    color: #42b983; /* Adjust color to fit your theme */
  }
  </style>
  