<template>
  <div class="add-opening">
    <h2>Add New Opening</h2>
    <form @submit.prevent="handleSubmit" class="opening-form">
      <div class="form-row">
        <label for="firm">Firm:</label>
        <input type="text" id="firm" v-model="newOpening.firm" required>
      </div>

      <div class="form-row">
        <label for="typeJob">Job Type:</label>
        <select id="typeJob" v-model="newOpening.type_job" @change="checkNewJobType" required>
          <option v-for="typeJob in jobTypes" :key="typeJob.id" :value="typeJob.type_name">{{ typeJob.type_name }}</option>
          <option value="add-new">Add New Job Type</option>
        </select>
      </div>

      <div v-if="newOpening.type_job === 'add-new'" class="form-row">
        <label for="newJobType">New Job Type:</label>
        <input type="text" id="newJobType" v-model="newOpening.newJobType" required>
      </div>

      <div class="form-row">
        <label for="result">Result:</label>
        <select id="result" v-model="newOpening.result" required>
          <option v-for="result in resultTypes" :key="result.id" :value="result.result_name">{{ result.result_name }}</option>
        </select>
      </div>

      <div class="form-row">
        <label for="url">URL:</label>
        <input type="text" id="url" v-model="newOpening.url">
      </div>

      <div class="form-row">
        <button type="submit">Add Opening</button>
      </div>
    </form>
  </div>
</template>

<script>
import axios from 'axios';
import { mapState } from 'vuex'; // Import mapState helper

export default {
  data() {
    return {
      newOpening: {
        firm: '',
        type_job: '',
        newJobType: '', // Additional field for new job type
        result: '',
        url: '',
        user_id: this.user_id // Initialize user_id from Vuex
      },
      resultTypes: [],
      jobTypes: []
    };
  },
  computed: {
    ...mapState({
      user_id: state => state.loggedInUser?.id // Map the user ID from Vuex
    })
  },
  watch: {
    user_id(newUserId) {
      // Watch for changes in user_id and update newOpening accordingly
      this.newOpening.user_id = newUserId;
    }
  },
  created() {
    this.fetchResultTypes();
    this.fetchJobTypes();
    this.newOpening.user_id = this.user_id; // Assign user_id from Vuex to the form
  },
  methods: {
    async fetchResultTypes() {
      try {
        const response = await axios.get('http://localhost:8080/getAllResultTypes');
        this.resultTypes = response.data;
      } catch (error) {
        console.error('Error fetching result types:', error);
      }
    },
    async fetchJobTypes() {
      try {
        const response = await axios.get('http://localhost:8080/getAllJobTypes');
        this.jobTypes = response.data;
      } catch (error) {
        console.error('Error fetching job types:', error);
      }
    },
    async handleSubmit() {
      try {
        if (this.newOpening.type_job === 'add-new') {
          // Add the new job type
          await axios.post('http://localhost:8080/addJobType', { name: this.newOpening.newJobType });
          this.newOpening.type_job = this.newOpening.newJobType; // Set the job type to the new one
        }

        // Add the opening
        const response = await axios.post('http://localhost:8080/addOpening', this.newOpening);
        this.$emit('opening-added', response.data);
        // Reset the form, except for user_id
        this.newOpening = {
          firm: '',
          type_job: '',
          newJobType: '', // Reset this as well
          result: '',
          url: '',
          user_id: this.user_id // Keep the user_id from Vuex
        };
      } catch (error) {
        console.error('Error adding opening:', error);
      }
    },
    checkNewJobType() {
      if (this.newOpening.type_job !== 'add-new') {
        this.newOpening.newJobType = ''; // Reset newJobType if not adding a new one
      }
    }
  }
};
</script>


<style scoped>
.add-opening {
  background-color: #1e1e1e;
  padding: 20px;
  border-radius: 10px;
  margin: 20px;
  text-align: left;
}

h2 {
  color: #e0e0e0;
}

.opening-form {
  display: flex;
  flex-wrap: wrap;
}

.form-row {
  flex: 1 0 300px; /* Adjust width as needed */
  margin-right: 10px;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}

.form-row label {
  min-width: 100px;
  margin-right: 10px;
  color: #e0e0e0;
}

input, select {
  flex: 1;
  padding: 5px;
  background-color: #333;
  color: #e0e0e0;
  border: 1px solid #444;
  border-radius: 4px;
}

input:focus, select:focus {
  outline: none;
  border-color: #555;
}

button {
  padding: 8px 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  cursor: pointer;
  border-radius: 4px;
}

button:hover {
  background-color: #45a049;
}
</style>
