import { createRouter, createWebHistory } from 'vue-router';
import OpeningsList from './components/OpeningsList.vue';
import JobOpeningDetail from './components/JobOpeningDetail.vue'; // Create this component

const routes = [
  {
    path: '/',
    name: 'openings',
    component: OpeningsList
  },
  {
    path: '/job/:id',
    name: 'job-detail',
    component: JobOpeningDetail,
    props: true // To pass route params as props
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;