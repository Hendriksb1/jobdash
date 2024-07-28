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
        <select id="typeJob" v-model="newOpening.type_job" required>
          <option v-for="typeJob in jobTypes" :key="typeJob.id" :value="typeJob.type_name">{{ typeJob.type_name }}</option>
        </select>
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

export default {
  data() {
    return {
      newOpening: {
        firm: '',
        type_job: '',
        result: '',
        url: ''
      },
      resultTypes: [],
      jobTypes: []
    };
  },
  created() {
    this.fetchResultTypes();
    this.fetchJobTypes();
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
        const response = await axios.post('http://localhost:8080/addOpening', this.newOpening);
        this.$emit('opening-added', response.data);
        this.newOpening = {
          firm: '',
          type_job: '',
          result: '',
          url: ''
        };
      } catch (error) {
        console.error('Error adding opening:', error);
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
