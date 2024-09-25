<template>
    <div class="settings-page">
      <nav class="sidebar">
        <ul>
          <li><a href="#">Profile Settings</a></li>
          <li><a href="#">Account Settings</a></li>
          <li><a href="#">Notification Preferences</a></li>
          <li><a href="#">Job Openings</a></li>
          <li><a href="#">Logout</a></li>
        </ul>
      </nav>
      
      <div class="main-content">
        <h1 class="job-title">{{ job.firm }}</h1>
        <div class="job-details">
          <!-- <p><strong>Firm:</strong> {{ job.firm }}</p> -->
          <p><strong>Job Type:</strong> {{ job.type_name }}</p>
          <p><strong>Result:</strong> {{ job.result_name }}</p>
          <p><strong>Application Date:</strong> {{ formatApplicationDate(job.application_date) }}</p>
          <p><strong>Link: <a class="job-url" :href="job.url" target="_blank">{{ job.url }}</a></strong></p>
        </div>
        
        <!-- New Textarea for Comments/Notes -->
        <div class="comment-section">
          <h2>Comments/Notes</h2>
          <div class="comment-card">
            <textarea 
              id="comments" 
              v-model="comments" 
              placeholder="Add your comments here..." 
              rows="6"
            ></textarea>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  import { mapState } from 'vuex';
  
  export default {
    props: ['id'],
    data() {
      return {
        job: {},
        comments: '', // Store the comments
      };
    },
    created() {
      this.fetchJobDetails();
    },
    computed: {
      ...mapState({
        user_id: state => state.loggedInUser?.id
      })
    },
    methods: {
      async fetchJobDetails() {
        try {
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
  
  <style scoped>
  html, body {
    height: 100%;
    margin: 0;
  }
  
  .settings-page {
    display: flex;
    height: 100vh; /* Full height of the viewport */
  }
  
  .sidebar {
    width: 200px;
    background-color: #1b1b1b;
    padding: 20px;
    color: #e0e0e0;
    height: 100%; /* Full height of the parent container */
  }
  
  .sidebar ul {
    list-style-type: none;
    padding: 0;
    text-align: left; /* Left align sidebar items */
  }
  
  .sidebar li {
    margin: 10px 0;
  }
  
  .sidebar a {
    color: #e0e0e0;
    text-decoration: none;
    font-weight: bold;
  }
  
  .sidebar a:hover {
    text-decoration: underline;
  }
  
  .main-content {
    max-width: 50%;

    flex: 1;
    padding: 20px;
    margin: 0 auto;
    /* background-color: #1e1e1e; */
    color: #e0e0e0;
    overflow-y: auto; /* Allows scrolling if content overflows */
  }

  .job-title, .job-details {
    text-align: left; /* Left align job details */
  }
  
  .job-title {
    color: #42b983;
    margin-bottom: 25px;
  }
  
  .job-details {
    text-align: left; /* Left align job details */
  }
  
  .job-details p {
    margin: 10px 0;
  }
  
  .job-url {
    color: #42b983;
    text-decoration: none;
    font-weight: bold;
  }
  
  .job-url:hover {
    text-decoration: underline;
  }
  
  .comment-section {
    margin-top: 20px;
  }
  
  .comment-section h2 {
    margin-bottom: 10px;
    font-weight: bold;
  }
  
  .comment-card {
    /* background-color: #2a2a2a; Card background color */
    border: 1px solid #42b983; /* Card border color */
    border-radius: 5px; /* Rounded corners */
    padding: 10px;
  }
  
  textarea {
    width: 100%;
    padding: 10px;
    border: none; /* Remove default border */
    background-color: transparent; /* Transparent background to match card */
    color: #e0e0e0; /* Text color */
    resize: none; /* Prevent resizing */
    font-family: inherit; /* Inherit font family */
  }
  
  textarea::placeholder {
    color: #a0a0a0; /* Placeholder text color */
  }
  </style>
  