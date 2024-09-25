<template>
  <div class="openings-list">
    <h2>Job Openings</h2>
    <table>
      <thead>
        <tr>
          <!-- <th @click="sortTable('id')">
            ID
            <span v-if="sortKey === 'id'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th> -->
          <th @click="sortTable('firm')">
            Firm
            <span v-if="sortKey === 'firm'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th @click="sortTable('type_job')">
            Job Type
            <span v-if="sortKey === 'type_job'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th @click="sortTable('result')">
            Result
            <span v-if="sortKey === 'result'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th @click="sortTable('application_date')">
            Application Date
            <span v-if="sortKey === 'application_date'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th @click="sortTable('url')">
            URL
            <span v-if="sortKey === 'url'" class="sort-indicator">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="opening in sortedOpenings" :key="opening.id"
          :class="getResultClass(opening.result)">
          <!-- <td>{{ opening.id }}</td> -->
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
            <button @click="navigateToDetail(opening.id)">More</button>
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
      jobTypes: [],
      sortKey: '',
      sortOrder: 'asc',
    };
  },
  computed: {
    sortedOpenings() {
      // Create a copy of the openings array to avoid mutating the original
      const sorted = [...this.openings];

      // Sort the copied array
      sorted.sort((a, b) => {
        let modifier = this.sortOrder === 'asc' ? 1 : -1;
        if (a[this.sortKey] < b[this.sortKey]) return -1 * modifier;
        if (a[this.sortKey] > b[this.sortKey]) return 1 * modifier;
        return 0;
      });

      return sorted;
    }
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
    console.log(this.$store.getters.user.id)
    this.fetchOpenings(this.$store.getters.user.id);
    this.fetchResultTypes();
    this.fetchJobTypes();
  },
  methods: {
    navigateToDetail(id) {
      console.log("emmiting event", id)
      this.$router.push({ name: 'job-detail', params: { id } });
      this.$emit('show-router-view', id); // Emit event with id
    },
    sortTable(key) {
      if (this.sortKey === key) {
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        this.sortKey = key;
        this.sortOrder = 'asc';
      }
    },
    async fetchOpenings(id) {
      try {
        const response = await axios.get(`http://localhost:8080/getAllOpenings/${id}`);
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
    externalLink(url) {
      window.open(url, '_blank');
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
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    },
    getResultClass(result) {
      switch (result.toLowerCase()) {
        case 'rejected':
          return 'rejected-row';
        case 'pending':
          return 'pending-row';
        case 'interview':
          return 'interview-row';
        case 'ghosted':
          return 'ghosted-row';
        default:
          return '';
      }
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
  overflow-x: auto; /* Enables horizontal scrolling if the table overflows */
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  color: #e0e0e0;
  table-layout: fixed; /* Ensures that the table doesn't exceed 100% width */
}

thead {
  background-color: #333;
}

th, td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #444;
  white-space: nowrap; /* Prevents text wrapping */
  overflow: hidden;
  text-overflow: ellipsis; /* Shows ellipsis if content overflows */
}

th {
  background-color: #444;
  color: #e0e0e0;
  cursor: pointer;
  position: relative;
}

th:hover {
  background-color: #555;
}

td {
  max-width: 250px; /* Sets a max width for table cells */
}

.sort-indicator {
  margin-left: 8px;
  font-size: 0.8em;
  color: #e0e0e0;
}

.rejected-row, .ghosted-row {
  background-color: #d6646459; /* Light red background */
}

.interview-row {
  background-color: #5cb85c59; /* Light green background */
}

tr:hover {
  background-color: #555;
}

a {
  color: #42b983;
}

button {
  padding: 5px 10px;
  background-color: #c84138;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 5px;
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
