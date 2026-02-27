import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

// Import ALL Material Web components globally
import '@material/web/all.js';

// M3 theme tokens
import './assets/styles/m3-theme.css';

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.mount('#app');
