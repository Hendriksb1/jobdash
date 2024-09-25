<template>
    <div>
      <h1>Job Opening: {{ job.id }}</h1>
      <p>Firm: {{ job.firm }}</p>
      <p>Job Type: {{ job.type_name }}</p> <!-- Ensure property names match your response -->
      <p>Result: {{ job.result_name }}</p> <!-- Ensure property names match your response -->
      <p>Application Date: {{ formatApplicationDate(job.application_date) }}</p>
      <a :href="job.url" target="_blank">{{ job.url }}</a>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  import { mapState } from 'vuex'; // Import mapState helper
  
  export default {
    props: ['id'], // Receiving the job id as a prop
    data() {
      return {
        job: {},
      };
    },
    created() {
      this.fetchJobDetails();
    },
    computed: {
      ...mapState({
        user_id: state => state.loggedInUser?.id // Map the user ID from Vuex
      })
    },
    methods: {
      async fetchJobDetails() {
        try {
          // Use a GET request with URL parameters
          const response = await axios.get(`http://localhost:8080/getOpening`, {
            params: {
              userId: this.user_id,
              jobId: this.id
            }
          });
          this.job = response.data;
        } catch (error) {
          console.error('Error fetching job details:', error);
        }
      },
      formatApplicationDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
      }
    }
  };
  </script>