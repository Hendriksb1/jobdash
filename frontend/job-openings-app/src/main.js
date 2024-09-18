import { createApp } from 'vue';
import App from './App.vue';
import store from './store'; // Ensure this path is correct

const app = createApp(App);
app.use(store);
app.mount('#app');