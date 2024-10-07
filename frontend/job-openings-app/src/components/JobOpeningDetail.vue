<template>
  <div class="settings-page">
    <div class="main-content">
      <h1 class="opening-title">{{ opening.firm }}</h1>
      <div class="opening-details">
        <p><strong>Opening Type:</strong> {{ opening.type_job }}</p>
        <p><strong>Result:</strong> {{ opening.result }}</p>
        <p><strong>Application Date:</strong> {{ formatApplicationDate(opening.application_date) }}</p>
        <p><strong>Link: <a class="opening-url" 
           :href="opening.url"
           target="_blank"
           >{{ opening.url }}</a>
        </strong></p>
      </div>

      <!-- New Textarea for Comment/Notes -->
      <div class="comment-section">
        <h2>Comment/Notes</h2>
        <div class="comment-card">
          <textarea 
              id="comment" 
              v-model="commentValue"
              placeholder="Add your comment here..." 
              rows="6"
          ></textarea>
        </div>
        <button @click="updateOpening(opening)">Update</button>
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
      opening: {
        comment: {
          String: ''
        }
      }, // Store the fetched opening details and ensure comment is initialized
    };
  },
  created() {
    this.fetchOpeningDetails();
  },
  computed: {
    ...mapState({
      user_id: state => state.loggedInUser?.id
    }),
    commentValue: {
      get() {
        return this.opening.comment ? this.opening.comment.String : ''; // Safely access comment.String
      },
      set(value) {
        if (this.opening.comment) {
          this.opening.comment.String = value;
        }
      }
    }
  },
  methods: {
    async fetchOpeningDetails() {
      console.log("fetch opening details");
      try {
        const response = await axios.get(`http://localhost:8080/getOpening`, {
          params: {
            userId: this.user_id,
            openingId: this.id
          }
        });
        console.log(response);
        this.opening = response.data;

        // Initialize comment if it does not exist
        if (!this.opening.comment) {
          this.$set(this.opening, 'comment', { String: '' });
        } else if (!this.opening.comment.Valid) {
          this.opening.comment.String = "";
        }

      } catch (error) {
        console.error('Error fetching opening details:', error);
      }
    },
    formatApplicationDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    },
    async updateOpening(opening) {
      try {
        const updatedOpening = {
          ...opening,
        };

        console.log('Updating with:', updatedOpening); // Log to check the payload
        const response = await axios.put(`http://localhost:8080/updateOpening/${opening.id}`, updatedOpening);
        console.log(response.data);
      } catch (error) {
        console.error('Error updating opening:', error);
      }
    }
  }
};
</script>

<style scoped>
html, body {
  height: 100%;
  margin: 0;
}

.main-content {
  max-width: 50%;
  flex: 1;
  padding: 20px;
  margin: 0 auto;
  color: #e0e0e0;
  overflow-y: auto;
}

.opening-title {
  color: #42b983;
  margin-bottom: 25px;
}

.opening-details {
  text-align: left;
}

.opening-details p {
  margin: 10px 0;
}

.opening-url {
  color: #42b983;
  text-decoration: none;
  font-weight: bold;
}

.opening-url:hover {
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
  border: 1px solid #42b983;
  border-radius: 5px;
  padding: 10px;
}

textarea {
  width: 100%;
  padding: 10px;
  border: none;
  background-color: transparent;
  color: #e0e0e0;
  resize: none;
  font-family: inherit;
}

textarea::placeholder {
  color: #a0a0a0;
}
</style>