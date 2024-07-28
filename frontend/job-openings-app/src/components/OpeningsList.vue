<template>
  <div class="openings-list">
    <h2>Job Openings</h2>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Firm</th>
          <th>Job Type</th>
          <th>Result</th>
          <th>Application Date</th>
          <th>URL</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="opening in openings" :key="opening.id">
          <td>{{ opening.id }}</td>
          <td @click="setEditable(opening, 'firm')" :contenteditable="isEditable(opening, 'firm')" @input="updateField(opening, 'firm', $event)">
            {{ opening.firm }}
          </td>
          <td>
            <select v-if="isEditable(opening, 'type_job')" v-model="opening.type_job" @change="markModified(opening)">
              <option v-for="typeJob in jobTypes" :key="typeJob.id" :value="typeJob.type_name">{{ typeJob.type_name }}</option>
            </select>
            <span v-else @click="setEditable(opening, 'type_job')">{{ opening.type_job }}</span>
          </td>
          <td>
            <select v-if="isEditable(opening, 'result')" v-model="opening.result" @change="markModified(opening)">
              <option v-for="result in resultTypes" :key="result.id" :value="result.result_name">{{ result.result_name }}</option>
            </select>
            <span v-else @click="setEditable(opening, 'result')">{{ opening.result }}</span>
          </td>
          <td>{{ formatApplicationDate(opening.application_date) }}</td>
          <td @click="setEditable(opening, 'url')" :contenteditable="isEditable(opening, 'url')" @input="updateField(opening, 'url', $event)">
            <a :href="opening.url" target="_blank">{{ opening.url }}</a>
          </td>
          <td>
            <button @click="deleteOpening(opening.id)">Delete</button>
            <button @click="updateOpening(opening)" :disabled="!isModified(opening)">Update</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: ['newOpening'],
  data() {
    return {
      openings: [],
      editableFields: {},
      modifiedOpenings: {},
      resultTypes: [],
      jobTypes: []
    };
  },
  watch: {
    newOpening: {
      handler(newOpening) {
        if (newOpening) {
          this.openings.push(newOpening);
        }
      },
      deep: true
    }
  },
  created() {
    this.fetchOpenings();
    this.fetchResultTypes();
    this.fetchJobTypes();
  },
  methods: {
    async fetchOpenings() {
      try {
        const response = await axios.get('http://localhost:8080/getAllOpenings');
        this.openings = response.data;
      } catch (error) {
        console.error('Error fetching openings:', error);
      }
    },
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
    async deleteOpening(id) {
      try {
        const response = await axios.delete(`http://localhost:8080/deleteOpening/${id}`);
        console.log(response.data);
        this.openings = this.openings.filter(opening => opening.id !== id);
      } catch (error) {
        console.error('Error deleting opening:', error);
      }
    },
    setEditable(opening, field) {
      if (!this.editableFields[opening.id]) {
        this.editableFields[opening.id] = {};
      }
      this.editableFields[opening.id][field] = true;
    },
    isEditable(opening, field) {
      return this.editableFields[opening.id] && this.editableFields[opening.id][field];
    },
    updateField(opening, field, event) {
      const value = event.target.innerText;
      if (opening[field] !== value) {
        opening[field] = value;
        this.markModified(opening);
      }
    },
    markModified(opening) {
      this.modifiedOpenings[opening.id] = true;
    },
    isModified(opening) {
      return this.modifiedOpenings[opening.id];
    },
    async updateOpening(opening) {
      try {
        const response = await axios.put(`http://localhost:8080/updateOpening/${opening.id}`, opening);
        console.log(response.data);
        this.modifiedOpenings[opening.id] = false;
        this.editableFields[opening.id] = {};
      } catch (error) {
        console.error('Error updating opening:', error);
      }
    },
    formatApplicationDate(dateString) {
      // Assuming dateString is in ISO format, e.g., "2024-07-15T12:30:00Z"
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    }
  }
};
</script>

<style scoped>
.openings-list {
  background-color: #1e1e1e;
  padding: 20px;
  border-radius: 10px;
  margin: 20px;
  text-align: left;
}

h2 {
  color: #e0e0e0;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  color: #e0e0e0;
}

thead {
  background-color: #333;
}

th, td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #444;
}

th {
  background-color: #444;
  color: #e0e0e0;
}

tr:hover {
  background-color: #555;
}

a {
  color: #42b983;
}

button {
  padding: 5px 10px;
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #dc3545;
}

button[disabled] {
  background-color: #ccc;
  cursor: not-allowed;
}

input, select {
  background-color: #333;
  color: #e0e0e0;
  border: 1px solid #444;
  padding: 5px;
  border-radius: 4px;
}

input:focus, select:focus {
  outline: none;
  border-color: #555;
}
</style>
